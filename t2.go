package main

import (
	"agent/src/agent/services"
	"errors"
	"flag"
	"fmt"
)

func main() {
	//mdzz,,只能使用bash -c 使用对于从定向
	//nohup xxx & echo $!>pid 来记录执行的pid
	//cmd := exec.Command("bash", "-c", "nohup redis-server >/dev/null 2>&1& echo $!>./pid")
	//err := cmd.Run()
	//if err != nil {
	//	log.Println(err.Error())
	//	log.Println(cmd.Args)
	//	os.Exit(1)
	//}
	//log.Println("process pid is ", cmd.Process.Pid)
	//log.Println("===main exit===")

	action := flag.String("s", "start", "start,stop,status,restart,help")
	flag.Parse()
	redis := services.RedisService{}
	var err error
	switch *action {
	case "start":
		err = redis.Start()
	case "stop":
		err = redis.Stop()
	case "status":
		pid := redis.GetPid()
		if pid > 0 {
			fmt.Println("redis正在运行")
		} else {
			err = errors.New("redis没有运行")
		}
	case "restart":
		err = redis.Restart()
	default:
		flag.Usage()
	}
	if err != nil {
		fmt.Println("运行错误", err.Error())
	}
}
