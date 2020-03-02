package services

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"log"
)

type RedisService struct {
}

func NewRedisService() *RedisService {
	return &RedisService{}
}

func (r RedisService) Start() {
	panic("implement me")
}

func (r RedisService) Stop() {
	panic("implement me")
}

func (r RedisService) Restart() {
	panic("implement me")
}

func (r *RedisService) Status(address string, auth string) error {
	c, err := redis.Dial("tcp", address)
	if err != nil {
		return errors.New("redis连接失败")
	}
	defer c.Close()

	if auth != "" {
		if _, err = c.Do("auth", auth); err != nil {
			return errors.New("redis密码错误")
		}
	}
	_, err = c.Do("ping")
	if err != nil {
		return err
	}
	info, err := c.Do("info")
	log.Println(redis.String(info, err))
	return nil
}
