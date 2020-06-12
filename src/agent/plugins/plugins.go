package plugins

import (
	"agent/src/agent/iface"
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
	"time"
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
			value := model.MetricValue{
				Metric:    "git",
				Timestamp: time.Now().Unix(),
				Value:     stdout.String(),
			}
			if isTimeout {
				value.Value = fmt.Sprintln("[ERROR] git timeout error:", err)
			}
			if err != nil {
				//如果异常退出,错误输出才有用,,
				//因为其他命令也可能使用错误输出
				errStr := stderr.String()
				if errStr != "" {
					value.Value = fmt.Sprintln("[ERROR] git update fail error:", errStr)
				}
			}
			//回应git更新成功,应该为日志
			//每执行一个操作后,应该将操作的成功或者失败的信息通知中控服务器
			a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
			pkt := g.NewPkt()
			pkt.Id = g.ActionNotice
			pkt.Data, _ = g.EncodeData([]*model.MetricValue{&value})
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
