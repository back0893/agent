package funcs

import (
	"agent/src/common/model"
	"errors"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
)

/**
net.port.listen 端口监听状态
*/

func ListenTcpPortMetrics(ports ...int64) ([]*model.MetricValue, error) {
	if len(ports) == 0 {
		return nil, errors.New("port empty")
	}
	ps, err := nux.TcpPorts()
	if err != nil {
		return nil, err
	}
	mvs := listenPort("net.tcp.port", ports, ps)
	return mvs, nil
}

func ListenUdpPortMetrics(ports ...int64) ([]*model.MetricValue, error) {
	if len(ports) == 0 {
		return nil, errors.New("port empty")
	}
	ps, err := nux.UdpPorts()
	if err != nil {
		return nil, err
	}
	mvs := listenPort("net.udp.port", ports, ps)
	return mvs, nil
}

func listenPort(metric string, ports []int64, listenPorts []int64) []*model.MetricValue {
	mvs := make([]*model.MetricValue, 0, 10)
	for _, port := range ports {
		if slice.ContainsInt64(listenPorts, port) {
			mvs = append(mvs, model.NewMetricValue(metric, port))
		}
	}
	return mvs
}
