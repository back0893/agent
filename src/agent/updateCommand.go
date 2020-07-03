package agent

import (
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/back0893/goTcp/utils"
)

var updateChan chan *model.UpdateInfo

func init() {
	updateChan = make(chan *model.UpdateInfo)
}
func GetUpdateChan() chan *model.UpdateInfo {
	return updateChan
}
func Undo() {
	filename := fmt.Sprintf("%s/agent", utils.GlobalConfig.GetString("root"))
	agent := NewUpdate(filename)
	agent.Undo()
}

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
func (uc *UpdateCommand) Do(url string) error {
	newFile := uc.GetNewFilename()
	currentFile := uc.GetFilename()
	oldFile := uc.GetOldIFilename()
	if err := g.Down(url, newFile); err != nil {
		uc.Undo()
		return err
	}
	if err := os.Rename(currentFile, oldFile); err != nil {
		uc.Undo()
		return err
	}
	if err := os.Rename(newFile, currentFile); err != nil {
		uc.Undo()
		return err
	}

	//agent退出

	return nil
}
func (uc *UpdateCommand) Undo() {
	newFile := uc.GetNewFilename()
	currentFile := uc.GetFilename()
	oldFile := uc.GetOldIFilename()
	_ = os.Rename(currentFile, newFile)
	_ = os.Rename(oldFile, currentFile)
	_ = os.Remove(newFile)
	//如果回退失败应该直接退出,并记录日志?
}

func AgentSelfUpdate(ctx context.Context) {
	for {
		select {
		case info := <-updateChan:
			Upgrade(info)
		case <-ctx.Done():
			break
		}
	}
}
func Upgrade(info *model.UpdateInfo) {
	agentPath := fmt.Sprintf("%s/agent", utils.GlobalConfig.GetString("root"))
	cfgPath := utils.GlobalConfig.GetString("cfgpath")
	log.Println(info.URL)
	success := false
	var err error
	switch info.Type {
	case 1:
		binUpdate := NewUpdate(agentPath)
		if err = binUpdate.Do(info.URL); err == nil {
			success = true
		}

	case 2:
		cfgUpdate := NewUpdate(cfgPath)
		if err = cfgUpdate.Do(info.URL); err == nil {
			success = true
		}
	}
	//通知中控升级成功或者失败
	agent := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
	pkt := g.NewPkt()
	pkt.Id = g.UPDATE
	data := model.UpdateResponse{LogID: info.LogID}
	if success {
		data.Status = true
	} else {
		data.Message = err.Error()
	}
	pkt.Data, _ = g.EncodeData(data)
	agent.GetCon().AsyncWrite(pkt, 5)
	//等待1s
	//自杀,等待重启
	time.Sleep(time.Second)
	if success {
		agent.Stop()
	} else {
		log.Println(err.Error())
	}

}
