package service

import (
	"awesome/datasource"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
type ServiceHandlerInterface interface {
	//处理服务的方法，在目前的程序中POST和GET请求都会映射到该方法上
	DoSrv(metestr string, inf ServiceHandlerInterface)
	//返回当前实现支持的动作和动作对应的处理函数
	getActionMap() map[string]SerivceActionHandler
	//返回请求报文，GET方法没有报文，只处理POST方法的报文
	getRBody() *ServiceRequestBody
	//根据元数据返回当前实例处理请求的数据源类，比如TableDataSource
	getServiceInterface(metestr string) (interface{}, error)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

//处理请求的方法类型
type SerivceActionHandler func(ids datasource.IDataSource, rBody *ServiceRequestBody)
type ServiceHandlerBase struct {
	Ctl       *beego.Controller
	ActionMap map[string]SerivceActionHandler
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func HasRightService(user string, serviceid string) (bool, error) {
	var maps []orm.Params
	o := orm.NewOrm()
	_, err := o.Raw("select * from G_USERSERVICE where USERID=?", user).Values(&maps)
	if err != nil {
		return false, err
	}
	if len(maps) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *ServiceHandlerBase) createErrorResponse(msg string) {
	CreateErrorResponse(msg,c.Ctl)
}
func (c *ServiceHandlerBase) createErrorResponseByError(err error) {
	CreateErrorResponseByError(err,c.Ctl)
}
func (c *ServiceHandlerBase) createErrorResult(msg string) {
	CreateErrorResult(msg,c.Ctl)
}
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *ServiceHandlerBase) setResult(msg string) {
	r := CreateRestResult(true)
	r["msg"] = msg
	c.Ctl.Data["json"] = r
}
func (c *ServiceHandlerBase) setResultSet(ds *datasource.DataResultSet) {
	r := CreateRestResult(true)
	if c.Ctl.Input().Get(REQUEST_PARAM_NOFIELDSINFO) != "" {
		r["data"] = ds.Data
	} else {
		r["resultset"] = ds
	}
	c.Ctl.Data["json"] = r
}



/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *ServiceHandlerBase) ServeJson() {
	c.Ctl.ServeJSON()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *ServiceHandlerBase) setPageParams(ids datasource.IDataSource) {
	psi, err := strconv.Atoi(c.Ctl.Input().Get(REQUEST_PARAM_PAGESIZE))
	pii, err2 := strconv.Atoi(c.Ctl.Input().Get(REQUEST_PARAM_PAGEINDEX))
	if err == nil && err2 == nil {
		ids.SetRowsLimit(psi)
		ids.SetRowsOffset(psi * (pii - 1))
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//转换字符串为指定的类型，转换不成功返回nil
func (c *ServiceHandlerBase) ConvertString2Type(value string, vtype string) (interface{}, error) {
	switch vtype {
	case datasource.Property_Datatype_INT:
		{
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			return i, nil
		}
	case datasource.Property_Datatype_DOU:
		{
			i, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			return i, nil
		}
	case datasource.Property_Datatype_STR:
		return value, nil
	case datasource.Property_Datatype_DATE:
		{

			theTime, err := time.Parse("2006-01-02", value)
			if err != nil {
				return nil, err
			}
			return theTime, nil
		}
	case datasource.Property_Datatype_TIME:
		{

			theTime, err := time.Parse("2006-01-02 15:04:05", value)
			if err != nil {
				theTime, err := time.Parse("2006-01-02", value)
				if err != nil {
					return nil, err
				}
				return theTime, nil
			}
			return theTime, nil
		}
	case datasource.Property_Datatype_ENUM:
		return value, nil
	case datasource.Property_Datatype_UNKN:
		return value, nil
	}
	return value, nil
}
