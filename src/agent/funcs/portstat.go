package funcs

import (
	"agent/src/agent/model"
	"errors"
	"github.com/toolkits/nux"
	"github.com/toolkits/slice"
)

/**
net.port.listen 端口监听状态
*/

func ListenTcpPortMetrics(ports ...int64) ([]*model.Port, error) {
	if len(ports) == 0 {
		return nil, errors.New("port empty")
	}
	ps, err := nux.TcpPorts()
	if err != nil {
		return nil, err
	}
	mvs := listenPort("tcp", ports, ps)
	return mvs, nil
}

func ListenUdpPortMetrics(ports ...int64) ([]*model.Port, error) {
	if len(ports) == 0 {
		return nil, errors.New("port empty")
	}
	ps, err := nux.UdpPorts()
	if err != nil {
		return nil, err
	}
	mvs := listenPort("udp", ports, ps)
	return mvs, nil
}

func listenPort(metric string, ports []int64, listenPorts []int64) []*model.Port {
	mvs := make([]*model.Port, 0, 10)
	for _, port := range ports {
		m := &model.Port{
			Type:   metric,
			Port:   port,
			Listen: false,
		}
		if slice.ContainsInt64(listenPorts, port) {
			m.Listen = true
		}
		mvs = append(mvs, m)
	}
	return mvs
}
