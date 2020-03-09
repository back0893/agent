package agent

import (
	"agent/src/g"
	"encoding/json"
	"github.com/back0893/goTcp/utils"
	"os"
)

/**
更新,需要有撤销动作当执行失败时..
*/
const (
	newFile     = "/xxx.new"
	oldFile     = "/xxx.old"
	currentFile = "/xxx"
)

type VersionInfo struct {
	version int    //版本
	path    string //如果需要更新,更新的地址
}
type UpdateCommand struct {
}

func (uc *UpdateCommand) Do() {
	url := utils.GlobalConfig.GetString("update.url")
	if url == "" {
		return
	}
	//todo 检测版本号
	data, err := g.Post(url, nil)
	//获得对应的版本信息失败
	if err != nil {
		return
	}
	versionInfo := &VersionInfo{}
	if err = json.Unmarshal(data, versionInfo); err != nil {
		return
	}

	g.Down(versionInfo.path, newFile)

	if err := os.Rename(currentFile, oldFile); err != nil {
		uc.Undo()
		return
	}
	if err := os.Rename(newFile, currentFile); err != nil {
		uc.Undo()
		return
	}
}
func (uc *UpdateCommand) Undo() {
	//如果回退失败应该直接退出,并记录日志?
	_ = os.Rename(currentFile, newFile)
	_ = os.Rename(oldFile, currentFile)
}
