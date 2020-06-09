package model

/**
插件分为2种
http下载,http下载使用url中token作为用户判断.
git拉取,拉取地址和用户作为用户判断
*/
type Plugins struct {
	Type     int //作为判断为git还是http
	Uri      string
	Name     string //插件名称
	Interval int    //执行时间间隔(秒)
	Branch   string //git的分或者tag名称
}
