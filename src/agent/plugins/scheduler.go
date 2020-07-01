package plugins

import (
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/g/model"
	"bytes"
	"encoding/json"
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

func PluginRun(plugin *Plugin) {

	timeout := plugin.Interval*1000 - 500
	if plugin.IsRepeat == false {
		timeout = 0
	}
	fpath := plugin.FilePath

	if !file.IsExist(fpath) {
		log.Println("no such plugin:", fpath)
		return
	}

	debug := utils.GlobalConfig.GetBool("debug")
	if debug {
		log.Println(fpath, "running...")
	}
	cmd := g.Command{
		Name:    fpath,
		Args:    []string{},
		Timeout: timeout,
	}
	cmd.Callback = func(stdout, stderr bytes.Buffer, err error, isTimeout bool) {
		var metrics []*model.MetricValue

		errStr := stderr.String()
		if errStr != "" {
			value := model.MetricValue{
				Metric:    "exec.fail",
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintln(file.Basename(plugin.FilePath), " err:", errStr),
			}
			metrics = append(metrics, &value)
		} else if isTimeout {
			// has be killed
			if debug {
				log.Println("[INFO] timeout and kill process", fpath, "successfully")
			}
			value := model.MetricValue{
				Metric:    "exec.fail",
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintln(file.Basename(plugin.FilePath), " timeout error:", err),
			}
			metrics = append(metrics, &value)
		} else if err != nil {
			value := model.MetricValue{
				Metric:    "exec.fail",
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintln(file.Basename(plugin.FilePath), " error:", err),
			}
			metrics = append(metrics, &value)
		} else {
			// exec successfully
			data := stdout.Bytes()
			if len(data) == 0 {
				if debug {
					log.Println("debug stdout empty")
				}
				return
			}
			log.Println(string(data))
			err = json.Unmarshal(data, &metrics)
			if err != nil {
				log.Print(err)
				return
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
	go func() {
		cmd.Run()
	}()
}
