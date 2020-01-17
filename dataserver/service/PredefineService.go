package service

import (
	"encoding/json"
	"fmt"
	"tongserver.dataserver/datasource"
)

type PredefineServiceHandler struct {
	IDSServiceHandler
	predefine *PredefineBody
}
type PredefineBody struct {
	ServiceRequestBody
	Ids        string
	Definetype string
}

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
func (c *PredefineServiceHandler) doAllData(sdef *ServiceDefine, ids datasource.IDataSource, rBody *ServiceRequestBody) {
	c.IDSServiceHandler.doQuery(sdef, ids, rBody)
}

func (c *PredefineServiceHandler) getActionMap() map[string]SerivceActionHandler {
	m := c.IDSServiceHandler.getActionMap()
	m[SrvAction_ALLDATA] = c.doAllData
	return m
}
