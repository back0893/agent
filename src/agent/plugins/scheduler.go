package plugins

import (
	"agent/src"
	"agent/src/agent/iface"
	"agent/src/g"
	cmd2 "agent/src/g/cmd"
	"agent/src/g/model"
	"bytes"
	"encoding/json"
	"github.com/back0893/goTcp/utils"
	"github.com/toolkits/file"
	"log"
	"path/filepath"
	"time"
)

type PluginScheduler struct {
	Ticker *time.Ticker
	Plugin *Plugin
	Quit   chan struct{}
}

func NewPluginScheduler(p *Plugin) *PluginScheduler {
	scheduler := PluginScheduler{
		Plugin: p,
		Ticker: time.NewTicker(time.Duration(p.Interval) * time.Second),
		Quit:   make(chan struct{}),
	}
	return &scheduler
}

func (this *PluginScheduler) Schedule() {
	go func() {
		for {
			select {
			case <-this.Ticker.C:
				PluginRun(this.Plugin)
			case <-this.Quit:
				this.Ticker.Stop()
				return
			}
		}
	}()
}

func (this *PluginScheduler) Stop() {
	close(this.Quit)
}

func PluginRun(plugin *Plugin) {

	timeout := plugin.Interval*1000 - 500
	fpath := filepath.Join(utils.GlobalConfig.GetString("plugin.dir"), plugin.FilePath)

	if !file.IsExist(fpath) {
		log.Println("no such plugin:", fpath)
		return
	}

	debug := utils.GlobalConfig.GetBool("debug")
	if debug {
		log.Println(fpath, "running...")
	}
	cmd := cmd2.Command{
		Name:    fpath,
		Args:    []string{},
		Timeout: timeout,
	}
	cmd.Callback = func(stdout, stderr bytes.Buffer, err error, isTimeout bool) {
		errStr := stderr.String()
		if errStr != "" {
			logFile := filepath.Join(utils.GlobalConfig.GetString("plugin.log"), plugin.FilePath+".stderr.log")
			if _, err = file.WriteString(logFile, errStr); err != nil {
				log.Printf("[ERROR] write log to %s fail, error: %s\n", logFile, err)
			}
		}

		if isTimeout {
			// has be killed
			if err == nil && debug {
				log.Println("[INFO] timeout and kill process", fpath, "successfully")
			}

			if err != nil {
				log.Println("[ERROR] kill process", fpath, "occur error:", err)
			}

			return
		}

		if err != nil {
			log.Println("[ERROR] exec plugin", fpath, "fail. error:", err)
			return
		}

		// exec successfully
		data := stdout.Bytes()
		if len(data) == 0 {
			if debug {
				log.Println("debug stdout empty")
			}
			return
		}

		var metrics []*model.MetricValue
		err = json.Unmarshal(data, &metrics)
		if err != nil {
			log.Print(err)
			return
		}

		a := utils.GlobalConfig.Get(g.AGENT).(iface.IAgent)
		pkt := src.NewPkt()
		if pkt.Data, err = g.EncodeData(metrics); err != nil {
			log.Print(err)
			return
		}
		if err := a.GetCon().Write(pkt); err != nil {
			log.Println(err)
		}
	}
	go func() {
		if err := cmd.Run(); err != nil {
			log.Println(err)
		}
	}()
}
