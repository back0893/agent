package g

import (
	"agent/src/g/model"
	"bytes"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/back0893/goTcp/iface"
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

func EncodeData(e ...interface{}) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	encoder := gob.NewEncoder(buffer)
	for _, tmp := range e {
		if err := encoder.Encode(tmp); err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

func DecodeData(data []byte, e ...interface{}) error {
	buffer := bytes.NewReader(data)
	decoder := gob.NewDecoder(buffer)
	for _, tmp := range e {
		if err := decoder.Decode(tmp); err != nil {
			return err
		}
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
	client := http.Client{Timeout: time.Minute * 10}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		//如果非200返回,说明是有问题的更新
		return errors.New("非200下载")
	}
	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
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

func GetInterval(args map[string]string, def time.Duration) time.Duration {
	v, ok := args["interval"]
	if ok {
		if m, err := strconv.Atoi(v); err == nil {
			return time.Duration(m)
		}
	}
	return def
}
func CSTLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		loc = time.FixedZone("Asia/Shanghai", 8*3600)
	}
	return loc
}
func CSTTime() time.Time {
	return time.Now().In(CSTLocation())
}

func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}

func UniqueInt64(nums []int64) []int64 {
	unique := make([]int64, 0)
	for _, num := range nums {
		if FindInInt64(unique, num) == -1 {
			unique = append(unique, num)
		}
	}
	log.Println(unique)
	return unique
}

//FindInInt64 寻找num在nums中 不存在-1
func FindInInt64(nums []int64, num int64) int {
	for i := 0; i < len(nums); i++ {
		if nums[i] == num {
			return i
		}
	}
	return -1
}
