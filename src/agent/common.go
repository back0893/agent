package agent

import (
	"github.com/back0893/goTcp/utils"
	"net"
	"strconv"
	"time"
)

func ConnectServer(cfg string) (*net.TCPConn, error) {
	utils.GlobalConfig.Load("json", cfg)
	host := net.JoinHostPort(utils.GlobalConfig.GetString("ip"), utils.GlobalConfig.GetString("port"))
	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return nil, err
	}
	con, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, err
	}
	return con, nil
}

func GetInterval(args map[string]string, def time.Duration) time.Duration {
	v, ok := args["interval"]
	if ok {
		if m, err := strconv.Atoi(v); err == nil {
			return time.Duration(m)
		}
	}
	return def
}
