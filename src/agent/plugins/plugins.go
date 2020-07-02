package plugins

import (
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"bytes"
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/file"
)

type Plugin struct {
	FilePath string
	Interval int
	IsRepeat bool
	MTime    int64
}

var (
	Plugins              = make(map[string]*Plugin)
	PluginsWithScheduler = make(map[string]*PluginScheduler)
)

func DelNoUsePlugins(newPlugins map[string]*Plugin) {
	for currKey, currPlugin := range Plugins {
		if newPlugin, ok := newPlugins[currKey]; !ok || currPlugin.MTime != newPlugin.MTime {
			deletePlugin(currKey)
		}
	}
}

func AddNewPlugins(newPlugins map[string]*Plugin) {
	for fpath, newPlugin := range newPlugins {
		if currPlugin, ok := Plugins[fpath]; ok {
			//存在但是插件被修改过,也需要停止后重启
			if newPlugin.MTime != currPlugin.MTime {
				deletePlugin(fpath)
			} else {
				continue
			}
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

func Git(dir string, repPlugins *model.Plugins, logID int32) {
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
		cmd.Callback = func(stdout, stderr *bytes.Buffer, err error, isTimeout bool) {
			var status int8 = 0
			var message string = ""
			if stderr != nil && stderr.String() != "" {
				message = string(stderr.Bytes())
			} else if isTimeout {
				// has be killed
				message = "git 执行超时"
			} else if err != nil {
				message = err.Error()
			} else {
				// exec successfully
				status = 1
				message = string(stdout.Bytes())
			}
			pkt := g.ComResponse(logID, status, message)

			a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
			if err := a.GetCon().Write(pkt); err != nil {
				log.Println(err)
			}
		}
		go func() {
			cmd.Run()
			fmt.Println("readList")
			desiredAll := ListPlugins(utils.GlobalConfig.GetString("plugin.dir"))
			DelNoUsePlugins(desiredAll)
			AddNewPlugins(desiredAll)
		}()
	}

}
