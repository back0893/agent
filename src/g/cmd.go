package g

import (
	"bytes"
	"log"
	"os/exec"
	"time"

	"github.com/toolkits/sys"
)

type Command struct {
	Name     string
	Args     []string
	Timeout  int
	Dir      string
	Callback func(stdout, stderr *bytes.Buffer, err error, isTimeout bool)
}

func (command *Command) Run() {
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Dir = command.Dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	log.Println(cmd.Args)
	//不设置成独立进程,方便在agent退出时,一起退出
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		command.Callback(&stdout, &stderr, err, false)
	}

	if command.Timeout <= 0 {
		//如果超时小于0那么不需要超时
		err := cmd.Wait()
		command.Callback(&stdout, &stderr, err, false)
	} else {
		err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(command.Timeout)*time.Millisecond)
		command.Callback(&stdout, &stderr, err, isTimeout)
	}
}
