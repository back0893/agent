package funcs

import (
	"agent/src/common/model"
	"github.com/toolkits/nux"
)

/**
硬盘相关
df.bytes.free 磁盘可用量
df.bytes.total：磁盘总大小
df.bytes.used：磁盘已用大小
*/

func DiskUseMetrics() (map[string][]*model.MetricValue, error) {
	mounts, err := nux.ListMountPoint()
	if err != nil {
		return nil, err
	}
	ms := make(map[string][]*model.MetricValue)
	for _, mount := range mounts {
		deviceUsage, err := nux.BuildDeviceUsage(mount[0], mount[1], mount[2])
		if err != nil {
			return nil, err
		}
		mvs := make([]*model.MetricValue, 0)
		mvs = append(mvs, model.NewMetricValue("df.bytes.total", deviceUsage.BlocksAll))
		mvs = append(mvs, model.NewMetricValue("df.bytes.free", deviceUsage.BlocksFree))
		mvs = append(mvs, model.NewMetricValue("df.bytes.used", deviceUsage.BlocksUsed))
		ms[deviceUsage.FsFile] = mvs
	}
	return ms, nil
}
