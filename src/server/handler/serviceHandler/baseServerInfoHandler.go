package serviceHandler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/db"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type BaseServerInfo struct {
}

func NewBaseServerInfo() *BaseServerInfo {
	return &BaseServerInfo{}
}
func (b BaseServerInfo) Handler(ctx context.Context, service *model.ServiceResponse, connection iface.IConnection) error {
	var cpu model.Cpu
	var mem model.Memory
	loadAvgs := make([]*model.LoadAvg, 0)
	var cpuNum int
	var cpuMhz string
	if err := g.DecodeData(service.Info, &cpu, &mem, &loadAvgs, &cpuNum, &cpuMhz); err != nil {
		log.Println("读取信息失败")
		return err
	}
	tmp, _ := connection.GetExtraData("auth")
	auth := tmp.(*model.Auth)
	db, _ := Db.DbConnections.Get("ep")
	ram_usage_ratio := g.Round(float64(mem.Used)/float64(mem.Total), 2)
	if _, err := db.Exec("insert cc_server_log (server_id,ram,cpu_usage_ratio,ram_usage_ratio,created_at) values (?,?,?,?,?)", auth.Id, float64(mem.Total)/(1024*1024), g.Round(cpu.Busy/100, 2), ram_usage_ratio, g.CSTTime()); err != nil {
		log.Println(err.Error())
	}
	if _, err := db.Exec("update cc_server set cpu_usage_ratio=?,ram_usage_ratio=?,cpu_num=?,cpu_mhz=? where id=?", g.Round(cpu.Busy/100, 2), ram_usage_ratio, cpuNum, cpuMhz, auth.Id); err != nil {
		log.Println(err.Error())
	}
	return nil
}
