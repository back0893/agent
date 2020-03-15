package funcs

import (
	model2 "agent/src/g/model"
	"github.com/toolkits/nux"
)

/**
机器负载
1,5,15Min内负载
*/

func LoadAvgMetrics() ([]*model2.LoadAvg, error) {
	load, err := nux.LoadAvg()
	if err != nil {
		return nil, err
	}
	return []*model2.LoadAvg{
		{
			Name: "1min",
			Load: load.Avg1min,
		},
		{
			Name: "5min",
			Load: load.Avg5min,
		},
		{
			Name: "15min",
			Load: load.Avg15min,
		},
	}, nil
}
