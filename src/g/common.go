package g

import (
	"agent/src/g/model"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/back0893/goTcp/iface"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func Mkdir(path string) error {
	return os.Mkdir(path, 0755)
}
func GetCon(s iface.IServer, username string) (con iface.IConnection, has bool) {
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

func SavePid(pidfile string) {
	pid := os.Getpid()
	file, err := os.Create(pidfile)
	if err != nil {
		return
	}
	defer file.Close()
	io.WriteString(file, strconv.Itoa(pid))
}

func ReadPid(pidfile string) int {
	file, err := os.Open(pidfile)
	if err != nil {
		return 0
	}
	defer file.Close()
	data, _ := ioutil.ReadAll(file)
	data = bytes.Trim(data, "\r\n")
	pid, _ := strconv.Atoi(string(data))
	return pid
}

func Status(pid int) bool {
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

func Post(url string, data interface{}) ([]byte, error) {
	client := http.Client{Timeout: time.Second * 5}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewReader(jsonStr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, errors.New("请求返回非200")
	}
	result, _ := ioutil.ReadAll(resp.Body)
	return result, nil
}

/**
file 是一个绝对路径
*/
func Down(url, file string) error {
	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fp, err := os.Create(file)
	if err != nil {
		return err
	}
	defer fp.Close()
	if _, err := io.Copy(fp, resp.Body); err != nil {
		//出现错误,删除
		os.Remove(file)
		return err
	}
	return nil
}
