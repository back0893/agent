package main

import (
	"agent/src/agent"
	"agent/src/g"
	"flag"
	"fmt"
	"github.com/back0893/goTcp/utils"
	"os"
	"syscall"
)

func start(cfg string) {
	agentClient, err := agent.NewAgent(cfg)
	if err != nil {
		panic(err)
	}
	agentClient.Start()
	g.SavePid("./pid")
	utils.GlobalConfig.Set(g.AGENT, agentClient)
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
	case "version":
		fmt.Printf("当前版本:%d\n", g.VERSION)
	default:
		flag.Usage()
	}

}
