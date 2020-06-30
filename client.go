package main

import (
	"agent/src/agent"
	httpHandler "agent/src/agent/http"
	"agent/src/agent/iface"
	"agent/src/g"
	"agent/src/http"
	"flag"
	"github.com/back0893/goTcp/utils"
	"log"
	"path/filepath"
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
	flag.StringVar(&action, "t", "start", "命令动作,start")
	flag.Parse()
	//获得当前的文件路径
	root, _ := filepath.Abs(".")
	utils.GlobalConfig.Set("root", root)
	utils.GlobalConfig.Set("cfgpath", cfg)
	switch action {
	case "start":
		start(cfg)
	default:
		flag.Usage()
	}

}
