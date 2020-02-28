package funcs

import (
	"agent/src/agent/model"
	"github.com/toolkits/nux"
)

/**
硬盘相关
df.bytes.free 磁盘可用量
df.bytes.total：磁盘总大小
df.bytes.used：磁盘已用大小
*/

func DiskUseMetrics() ([]*model.Disk, error) {
	mounts, err := nux.ListMountPoint()
	if err != nil {
		return nil, err
	}
	disks := make([]*model.Disk, 0)
	for _, mount := range mounts {
		deviceUsage, err := nux.BuildDeviceUsage(mount[0], mount[1], mount[2])
		if err != nil {
			return nil, err
		}
		disk := &model.Disk{
			FsFile: deviceUsage.FsFile,
			Free:   deviceUsage.BlocksFree,
			Total:  deviceUsage.BlocksAll,
			Used:   deviceUsage.BlocksUsed,
		}
		disks = append(disks, disk)
	}

	return disks, nil
}
