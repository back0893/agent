package plugins

import (
	"agent/src/g"
	"agent/src/g/model"
	"bytes"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/file"
	"log"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type Plugin struct {
	FilePath string
	Interval int
	IsRepeat bool
}

var (
	Plugins              = make(map[string]*Plugin)
	PluginsWithScheduler = make(map[string]*PluginScheduler)
)

func DelNoUsePlugins(newPlugins map[string]*Plugin) {
	for currKey := range Plugins {
		if _, ok := newPlugins[currKey]; !ok {
			deletePlugin(currKey)
		}
	}
}

func AddNewPlugins(newPlugins map[string]*Plugin) {
	for fpath, newPlugin := range newPlugins {
		if _, ok := Plugins[fpath]; ok {
			continue
		}

		Plugins[fpath] = newPlugin
		sch := NewPluginScheduler(newPlugin)
		PluginsWithScheduler[fpath] = sch
		sch.Schedule()
	}
}

func ClearAllPlugins() {
	for k := range Plugins {
		deletePlugin(k)
	}
}

func deletePlugin(key string) {
	v, ok := PluginsWithScheduler[key]
	if ok {
		v.Stop()
		delete(PluginsWithScheduler, key)
	}
	delete(Plugins, key)
}

func Git(dir string, repPlugins *model.Plugins) {
	//没有插件的git地址,说明灭有配置插件
	if len(repPlugins.Uri) == 0 {
		ClearAllPlugins()
	}
	file.InsureDir(dir)

	for _, uri := range repPlugins.Uri {
		//去掉尾部的扩展名.git
		name := strings.Replace(path.Base(uri), ".git", "", 1)
		dirPath := filepath.Join(dir, name)
		var cmd g.Command
		fmt.Println(dirPath)
		if file.IsExist(dirPath) {
			cmd = g.Command{
				Name:    "git",
				Args:    []string{"pull"},
				Timeout: 60 * 1000,
				Dir:     dirPath,
			}
		} else {
			cmd = g.Command{
				Name:    "git",
				Args:    []string{"clone", uri},
				Timeout: 120 * 1000,
				Dir:     dir,
			}
		}
		cmd.Callback = func(stdout, stderr bytes.Buffer, err error, isTimeout bool) {
			errStr := stderr.String()
			if errStr != "" {
				logFile := filepath.Join(utils.GlobalConfig.GetString("plugin.log"), "git"+".stderr.log")
				if _, err = file.WriteString(logFile, errStr); err != nil {
					log.Printf("[ERROR] write log to %s fail, error: %s\n", logFile, err)
				}
			}

			if isTimeout {
				// has be killed
				if err == nil {
					log.Println("[INFO] git timeout and kill process")
				}

				return
			}

			if err != nil {
				log.Println("[ERROR] exec git fail. error:", err)
				return
			}

			//回应git更新成功,应该为日志
			//a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
			//pkt := g.NewPkt()
			//
			////todo 回应git更新成功
			////id没有想好
			//pkt.Id = g.MinePluginsResponse
			//if err := a.GetCon().Write(pkt); err != nil {
			//	log.Println(err)
			//}
		}
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			if err := cmd.Run(); err != nil {
				log.Println(err)
			}
			fmt.Println("readList")
			desiredAll := ListPlugins("")
			DelNoUsePlugins(desiredAll)
			AddNewPlugins(desiredAll)
			wg.Done()
		}()
		wg.Wait()
	}

}
