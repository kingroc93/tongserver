package service

import (
	"fmt"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
	"time"
	"tongserver.dataserver/datasource"

	"tongserver.dataserver/utils"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	//返回全部数据
	SrvActionALLDATA string = "all"
	//查询动作
	SrvActionQUERY string = "query"
	//根据主键返回
	SrvActionGET string = "get"
	//返回缓存
	SrvActionCACHE string = "cache"
	//根据字段值返回
	SrvActionBYFIELD string = "byfield"
	//返回服务元数据
	SrvActionMETA string = "meta"
	//删除操作
	SrvActionDELETE string = "delete"
	//更新操作
	SrvActionUPDATE string = "update"
	//插入操作
	SrvActionINSERT string = "insert"

	//以下三个常量均为通过QueryString传入的参数名
	//针对查询自动分页中每页记录数
	RequestParamPagesize string = "_pagesize"
	//针对查询自动分页中的页索引
	RequestParamPageindex string = "_pageindex"
	//是否返回字段元数据，默认为返回
	RequestParamNofieldsinfo string = "_nofield"
	// 响应的风格，默认是数组风格array，可以设定为map风格
	ResponseStyle string = "_repstyle"
	//当前请求不执行而是只返回SQL语句，仅针对IDS类型的服务有效
	RequestParamSQL string = "_sql"
	//当前请求的响应信息不直接返回
	//该参数只对query、all两个操作起作用
	RequestParamCache      string = "_cache"
	RequestParamCachebykey string = "_cachekey"
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
	getServiceInterface(meta map[string]interface{}, sdef *SDefine) (interface{}, error)
}

// 请求响应的接口
type RequestResponseHandler interface {
	CreateResponseData(style int, data interface{})
	GetParam(name string) string
	GetRequestBody() (*SRequestBody, error)
}

// SerivceActionHandler 处理请求的方法类型
type SerivceActionHandler func(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody)

// SHandlerBase 服务处理句柄基类
type SHandlerBase struct {
	//Ctl           *beego.Controller
	RRHandler     RequestResponseHandler
	ActionMap     map[string]SerivceActionHandler
	CurrentUserId string
}

