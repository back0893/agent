package funcs

import (
	model2 "agent/src/g/model"
	"github.com/toolkits/nux"
)

/**
硬盘相关
df.bytes.free 磁盘可用量
df.bytes.total：磁盘总大小
df.bytes.used：磁盘已用大小
*/

func DiskUseMetrics() ([]*model2.Disk, error) {
	mounts, err := nux.ListMountPoint()
	if err != nil {
		return nil, err
	}
	disks := make([]*model2.Disk, 0)
	for _, mount := range mounts {
		deviceUsage, err := nux.BuildDeviceUsage(mount[0], mount[1], mount[2])
		if err != nil {
			return nil, err
		}
		disk := &model2.Disk{
			FsFile: deviceUsage.FsFile,
			Free:   deviceUsage.BlocksFree,
			Total:  deviceUsage.BlocksAll,
			Used:   deviceUsage.BlocksUsed,
		}
		disks = append(disks, disk)
	}

	return disks, nil
}
