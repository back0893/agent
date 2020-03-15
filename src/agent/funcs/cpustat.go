package funcs

import (
	model2 "agent/src/g/model"
	"github.com/toolkits/nux"
	"sync"
)

/**
监控cpu使用
cpu.idel
cpu.busy =100-cpu.idel
*/
const historyCount int = 2

var (
	procStatHistory [historyCount]*nux.ProcStat
	psLock          = new(sync.RWMutex)
)

func deltaTotal() uint64 {
	if procStatHistory[1] == nil {
		return 0
	}
	return procStatHistory[0].Cpu.Total - procStatHistory[1].Cpu.Total
}

func CpuIdle() float64 {
	psLock.RLock()
	defer psLock.RUnlock()
	dt := deltaTotal()
	if dt == 0 {
		return 0.0
	}
	invQuotient := 100.00 / float64(dt)
	return float64(procStatHistory[0].Cpu.Idle-procStatHistory[1].Cpu.Idle) * invQuotient
}
func CpuPrepared() bool {
	psLock.RLock()
	defer psLock.RUnlock()
	return procStatHistory[1] != nil
}

func UpdateCpuStat() error {
	ps, err := nux.CurrentProcStat()
	if err != nil {
		return err
	}
	psLock.Lock()
	defer psLock.Unlock()
	procStatHistory[1] = procStatHistory[0]
	procStatHistory[0] = ps
	return nil
}

func CpuMetrics() *model2.Cpu {
	cpuIdle := CpuIdle()
	return &model2.Cpu{
		Idle: cpuIdle,
		Busy: 100 - cpuIdle,
	}
}