func (c *SHandlerBase) createErrorResponse(msg string) {
	r := utils.CreateRestResult(false)
	r["msg"] = msg
	c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据元数据返回处理服务的接口
func (c *SHandlerBase) getServiceInterface(meta map[string]interface{}, sdef *SDefine) (interface{}, error) {
	idstr := meta["ids"].(string)
	if strings.Index(idstr, ".") == -1 {
		idstr = sdef.ProjectId + "." + idstr
	}
	return datasource.CreateIDSFromName(idstr)
}

// DoSrv 处理服务请求的入口
func (c *SHandlerBase) DoSrv(sdef *SDefine, inf SHandlerInterface) {
	//////////////////////////////////////////////////////////////////////////
	//调用传入的接口中的方法实现下面的功能,因为需要通过不同的接口实现来实现不同的行为
	//meta := make(map[string]interface{})
	//err2 := json.Unmarshal([]byte(metestr), &meta)
	meta, err2 := utils.ParseJSONStr2Map(sdef.Meta)
	if err2 != nil {
		c.createErrorResponse("meta信息不正确,应为JSON格式")
		return
	}
	obj, err := inf.getServiceInterface(meta, sdef)
	if err != nil {
		c.createErrorResponse(err.Error())
		return
	}
	rBody := inf.getRBody()
	//////////////////////////////////////////////////////////////////////////
	ids, ok := obj.(datasource.IDataSource)
	if !ok {
		c.createErrorResponse("请求的服务没有实现IDataSource接口")
		return
	}
	act := c.RRHandler.GetParam(":action") //c.Ctl.Ctx.Input.Param(":action")
	amap := inf.getActionMap()
	f, ok := amap[act]
	if !ok {
		c.createErrorResponse("请求的动作当前服务没有实现")
		return
	}
	f(sdef, meta, ids, rBody)
}

func (c *SHandlerBase) getActionMap() map[string]SerivceActionHandler {
	return map[string]SerivceActionHandler{
		SrvActionMETA:  c.doGetMeta,
		SrvActionCACHE: c.doGetCache,
	}
}

func (c *SHandlerBase) getCache(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) (*utils.RestResult, error) {
	key := c.RRHandler.GetParam(RequestParamCachebykey) //c.Ctl.Input().Get(RequestParamCachebykey)
	if key == "" {
		return nil, fmt.Errorf(RequestParamCachebykey + "不得为空")
	}
	obj := utils.DataSetResultCache.Get(key)
	if obj == nil {
		return nil, fmt.Errorf("没有找到请求的缓存信息")
	}
	r, ok := obj.(utils.RestResult)
	if !ok {
		return nil, fmt.Errorf("缓存对象类型非法")
	}
	times := r["cachetimes"].(int)
	d := r["duration"].(int)
	if times > 0 {
		times = times - 1
	}
	r["cachetimes"] = times
	if times == 0 {
		err := utils.DataSetResultCache.Delete(key)
		if err != nil {
			r["result"] = false
			r["msg"] = "删除缓存时发生错误：" + err.Error()
			return &r, fmt.Errorf("删除缓存时发生错误：" + err.Error())
		}
	} else {
		err := utils.DataSetResultCache.Put(key, obj, time.Duration(d)*time.Second)
		if err != nil {
			r["result"] = false
			r["msg"] = "加入缓存时发生错误：" + err.Error()
			return &r, fmt.Errorf("加入缓存时发生错误：" + err.Error())
		}
	}
	return &r, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回缓存的结果数据
func (c *SHandlerBase) doGetCache(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	r, err := c.getCache(sdef, ids, rBody)
	if r != nil {
		(*r)["result"] = true
	}
	if err != nil {
		(*r)["msg"] = err.Error()
	}
	c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回服务元数据
func (c *SHandlerBase) doGetMeta(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	r := utils.CreateRestResult(true)
	sd := make(map[string]interface{})
	r["servicedefine"] = sd
	sd["Context"] = sdef.Context
	sd["BodyType"] = sdef.BodyType
	sd["ServiceType"] = sdef.ServiceType
	sd["Namespace"] = sdef.Namespace
	sd["Enabled"] = sdef.Enabled
	sd["MsgLog"] = sdef.MsgLog
	sd["Security"] = sdef.Security
	sd["Meta"] = meta

	imp := []string{"IDataSource"}
	if inf, ok := ids.(datasource.ICriteriaDataSource); ok {
		imp = append(imp, "ICriteriaDataSource")
		sd["Fields"] = inf.GetFields()
		sd["KeyFields"] = inf.GetKeyFields()
	}
	if _, ok := ids.(datasource.IFilterAdder); ok {
		imp = append(imp, "IFilterAdder")
	}
	if _, ok := ids.(datasource.IAggregativeAdder); ok {
		imp = append(imp, "IAggregativeAdder")
	}
	if _, ok := ids.(datasource.IWriteableDataSource); ok {
		imp = append(imp, "IWriteableDataSource")
	}
	r["ids"] = imp

	//c.Ctl.Data["json"] = r
	//c.ServeJSON()
	c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
}

// setResultSet 设定结果集
func (c *SHandlerBase) setResultSet(ds *datasource.DataResultSet) {
	//c.Ctl.Input().Get(RequestParamCache /**_cache**/)
	if c.RRHandler.GetParam(RequestParamCache /**_cache**/) != "" {
		// 处理缓存请求 [缓存时间]_[最大请求次数]  10_1  缓存的结果集请求一次即删除，
		// 最长保存10秒钟，“缓存时间”为0时表示使用系统定义的默认缓存时间，为30s
		// 缓存的结果集随时都有可能消失
		cs := c.RRHandler.GetParam(RequestParamCache) //c.Ctl.Input().Get(RequestParamCache)
		css := strings.Split(cs, "_")
		r := utils.CreateRestResult(true)
		if len(css) != 2 {
			c.createErrorResponse("缓存参数" + RequestParamCache + "必须为 [缓存时间]_[最大请求次数] 的形式")
			return
		}
		t, ok := strconv.Atoi(css[0])
		if ok != nil {
			c.createErrorResponse("缓存时间非法")
			return
		}
		t2, ok := strconv.Atoi(css[1])
		if ok != nil {
			c.createErrorResponse("最大请求次数非法")
			return
		}
		if t < 0 || t2 < 0 {
			c.createErrorResponse("非法的最大请求次数或缓存时间")
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
			c.createErrorResponse("加入缓存时发生错误：" + err.Error())
			return
		}
		c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
		return
	}
	r := utils.CreateRestResult(true)
	//if c.Ctl.Input().Get(ResponseStyle) != "map" {
	if c.RRHandler.GetParam(ResponseStyle) != "map" {
		//if c.Ctl.Input().Get(RequestParamNofieldsinfo) != "" {
		if c.RRHandler.GetParam(RequestParamNofieldsinfo) != "" {
			r["data"] = ds.Data
		} else {
			r["resultset"] = ds
		}
	} else {
		result := make([]map[string]interface{}, len(ds.Data), len(ds.Data))
		for i, d := range ds.Data {
			item := make(map[string]interface{})
			for k, v := range ds.Fields {
				item[k] = d[v.Index]
			}
			result[i] = item
		}
		//if c.Ctl.Input().Get(RequestParamNofieldsinfo) != "" {
		if c.RRHandler.GetParam(RequestParamNofieldsinfo) != "" {
			r["data"] = result
		} else {
			rsd := make(map[string]interface{})
			rsd["Fields"] = ds.Fields
			rsd["Data"] = result
			rsd["Meta"] = ds.Meta
			r["resultset"] = rsd
		}
	}
	c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
}

// setPageParams 设定从querystring传入的公共参数
func (c *SHandlerBase) setPageParams(ids datasource.IDataSource) {
	psi, err := strconv.Atoi(c.RRHandler.GetParam(RequestParamPagesize))
	pii, err2 := strconv.Atoi(c.RRHandler.GetParam(RequestParamPageindex))
	if err == nil {
		ids.SetRowsLimit(psi)
		if err2 == nil {
			ids.SetRowsOffset(psi * (pii - 1))
		}
	}
}

// ConvertString2Type 转换字符串为指定的类型，转换不成功返回nil
func (c *SHandlerBase) ConvertString2Type(value string, vtype string) (interface{}, error) {
	return datasource.ConvertString2Type(value, vtype)
}
