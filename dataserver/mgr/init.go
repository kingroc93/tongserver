package mgr

import (
	"github.com/astaxie/beego"
	"strings"
	"tongserver.dataserver/service"
)

var JedaSrvContainer = make(service.SDefineContainerType)

func createService(cnt string, idsname string) *service.SDefine {
	ps := strings.Split(cnt, ".")
	projectid := "00000000-0000-0000-0000-000000000000"
	srv := &service.SDefine{
		ServiceId:   "JEDA_SRV_" + cnt,
		Context:     ps[1],
		BodyType:    "body",
		ServiceType: "",
		Namespace:   ps[0],
		Enabled:     true,
		MsgLog:      false,
		Security:    true,
		Meta:        "{\"ids\": \"" + idsname + "\"}",
		ProjectId:   projectid}
	return srv
}
func createDefaultService() {
	//User Service
	JedaSrvContainer["jeda.meta"] = createService("jeda.user", "default.mgr.G_META")
	JedaSrvContainer["jeda.metaitem"] = createService("jeda.user", "default.mgr.G_META_ITEM")
	JedaSrvContainer["jeda.user"] = createService("jeda.user", "default.mgr.JEDA_USER")
	JedaSrvContainer["jeda.role"] = createService("jeda.user", "default.mgr.JEDA_ROLE")
	JedaSrvContainer["jeda.roleuser"] = createService("jeda.user", "default.mgr.JEDA_ROLE_USER")
	JedaSrvContainer["jeda.org"] = createService("jeda.user", "default.mgr.JEDA_ORG")
	JedaSrvContainer["jeda.menu"] = createService("jeda.user", "default.mgr.JEDA_MENU")
	JedaSrvContainer["jeda.userservice"] = createService("jeda.user", "default.mgr.G_USERSERVICE")
	JedaSrvContainer["jeda.userproject"] = createService("jeda.user", "default.mgr.G_USERPROJECT")
	JedaSrvContainer["jeda.service"] = createService("jeda.user", "default.mgr.G_SERVICE")
	JedaSrvContainer["jeda.project"] = createService("jeda.user", "default.mgr.G_PROJECT")
	JedaSrvContainer["jeda.ids"] = createService("jeda.user", "default.mgr.G_IDS")
	JedaSrvContainer["jeda.databaseurl"] = createService("jeda.user", "default.mgr.G_DATABASEURL")
}

func init() {
	createDefaultService()
	beego.Router("/mgr/?:context/?:action", &JedaController{}, "get,post:DoSrv")
	//JWT Request
	beego.Router("/token/verify", &SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &SecurityController{}, "post:CreateToken")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})

}
