package g

const (
	HeaderLength = 21 //包的固定长度
	PING         = 1  //心跳
	STOP         = 2  //停止
	UPDATE       = 3  //更新
	Response     = 99 //通用回应

	Auth       = 100 //身份识别
	CPU        = 101 //cpu使用
	HHD        = 102 //硬盘使用
	MEM        = 103 //内存使用
	LoadAvg    = 104 //负载
	PortListen = 105 //端口监听情况

	Service         = 201 //对于service的命令
	ServiceResponse = 301 //service执行后的回应

	//agent的context传递key
	AGENT string = "agent"
)
