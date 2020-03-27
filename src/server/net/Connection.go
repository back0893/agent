package net

import (
	"github.com/back0893/goTcp/net"
	"github.com/back0893/goTcp/utils"
	"time"
)

type Connection struct {
	net.Connection
}

func (c *Connection) UpdateTimeOut() {
	interval := utils.GlobalConfig.GetInt("heartTimeOut")
	if interval == 0 {
		interval = 15
	}
	timeOut := time.Duration(interval)
	//每次心跳为connect设置新的过期时间,如果写入或者读取超过就会触发timeout的错误
	_ = c.GetRawCon().SetDeadline(time.Now().Add(time.Second * timeOut))
}
