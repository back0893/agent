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
	"flag"
	"fmt"
	"os"
	"syscall"
)

func selfStart(cfg string) {
	updater, err := agent.NewAgent(cfg)
	if err != nil {
		panic(err)
	}
	src.SavePid("./pid")
	updater.Start()
	updater.Wait()
}

func stop(name string) {
	//todo 应该直接执行对应的agent的可执行文件
	pid := src.ReadPid("./pid")
	_ = syscall.Kill(pid, syscall.SIGKILL)
}

func start(name, cfg string) {
	//todo 应该直接执行对应的agent的可执行文件
	selfStart(cfg)
}
func update(name string) {
	//todo 请求更新的配置,下载,替换
	os.Rename("./update", "./update.old")
	os.Rename("./update.new", "./update")
}
func status(name string) {
	isRun := src.Status(src.ReadPid("./pid"))
	str := ""
	if isRun {
		str = fmt.Sprintf("%s正在运行", name)
	} else {
		str = fmt.Sprintf("%s未运行", name)
	}
	fmt.Print(str)
}

func main() {
	var action string
	var agentName string
	var cfg string
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置,只有start时才有用")
	flag.StringVar(&action, "t", "start", "命令动作,start|stop|update|check")
	flag.StringVar(&agentName, "u", "", "操作的agent的名称,默认为自己")
	flag.Parse()

	switch action {
	case "start":
		start(agentName, cfg)
	case "stop":
		stop(agentName)
	case "update":
		update(agentName)
	//todo 更新完成后应该停止并且启动对应的agent
	case "status":
		status(agentName)
	case "check":
		//从中心服务器下载最新的
	}
}
