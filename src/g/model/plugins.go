package model

/**
插件分为2种
git拉取,拉取地址和用户作为用户判断
使用gitlab的api创建仓库和新增删除文件
可是
*/
type Plugins struct {
	Uri []string
}

/**
插件的标准返回
*/
type MetricValue struct {
	Id        int32       `json:"id"`        //插件编号
	Data      interface{} `json:"value"`     //插件数据
	Timestamp int64       `json:"timestamp"` //插件返回时间
}
