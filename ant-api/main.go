package main

import (
	_ "ant-api/routers"
	_ "ant-api/workers"
	_ "github.com/astaxie/beego/session/redis"

	"ant-api/auth"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	"github.com/wonderivan/logger"
)

func main() {

	// 允许跨域访问
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			AllowCredentials: true,

	}))

	// 登录验证
	beego.InsertFilter("*", beego.BeforeRouter, auth.FilterAuth)

	// 日志配置
	logger.SetLogger(beego.AppConfig.String("LogJson"))

	// 开启dev模式
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}

	beego.Run()

}
