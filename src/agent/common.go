package agent

import (
	"github.com/back0893/goTcp/utils"
	"net"
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
