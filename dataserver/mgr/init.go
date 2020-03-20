package mgr

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"strings"
	"tongserver.dataserver/service"
)

var JedaSrvContainer = make(service.SDefineContainerType)

// 创建用于系统管理的默认服务
func createService(cnt string, meta map[string]interface{}) *service.SDefine {
	jsonStr, _ := json.Marshal(meta)

	ps := strings.Split(cnt, ".")
	projectid := "00000000-0000-0000-0000-000000000000"
	sec := beego.BConfig.RunMode == beego.PROD
	srv := &service.SDefine{
		ServiceId:   "JEDA_SRV_" + cnt,
		Context:     ps[1],
		BodyType:    "body",
		ServiceType: "",
		Namespace:   ps[0],
		Enabled:     true,
		MsgLog:      false,
		Security:    sec || true,
		Meta:        string(jsonStr),
		ProjectId:   projectid}
	return srv
}

// 创建用于系统管理的默认服务
func createDefaultService() {
	//User Service
	JedaSrvContainer["jeda.meta"] = createService("jeda.meta", map[string]interface{}{
		"ids": "default.mgr.G_META",
		"userfilter": map[string]string{
			"filterkey": "PROJECTID",
			"opera":     "in",
			"ids":       "default.mgr.G_USERPROJECT",
			"userfield": "USERID",
			"joinfield": "PROJECTNAME"},
	})

	JedaSrvContainer["jeda.metaitem"] = createService("jeda.metaitem", map[string]interface{}{
		"ids": "default.mgr.G_META_ITEM"})

	JedaSrvContainer["jeda.user"] = createService("jeda.user", map[string]interface{}{
		"ids": "default.mgr.JEDA_USER"})

	JedaSrvContainer["jeda.role"] = createService("jeda.role", map[string]interface{}{
		"ids": "default.mgr.JEDA_ROLE"})

	JedaSrvContainer["jeda.roleuser"] = createService("jeda.roleuser", map[string]interface{}{
		"ids": "default.mgr.JEDA_ROLE_USER"})

	JedaSrvContainer["jeda.org"] = createService("jeda.org", map[string]interface{}{
		"ids": "default.mgr.JEDA_ORG"})

	JedaSrvContainer["jeda.menu"] = createService("jeda.menu", map[string]interface{}{
		"ids": "default.mgr.JEDA_MENU"})

	JedaSrvContainer["jeda.userservice"] = createService("jeda.userservice", map[string]interface{}{
		"ids": "default.mgr.G_USERSERVICE"})

	JedaSrvContainer["jeda.userproject"] = createService("jeda.userproject", map[string]interface{}{
		"ids": "default.mgr.G_USERPROJECT"})

	JedaSrvContainer["jeda.service"] = createService("jeda.service", map[string]interface{}{
		"ids": "default.mgr.G_SERVICE"})

	JedaSrvContainer["jeda.project"] = createService("jeda.user", map[string]interface{}{
		"ids": "default.mgr.G_PROJECT"})

	JedaSrvContainer["jeda.ids"] = createService("jeda.user", map[string]interface{}{
		"ids": "default.mgr.G_IDS"})

	JedaSrvContainer["jeda.databaseurl"] = createService("jeda.user", map[string]interface{}{
		"ids": "default.mgr.G_DATABASEURL"})
}

func init() {
	createDefaultService()
	beego.Router("/mgr/?:context/?:action", &JedaController{}, "get,post:DoSrv")
	beego.Router("/meta/?:context", &JedaController{}, "get,post:GetMeta")
	beego.Router("/jeda/user/?:cat", &JedaController{}, "get,post:GetCurrentUserInfo")
	//JWT Request
	beego.Router("/token/verify", &SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &SecurityController{}, "post:CreateToken")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})

}
