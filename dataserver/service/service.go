package service

import (
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/utils"
)

const (
	// SrvTypeIds 基于TableDataSource类的服务
	SrvTypeIds string = "IDS"
	// SrvTypePredef 预定义服务
	SrvTypePredef string = "PREDEF"
	// SrvValueKey value-key形式的服务
	SrvValueKey string = "VK"
	// SrvTypeSrvflow 基于服务流程的服务
	SrvTypeSrvflow string = "SRVFLOW"
)

// createServiceHandlerInterfaceFun 创建服务处理句柄的函数类型
type createServiceHandlerInterfaceFun func(*beego.Controller, string) SHandlerInterface

// SHandlerContainer 服务处理句柄容器
var SHandlerContainer = make(map[string]createServiceHandlerInterfaceFun)

// SDefine 服务定义结构体
type SDefine struct {
	// ServiceId 服务ID GUID类型
	ServiceId string
	// Context 上下文
	Context string
	// BodyType 报文类型
	BodyType string
	// ServiceType 服务类型
	ServiceType string
	// Meta 服务元数据
	Meta string
	// Namespace 命名空间
	Namespace string
	// Enabled 是否可用
	Enabled bool
	// MsgLog 是否开启消息日志
	MsgLog bool
	// Security 是否开启安全认证
	Security bool
	// 项目id
	ProjectId string
}

// SDefineContainerType 服务定义类型
type SDefineContainerType map[string]*SDefine

// SController 服务控制器基类
type SController struct {
	beego.Controller
}

// SDefineContainer 服务定义容器
var SDefineContainer = make(SDefineContainerType)

// reloadSrvMetaFromDatabase 从数据库中加载服务定义信息
func (c *SController) reloadSrvMetaFromDatabase(cnt string) (*SDefine, error) {
	ps := strings.Split(cnt, ".")
	if len(ps) != 2 && len(ps) != 1 {
		return nil, fmt.Errorf("上下文格式不正确")
	}
	sdef, ok := SDefineContainer[cnt]

	if ok {
		return sdef, nil
	}
	ds := datasource.CreateTableDataSource("GSERVICE", "default", "G_SERVICE")
	if len(ps) == 2 {
		ds.AddCriteria("NAMESPACE", datasource.OperEq, ps[0]).AndCriteria("CONTEXT", datasource.OperEq, ps[1])
	}
	if len(ps) == 1 {
		ds.AddCriteria("CONTEXT", datasource.OperEq, ps[0]).AndCriteria("CONTEXT", datasource.OperEq, "")
	}
	rs, err := ds.DoFilter()
	if err != nil {
		return nil, err
	}
	if len(rs.Data) == 0 {
		return nil, fmt.Errorf("没有找到对应的服务" + cnt)
	}

	srv := &SDefine{
		ServiceId:   rs.Data[0][rs.Fields["ID"].Index].(string),
		Context:     rs.Data[0][rs.Fields["CONTEXT"].Index].(string),
		BodyType:    rs.Data[0][rs.Fields["BODYTYPE"].Index].(string),
		ServiceType: rs.Data[0][rs.Fields["SERVICETYPE"].Index].(string),
		Namespace:   rs.Data[0][rs.Fields["NAMESPACE"].Index].(string),
		Enabled:     rs.Data[0][rs.Fields["ENABLED"].Index].(int32) == 1,
		MsgLog:      rs.Data[0][rs.Fields["MSGLOG"].Index].(int32) == 1,
		Security:    rs.Data[0][rs.Fields["SECURITY"].Index].(int32) == 1,
		Meta:        rs.Data[0][rs.Fields["META"].Index].(string),
		ProjectId:   rs.Data[0][rs.Fields["PROJECTID"].Index].(string)}
	SDefineContainer[cnt] = srv
	return srv, nil
}

// DoSrv 处理请求
func (c *SController) DoSrv() {

	//获取上下文
	cnt := c.Ctx.Input.Param(":context")
	//根据上下文获取服务定义信息
	//默认是从数据库获取
	sdef, err := c.reloadSrvMetaFromDatabase(cnt)
	if err != nil {
		utils.CreateErrorResponse(err.Error(), &c.Controller)
		return
	}
	if !sdef.Enabled {
		utils.CreateErrorResponse("请求的服务未启用", &c.Controller)
		return
	}
	userid := ""
	if sdef.Security {
		// 处理访问控制
		userid, err = GetISevurityServiceInstance().VerifyToken(&c.Controller)
		if err != nil {
			utils.CreateErrorResponse(err.Error(), &c.Controller)
			return
		}
		if !GetISevurityServiceInstance().VerifyService(userid, sdef.ServiceId, 0) {
			utils.CreateErrorResponse("未授权的请求", &c.Controller)
			return
		}
	}
	handler, ok := SHandlerContainer[sdef.ServiceType]
	if !ok {
		utils.CreateErrorResponse("没有找到"+sdef.ServiceType+"定义的服务接口处理程序", &c.Controller)
		return
	}
	h := handler(&c.Controller, userid)
	h.DoSrv(sdef, h)
}

// init 初始化
func init() {
	SHandlerContainer[SrvTypeIds] = func(c *beego.Controller, caller string) SHandlerInterface {
		return &IDSServiceHandler{SHandlerBase{Ctl: c, CurrentUserId: caller}}
	}
	SHandlerContainer[SrvTypePredef] = func(c *beego.Controller, caller string) SHandlerInterface {
		return &PredefineServiceHandler{IDSServiceHandler: IDSServiceHandler{SHandlerBase{Ctl: c, CurrentUserId: caller}}}
	}
	SHandlerContainer[SrvValueKey] = func(c *beego.Controller, caller string) SHandlerInterface {
		return &ValueKeyService{SHandlerBase{Ctl: c, CurrentUserId: caller}}
	}
}
