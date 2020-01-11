package mgr

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/jeda/menu", &JedaController{}, "get:GetMenu")
	beego.Router("/jeda/srvs", &JedaController{}, "get:GetSrvs")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})
	//JWT Request
	beego.Router("/token/verify", &SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &SecurityController{}, "post:CreateToken")
	beego.Router("/token/refash", &SecurityController{}, "get:RefashToken")
}
