package serviceHandler

import (
	"agent/src/g"
	"agent/src/g/model"
	"agent/src/server/Db"
	serverModel "agent/src/server/model"
	"context"
	"github.com/back0893/goTcp/iface"
	"log"
)

type HHDHandler struct {
}

func NewHHDHandler() *HHDHandler {
	return &HHDHandler{}
}
func (H HHDHandler) Handler(ctx context.Context, service *model.ServiceResponse, connection iface.IConnection) error {
	disks := make([]*model.Disk, 0)
	if err := g.DecodeData(service.Info, &disks); err != nil {
		log.Println("读取硬盘信息失败")
		return err
	}
	tmp, _ := connection.GetExtraData("auth")
	auth := tmp.(*model.Auth)
	db, _ := Db.DbConnections.Get("ep")
	for _, disk := range disks {
		var serverDisk serverModel.ServerDisk
		if err := db.Get(&serverDisk, "select id,name,gb,server_id from cc_server_disk where server_id=? and name=?", auth.Id, disk.FsFile); err != nil {
		}
		total := float64(disk.Total) / (1024 * 1024)
		if serverDisk.Id == 0 {
			if re, err := db.Exec("insert cc_server_disk (`name`,`gb`,`server_id`) values(?,?,?)", disk.FsFile, g.Round(total/1024, 2), auth.Id); err != nil {
				continue
			} else {
				serverDisk.Id, _ = re.LastInsertId()
			}
		} else {
			if _, err := db.Exec("update cc_server_disk  set gb=? where server_id=? and `name`=?", g.Round(total/1024, 2), auth.Id, disk.FsFile); err != nil {
				continue
			}
		}
		created_at := g.CSTTime()
		if _, err := db.Exec("insert cc_server_disk_log (`disk_id`,`usage_ratio`,`created_at`) values (?,?,?)", serverDisk.Id, g.Round(disk.Percent, 2), created_at); err != nil {
			continue
		}
	}
	return nil
}
