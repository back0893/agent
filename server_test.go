package main

import (
	"agent/src/agent/plugins"
	"agent/src/g"
	"agent/src/g/model"
	"bufio"
	"encoding/binary"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
	"net"
	"testing"
)

func read(buffer *bufio.Reader) (*g.Packet, error) {
	pkt := g.NewPkt()
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Length); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Version); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Timestamp); err != nil {
		return nil, err
	}
	if err := binary.Read(buffer, binary.BigEndian, &pkt.Id); err != nil {
		return nil, err
	}
	length := pkt.Length - 21
	pkt.Data = make([]byte, length)
	if err := binary.Read(buffer, binary.BigEndian, pkt.Data); err != nil {
		return nil, err
	}
	return pkt, nil
}
func newConn() *net.TCPConn {
	add, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8001")
	if err != nil {
		panic(err)
	}
	con, err := net.DialTCP("tcp", nil, add)
	if err != nil {
		panic(err)
	}
	return con
}
func TestServerCon(t *testing.T) {
	con := newConn()
	pkt := g.NewPkt()
	pkt.Id = g.MinePlugins
	data, _ := pkt.Serialize()
	log.Println("sned")
	con.Write(data)
	log.Println("read")
	res, err := read(bufio.NewReader(con))
	if err != nil {
		panic(err)
	}
	repPlugins := model.Plugins{}
	fmt.Println(res.Id, res.Length, len(res.Data))
	if err := g.DecodeData(res.Data, &repPlugins); err != nil {
		panic(err)
	}
}
func TestGit(t *testing.T) {
	g.LoadInit("./client.json")
	plugin := model.Plugins{
		Uri: []string{"git@gitee.com:liuzy1988/SealKit.git"},
	}
	plugins.Git(utils.GlobalConfig.GetString("plugin.dir"), &plugin)
	for str, schedule := range plugins.PluginsWithScheduler {
		fmt.Println(str, *schedule.Plugin)
	}
}
func TestReadList(t *testing.T) {
	g.LoadInit("./client.json")

	log.Println("read list")
	desiredAll := plugins.ListPlugins("")
	for key, value := range desiredAll {
		fmt.Println(key, value.FilePath)
	}
}
