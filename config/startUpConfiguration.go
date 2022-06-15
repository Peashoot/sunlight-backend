package config

// StartUpConfiguration 服务启动配置项
type StartUpConfiguration struct {
	ListenPort int                       `json:"port"`    // 监听端口
	LogPath    string                    `json:"logPath"` // 日志路径
	Mysql      StartUpDBConfiguration    `json:"mysql"`   // mysql配置
	Redis      StartUpCacheConfiguration `json:"redis"`   // redis配置
}

// StartUpDBConfiguration mysql配置参数
type StartUpDBConfiguration struct {
	Host           string `json:"host"`     // 连接地址
	Port           int    `json:"port"`     // 连接端口
	Password       string `json:"password"` // 连接密码
	Username       string `json:"username"` // 连接用户名
	DBName         string `json:"db"`       // 连接数据库名
	ConnectTimeout int    `json:"timeout"`  // 连接超时时间（单位：秒）
}

// StartUpCacheConfiguration Redis配置参数
type StartUpCacheConfiguration struct {
	Host     string `json:"host"`     // 连接地址
	Port     uint16 `json:"port"`     // 连接端口
	Password string `json:"password"` // 连接密码
	DBIndex  int    `json:"db"`       // 连接数据库ID
}
