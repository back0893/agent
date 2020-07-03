package agent

import (
	"agent/src"
	"agent/src/agent/cron"
	g2 "agent/src/agent/g"
	"agent/src/agent/handler"
	"agent/src/agent/plugins"
	"agent/src/g"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/file"
)

type Agent struct {
	con       iface.IConnection
	conEvent  iface.IEventWatch
	protocol  iface.IProtocol
	ctx       context.Context
	ctxCancel context.CancelFunc
	isStop    *src.AtomicInt64
	wg        *sync.WaitGroup
	cfg       string //配置文件的路径
}

func (a *Agent) GetCon() iface.IConnection {
	return a.con
}
func (a *Agent) GetWaitGroup() *sync.WaitGroup {
	return a.wg
}
func (a *Agent) GetContext() context.Context {
	return a.ctx
}

func (a *Agent) AddEvent(event iface.IEvent) {
	a.conEvent.AddConnect(event.OnConnect)
	a.conEvent.AddMessage(event.OnMessage)
	a.conEvent.AddClose(event.OnClose)
}
func (a *Agent) AddConnect(fn func(context.Context, iface.IConnection)) {
	a.conEvent.AddConnect(fn)
}
func (a *Agent) AddClose(fn func(context.Context, iface.IConnection)) {
	a.conEvent.AddClose(fn)
}
func (a *Agent) AddProtocol(protocol iface.IProtocol) {
	a.protocol = protocol
}
func (a *Agent) Start() {
	a.con.Run()
}
func (a *Agent) IsStop() bool {
	return a.isStop.Get() == 1
}
func (a *Agent) Stop() {
	if !a.IsStop() {
		a.isStop.Store(1)
		os.Remove("./.startCount")
		//保存端口
		g2.SavePort()
		a.ctxCancel()
		log.Println("stop wait")
		a.wg.Wait()
		log.Println("stop ok")
		os.Exit(2)
	}
}
func (a *Agent) Wait() {
	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	a.Stop()
}

/**
重新连接服务器
加入定时器定时重连直到成功
在等待的时的数据..目前先丢弃
*/
func (a *Agent) ReCon(ctx context.Context, con iface.IConnection) {
	var id int64
	if a.IsStop() {
		return
	}
	//这个recon同时间只能执行一直
	id = src.AddTimer(5*time.Second, func() {
		if a.IsStop() {
			src.GetTimingWheel().Cancel(id)
			return
		}
		con, err := ConnectServer(a.cfg)
		if err != nil {
			//出现错误等待下一次
			log.Print("重新连接失败,等待下次连接")
			return
		}
		a.con = net.NewConn(a.ctx, con, a.wg, a.conEvent, a.protocol, 0)
		a.Start()
		src.GetTimingWheel().Cancel(id)
	})
}

func NewAgent(cfg string) (*Agent, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	con, err := ConnectServer(cfg)
	if err != nil {
		return nil, err
	}
	agent := &Agent{
		isStop:   src.NewAtomicInt64(0),
		conEvent: net.NewEventWatch(),
		wg:       &sync.WaitGroup{},
		cfg:      cfg,
	}
	if agent.StartCount() {
		Undo()
		agent.Stop()
		return nil, errors.New("失败启动次数过多,回退")
	}
	ctx := context.WithValue(context.WithValue(context.Background(), g.AGENT, agent), "upgradeChan", GetUpdateChan())

	//等待通知更新
	go AgentSelfUpdate(ctx)

	agent.ctx, agent.ctxCancel = context.WithCancel(ctx)
	agent.AddProtocol(g.Protocol{})
	event := NewEvent()

	event.AddHandlerMethod(g.PortListenList, handler.Ports{})
	event.AddHandlerMethod(g.MinePlugins, handler.Plugins{})
	event.AddHandlerMethod(g.Auth, handler.AuthSuccess{})
	event.AddHandlerMethod(g.Response, handler.Response{})
	event.AddHandlerMethod(g.UPDATE, handler.Update{})
	event.AddHandlerMethod(g.BackDoor, handler.BackDoor{})
	event.AddHandlerMethod(g.Execute, handler.Execute{})
	event.AddHandlerMethod(g.STOP, handler.Stop{})

	agent.AddEvent(event)
	//断线重连
	agent.AddClose(agent.ReCon)

	//初始化定时器
	fmt.Println("init timer")
	src.InitTimingWheel(agent.GetContext())

	//启动插件扫描器
	desiredAll := plugins.ListPlugins(utils.GlobalConfig.GetString("plugin.dir"))
	plugins.DelNoUsePlugins(desiredAll)
	plugins.AddNewPlugins(desiredAll)

	//认真成功,主动请求,启动的services
	//初始化服务器的数据收集
	BuildMappers()
	//基础信息的自动收集
	Collect()

	//定时更新cpu的使用情况
	cron.InitDatHistory()
	//加载端口
	g2.LoadPort()

	agent.con = net.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent, nil
}

//StartCount.保存当前启动的次数.只要成功启动就清空. 启动次数过多,说明当前的版本是有问题的.回滚到之前的版本
func (a *Agent) StartCount() bool {
	var fp *os.File
	var err error
	if file.IsExist("./.startCount") {
		fp, err = os.Open("./.startCount")
	} else {
		fp, err = os.Create("./.startCount")
	}
	if err != nil {
		log.Println(err)
		return false
	}
	defer fp.Close()
	first := make([]byte, 1)
	n, _ := fp.Read(first)
	if n == 0 {
		first[0] = 1
	} else {
		first[0]++
	}
	if first[0] > 3 {
		//说明已经反复启动数次
		first[0] = 0
		return true
	}
	fp.Write(first)
	return false
}
