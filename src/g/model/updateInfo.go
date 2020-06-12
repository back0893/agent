package model

type UpdateInfo struct {
	Url     string //更新地址
	Type    int    //更新的类型 0=>全部 1=>只更新agent 2=>只更新配置文件
	Version string //版本tag
}
