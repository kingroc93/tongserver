package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"tongserver.dataserver/datasource"
)

const (
	//基于TableDataSource类的服务
	SRV_TYPE_IDS string = "IDS"
	//预定义服务
	SRV_TYPE_PREDEF string = "PREDEF"
	//基于RnmuSource类的服务
	SRV_TYPE_ENMU string = "ENMU"
	//基于服务流程的服务
	SRV_TYPE_SRVFLOW string = "SRVFLOW"
)

type ServiceDefine struct {
	//上下文
	Context string
	//报文类型
	BodyType string
	//服务类型
	ServiceType string
	//服务元数据
	Meta string
	//命名空间
	Namespace string
	//是否可用
	Enabled bool
	//是否开启消息日志
	MsgLog bool
	//是否开启安全认证
	Security bool
}

type ServiceDefineContainerType map[string]*ServiceDefine

var ServiceDefineContainer = make(ServiceDefineContainerType)

type ServiceController struct {
	beego.Controller
}

////////////////////////////////////////////////////////////////////////////////////////////////
//从数据库中加载服务定义信息
func (c *ServiceController) reloadSrvMetaFromDatabase(cnt string) (*ServiceDefine, error) {
	ps := strings.Split(cnt, ".")
	if len(ps) != 2 && len(ps) != 1 {
		return nil, fmt.Errorf("上下文格式不正确")
	}
	sdef, ok := ServiceDefineContainer[cnt]
	if ok {
		return sdef, nil
	}
	ds := datasource.CreateTableDataSource("GSERVICE", "default", "G_SERVICE")
	if len(ps) == 2 {
		ds.AddCriteria("NAMESPACE", datasource.OPER_EQ, ps[0]).AndCriteria("CONTEXT", datasource.OPER_EQ, ps[1])
	}
	if len(ps) == 1 {
		ds.AddCriteria("CONTEXT", datasource.OPER_EQ, ps[0]).AndCriteria("CONTEXT", datasource.OPER_EQ, "")
	}
	rs, err := ds.DoFilter()
	if err != nil {
		return nil, err
	}
	if len(rs.Data) == 0 {
		return nil, fmt.Errorf("没有找到对应的服务" + cnt)
	}

	srv := &ServiceDefine{
		Context:     rs.Data[0][rs.Fields["CONTEXT"].Index].(string),
		BodyType:    rs.Data[0][rs.Fields["BODYTYPE"].Index].(string),
		ServiceType: rs.Data[0][rs.Fields["SERVICETYPE"].Index].(string),
		Namespace:   rs.Data[0][rs.Fields["NAMESPACE"].Index].(string),
		Enabled:     rs.Data[0][rs.Fields["ENABLED"].Index].(int32) == 1,
		MsgLog:      rs.Data[0][rs.Fields["MSGLOG"].Index].(int32) == 1,
		Security:    rs.Data[0][rs.Fields["SECURITY"].Index].(int32) == 1,
		Meta:        rs.Data[0][rs.Fields["META"].Index].(string)}
	ServiceDefineContainer[cnt] = srv
	return srv, nil
}

////////////////////////////////////////////////////////////////////////////////////////////////
// 处理请求
func (c *ServiceController) DoSrv() {
	//获取上下文
	cnt := c.Ctx.Input.Param(":context")
	//根据上下文获取服务定义信息
	//默认是从数据库获取
	sdef, err := c.reloadSrvMetaFromDatabase(cnt)
	if err != nil {
		r := CreateRestResult(false)
		r["msg"] = err.Error()
		c.Data["json"] = r
		c.ServeJSON()
	}
	//IDS类型的服务，即使用IDataSource接口作为服务处理器的服务
	if sdef.ServiceType == SRV_TYPE_IDS {
		ctl := &IDSServiceHandler{
			ServiceHandlerBase{
				Ctl: &c.Controller,
			}}
		ctl.DoSrv(sdef.Meta, ctl)
		return
	}
	//IDS服务会通过POST请求某些服务是预定义的处理过程
	if sdef.ServiceType == SRV_TYPE_PREDEF {
		ctl := &PredefineServiceHandler{
			IDSServiceHandler: IDSServiceHandler{
				ServiceHandlerBase{
					Ctl: &c.Controller,
				}}}
		ctl.DoSrv(sdef.Meta, ctl)
		return
	}
}

func init() {

}
