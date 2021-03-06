package g

const (
	HeaderLength = 21 //包的固定长度
	AuthSuccess  = 1
	AuthFail     = 2
	PING         = 3 //心跳
	STOP         = 4 //停止
	UPDATE       = 5 //更新
	Services     = 6
	Response     = 99 //通用回应

	Auth           = 100 //身份识别
	BaseServerInfo = 101 //cpu使用
	HHD            = 102 //硬盘使用
	PortListen     = 104 //端口监听情况
	ServicesList   = 105 //下发默认启动服务

	Service         = 201 //对于service的命令
	ServiceResponse = 301 //service执行后的回应

	//agent的context传递key
	AGENT string = "agent"
	//全局产量,当前agent的版本
	VERSION = 4

	/**
	下发的service的id映射服务
	*/
	REDISSERVICE = 1
)
