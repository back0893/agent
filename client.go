package main

import (
	"agent/src"
	"agent/src/agent/funcs"
	"bytes"
	"context"
	"encoding/gob"
	"github.com/back0893/goTcp/iface"
	net2 "github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func EncodeData(e interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(e); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}
func SendHeart(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.PING
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendCPU(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.CPU
	cpu := funcs.CpuMetrics()

	data, err := EncodeData(cpu)
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data = data
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendHHD(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.HHD

	disks, err := funcs.DiskUseMetrics()
	if err != nil {
		log.Print(err)
		return
	}

	pkt.Data, err = EncodeData(disks)

	if err != nil {
		log.Println(err)
		return
	}

	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendMem(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.MEM
	memory, err := funcs.MemMetrics()
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data, err = EncodeData(memory)
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendLoadAvg(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.LoadAvg
	loadAvg, err := funcs.LoadAvgMetrics()
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data, err = EncodeData(loadAvg)
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendPort(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.PortListen
	loadAvg, err := funcs.ListenTcpPortMetrics(80, 10086)
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data, err = EncodeData(loadAvg)
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}

type Agent struct {
	con       iface.IConnection
	conEvent  iface.IEventWatch
	protocol  iface.IProtocol
	ctx       context.Context
	ctxCancel context.CancelFunc
	isStop    *src.AtomicInt64
	wg        *sync.WaitGroup
}

type AgentEvent struct{}

func (a AgentEvent) OnConnect(ctx context.Context, connection iface.IConnection) {
	log.Println("接连成功时")
}

func (a AgentEvent) OnMessage(ctx context.Context, packet iface.IPacket, connection iface.IConnection) {
	pkt := packet.(*src.Packet)
	log.Println("接受的回应id=>", pkt.Id)
}

func (a AgentEvent) OnClose(ctx context.Context, connection iface.IConnection) {
	log.Println("接连关闭")
}

func (a *Agent) AddEvent(event iface.IEvent) {
	a.conEvent.AddConnect(event.OnConnect)
	a.conEvent.AddMessage(event.OnMessage)
	a.conEvent.AddClose(event.OnClose)

}
func (a *Agent) AddProtocol(protocol iface.IProtocol) {
	a.protocol = protocol
}
func (a *Agent) Start() {
	go a.con.Run()
}
func (a *Agent) IsStop() bool {
	return a.isStop.Get() == 1
}
func (a *Agent) Stop() {
	if a.IsStop() {
		a.isStop.Store(0)
		a.con.Close()
		a.wg.Wait()
	}
}
func NewAgent(con *net.TCPConn, event iface.IEvent, protocol iface.IProtocol) *Agent {
	agent := &Agent{
		isStop:   src.NewAtomicInt64(1),
		conEvent: net2.NewEventWatch(),
		wg:       &sync.WaitGroup{},
	}
	agent.AddProtocol(protocol)
	agent.AddEvent(event)
	agent.ctx, agent.ctxCancel = context.WithCancel(context.Background())

	agent.con = net2.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent
}
func main() {
	utils.GlobalConfig.Load("json", "./client.json")
	host := net.JoinHostPort(utils.GlobalConfig.GetString("ip"), utils.GlobalConfig.GetString("port"))
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		panic(err)
	}
	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		panic(err)
	}

	agent := NewAgent(con, AgentEvent{}, src.Protocol{})
	src.InitTimingWheel(agent.ctx)

	//心跳单独实现.
	src.AddTimer(2*time.Second, func() {
		SendHeart(agent.con)
	})
	//目前定时汇报cpu,内存,硬盘使用情况
	src.AddTimer(3*time.Second, func() {
		SendCPU(agent.con)
		SendMem(agent.con)
		SendHHD(agent.con)
		SendLoadAvg(agent.con)
	})

	//todo agent的断线重连

	go agent.Start()

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	agent.Stop()
}
