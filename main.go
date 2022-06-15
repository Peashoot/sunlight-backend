package main

import (
	"fmt"

	"github.com/kataras/iris/v12"
	"github.com/peashoot/sunlight/api"
	"github.com/peashoot/sunlight/auth"
	"github.com/peashoot/sunlight/config"
	"github.com/peashoot/sunlight/log"
)

func main() {
	app := iris.New()
	// 日志配置
	log.Init(app)
	// JWT认证配置
	auth.Init(app)
	// 路由配置
	api.Init(app)
	// 启动服务
	app.Listen(fmt.Sprintf(":%d", config.ListenPort))
}
