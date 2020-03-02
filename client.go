package main

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/agent/model"
	"bytes"
	"context"
	"encoding/gob"
	"flag"
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

const (
	AGENT string = "agent"
)

var (
	cfg string
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
	//更新cpu状态
	funcs.UpdateCpuStat()
	if funcs.CpuPrepared() == false {
		//如果cpu状态还未准备好久不发送
		return
	}

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
	listenPorts := utils.GlobalConfig.GetIntSlice("listenPort")
	lp := make([]int64, 0)
	for _, val := range listenPorts {
		lp = append(lp, int64(val))
	}
	ports, err := funcs.ListenTcpPortMetrics(lp...)
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data, err = EncodeData(ports)
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
	//这个时候发送身份识别
	pkt := src.NewPkt()
	pkt.Id = src.Auth
	authModel := model.Auth{
		Username: utils.GlobalConfig.GetString("username"),
		Password: utils.GlobalConfig.GetString("password"),
	}
	pkt.Data, _ = EncodeData(authModel)
	if err := connection.Write(pkt); err != nil {
		log.Println(err)
	}
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
		con, err := ConnectServer()
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
func NewAgent(con *net.TCPConn, event iface.IEvent, protocol iface.IProtocol) *Agent {
	agent := &Agent{
		isStop:   src.NewAtomicInt64(0),
		conEvent: net2.NewEventWatch(),
		wg:       &sync.WaitGroup{},
	}
	agent.AddProtocol(protocol)
	agent.AddEvent(event)
	agent.ctx, agent.ctxCancel = context.WithCancel(context.Background())

	agent.con = net2.NewConn(agent.ctx, con, agent.wg, agent.conEvent, agent.protocol, 0)

	return agent
}

func ConnectServer() (*net.TCPConn, error) {
	utils.GlobalConfig.Load("json", cfg)
	host := net.JoinHostPort(utils.GlobalConfig.GetString("ip"), utils.GlobalConfig.GetString("port"))
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}
	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return con, nil
}
func main() {
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置json")
	flag.Parse()

	con, err := ConnectServer()
	if err != nil {
		panic(err)
	}
	agent := NewAgent(con, AgentEvent{}, src.Protocol{})

	//断线重连
	agent.AddClose(agent.ReCon)

	src.InitTimingWheel(agent.ctx)

	//心跳单独实现.
	headrtBeat := utils.GlobalConfig.GetInt("heartBeat")
	src.AddTimer(time.Second*time.Duration(headrtBeat), func() {
		SendHeart(agent.con)
	})
	//目前定时汇报cpu,内存,硬盘使用情况
	src.AddTimer(5*time.Second, func() {
		SendCPU(agent.con)
		SendMem(agent.con)
		SendHHD(agent.con)
		SendLoadAvg(agent.con)
		SendPort(agent.con)
	})

	agent.Start()

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	agent.Stop()
}
