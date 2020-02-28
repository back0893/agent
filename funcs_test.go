package main

import (
	"agent/src/agent/funcs"
	"testing"
)

func TestCpu(t *testing.T) {
	funcs.UpdateCpuStat()
	funcs.UpdateCpuStat()
	ms := funcs.CpuMetrics()
	for _, m := range ms {
		t.Log(m.Value)
	}
}
func TestDisk(t *testing.T) {
	disk, _ := funcs.DiskUseMetrics()
	for fsFile, ms := range disk {
		t.Logf("%s has %d free %d,use %d", fsFile, ms[0].Value, ms[1].Value, ms[2].Value)
	}
}
func TestMem(t *testing.T) {
	mem, _ := funcs.MemMetrics()
	t.Logf(" total %d,use %d", mem[0].Value, mem[1].Value)
}

func TestLoadAvg(t *testing.T) {
	avg, _ := funcs.LoadAvgMetrics()
	for _, a := range avg {
		t.Log(a.Value)
	}
}
func TestPort(t *testing.T) {
	avg, err := funcs.ListenTcpPortMetrics(8412, 8001, 8000)
	if err != nil {
		t.Error(err)
	}
	for _, a := range avg {
		t.Log(a.Value)
	}
}
