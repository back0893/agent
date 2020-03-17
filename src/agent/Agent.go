package agent

import (
	"agent/src"
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Agent struct {
	con          iface.IConnection
	conEvent     iface.IEventWatch
	protocol     iface.IProtocol
	ctx          context.Context
	ctxCancel    context.CancelFunc
	isStop       *src.AtomicInt64
	wg           *sync.WaitGroup
	cfg          string //配置文件的路径
	servicesList *ServicesList
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
	con, err := ConnectServer(cfg)
	if err != nil {
		return nil, err
	}
	agent := &Agent{
		isStop:       src.NewAtomicInt64(0),
		conEvent:     net.NewEventWatch(),
		wg:           &sync.WaitGroup{},
		cfg:          cfg,
		servicesList: NewServicesList(),
	}

	agent.ctx, agent.ctxCancel = context.WithCancel(context.WithValue(context.Background(), g.AGENT, agent))
	agent.AddProtocol(src.Protocol{})
	agent.AddEvent(Event{})

	//断线重连
	agent.AddClose(agent.ReCon)

	//初始化定时器
	src.InitTimingWheel(agent.GetContext())

	//初始化services列表
	agent.servicesList.WakeUp()
	//新增在结束时,保存
	agent.AddClose(func(ctx context.Context, connection iface.IConnection) {
		ctx.Value(g.AGENT).(*Agent).servicesList.Sleep()
	})
	//新增一个服务的定时监控
	src.AddTimer(5*time.Second, agent.servicesList.Listen)

	agent.con = net.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent, nil
}
