package server

import (
	"github.com/back0893/goTcp/utils"
	"net"
	"time"
)

func SetTimeOut(connection net.Conn) {
	timeOut := time.Duration(utils.GlobalConfig.GetInt("heartTimeOut"))
	//每次心跳为connect设置新的过期时间,如果写入或者读取超过就会触发timeout的错误
	_ = connection.SetDeadline(time.Now().Add(time.Second * timeOut))
}
