package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"reflect"
	"strings"
	"tongserver.dataserver/activity"
	"tongserver.dataserver/utils"
	"tongserver.dataserver/utils/mapstructure"
)

// 内部服务活动
// 在流程中调用内部的服务
type InnerServiceActivity struct {
	activity.Activity
	cnt         string
	rv          string
	body        map[string]interface{}
	params      map[string]string
	resultData  interface{}
	resultStyle int
	context     activity.IContext
}

func (c *InnerServiceActivity) Execute(flowcontext activity.IContext) error {
	c.context = flowcontext
	cnt, err := activity.ReplaceExpressionLStr(flowcontext, c.cnt)
	if err != nil {
		logs.Error("InnerServiceActivity:Execute:解析EL表达式错误")
		return err
	}
	sdef, err := GetSrvMetaFromPath(cnt)
	if err != nil {
		return err
	}
	if !sdef.Enabled {
		return fmt.Errorf("请求的服务%s未启用", cnt)
	}
	userid := ""
	if sdef.Security {
		user := flowcontext.GetVarbiableByName("userid")
		if user == nil {
			return fmt.Errorf("服务需要userid参数")
		}
		userid = user.(string)
		if !GetISevurityServiceInstance().VerifyService(userid, sdef.ServiceId, 0) {
			return fmt.Errorf("未授权的请求,用户id:%s,服务id:%s", userid, sdef.ServiceId)
		}
	}
	handler, ok := SHandlerContainer[sdef.ServiceType]
	if !ok {
		return fmt.Errorf("没有找到" + sdef.ServiceType + "定义的服务接口处理程序")
	}
	h := handler(c, userid)
	h.DoSrv(sdef, h)
	if c.rv == "" {
		c.rv = cnt + "_result"
	}
	flowcontext.SetVarbiable(c.rv, c.GetResponseData())
	return nil
}
func (c *InnerServiceActivity) GetResponseData() interface{} {
	return c.resultData
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// type RequestResponseHandler interface {
//	CreateResponseData(style int, data interface{})
//	GetParam(name string) string
//	GetRequestBody() (*SRequestBody, error)
//}
// 实现RequestResponseHandler接口
func (c *InnerServiceActivity) CreateResponseData(style int, data interface{}) {
	c.resultStyle = style
	c.resultData = data
}

func (c *InnerServiceActivity) GetParam(name string) string {
	s, err := activity.ReplaceExpressionLStr(c.context, c.params[name])
	if err != nil {
		logs.Error("处理EL表达式发生错误，%s", err.Error())
		return c.params[name]
	} else {
		return s
	}
}

// 不完整的替换EL表达式方法，只处理了Array、Slice、Map、String这几种类型
// Array、Slice会重置为[]interface{}
// Map 会重置为map[string]interface{}
func (c *InnerServiceActivity) replaceBodyEL(v interface{}) interface{} {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Map:
		{
			s := reflect.ValueOf(v)
			pvs := make(map[string]interface{})
			vs := s.MapKeys()
			for _, k := range vs {
				nv := c.replaceBodyEL(s.MapIndex(k).Interface())
				pvs[k.Interface().(string)] = nv
			}
			return pvs
		}
	case reflect.Slice, reflect.Array:
		{
			s := reflect.ValueOf(v)
			pvs := make([]interface{}, s.Len(), s.Len())
			for i := 0; i < s.Len(); i++ {
				pvs[i] = c.replaceBodyEL(s.Index(i).Interface())
			}
			return pvs
		}
	case reflect.String:
		{
			s := strings.TrimSpace(v.(string))
			if s[:2] == "${" {
				//一上来就是EL表达式的
				eL, err := activity.ReplaceExpressionL(c.context, s)
				if err != nil {
					logs.Error("处理EL表达式发生错误，%s", err.Error())
					return v
				}
				return eL
			} else {
				//在字符串中间的EL表达式
				eL, err := activity.ReplaceExpressionLStr(c.context, v.(string))
				if err != nil {
					logs.Error("处理EL表达式发生错误，%s", err.Error())
					return v
				}
				return eL
			}
		}
	default:
		return v
	}
}

func (c *InnerServiceActivity) GetRequestBody() (*SRequestBody, error) {
	r := &SRequestBody{}
	m := make(map[string]interface{})
	for k, v := range c.body {
		m[k] = c.replaceBodyEL(v)
	}
	err := mapstructure.Decode(m, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// 实现RequestResponseHandler接口结束
//////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// 	{
//		"style" : "innerservice",
//		"cnt":"",
//		"params":{
//			"name":"value"
//		}
//		"rbody":{}
//	}
func CreateInnerServiceActivity(acti *activity.Activity) (activity.IActivity, error) {
	cnt := acti.GetDef()["cnt"]
	if cnt == nil || cnt.(string) == "" {
		return nil, fmt.Errorf("CreateInnerServiceActivity:cnt属性不能为空")
	}

	pm := utils.GetMapFromMap(acti.GetDef(), "params")
	if pm == nil {
		return nil, fmt.Errorf("CreateInnerServiceActivity:param属性不能为空")
	}
	rb := utils.GetMapFromMap(acti.GetDef(), "rbody")
	if rb == nil {
		return nil, fmt.Errorf("CreateInnerServiceActivity:rbody属性不能为空")
	}

	act := &InnerServiceActivity{Activity: *acti}
	s, ok := acti.GetDef()["resultvariable"]
	if ok {
		act.rv = s.(string)
	} else {
		act.rv = ""
	}
	act.params = make(map[string]string)
	for k, v := range pm {
		act.params[k] = v.(string)
	}
	act.body = rb
	act.cnt = cnt.(string)
	return act, nil
}
