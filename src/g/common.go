package g

import (
	"agent/src/agent/model"
	"github.com/back0893/goTcp/iface"
	"github.com/back0893/goTcp/net"
	"os"
)

func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}

func GetCon(s *net.Server, username string) (con iface.IConnection, has bool) {
	s.GetConnections().Range(func(key, value interface{}) bool {
		con = value.(iface.IConnection)
		data, ok := con.GetExtraData("auth")
		if ok == false {
			return true
		}
		auth := data.(*model.Auth)
		if auth.Username == username {
			has = true
			return false
		}
		return true
	})
	return con, has
}
