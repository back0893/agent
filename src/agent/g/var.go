package g

import (
	"encoding/json"
	"log"
	"os"
	"sync"
)

var (
	portsListen     []int64
	porstListenLock = new(sync.RWMutex)
)

func SetPortListen(ports []int64) {
	porstListenLock.RLock()
	defer porstListenLock.RUnlock()
	portsListen = ports
}
func GetPortListen() []int64 {
	porstListenLock.Lock()
	defer porstListenLock.Unlock()
	return portsListen
}

//SavePort 保存端口
func SavePort() {
	porstListenLock.RLock()
	defer porstListenLock.RUnlock()
	if len(portsListen) == 0 {
		return
	}
	fp, err := os.Create("ports.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()
	encoder := json.NewEncoder(fp)
	encoder.Encode(portsListen)
}

//LoadPort 加载端口
func LoadPort() {
	porstListenLock.RLock()
	defer porstListenLock.RUnlock()
	fp, err := os.Open("ports.json")
	if err != nil {
		log.Println(err)
		return
	}
	defer fp.Close()
	decoder := json.NewDecoder(fp)
	if err := decoder.Decode(&portsListen); err != nil {
		log.Println(err)
		return
	}
}
