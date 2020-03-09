package main

import (
	"agent/src"
	"agent/src/agent"
	"agent/src/agent/cron"
	"agent/src/g"
	"flag"
	"os"
	"syscall"
	"time"
)

func start(cfg string) {
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
	g.SavePid("./pid")
	agentClient.Wait()
}
func stop() {
	pid := g.ReadPid("./pid")
	_ = syscall.Kill(pid, syscall.SIGKILL)
}

func update() {
	//todo 从中心服务器检测是否需要更新
	//todo 请求更新的配置,下载,替换
	os.Rename("./update", "./update.old")
	os.Rename("./update.new", "./update")
}
func main() {
	var action string
	var cfg string
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置,只有start时才有用")
	flag.StringVar(&action, "t", "start", "命令动作,start|stop|update")
	flag.Parse()

	switch action {
	case "start":
		start(cfg)
	case "stop":
		stop()
	case "update":
		update()
	}
}
