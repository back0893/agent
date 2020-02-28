package model

type MetricValue struct {
	Metric string
	Value  interface{}
}

func NewMetricValue(metric string, value interface{}) *MetricValue {
	return &MetricValue{
		Metric: metric,
		Value:  value,
	}
}
