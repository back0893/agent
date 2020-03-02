package services

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/gomodule/redigo/redis"
	"os/exec"
)

type RedisService struct {
}

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (r RedisService) Start() error {
	cmd := exec.Command("bash", "-c", "nohup redis-server >/dev/null 2>&1 & echo &!>./redisPid")
	return cmd.Run()
}

func (r RedisService) Stop() {
	//防止因为人工或者其他原因导致redis的pid改变
}

func (r RedisService) Restart() {
	panic("implement me")
}

func (r *RedisService) Status(address string, auth string) (map[string]string, error) {
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
