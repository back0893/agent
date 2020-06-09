package cron

import (
	"agent/src"
	"agent/src/agent/funcs"
	"time"
)

func InitDatHistory() {
	src.AddTimer(5*time.Second, func() {
		funcs.UpdateCpuStat()
	})
}
