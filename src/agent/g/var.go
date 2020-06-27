package g

import "sync"

var (
	portsListen     []int64
	porstListenLock = new(sync.RWMutex)
)

func SetPortListen(ports []int64) {
	porstListenLock.RLock()
	defer porstListenLock.RUnlock()
	portsListen = append(portsListen, ports...)
}
func GetPortListen() []int64 {
	porstListenLock.Lock()
	defer porstListenLock.Unlock()
	return portsListen
}
