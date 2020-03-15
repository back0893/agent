package funcs

import (
	model2 "agent/src/g/model"
	"github.com/toolkits/nux"
)

/**
内存使用
mem.memtotal：内存总大小
mem.memused：使用了多少内存
*/

func MemMetrics() (*model2.Memory, error) {
	m, err := nux.MemInfo()
	if err != nil {
		return nil, err
	}
	free := m.MemFree
	if m.MemAvailable > 0 {
		free = m.MemAvailable
	}
	return &model2.Memory{
		Total: m.MemTotal,
		Used:  m.MemTotal - free,
		Free:  free,
	}, nil
}
