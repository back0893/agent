package model

//UpdateInfo 更新信息
type UpdateInfo struct {
	URL  string //更新地址
	Type int    //更新的类型 0=>全部 1=>只更新agent 2=>只更新配置文件
}

//UpdateResponse 更新回应
type UpdateResponse struct {
	//Status 更新成功
	Status  bool
	Message string
}
