package service

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/utils"
)

// SHandlerInterface 服务处理接口
type SHandlerInterface interface {
	//处理服务的方法，在目前的程序中POST和GET请求都会映射到该方法上
	DoSrv(sdef *SDefine, inf SHandlerInterface)
	//返回当前实现支持的动作和动作对应的处理函数
	getActionMap() map[string]SerivceActionHandler
	//返回请求报文，GET方法没有报文，只处理POST方法的报文
	getRBody() *SRequestBody
	//根据元数据返回当前实例处理请求的数据源类，比如TableDataSource
	getServiceInterface(metestr string) (interface{}, error)
}

// SerivceActionHandler 处理请求的方法类型
type SerivceActionHandler func(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody)

// SHandlerBase 服务处理句柄基类
type SHandlerBase struct {
	Ctl       *beego.Controller
	ActionMap map[string]SerivceActionHandler
}

// HasRightService 判断是否有权限
func HasRightService(user string, serviceid string) (bool, error) {
	var maps []orm.Params
	o := orm.NewOrm()
	_, err := o.Raw("select * from G_USERSERVICE where USERID=?", user).Values(&maps)
	if err != nil {
		return false, err
	}
	if len(maps) == 0 {
		return false, nil
	}
	return true, nil
}

// createErrorResponse 设定失败结果
func (c *SHandlerBase) createErrorResponse(msg string) {
	CreateErrorResponse(msg, c.Ctl)
}

// createErrorResponseByError 根据error设定失败结果
func (c *SHandlerBase) createErrorResponseByError(err error) {
	CreateErrorResponseByError(err, c.Ctl)
}

// createErrorResult 设定失败结果
func (c *SHandlerBase) createErrorResult(msg string) {
	CreateErrorResponse(msg, c.Ctl)
}

// setResult 设定请求成功的返回结果
func (c *SHandlerBase) setResult(msg string) {
	r := CreateRestResult(true)
	r["msg"] = msg
	c.Ctl.Data["json"] = r
}

// setResultSet 设定结果集
func (c *SHandlerBase) setResultSet(ds *datasource.DataResultSet) {

	r := CreateRestResult(true)
	if c.Ctl.Input().Get(RequestParamNofieldsinfo) != "" {
		r["data"] = ds.Data
	} else {
		r["resultset"] = ds
	}

	if c.Ctl.Input().Get(RequestParamCache /**_cache**/) != "" {
		// 处理缓存请求 [缓存时间]_[最大请求次数]  10_1  缓存的结果集请求一次即删除，
		// 最长保存10秒钟，“缓存时间”为0时表示使用系统定义的默认缓存时间，为30s
		// 缓存的结果集随时都有可能消失
		cs := c.Ctl.Input().Get(RequestParamCache)
		css := strings.Split(cs, "_")
		if len(css) != 2 {
			r["result"] = false
			r["msg"] = "缓存参数" + RequestParamCache + "必须为 [缓存时间]_[最大请求次数] 的形式"
			c.Ctl.Data["json"] = r
			return
		}
		t, ok := strconv.Atoi(css[0])
		if ok != nil {
			r["result"] = false
			r["msg"] = "缓存时间非法"
			c.Ctl.Data["json"] = r
			return
		}
		t2, ok := strconv.Atoi(css[1])
		if ok != nil {
			r["result"] = false
			r["msg"] = "最大请求次数非法"
			c.Ctl.Data["json"] = r
			return
		}
		if t < 0 || t2 < 0 {
			r["result"] = false
			r["msg"] = "非法的最大请求次数或缓存时间"
			c.Ctl.Data["json"] = r
			return
		}

		keys := uuid.NewV4().String()
		r["cacheid"] = keys
		if t == 0 {
			t = 10
		}
		if t2 == 0 {
			t2 = -1
		}
		r["cachetimes"] = t2
		r["duration"] = t
		err := utils.DataSetResultCache.Put(keys, r, time.Duration(t)*time.Second)
		if err != nil {
			r["result"] = false
			r["msg"] = "加入缓存时发生错误：" + err.Error()
			c.Ctl.Data["json"] = r
			return
		}
	}

	c.Ctl.Data["json"] = r
}

// ServeJSON call the c.Ctl.ServeJSON()
func (c *SHandlerBase) ServeJSON() {
	c.Ctl.ServeJSON()
}

// setPageParams 设定从querystring传入的公共参数
func (c *SHandlerBase) setPageParams(ids datasource.IDataSource) {
	psi, err := strconv.Atoi(c.Ctl.Input().Get(RequestParamPagesize))
	pii, err2 := strconv.Atoi(c.Ctl.Input().Get(RequestParamPageindex))
	if err == nil {
		ids.SetRowsLimit(psi)
		if err2 == nil {
			ids.SetRowsOffset(psi * (pii - 1))
		}
	}
}

// ConvertString2Type 转换字符串为指定的类型，转换不成功返回nil
func (c *SHandlerBase) ConvertString2Type(value string, vtype string) (interface{}, error) {
	switch vtype {
	case datasource.PropertyDatatypeInt:
		{
			i, err := strconv.Atoi(value)
			if err != nil {
				return nil, err
			}
			return i, nil
		}
	case datasource.PropertyDatatypeDou:
		{
			i, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return nil, err
			}
			return i, nil
		}
	case datasource.PropertyDatatypeStr:
		return value, nil

	case datasource.PropertyDatatypeDate:
		{
			theTime, err := time.Parse("2006-01-02", value)
			if err != nil {
				return nil, err
			}
			return theTime, nil
		}
	case datasource.PropertyDatatypeTime:
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
	case datasource.PropertyDatatypeEnum:
		return value, nil
	case datasource.PropertyDatatypeUnkn:
		return value, nil
	}
	return value, nil
}
