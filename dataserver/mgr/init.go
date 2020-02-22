package mgr

import (
	"github.com/astaxie/beego"
)

// RestResult 请求返回的数据类型
type RestResult map[string]interface{}

// CreateRestResult 创建返回数据实例
func CreateRestResult(success bool) RestResult {
	var result = make(RestResult)
	result["result"] = success
	return result
}

func init() {
	// jeda manage
	beego.Router("/jeda/menu", &JedaController{}, "get:GetMenu")
	beego.Router("/jeda/srvs", &JedaController{}, "get:GetSrvs")
	beego.Router("/jeda/ids", &JedaController{}, "get:GetIdsList")
	beego.Router("/jeda/reloadmeta", &JedaController{}, "get:ReloadMetaData")
	beego.Router("/jeda/testdbconn", &JedaController{}, "get:Testdbconn")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})
	//JWT Request
	beego.Router("/jeda/token/verify", &SecurityController{}, "post:VerifyToken")
	beego.Router("/jeda/token/create", &SecurityController{}, "post:CreateToken")

}
