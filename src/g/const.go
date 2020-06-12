package g

const (
	HeaderLength = 21 //包的固定长度
	AuthSuccess  = 1  //认真成功
	AuthFail     = 2  //认真失败
	STOP         = 4  //停止
	UPDATE       = 5  //更新

	//客户端请求id
	Services        = 6
	PING            = 3   //心跳
	Auth            = 100 //身份识别
	BaseServerInfo  = 101 //cpu使用
	HHD             = 102 //硬盘使用
	PortListen      = 104 //端口监听情况
	PortListenList  = 106 //请求监控的port
	ProcessNumList  = 107 //请求监控进程id
	MinePlugins     = 108 //请求当前agent配置的插件
	ServiceResponse = 109
	ActionNotice    = 110 //插件或者活动的执行后的

	//服务器端回应id
	Response               = 99  //通用回应
	PortListenListResponse = 306 //需要监控的port
	ProcessNumListResponse = 307 //需要监控进程id
	MinePluginsResponse    = 308 //当前agent配置的插件

	//agent的context传递key
	AGENT string = "agent"
)
