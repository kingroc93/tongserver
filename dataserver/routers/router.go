package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"tongserver.dataserver/mgr"
	"tongserver.dataserver/service"
)

func RegisterRoutes() {

	//JWT Request
	logs.Info("启用令牌创建及验证服务")
	logs.Info("    /token/verify")
	logs.Info("    /token/create")
	logs.Info("    /jeda/user/?:cat")
	beego.Router("/token/verify", &mgr.SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &mgr.SecurityController{}, "post:CreateToken")
	beego.Router("/jeda/user/?:cat", &mgr.JedaController{}, "get,post:GetCurrentUserInfo")

	// 所有服务请求的入口函数
	beego.Router("/services/?:context/?:action", &service.SController{}, "get,post:DoSrv")

}
