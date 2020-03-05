package main

/**
这是更新,和部署agent
是所有的agent的根节点
目前这个跟节点不能自动升级,只能手动更新
同时这个也是一个agent
*/

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



func main(){
	var cfg string
	var
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置json")
	flag.Parse()

	con, err := agent.ConnectServer(cfg)
	if err != nil {
		panic(err)
	}
	updater := agent.NewAgent(con, agent.Event{}, src.Protocol{}, cfg)

	//断线重连
	updater.AddClose(updater.ReCon)
	//心跳单独实现.
	headrtBeat := utils.GlobalConfig.GetInt("heartBeat")
	src.AddTimer(time.Second*time.Duration(headrtBeat), func() {
		cron.SendHeart(updater.GetCon())
	})

	src.InitTimingWheel(updater.GetContext())
	updater.Start()

	log.Println("接受停止或者ctrl-c停止")
	chSign := make(chan os.Signal)
	signal.Notify(chSign, syscall.SIGINT, syscall.SIGTERM)
	log.Println("接受到信号:", <-chSign)
	updater.Stop()
}
