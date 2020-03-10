package agent

import (
	"agent/src/agent/model"
	"agent/src/g"
	"os"
)

/**
更新,需要有撤销动作当执行失败时..
*/
const (
	newFile     = "./hi2.new"
	oldFile     = "./hi2.old"
	currentFile = "./hi2"
)

type UpdateCommand struct {
}

func NewUpdate() *UpdateCommand {
	return &UpdateCommand{}
}

func (uc *UpdateCommand) Do(Info *model.UpdateInfo) error {
	//版本小于当前的版本号
	//if Info.Version<utils.GlobalConfig.GetInt("Version"){
	//	return nil
	//}
	g.Down(Info.Url, newFile)
	if err := os.Rename(currentFile, oldFile); err != nil {
		uc.Undo()
		return err
	}
	if err := os.Rename(newFile, currentFile); err != nil {
		uc.Undo()
		return err
	}
	return nil
}
func (uc *UpdateCommand) Undo() {
	//如果回退失败应该直接退出,并记录日志?
	_ = os.Rename(currentFile, newFile)
	_ = os.Rename(oldFile, currentFile)
}
