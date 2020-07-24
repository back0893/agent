package agent

import (
	"agent/src"
	"agent/src/agent/handler"
	"agent/src/g"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
)

type ControllerAgent struct {
	con       iface.IConnection
	conEvent  iface.IEventWatch
	protocol  iface.IProtocol
	ctx       context.Context
	ctxCancel context.CancelFunc
	isStop    *src.AtomicInt64
	wg        *sync.WaitGroup
	cfg       string //配置文件的路径
}

func (ca *ControllerAgent) GetCon() iface.IConnection {
	return ca.con
}

func (ca *ControllerAgent) GetWaitGroup() *sync.WaitGroup {
	return ca.wg
}

func (ca *ControllerAgent) GetContext() context.Context {
	return ca.ctx
}

func (ca *ControllerAgent) AddEvent(event iface.IEvent) {
	ca.conEvent.AddConnect(event.OnConnect)
	ca.conEvent.AddMessage(event.OnMessage)
	ca.conEvent.AddClose(event.OnClose)
}

func (ca *ControllerAgent) AddConnect(fn func(context.Context, iface.IConnection)) {
	ca.conEvent.AddConnect(fn)
}

func (ca *ControllerAgent) AddClose(fn func(context.Context, iface.IConnection)) {
	ca.conEvent.AddClose(fn)
}

func (ca *ControllerAgent) AddProtocol(protocol iface.IProtocol) {
	ca.protocol = protocol
}

func (ca *ControllerAgent) Start() {
	panic("not implemented") // TODO: Implement
}

func (ca *ControllerAgent) IsStop() bool {
	return ca.isStop.Get() == 1
}

func (ca *ControllerAgent) Stop() {
	if !ca.IsStop() {
		ca.isStop.Store(1)
		ca.con.Close()
		ca.ctxCancel()
		//平滑退出
		ca.wg.Wait()
		os.Exit(0)
	}
}

func (ca *ControllerAgent) Wait() {
	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	ca.Stop()
}

func (ca *ControllerAgent) ReCon(ctx context.Context, con iface.IConnection) {
	var id int64
	if ca.IsStop() {
		return
	}
	//这个recon同时间只能执行一直
	id = src.AddTimer(5*time.Second, func() {
		if ca.IsStop() {
			src.GetTimingWheel().Cancel(id)
			return
		}
		con, err := ConnectServer(ca.cfg)
		if err != nil {
			//出现错误等待下一次
			log.Print("重新连接失败,等待下次连接")
			return
		}
		ca.con = net.NewConn(ca.ctx, con, ca.wg, ca.conEvent, ca.protocol, 0)
		ca.Start()
		src.GetTimingWheel().Cancel(id)
	})
}
func (ca *ControllerAgent) ChildAgent() {
	//扫描一个文件
	//获得保存agent的文件夹
	//判断对应agent的文件夹存在.不存在尝试下载,下载失败,上报
	//进入agent对应的文件夹,获得应执行的版本.不存在尝试下载,下载失败,上报
	//尝试启动.失败次数过多,回退版本
	//启动http服务器,接收手工http消息启动,或者直接kill -9 退出子进程

}
func NewControllerAgent(cfg string) (*ControllerAgent, error) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	con, err := ConnectServer(cfg)
	if err != nil {
		return nil, err
	}
	agent := &ControllerAgent{
		isStop:   src.NewAtomicInt64(0),
		conEvent: net.NewEventWatch(),
		wg:       &sync.WaitGroup{},
		cfg:      cfg,
	}
	ctx := context.WithValue(context.WithValue(context.Background(), g.AGENT, agent), "upgradeChan", GetUpdateChan())

	//等待通知更新
	go AgentSelfUpdate(ctx)

	agent.ctx, agent.ctxCancel = context.WithCancel(ctx)
	agent.AddProtocol(g.Protocol{})
	event := NewEvent()

	event.AddHandlerMethod(g.UPDATE, handler.Update{})

	agent.AddEvent(event)
	//断线重连
	agent.AddClose(agent.ReCon)

	//初始化定时器
	fmt.Println("init timer")
	src.InitTimingWheel(agent.GetContext())

	agent.con = net.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent, nil
}
