package main

import (
	"agent/src"
	agentClient "agent/src/agent"
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

	con, err := agentClient.ConnectServer(cfg)
	if err != nil {
		panic(err)
	}
	agent := agentClient.NewAgent(con, agentClient.Event{}, src.Protocol{}, cfg)

	//断线重连
	agent.AddClose(agent.ReCon)

	src.InitTimingWheel(agent.GetContext())

	//心跳单独实现.
	headrtBeat := utils.GlobalConfig.GetInt("heartBeat")
	src.AddTimer(time.Second*time.Duration(headrtBeat), func() {
		cron.SendHeart(agent.GetCon())
	})
	go agent.RunTask()
	agent.Start()

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	agent.Stop()
}
