package mgr

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/skip2/go-qrcode"
	"strconv"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/service"
	"tongserver.dataserver/utils"
)

// JedaController 后台管理控制器
type JedaController struct {
	beego.Controller
	ControllerWithVerify
}

// QrcodeController 生成二维码的控制器
type QrcodeController struct {
	beego.Controller
	ControllerWithVerify
}

// reloadMetaFun 重新加载元数据的函数句柄类型
type reloadMetaFun func() error

// metaFuns 用于加载元数据的函数列表
var metaFuns = make(map[string]reloadMetaFun)

// AddMetaFuns 添加加载元数据的函数句柄
func AddMetaFuns(name string, f reloadMetaFun) {
	metaFuns[name] = f
}

// Get 生成二维码
func (c *QrcodeController) Get() {
	cnt := c.Input().Get("c")
	bs := c.Input().Get("t")
	size, errw := strconv.Atoi(c.Input().Get("s"))
	if errw != nil {
		size = 256
	}
	c.Ctx.Output.ContentType("png")
	var png []byte
	if bs == "64" {
		cnt = utils.DecodeURLBase64(cnt)
	}
	png, _ = qrcode.Encode(cnt, qrcode.Medium, size)
	c.Ctx.Output.Body(png)
}

// GetSrvs 返回所有服务
func (c *JedaController) GetSrvs() {

	c.Data["json"] = datasource.IDSContainer
	c.ServeJSON()
}

// ReloadMetaData 重新加载系统元数据
func ReloadMetaData() error {
	for _, f := range metaFuns {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

// ReloadMetaData 重新加载元数据
func (c *JedaController) ReloadMetaData() {
	if !c.ControllerWithVerify.Verifty(&c.Controller) {
		return
	}
	err := ReloadMetaData()
	if err != nil {
		logs.Debug("加载系统元数据时发生错误,%v", err.Error())
		utils.CreateErrorResponseByError(err, &c.Controller)
	}
	r := utils.CreateRestResult(true)
	r["msg"] = "重新加载成功"
	c.ServeJSON()
}

// commonCheckGetSrvList commonCheckGetSrvList
func (c *JedaController) commonCheckGetSrvList() bool {
	f := metaFuns["ids"]
	if f == nil {
		r := utils.CreateRestResult(false)
		r["msg"] = "没有找到名称为ids的元数据加载函数"
		c.Data["json"] = r
		c.ServeJSON()
		return true
	}
	err := f()
	if err != nil {
		r := utils.CreateRestResult(false)
		r["msg"] = err.Error()
		c.Data["json"] = r
		c.ServeJSON()
		return true
	}
	return false
}

func (c *JedaController) Testdbconn() {
	if !c.ControllerWithVerify.Verifty(&c.Controller) {
		return
	}
}

// GetIdsList 返回数据源列表
func (c *JedaController) GetIdsList() {
	if !c.ControllerWithVerify.Verifty(&c.Controller) {
		return
	}
	if c.commonCheckGetSrvList() {
		return
	}
	ids := datasource.IDSContainer
	var result = &datasource.DataResultSet{}
	result.Fields = make(datasource.FieldDescType)
	result.Fields["ID"] = &datasource.FieldDesc{
		FieldType: datasource.PropertyDatatypeStr,
		Index:     0,
		Meta:      &map[string]string{"CAP": "编号"}}
	result.Fields["IdsName"] = &datasource.FieldDesc{
		FieldType: datasource.PropertyDatatypeStr,
		Index:     1,
		Meta:      &map[string]string{"CAP": "数据源名称"}}
	result.Fields["TableName"] = &datasource.FieldDesc{
		FieldType: datasource.PropertyDatatypeStr,
		Index:     2,
		Meta:      &map[string]string{"CAP": "表名"}}
	result.Fields["DbAlias"] = &datasource.FieldDesc{
		FieldType: datasource.PropertyDatatypeStr,
		Index:     3,
		Meta:      &map[string]string{"CAP": "数据库别名"}}

	result.Fields["Writeable"] = &datasource.FieldDesc{
		FieldType: datasource.PropertyDatatypeStr,
		Index:     4,
		Meta:      &map[string]string{"CAP": "是否可写"}}
	result.Data = make([][]interface{}, 0, len(ids))
	for k, v := range ids {
		if v["inf"].(string) != "CreateTableDataSource" && v["inf"].(string) != "CreateWriteableTableDataSource" {
			continue
		}
		row := make([]interface{}, 5, 5)
		row[0] = k
		row[1] = v["name"].(string)
		row[2] = v["tablename"].(string)
		row[3] = v["dbalias"].(string)
		if v["inf"].(string) != "CreateWriteableTableDataSource" {
			row[4] = "true"
		}
		if v["inf"].(string) != "CreateTableDataSource" {
			row[4] = "true"
		}
		result.Data = append(result.Data, row)
	}
	r := utils.CreateRestResult(true)
	r["resultset"] = result
	c.Data["json"] = r
	c.ServeJSON()
}

// GetMenu 返回菜单信息
func (c *JedaController) GetMenu() {

}

// DoSrv
func (c *JedaController) DoSrv() {
	//获取上下文
	cnt := c.Ctx.Input.Param(":context")
	//根据上下文获取服务定义信息
	//默认是从数据库获取
	sdef := JedaSrvContainer[cnt]
	if sdef == nil {
		utils.CreateErrorResponse("没有找到请求的服务,"+cnt, &c.Controller)
		return
	}
	if sdef.Security {
		// 处理访问控制
		userid, err := service.GetISevurityServiceInstance().VerifyToken(&c.Controller)
		if err != nil {
			utils.CreateErrorResponse(err.Error(), &c.Controller)
			return
		}
		if !service.GetISevurityServiceInstance().VerifyService(userid, sdef.ServiceId, 0) {
			utils.CreateErrorResponse("未授权的请求", &c.Controller)
			return
		}
	}
	h := &service.IDSServiceHandler{service.SHandlerBase{Ctl: &c.Controller}}
	h.DoSrv(sdef, h)
}
