package services

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

type RedisService struct {
}

func (r RedisService) Status() bool {
	pid := r.GetPid()
	if pid == 0 {
		return false
	}
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

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (r RedisService) GetPid() int {
	file, err := os.Open("./pid")
	if err != nil {
		return 0
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return 0
	}
	pid, err := strconv.Atoi(string(bytes.Trim(data, "\n\r")))
	if err != nil {
		return 0
	}
	return pid
}
func (r RedisService) Start() error {
	if r.Status() {
		return errors.New("redis已经运行")
	}
	cmd := exec.Command("bash", "-c", "nohup redis-server >/dev/null 2>&1& echo $!>./pid")
	return cmd.Run()
}

func (r RedisService) Stop() error {
	pid := r.GetPid()
	if pid == 0 {
		return errors.New("redis灭有在运行")
	}
	syscall.Kill(pid, syscall.SIGKILL)
	//参数pid
	os.Remove("./pid")
	return nil
}

func (r RedisService) Restart() error {
	if err := r.Stop(); err != nil {
		return err
	}
	if err := r.Start(); err != nil {
		return err
	}
	return nil
}

func (r *RedisService) info(address string, auth string) (map[string]string, error) {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		return nil, errors.New("redis连接失败")
	}
	defer c.Close()

	if auth != "" {
		if _, err = c.Do("auth", auth); err != nil {
			return nil, errors.New("redis密码错误")
		}
	}
	_, err = c.Do("ping")
	if err != nil {
		return nil, err
	}
	info, err := c.Do("info")
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(info.([]byte))
	scaner := bufio.NewScanner(reader)
	scaner.Split(bufio.ScanLines)
	infoMap := make(map[string]string)
	for scaner.Scan() {
		line := scaner.Bytes()
		p := bytes.Split(line, []byte{':'})
		if len(p) != 2 {
			continue
		}
		infoMap[string(p[0])] = string(p[1])
	}
	return infoMap, nil
}
