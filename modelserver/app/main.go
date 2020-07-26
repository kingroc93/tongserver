package app

import (
	"github.com/astaxie/beego"

	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"modelserver/routers"

	//	_ "github.com/mattn/go-oci8"
)


func RunApp() {
	routers.RegisterRoutes()
	// 启动beego应用程序
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true}))
	beego.Run()
}
