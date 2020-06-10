package cmd

import (
	"bytes"
	"github.com/toolkits/sys"
	"os/exec"
	"syscall"
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
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	if err := cmd.Start(); err != nil {
		return err
	}

	err, isTimeout := sys.CmdRunWithTimeout(cmd, time.Duration(command.Timeout)*time.Millisecond)
	command.Callback(stdout, stderr, err, isTimeout)
	return nil
}
