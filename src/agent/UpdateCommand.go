package agent

import (
	"agent/src/agent/model"
	"agent/src/g"
	"fmt"
	"os"
)

/**
更新,需要有撤销动作当执行失败时..
*/
type UpdateCommand struct {
	filename string
}

//新文件名
func (uc UpdateCommand) GetNewFilename() string {
	return fmt.Sprintf("%s.new", uc.filename)
}

//旧文件名
func (uc UpdateCommand) GetOldIFilename() string {
	return fmt.Sprintf("%s.old", uc.filename)
}
func (uc UpdateCommand) GetFilename() string {
	return uc.filename
}
func NewUpdate(filename string) *UpdateCommand {
	return &UpdateCommand{
		filename: filename,
	}
}

func (uc *UpdateCommand) Do(Info *model.UpdateInfo) error {
	//版本小于当前的版本号
	//if Info.Version<utils.GlobalConfig.GetInt("Version"){
	//	return nil
	//}
	newFile := uc.GetNewFilename()
	currentFile := uc.GetFilename()
	oldFile := uc.GetOldIFilename()
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
	newFile := uc.GetNewFilename()
	currentFile := uc.GetFilename()
	oldFile := uc.GetOldIFilename()
	_ = os.Rename(currentFile, newFile)
	_ = os.Rename(oldFile, currentFile)
	//如果回退失败应该直接退出,并记录日志?
}
