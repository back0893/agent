package main

import (
	"context"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func startServer() (*exec.Cmd, context.CancelFunc) {
	ctx, canncel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "./server")
	if err := cmd.Start(); err != nil {
		panic(err)
	}
	return cmd, canncel
}
func main() {
	cmd, cancel := startServer()
	re := make(chan struct{})
	go func() {
		for {
			//子进程意外退出
			if err := cmd.Wait(); err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					if exitErr.ProcessState.ExitCode() == 0 || exitErr.ProcessState.ExitCode() == -1 {
						log.Println("正常退出1")
						re <- struct{}{}
						return
					}
				}
				cmd, cancel = startServer()
				log.Println("意外退出")
				return
			} else {
				log.Println("正常退出2")
				re <- struct{}{}
				return
			}

		}
	}()
	//接收到sign后重启server
	cha := make(chan os.Signal)
	signal.Notify(cha, syscall.SIGALRM)
	log.Println(syscall.Getpid())
	for {
		select {
		case s := <-cha:
			switch s {
			case syscall.SIGALRM:
				cancel()
				//这里等待,父进程会马上退出,导致子进程灭有接收到关闭信号,而被init接管
				<-re
				log.Println("退出")
				return
			}
		}
	}
}
