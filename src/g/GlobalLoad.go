package g

import (
	"fmt"
	"github.com/back0893/goTcp/utils"
	"github.com/spf13/cast"
	"log"
	"os"
	"strings"
	"time"
)

func runtimeDir() (string, error) {
	var path string
	p := utils.GlobalConfig.Get("runtime")
	if p == nil {
		path = "./runtime"
	} else {
		path = cast.ToString(p)
	}
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) == false {
			err := Mkdir(path)
			if err != nil {
				return "", err
			}
		} else {
			return "", err
		}
	}
	filePath := fmt.Sprintf("%s/%s.log", strings.TrimRight(path, "/"), time.Now().Format("2006-01-02"))
	return filePath, nil
}

func setLogWrite() {
	file, err := runtimeDir()
	if err != nil {
		log.Println("无法接入日志")
		return
	}
	fileWrite, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("无法接入日志")
		return
	}
	log.SetOutput(fileWrite)
}
func LoadInit(file string) {
	utils.GlobalConfig.Load("json", file)
	setLogWrite()
}
