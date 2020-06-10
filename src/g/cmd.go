package g

import (
	"bytes"
	"github.com/toolkits/sys"
	"log"
	"os/exec"
	"time"
)

type Command struct {
	Name     string
	Args     []string
	Timeout  int
	Dir      string
	Callback func(stdout, stderr bytes.Buffer, err error, isTimeout bool)
}

func (command *Command) Run() error {
	cmd := exec.Command(command.Name, command.Args...)
	cmd.Dir = command.Dir
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	log.Println(cmd.Args)
	//cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}

	if command.Timeout <= 0 {
		//如果超时小于0那么不需要超时
		err := cmd.Wait()
		command.Callback(stdout, stderr, err, false)
	} else {
		err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(command.Timeout)*time.Millisecond)
		command.Callback(stdout, stderr, err, isTimeout)
	}
	return nil
}
