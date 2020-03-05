package g

import (
	"agent/src/agent/model"
	"bytes"
	"encoding/gob"
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

func EncodeData(e interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	if err := encoder.Encode(e); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func DecodeData(data []byte, e interface{}) error {
	buffer := bytes.NewReader(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(e); err != nil {
		return err
	}
	return nil
}
