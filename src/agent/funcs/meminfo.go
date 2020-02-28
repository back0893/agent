package funcs

import (
	"agent/src/common/model"
	"github.com/toolkits/nux"
)

/**
内存使用
mem.memtotal：内存总大小
mem.memused：使用了多少内存
*/

func MemMetrics() ([]*model.MetricValue, error) {
	m, err := nux.MemInfo()
	if err != nil {
		return nil, err
	}

	return []*model.MetricValue{
		model.NewMetricValue("mem.memtotal", m.MemFree),
		model.NewMetricValue("mem.memused", m.MemTotal-m.MemFree),
	}, nil
}
