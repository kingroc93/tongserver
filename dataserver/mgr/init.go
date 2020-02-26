package mgr

import (
	"github.com/astaxie/beego"
	"strings"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/service"
)

var JedaSrvContainer = make(service.SDefineContainerType)

// createDefaultDataIDs 注册默认数据源，这些数据源用于系统管理
func CreateDefaultDataIDs() {
	var meta map[string]interface{}
	tablesname := []string{"G_META", "G_META_ITEM", "JEDA_USER", "JEDA_ROLE", "JEDA_ROLE_USER", "JEDA_ORG", "JEDA_MENU",
		"G_USERSERVICE", "G_USERPROJECT", "G_SERVICE", "G_PROJECT", "G_IDS", "G_DATABASEURL"}
	for _, name := range tablesname {
		meta = map[string]interface{}{
			"name":      "default.mgr." + name,
			"inf":       "CreateWriteableTableDataSource",
			"dbalias":   "default",
			"tablename": name}
		datasource.IDSContainer[meta["name"].(string)] = meta
	}
}
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

	//Role Service
	//Org Service
	//DataBaseUrl Service
	//Meta Service
	//MetaItem Service
	//UserProject Service
	//UserService Service

}

func init() {
	createDefaultService()
	// jeda manage
	//beego.Router("/jeda/menu", &JedaController{}, "get:GetMenu")
	//beego.Router("/jeda/srvs", &JedaController{}, "get:GetSrvs")
	//beego.Router("/jeda/ids", &JedaController{}, "get:GetIdsList")
	//beego.Router("/jeda/reloadmeta", &JedaController{}, "get:ReloadMetaData")
	//beego.Router("/jeda/testdbconn", &JedaController{}, "get:Testdbconn")

	beego.Router("/jeda/?:context/?:action", &JedaController{}, "get,post:DoSrv")

	//JWT Request
	beego.Router("/token/verify", &SecurityController{}, "post:VerifyToken")
	beego.Router("/token/create", &SecurityController{}, "post:CreateToken")
	//qrcode
	beego.Router("/qrcode", &QrcodeController{})

}
