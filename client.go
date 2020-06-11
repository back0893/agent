package main

import (
	"agent/src/agent"
	httpHandler "agent/src/agent/http"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/http"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"log"
	"os"
)

func start(cfg string) {
	agentClient, err := agent.NewAgent(cfg)
	if err != nil {
		panic(err)
	}
	agentClient.Start()
	utils.GlobalConfig.Set(g.AGENT, agentClient)
	go startHttp(agentClient)
	agentClient.Wait()
}

func startHttp(agent iface.IAgent) {
	if !utils.GlobalConfig.GetBool("http.enabled") {
		return
	}
	addr := utils.GlobalConfig.GetString("http.listen")
	if addr == "" {
		return
	}
	log.Println("http start ,listening", addr)
	httpServer := http.NewServer(addr)
	httpServer.AddHandler("/push", httpHandler.WrapperTransfer(agent))
	log.Fatalln(httpServer.Run())
}

func main() {
	var action string
	var cfg string
	flag.StringVar(&cfg, "c", "./client.json", "加载的配置,只有start时才有用")
	flag.StringVar(&action, "t", "start", "命令动作,start|stop")
	flag.Parse()
	//获得当前的文件名,以便更新
	utils.GlobalConfig.Set("filename", os.Args[0])

	switch action {
	case "start":
		start(cfg)
	case "version":
		fmt.Printf("当前版本:%d\n", g.VERSION)
	default:
		flag.Usage()
	}

}
