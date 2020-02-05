package mgr

import (
	"github.com/astaxie/beego"
	"tongserver.dataserver/service"
)

func init() {
	// jeda manage
	beego.Router("/jeda/menu", &JedaController{}, "get:GetMenu")
	beego.Router("/jeda/srvs", &JedaController{}, "get:GetSrvs")
	beego.Router("/jeda/ids", &JedaController{}, "get:GetIdsList")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})
	//JWT Request
	beego.Router("/token/verify", &service.SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &service.SecurityController{}, "post:CreateToken")

}
