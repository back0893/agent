package funcs

import (
	"agent/src/agent/model"
	"github.com/toolkits/nux"
)

/**
内存使用
mem.memtotal：内存总大小
mem.memused：使用了多少内存
*/

func MemMetrics() (*model.Memory, error) {
	m, err := nux.MemInfo()
	if err != nil {
		return nil, err
	}
	free := m.MemFree
	if m.MemAvailable > 0 {
		free = m.MemAvailable
	}
	return &model.Memory{
		Total: m.MemTotal,
		Used:  m.MemTotal - free,
		Free:  free,
	}, nil
}
