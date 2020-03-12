package main

import (
	"agent/src"
	"agent/src/agent"
	"agent/src/agent/cron"
	"agent/src/g"
	"flag"
	"github.com/back0893/goTcp/utils"
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
		//todo 将基础信息上班规整层一个服务
		cron.SendMem(agentClient.GetCon())
		//cron.SendPort(agentClient.GetCon())
		//cron.SendCPU(agentClient.GetCon())
		//cron.SendHHD(agentClient.GetCon())
		//cron.SendLoadAvg(agentClient.GetCon())
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
func main() {
	var action string
	var cfg string
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置,只有start时才有用")
	flag.StringVar(&action, "t", "start", "命令动作,start|stop")
	flag.Parse()
	//获得当前的文件名,以便更新
	utils.GlobalConfig.Set("filename", os.Args[0])

	switch action {
	case "start":
		start(cfg)
	case "stop":
		stop()
	default:
		flag.Usage()
	}

}
