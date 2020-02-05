package service

import (
	"encoding/json"
	"fmt"
	"tongserver.dataserver/datasource"
)

// PredefineServiceHandler 预定义的服务，可以将请求的报文保存在数据库中，形成预定义的服务4
// 在调用时无需提交整个报文，可以不提交报文或提交部分报文，系统合并提交的报文和预定的报文
// 预定义的报文中，针对条件的值value属性可以使用:?作为占位符，通过QueryString传入参数
type PredefineServiceHandler struct {
	IDSServiceHandler
	// 保存预定义的报文信息
	predefine *PredefineBody
}

// PredefineBody 预定义的报文结构体
type PredefineBody struct {
	SRequestBody
	Ids        string
	Definetype string
}

// DoSrv 预定义服务处理请求
func (c *PredefineServiceHandler) DoSrv(sdef *SDefine, inf SHandlerInterface) {
	c.predefine = &PredefineBody{}
	metestr := sdef.Meta
	err2 := json.Unmarshal([]byte(metestr), c.predefine)
	if err2 != nil {
		c.createErrorResponse("meta信息不正确,应为JSON格式")
		return
	}
	c.IDSServiceHandler.DoSrv(sdef, inf)
}

// merageRbody 合并请求报文,传入的请求报文和预定义的报文进行合并
func (c *PredefineServiceHandler) merageRbody(rBody *SRequestBody) *SRequestBody {
	b := c.predefine.SRequestBody
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

// getRBody 返回请求体
func (c *PredefineServiceHandler) getRBody() *SRequestBody {
	if c.Ctl.Ctx.Request.Method == "POST" {
		rBody := &SRequestBody{}
		if len(c.Ctl.Ctx.Input.RequestBody) != 0 {
			err := json.Unmarshal([]byte(c.Ctl.Ctx.Input.RequestBody), rBody)
			if err != nil {
				c.createErrorResponse("解析报文时发生错误" + err.Error())
			}
		}
		for i, cri := range c.predefine.SRequestBody.Criteria {
			if cri.Value == ":?" {
				c.predefine.SRequestBody.Criteria[i].Value = c.Ctl.Input().Get(cri.Field)
			}
		}
		return c.merageRbody(rBody)
	}
	for i, cri := range c.predefine.SRequestBody.Criteria {
		if cri.Value == ":?" {
			c.predefine.SRequestBody.Criteria[i].Value = c.Ctl.Input().Get(cri.Field)
		}
	}
	return &c.predefine.SRequestBody

}

// getServiceInterface 返回该服务需要的数据源接口
func (c *PredefineServiceHandler) getServiceInterface(sdef *SDefine) (interface{}, error) {
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

// doAllData 返回所有数据
func (c *PredefineServiceHandler) doAllData(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
	c.IDSServiceHandler.doQuery(sdef, ids, rBody)
}

// doGetMeta 返回元数据
func (c *PredefineServiceHandler) doGetMeta(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
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
	c.ServeJSON()
}

// getActionMap 返回动作映射表
func (c *PredefineServiceHandler) getActionMap() map[string]SerivceActionHandler {
	m := c.IDSServiceHandler.getActionMap()
	m[SrvActionALLDATA] = c.doAllData
	m[SrvActionMETA] = c.doGetMeta
	return m
}
