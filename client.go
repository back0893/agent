package main

import (
	"agent/src"
	"bytes"
	"context"
	"encoding/binary"
	"github.com/back0893/goTcp/iface"
	net2 "github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/nux"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

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
	buffer := bytes.NewBuffer([]byte{})
	num := nux.NumCpu()
	binary.Write(buffer, binary.BigEndian, int32(num))
	pkt.Data = buffer.Bytes()
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendHHD(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.HHD
	buffer := bytes.NewBuffer([]byte{})
	disks, err := nux.ListDiskStats()
	if err != nil {
		log.Println(err)
		return
	}
	diskName := make([]string, 0)
	for _, disk := range disks {
		diskName = append(diskName, disk.Device)
	}
	str := strings.Join(diskName, "\n")
	buffer.WriteString(str)
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
func SendMem(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = src.MEM
	buffer := bytes.NewBuffer([]byte{})
	info, err := nux.MemInfo()
	if err != nil {
		log.Println(err)
		return
	}
	binary.Write(buffer, binary.BigEndian, info.MemTotal)
	binary.Write(buffer, binary.BigEndian, info.MemTotal-info.MemFree)
	pkt.Data = buffer.Bytes()
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
	tw        *src.TimingWheel
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
	log.Println("接连成功关闭")
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
	go a.tw.Start()
}
func (a *Agent) IsStop() bool {
	return a.isStop.Get() == 0
}
func (a *Agent) Stop() {
	if a.IsStop() {
		a.isStop.Store(1)
		a.con.Close()
		a.tw.Stop()
		a.wg.Wait()
	}
}
func (a *Agent) AddTimer(interval time.Duration, fn func()) int64 {
	return a.tw.AddTimer(time.Now(), interval, fn)
}
func (a *Agent) TimerAt(when time.Time, fn func()) int64 {
	return a.tw.AddTimer(when, 0, fn)
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
	agent.tw = src.NewTimingWheel(agent.ctx)

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

	//心跳单独实现.
	agent.AddTimer(2*time.Second, func() {
		SendHeart(agent.con)
	})
	//目前定时汇报cpu,内存,硬盘使用情况
	agent.AddTimer(3*time.Second, func() {
		SendCPU(agent.con)
	})
	agent.AddTimer(4*time.Second, func() {
		SendMem(agent.con)
	})

	go agent.Start()

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	agent.Stop()
}
