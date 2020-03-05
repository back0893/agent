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
	"bytes"
	"flag"
	"github.com/back0893/goTcp/utils"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func recordPid(pidfile string) {
	pid := os.Getpid()
	file, _ := os.Create(pidfile)
	defer file.Close()
	io.WriteString(file, strconv.Itoa(pid))
}
func selfStart(cfg string) {
	con, err := agent.ConnectServer(cfg)
	if err != nil {
		panic(err)
	}
	updater := agent.NewAgent(con, agent.Event{}, src.Protocol{}, cfg)
	src.InitTimingWheel(updater.GetContext())
	//断线重连
	updater.AddClose(updater.ReCon)
	//心跳单独实现.
	heartBeat := utils.GlobalConfig.GetInt("heartBeat")
	src.AddTimer(time.Second*time.Duration(heartBeat), func() {
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
func stop(name string) {
	//todo 应该直接执行对应的agent的可执行文件
	file, _ := os.Open("./pid.pid")
	data, _ := ioutil.ReadAll(file)
	data = bytes.Trim(data, "\r\n")
	pid, _ := strconv.Atoi(string(data))
	syscall.Kill(pid, syscall.SIGKILL)
}
func start(name, cfg string) {
	//todo 应该直接执行对应的agent的可执行文件
	recordPid("./pid.pid")
	selfStart(cfg)
}
func update(name string) {
	os.Rename("./update", "./update.old")
	os.Rename("./update.new", "./update")
}
func main() {
	var status string
	var agentName string
	var cfg string
	flag.StringVar(&cfg, "c", "./app.json", "加载的配置,只有start时才有用")
	flag.StringVar(&status, "s", "start", "命令动作,start|stop|update")
	flag.StringVar(&agentName, "u", "", "更新agent的名称")
	flag.Parse()

	switch status {
	case "start":
		start(agentName, cfg)
	case "stop":
		stop(agentName)
	case "update":
		update(agentName)
		//todo 更新完成后应该停止并且启动对应的agent
	}
}
