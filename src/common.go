package src

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
)

/**
一些公共的函数
*/
func SavePid(pidfile string) {
	pid := os.Getpid()
	file, err := os.Create(pidfile)
	if err != nil {
		return
	}
	defer file.Close()
	io.WriteString(file, strconv.Itoa(pid))
}

func ReadPid(pidfile string) int {
	file, err := os.Open(pidfile)
	if err != nil {
		return 0
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	data = bytes.Trim(data, "\r\n")
	pid, _ := strconv.Atoi(string(data))
	return pid
}

func Status(pid int) bool {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("ps -p %d |grep -v \"PID TTY\"|wc -l", pid))
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	out = bytes.Trim(out, "\r\n")
	wc, _ := strconv.Atoi(string(out))
	if wc > 0 {
		return true
	}
	return false
}
