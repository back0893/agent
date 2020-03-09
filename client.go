package main

import (
	"agent/src"
	"agent/src/agent"
	"agent/src/agent/cron"
	"flag"
	"time"
)

func main() {
	var cfg string
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置json")
	flag.Parse()

	agentClient, err := agent.NewAgent(cfg)
	if err != nil {
		panic(err)
	}

	src.AddTimer(time.Second*5, func() {
		cron.SendMem(agentClient.GetCon())
		cron.SendPort(agentClient.GetCon())
		cron.SendCPU(agentClient.GetCon())
		cron.SendHHD(agentClient.GetCon())
		cron.SendLoadAvg(agentClient.GetCon())
	})
	go agentClient.RunTask()
	agentClient.Start()
	agentClient.Wait()
}
