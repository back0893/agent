package plugins

import (
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/file"
)

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
	once   sync.Once
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{
		Plugin: p,
		Quit:   make(chan struct{}),
	}
	if p.Interval > 0 {
		scheduler.Ticker = time.NewTicker(time.Duration(p.Interval) * time.Second)
	} else {
		scheduler.Ticker = time.NewTicker(1 * time.Second)
	}
	return &scheduler
}

func (this *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				PluginRun(this.Plugin)
				if !this.Plugin.IsRepeat {
					this.Stop()
				}
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *PluginScheduler) Stop() {
	this.once.Do(func() {
		close(this.Quit)
	})
}

//PluginRun 插件的定时执行
func PluginRun(plugin *Plugin) {
	timeout := plugin.Interval*1000 - 500
	if plugin.IsRepeat == false {
		timeout = 0
	}
	fpath := plugin.FilePath
	cmd := g.Command{
		Name:    fpath,
		Args:    []string{},
		Timeout: timeout,
	}
	debug := utils.GlobalConfig.GetBool("debug")
	if debug {
		log.Println(fpath, "running...")
	}
	cmd.Callback = func(stdout, stderr *bytes.Buffer, err error, isTimeout bool) {
		var metrics []*model.MetricValue

		if stderr != nil && stderr.String() != "" {
			value := model.MetricValue{
				Metric:    "plugin.fail",
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintln(file.Basename(plugin.FilePath), " err:", stderr.String()),
			}
			metrics = append(metrics, &value)
		} else if isTimeout {
			// has be killed
			if debug {
				log.Println("[INFO] timeout and kill process", fpath, "successfully")
			}
			value := model.MetricValue{
				Metric:    "plugin.fail",
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintln(file.Basename(plugin.FilePath), " timeout error:", err),
			}
			metrics = append(metrics, &value)
		} else if err != nil {
			value := model.MetricValue{
				Metric:    "plugin.fail",
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintln(file.Basename(plugin.FilePath), " error:", err),
			}
			metrics = append(metrics, &value)
		} else {
			// exec successfully
			if stdout != nil {
				data := stdout.Bytes()
				if len(data) == 0 {
					if debug {
						log.Println("debug stdout empty")
					}
					value := model.MetricValue{
						Metric:    "plugin.success",
						Timestamp: time.Now().Unix(),
						Value:     "stdout empty",
					}
					metrics = append(metrics, &value)
				} else {
					log.Println(string(data))
					err = json.Unmarshal(data, &metrics)
					if err != nil {
						value := model.MetricValue{
							Metric:    "plugin.success",
							Timestamp: time.Now().Unix(),
							Value:     string(data),
						}
						metrics = append(metrics, &value)
					}
				}
			}
		}
		pkt := g.NewPkt()
		pkt.Id = g.ActionNotice
		if pkt.Data, err = g.EncodeData(metrics); err != nil {
			log.Print(err)
			return
		}
		a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
		if err := a.GetCon().Write(pkt); err != nil {
			log.Println(err)
		}
	}

	if !file.IsExist(fpath) {
		cmd.Callback(nil, nil, errors.New(fmt.Sprintf("%s 不存在", plugin.FilePath)), false)
		return
	}

	go func() {
		cmd.Run()
	}()
}

//PluginExecute 插件的独立执行
func PluginExecute(plugin *Plugin, fn func(stdout, stderr *bytes.Buffer, err error, isTimeout bool)) {
	fpath := plugin.FilePath
	if !file.IsExist(fpath) {
		fn(nil, nil, fmt.Errorf("%s 不存在", plugin.FilePath), false)
		return
	}

	debug := utils.GlobalConfig.GetBool("debug")
	if debug {
		log.Println(fpath, "running...")
	}
	cmd := g.Command{
		Name:    fpath,
		Args:    []string{},
		Timeout: plugin.Interval,
	}
	cmd.Callback = fn
	go func() {
		cmd.Run()
	}()
}
