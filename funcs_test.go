package main

import (
	"agent/src/agent/funcs"
	"testing"
)

func TestCpu(t *testing.T) {
	funcs.UpdateCpuStat()
	funcs.UpdateCpuStat()
	ms := funcs.CpuMetrics()
	t.Logf("idel=>%.2f,busy=>%.2f", ms.Idle, ms.Busy)
}
func TestDisk(t *testing.T) {
	disk, _ := funcs.DiskUseMetrics()
	for _, ms := range disk {
		t.Logf("%s has %d free %d,use %d", ms.FsFile, ms.Total, ms.Free, ms.Used)
	}
}
func TestMem(t *testing.T) {
	mem, _ := funcs.MemMetrics()
	t.Logf(" total %d,use %d", mem.Total, mem.Used)
}

func TestLoadAvg(t *testing.T) {
	avg, _ := funcs.LoadAvgMetrics()
	for _, a := range avg {
		t.Log(a.Load)
	}
}
func TestPort(t *testing.T) {
	avg, err := funcs.ListenTcpPortMetrics(8412, 8001, 8000)
	if err != nil {
		t.Error(err)
	}
	for _, a := range avg {
		t.Log(a.Port)
	}
}
