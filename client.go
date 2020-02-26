package main

import (
	"agent/src"
	"bufio"
	"bytes"
	"encoding/binary"
	"github.com/back0893/goTcp/iface"
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

type connection struct {
	rawCon net.Conn
	buffer *bufio.Reader
}

func (con connection) GetId() uint32 {
	return 0
}

func (con *connection) GetRawCon() net.Conn {
	return con.rawCon
}

func (con *connection) GetBuffer() *bufio.Reader {
	return con.buffer
}

func (con *connection) Write(p iface.IPacket) error {
	if data, err := p.Serialize(); err != nil {
		return err
	} else {
		if _, err := con.rawCon.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func (con *connection) AsyncWrite(p iface.IPacket, timeout time.Duration) error {
	return con.Write(p)
}

func (con connection) GetExtraData(key interface{}) (interface{}, bool) {
	return nil, false
}

func (con connection) SetExtraData(key interface{}, value interface{}) {

}

func (con connection) GetExtraMap() *sync.Map {
	return nil
}

func (con *connection) Close() {
	con.rawCon.Close()
}

func (con connection) IsClosed() bool {
	return false
}

func SendHeart(conn iface.IConnection) {
	interval := utils.GlobalConfig.GetInt("heartInterval")
	ticker := time.NewTicker(time.Duration(interval) * time.Second)
	pkt := src.NewPkt()
	pkt.Id = src.PING
	for {
		select {
		case <-ticker.C:
			{
				if err := conn.Write(pkt); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
func SendCPU(conn iface.IConnection) {
	ticker := time.NewTicker(time.Second * 1)
	pkt := src.NewPkt()
	pkt.Id = src.CPU
	for {
		select {
		case <-ticker.C:
			{
				buffer := bytes.NewBuffer([]byte{})
				num := nux.NumCpu()
				binary.Write(buffer, binary.BigEndian, int32(num))
				pkt.Data = buffer.Bytes()
				if err := conn.Write(pkt); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
func SendHHD(conn iface.IConnection) {
	ticker := time.NewTicker(time.Second * 2)
	pkt := src.NewPkt()
	pkt.Id = src.HHD
	for {
		select {
		case <-ticker.C:
			{
				buffer := bytes.NewBuffer([]byte{})
				disks, err := nux.ListDiskStats()
				if err != nil {
					log.Println(err)
					continue
				}
				diskName := make([]string, 0)
				for _, disk := range disks {
					diskName = append(diskName, disk.Device)
				}
				str := strings.Join(diskName, "\n")
				buffer.WriteString(str)
				if err := conn.Write(pkt); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
func SendMem(conn iface.IConnection) {
	ticker := time.NewTicker(time.Second * 3)
	pkt := src.NewPkt()
	pkt.Id = src.MEM
	for {
		select {
		case <-ticker.C:
			{
				buffer := bytes.NewBuffer([]byte{})
				info, err := nux.MemInfo()
				if err != nil {
					log.Println(err)
					break
				}
				binary.Write(buffer, binary.BigEndian, info.MemTotal)
				binary.Write(buffer, binary.BigEndian, info.MemTotal-info.MemFree)
				pkt.Data = buffer.Bytes()
				if err := conn.Write(pkt); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
func Rev(conn iface.IConnection) {
	protocol := src.Protocol{}
	for {
		packet, err := protocol.UnPack(conn)
		if err != nil {
			log.Panic(err)
		}
		pkt := packet.(*src.Packet)
		log.Println("接受的回应id=>", pkt.Id)

	}
}

func main() {
	utils.GlobalConfig.Load("json", "./client.json")
	host := net.JoinHostPort(utils.GlobalConfig.GetString("ip"), utils.GlobalConfig.GetString("port"))
	con, err := net.Dial("tcp", host)
	if err != nil {
		panic(err)
	}
	conn := &connection{
		rawCon: con,
		buffer: bufio.NewReader(con),
	}

	//心跳单独实现.
	//go SendHeart(conn)
	//目前定时汇报cpu,内存,硬盘使用情况
	go SendCPU(conn)
	go SendMem(conn)

	//读取服务端发送来的消息
	go Rev(conn)

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)

}
