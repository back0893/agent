package agent

import (
	"agent/src"
	iface2 "agent/src/agent/iface"
	"agent/src/agent/model"
	"agent/src/agent/services"
	"agent/src/g"
	"context"
	"github.com/back0893/goTcp/iface"
	net2 "github.com/back0893/goTcp/net"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

type Agent struct {
	con       iface.IConnection
	conEvent  iface.IEventWatch
	protocol  iface.IProtocol
	ctx       context.Context
	ctxCancel context.CancelFunc
	isStop    *src.AtomicInt64
	wg        *sync.WaitGroup
	taskQueue *src.TaskQueue
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
		a.ctxCancel()
		log.Println("stop wait")
		a.wg.Wait()
		log.Println("stop ok")
		os.Exit(2)
	}
}
func (a *Agent) RunTask() {
	//读取taskQueue,执行相应的操作
	go func() {
		var service iface2.IService
		var task *model.Service
		for {
			task = a.taskQueue.Pop()
			switch strings.ToLower(task.Service) {
			case "redis":
				service = services.NewRedisService(a)
			case "heart":
				service = services.NewHeartBeatService(a)
			case "loadavg":
				service = services.NewLoadAvgServiceService(a)
			case "memory":
				service = services.NewMemoryService(a)
			case "hhd":
				service = services.NewHHDService(a)
			case "port":
				service = services.NewPortService(a)
			case "cpu":
				service = services.NewCPUService(a)
			default:
				continue
			}
			service.Action(task.Action, task.Args)
		}

	}()
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
		a.con = net2.NewConn(a.ctx, con, a.wg, a.conEvent, a.protocol, 0)
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
		isStop:    src.NewAtomicInt64(0),
		conEvent:  net2.NewEventWatch(),
		wg:        &sync.WaitGroup{},
		taskQueue: src.NewTaskQueue(),
		cfg:       cfg,
	}

	agent.ctx, agent.ctxCancel = context.WithCancel(context.WithValue(context.Background(), g.AGENT, agent))
	agent.AddProtocol(src.Protocol{})
	agent.AddEvent(Event{})

	//断线重连
	agent.AddClose(agent.ReCon)

	src.InitTimingWheel(agent.GetContext())

	agent.con = net2.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent, nil
}
