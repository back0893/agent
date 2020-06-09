package model

import (
	"agent/src/g"
	"errors"
	"fmt"
	"github.com/toolkits/file"
	"os/exec"
)

/**
插件分为2种
http下载,http下载使用url中token作为用户判断.
git拉取,拉取地址和用户作为用户判断
*/
type Plugins struct {
	Type     int //作为判断为git还是http
	Uri      string
	Name     string //插件名称
	Interval int    //执行时间间隔(秒)
	Branch   string //git的分或者tag名称
}

//规则化命名 name_interval 名称_执行间隔
func (p *Plugins) FormatName() string {
	return fmt.Sprintf("%s_%d", p.Name, p.Interval)
}
func (p *Plugins) http(dir string) error {
	filePath := fmt.Sprintf("%s/%s", dir, p.FormatName())
	return g.Down(p.Uri, filePath)
}
func (p *Plugins) git(dir string) error {
	if file.IsExist(dir) {
		cmd := exec.Command("git", "pull")
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {

		}
	} else {
		var cmd *exec.Cmd
		if p.Branch != "" {
			cmd = exec.Command("git", "clone", "-branch", p.Branch, p.Uri, p.FormatName())
		} else {
			cmd = exec.Command("git", "clone", p.Uri, p.FormatName())
		}
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {

		}
	}
	return nil
}
func (p *Plugins) Down(dir string) error {
	switch p.Type {
	case 1:
		return p.git(dir)
	case 2:
		return p.http(dir)
	}
	return errors.New("tag不存在")
}
