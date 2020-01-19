package service

import (
	"encoding/json"
	"fmt"
	"tongserver.dataserver/datasource"
)

// 预定义的服务，可以将请求的报文保存在数据库中，形成预定义的服务4
// 在调用时无需提交整个报文，可以不提交报文或提交部分报文，系统合并提交的报文和预定的报文
// 预定义的报文中，针对条件的值value属性可以使用:?作为占位符，通过QueryString传入参数
type PredefineServiceHandler struct {
	IDSServiceHandler
	// 保存预定义的报文信息
	predefine *PredefineBody
}

// 预定义的报文结构体
type PredefineBody struct {
	ServiceRequestBody
	Ids        string
	Definetype string
}

// 预定义服务处理请求
func (c *PredefineServiceHandler) DoSrv(sdef *ServiceDefine, inf ServiceHandlerInterface) {
	c.predefine = &PredefineBody{}
	metestr := sdef.Meta
	err2 := json.Unmarshal([]byte(metestr), c.predefine)
	if err2 != nil {
		c.createErrorResponse("meta信息不正确,应为JSON格式")
		return
	}
	c.IDSServiceHandler.DoSrv(sdef, inf)
}
func (c *PredefineServiceHandler) merageRbody(rBody *ServiceRequestBody) *ServiceRequestBody {
	b := c.predefine.ServiceRequestBody
	rBody.Criteria = append(rBody.Criteria, b.Criteria...)
	rBody.Aggre = append(rBody.Aggre, b.Aggre...)
	rBody.Bulldozer = append(rBody.Bulldozer, b.Bulldozer...)
	rBody.PostAction = append(rBody.PostAction, b.PostAction...)
	if rBody.OrderBy == "" {
		rBody.OrderBy = b.OrderBy
	}
	if rBody.InnerJoin == "" {
		rBody.InnerJoin = b.InnerJoin
	}
	return rBody
}

// 返回请求体
func (c *PredefineServiceHandler) getRBody() *ServiceRequestBody {
	if c.Ctl.Ctx.Request.Method == "POST" {
		rBody := &ServiceRequestBody{}
		if len(c.Ctl.Ctx.Input.RequestBody) != 0 {
			err := json.Unmarshal([]byte(c.Ctl.Ctx.Input.RequestBody), rBody)
			if err != nil {
				c.createErrorResponse("解析报文时发生错误" + err.Error())
			}
		}
		for i, cri := range c.predefine.ServiceRequestBody.Criteria {
			if cri.Value == ":?" {
				c.predefine.ServiceRequestBody.Criteria[i].Value = c.Ctl.Input().Get(cri.Field)
			}
		}
		return c.merageRbody(rBody)
	} else {
		for i, cri := range c.predefine.ServiceRequestBody.Criteria {
			if cri.Value == ":?" {
				c.predefine.ServiceRequestBody.Criteria[i].Value = c.Ctl.Input().Get(cri.Field)
			}
		}
		return &c.predefine.ServiceRequestBody
	}
}

//返回该服务需要的数据源接口
func (c *PredefineServiceHandler) getServiceInterface(metestr string) (interface{}, error) {
	if c.predefine.Definetype == "ids" {
		param := datasource.IDSContainer[c.predefine.Ids]
		obj := datasource.CreateIDSFromParam(param)
		if obj == nil {
			return nil, fmt.Errorf(c.predefine.Ids + "没有找到对应的处理程序")
		}
		return obj, nil
	}
	return nil, fmt.Errorf(c.predefine.Ids)
}

// 返回所有数据
func (c *PredefineServiceHandler) doAllData(sdef *ServiceDefine, ids datasource.IDataSource, rBody *ServiceRequestBody) {
	c.IDSServiceHandler.doQuery(sdef, ids, rBody)
}

// 返回元数据
func (c *PredefineServiceHandler) doGetMeta(sdef *ServiceDefine, ids datasource.IDataSource, rBody *ServiceRequestBody) {
	r := CreateRestResult(true)
	sd := make(map[string]interface{})
	r["servicedefine"] = sd
	sd["Context"] = sdef.Context
	sd["BodyType"] = sdef.BodyType
	sd["ServiceType"] = sdef.ServiceType
	sd["Namespace"] = sdef.Namespace
	sd["Enabled"] = sdef.Enabled
	sd["MsgLog"] = sdef.MsgLog
	sd["Security"] = sdef.Security
	meta := make(map[string]interface{})
	err2 := json.Unmarshal([]byte(sdef.Meta), &meta)
	if err2 == nil {
		sd["Meta"] = meta
	} else {
		sd["Meta"] = sdef.Meta
	}
	c.Ctl.Data["json"] = r
	c.ServeJson()
}

//返回动作映射表
func (c *PredefineServiceHandler) getActionMap() map[string]SerivceActionHandler {
	m := c.IDSServiceHandler.getActionMap()
	m[SrvActionALLDATA] = c.doAllData
	m[SrvActionMETA] = c.doGetMeta
	return m
}
