package agent

import (
	"agent/src"
	"agent/src/agent/services"
	"agent/src/g"
	"context"
	"fmt"
	"github.com/back0893/goTcp/iface"
	net2 "github.com/back0893/goTcp/net"
	"log"
	"net"
	"sync"
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
	if a.IsStop() {
		a.isStop.Store(1)
		a.ctxCancel()
		a.wg.Wait()
	}
}
func (a *Agent) RunTask() {
	//读取taskQueue,执行相应的操作
	go func() {
		for {
			service := a.taskQueue.Pop()
			pkt := src.NewPkt()
			pkt.Id = g.ServiceResponse
			var str = "未知命令"
			switch service.Service {
			case "redis":
				redis := services.NewRedisService()
				switch service.Action {
				case "start":
					err := redis.Start()
					if err != nil {
						str = "redis启动失败"
					} else {
						str = "redis启动成功"
					}
				case "stop":
					err := redis.Stop()
					if err != nil {
						str = "redis停止失败"
					} else {
						str = "redis停止成功"
					}
				case "status":
					status := redis.Status()
					if status {
						str = "redis正在运行"
					} else {
						str = "redis没有运行"
					}
				case "restart":
					err := redis.Restart()
					if err != nil {
						str = "redis重启失败"
					} else {
						str = "redis重启成功"
					}
				}
			default:
			}
			fmt.Println(str)
			pkt.Data = []byte(str)
			err := a.con.Write(pkt)
			if err != nil {
				//todo 发送失败..应该有后续操作
			}
		}
	}()
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

func NewAgent(con *net.TCPConn, event iface.IEvent, protocol iface.IProtocol, cfg string) *Agent {
	agent := &Agent{
		isStop:    src.NewAtomicInt64(0),
		conEvent:  net2.NewEventWatch(),
		wg:        &sync.WaitGroup{},
		taskQueue: src.NewTaskQueue(),
		cfg:       cfg,
	}
	agent.AddProtocol(protocol)
	agent.AddEvent(event)
	agent.ctx, agent.ctxCancel = context.WithCancel(context.WithValue(context.Background(), g.AGENT, agent))

	agent.con = net2.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent
}
