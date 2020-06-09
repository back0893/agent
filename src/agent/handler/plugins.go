package handler

import (
	"agent/src"
	"agent/src/g"
	"agent/src/g/model"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"github.com/toolkits/file"
	"log"
	"os/exec"
)

type Plugins struct {
}

//规则化命名 name_interval 名称_执行间隔
func (p *Plugins) FormatName(name string, interval int) string {
	return fmt.Sprintf("%s_%d", name, interval)
}
func (p *Plugins) http(uri, filePath string) error {
	return g.Down(uri, filePath)
}
func (p *Plugins) git(dir string, pkt *model.Plugins) error {
	if file.IsExist(dir) {
		cmd := exec.Command("git", "pull")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {

		}
	} else {
		var cmd *exec.Cmd
		if pkt.Branch != "" {
			cmd = exec.Command("git", "clone", "-branch", pkt.Branch, pkt.Uri, p.FormatName(pkt.Name, pkt.Interval))
		} else {
			cmd = exec.Command("git", "clone", pkt.Uri, p.FormatName(pkt.Name, pkt.Interval))
		}
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {

		}
	}
	return nil
}
func (p *Plugins) Down(dir string, pkt *model.Plugins) error {
	switch pkt.Type {
	case 1:
		return p.git(dir, pkt)
	case 2:
		filePath := fmt.Sprintf("%s/%s", dir, p.FormatName(pkt.Name, pkt.Interval))
		return p.http(pkt.Uri, filePath)
	}
	return errors.New("tag不存在")
}

func (p Plugins) Handler(ctx context.Context, packet *src.Packet, connection iface.IConnection) {
	decoder := gob.NewDecoder(bytes.NewReader(packet.Data))
	plugins := model.Plugins{}
	if err := decoder.Decode(&plugins); err != nil {
		log.Println(err)
		return
	}
	if err := p.Down("", &plugins); err != nil {
		log.Println(err)
		return
	}
}
