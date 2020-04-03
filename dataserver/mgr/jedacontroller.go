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
	service.ServiceControllerBase
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
	if _, ok := c.ControllerWithVerify.Verifty(&c.Controller); !ok {
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

// 测试链接
func (c *JedaController) Testdbconn() {
	if _, ok := c.ControllerWithVerify.Verifty(&c.Controller); !ok {
		return
	}
}

// GetIdsList 返回数据源列表
func (c *JedaController) GetIdsList() {
	if _, ok := c.ControllerWithVerify.Verifty(&c.Controller); !ok {
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

// 根据SQL语句和参数返回map
func (c *JedaController) renderSQL(sql string, ps ...interface{}) ([]interface{}, error) {
	sqld := datasource.CreateSQLDataSource("", "default", sql)
	sqld.ParamsValues = ps
	rs, err := sqld.GetAllData()
	if err != nil {
		return nil, err
	}
	mr := make([]interface{}, len(rs.Data), len(rs.Data))
	for i := 0; i < len(rs.Data); i++ {
		m := c.convertRset2map(rs, i)
		mr[i] = m
	}
	return mr, nil
}

// 填充用户基本信息
func (c *JedaController) renderUserinfo(userid string) {
	obj := datasource.CreateIDSFromParam(datasource.IDSContainer["default.mgr.JEDA_USER"])

	if obj == nil {
		utils.CreateErrorResponse("没有找到jeda.user数据源", &c.Controller)
		return
	}
	ids := obj.(datasource.ICriteriaDataSource)
	rs, err := ids.QueryDataByFieldValues(&map[string]interface{}{"USER_ID": userid})
	if err != nil {
		utils.CreateErrorResponse(err.Error(), &c.Controller)
		return
	}
	if len(rs.Data) == 0 {
		utils.CreateErrorResponse("没有找到用户信息"+userid, &c.Controller)
		return
	}
	rest := c.convertRset2map(rs, 0)

	sqld := datasource.CreateSQLDataSource("", "default",
		"SELECT a.* FROM JEDA_ROLE a inner join JEDA_ROLE_USER b on a.ROLE_ID=b.ROLE_ID and b.USER_ID=?")
	sqld.ParamsValues = []interface{}{userid}
	rs, err = sqld.GetAllData()
	if err != nil {
		utils.CreateErrorResponse(err.Error(), &c.Controller)
		return
	}
	if len(rs.Data) != 0 {
		mr := make([]interface{}, len(rs.Data), len(rs.Data))
		for i := 0; i < len(rs.Data); i++ {
			m := c.convertRset2map(rs, i)
			mr[i] = m
		}
		rest["roleset"] = mr
	}
	r := utils.CreateRestResult(true)
	r["resultset"] = rest
	c.Data["json"] = r
	c.ServeJSON()
}

// 返回用户在当前系统中的详细信息
func (c *JedaController) GetCurrentUserInfo() {
	userid, err := service.GetISevurityServiceInstance().VerifyToken(&c.Controller)
	if err != nil {
		utils.CreateErrorResponse(err.Error(), &c.Controller)
		return
	}
	cnt := c.Ctx.Input.Param(":cat")
	switch cnt {
	case "info":
		c.renderUserinfo(userid)
	case "role":
		{
			mr, err := c.renderSQL("SELECT a.* FROM JEDA_ROLE a inner join JEDA_ROLE_USER b on a.ROLE_ID=b.ROLE_ID and b.USER_ID=?", userid)
			if err != nil {
				utils.CreateErrorResponse(err.Error(), &c.Controller)
				return
			}
			r := utils.CreateRestResult(true)
			r["resultset"] = mr
			c.Data["json"] = r
			c.ServeJSON()
		}
	case "service":
		{
			mr, err := c.renderSQL("SELECT * FROM G_SERVICE a inner join G_USERSERVICE b on a.ID=b.SERVICEID inner join JEDA_ROLE_USER c on c.ROLE_ID=b.ROLEID and c.USER_ID=?", userid)
			if err != nil {
				utils.CreateErrorResponse(err.Error(), &c.Controller)
				return
			}
			r := utils.CreateRestResult(true)
			r["resultset"] = mr
			c.Data["json"] = r
			c.ServeJSON()
		}
	case "project":
		{
			mr, err := c.renderSQL("SELECT a.* FROM G_PROJECT a inner join G_USERPROJECT b on a.ID=b.PROJECTID and b.USERID=?", userid)
			if err != nil {
				utils.CreateErrorResponse(err.Error(), &c.Controller)
				return
			}
			r := utils.CreateRestResult(true)
			r["resultset"] = mr
			c.Data["json"] = r
			c.ServeJSON()
		}

	}
}

// 将结果集由数组的形式转换为map的形式
func (c *JedaController) convertRset2map(ds *datasource.DataResultSet, index int) map[string]interface{} {
	d := ds.Data[index]
	item := make(map[string]interface{})
	for k, v := range ds.Fields {
		item[k] = d[v.Index]
	}
	return item
}

// 判断用户是否可以访问服务
func (c *JedaController) verifyUserAccess(srvid string) (string, bool) {
	// 处理访问控制
	userid, err := service.GetISevurityServiceInstance().VerifyToken(&c.Controller)
	if err != nil {
		utils.CreateErrorResponse(err.Error(), &c.Controller)
		return "", false
	}
	if !service.GetISevurityServiceInstance().VerifyService(userid, srvid, 0) {
		utils.CreateErrorResponse("未授权的请求", &c.Controller)
		return "", false
	}
	return userid, true
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
	userid := ""
	ok := false
	if sdef.Security {
		// 处理访问控制
		if userid, ok = c.verifyUserAccess(sdef.ServiceId); !ok {
			return
		}
	}
	h := &service.IDSServiceHandler{service.SHandlerBase{RRHandler: c, CurrentUserId: userid}}
	h.DoSrv(sdef, h)
	c.ServeJSON()
}
