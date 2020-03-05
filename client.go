package main

import (
	"agent/src"
	"agent/src/agent"
	"agent/src/agent/cron"
	"flag"
	"github.com/back0893/goTcp/utils"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var cfg string
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置json")
	flag.Parse()

	con, err := agent.ConnectServer(cfg)
	if err != nil {
		panic(err)
	}
	agentClient := agent.NewAgent(con, agent.Event{}, src.Protocol{}, cfg)

	//断线重连
	agentClient.AddClose(agentClient.ReCon)

	src.InitTimingWheel(agentClient.GetContext())

	//心跳单独实现.
	headrtBeat := utils.GlobalConfig.GetInt("heartBeat")
	src.AddTimer(time.Second*time.Duration(headrtBeat), func() {
		cron.SendHeart(agentClient.GetCon())
	})
	src.AddTimer(time.Second*5, func() {
		cron.SendMem(agentClient.GetCon())
		cron.SendPort(agentClient.GetCon())
		cron.SendCPU(agentClient.GetCon())
		cron.SendHHD(agentClient.GetCon())
		cron.SendLoadAvg(agentClient.GetCon())
	})
	go agentClient.RunTask()
	agentClient.Start()

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	agentClient.Stop()
}
