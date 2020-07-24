package agent

import (
	"agent/src/g"
	"net"

	"github.com/back0893/goTcp/utils"
)

func ConnectServer(cfg string) (*net.TCPConn, error) {
	utils.GlobalConfig.Load("json", cfg)
	return g.ConnectTcp(utils.GlobalConfig.GetString("ip"), utils.GlobalConfig.GetString("port"))
}
