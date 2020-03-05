package cron

import (
	"agent/src"
	"agent/src/agent/funcs"
	"agent/src/g"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/utils"
	"log"
)

func SendHeart(conn iface.IConnection) {
	pkt := src.NewPkt()
	pkt.Id = g.PING
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
	pkt.Id = g.CPU
	//这里cpu的范围区间才能计算,所以需要一个定时器来定时查询
	src.AddTimer(5, func() {
		_ = funcs.UpdateCpuStat()
	})
	cpu := funcs.CpuMetrics()

	data, err := g.EncodeData(cpu)
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
	pkt.Id = g.HHD

	disks, err := funcs.DiskUseMetrics()
	if err != nil {
		log.Print(err)
		return
	}

	pkt.Data, err = g.EncodeData(disks)

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
	pkt.Id = g.MEM
	memory, err := funcs.MemMetrics()
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data, err = g.EncodeData(memory)
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
	pkt.Id = g.LoadAvg
	loadAvg, err := funcs.LoadAvgMetrics()
	if err != nil {
		log.Println(err)
		return
	}
	pkt.Data, err = g.EncodeData(loadAvg)
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
	pkt.Id = g.PortListen
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
	pkt.Data, err = g.EncodeData(ports)
	if err != nil {
		log.Println(err)
		return
	}
	if err := conn.Write(pkt); err != nil {
		log.Println(err)
	}
}
