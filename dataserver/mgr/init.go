package mgr

import (
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/jeda/user/?:cat", &JedaController{}, "get,post:GetCurrentUserInfo")
	//JWT Request
	beego.Router("/token/verify", &SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &SecurityController{}, "post:CreateToken")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})
}
