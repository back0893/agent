package funcs

import (
	"agent/src/common/model"
	"github.com/toolkits/nux"
)

/**
机器负载
1,5,15Min内负载
*/

func LoadAvgMetrics() ([]*model.MetricValue, error) {
	load, err := nux.LoadAvg()
	if err != nil {
		return nil, err
	}
	return []*model.MetricValue{
		model.NewMetricValue("load.1min", load.Avg1min),
		model.NewMetricValue("load.5min", load.Avg5min),
		model.NewMetricValue("load.15min", load.Avg15min),
	}, nil
}
