package main

import (
	"log"
	"os"
	"os/exec"
)

func main() {
	//mdzz,,只能使用bash -c 使用对于从定向
	//nohup xxx & echo $!>pid 来记录执行的pid
	cmd := exec.Command("bash", "-c", "nohup redis-server >/dev/null 2>&1& echo $!>./pid")
	err := cmd.Run()
	if err != nil {
		log.Println(err.Error())
		log.Println(cmd.Args)
		os.Exit(1)
	}
	log.Println("process pid is ", cmd.Process.Pid)
	log.Println("===main exit===")

}
