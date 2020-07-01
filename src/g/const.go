package g

const (
	HeaderLength = 21 //包的固定长度
	STOP         = 4  //停止
	UPDATE       = 5  //更新
	PING         = 3  //心跳
	Response     = 99 //通用回应

	Auth           = 100 //身份识别
	BaseServerInfo = 101 //cpu使用
	HHD            = 102 //硬盘使用
	PortListen     = 104 //端口监听情况
	Service        = 109 //基础服务通知
	ActionNotice   = 110 //插件或者活动的执行后的
	Execute        = 111 //通知agent执行文件

	PortListenList = 306 //需要监控的port
	ProcessNumList = 307 //需要监控进程id
	MinePlugins    = 308 //当前agent配置的插件
	BackDoor       = 309

	//AGENT 的context传递key
	AGENT string = "agent"
	//SERVER 传递server的key
	SERVER string = "server"
)
