package service

import (
	"fmt"
	"tongserver.dataserver/activity"
	"tongserver.dataserver/utils"
	"tongserver.dataserver/utils/mapstructure"
)

// "InnerService": {
// 		"style": "InnerService",
// 		"url":"${}",
// 		"params": {},
// }
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
}

func (c *InnerServiceActivity) Execute(flowcontext activity.IContext) error {
	sdef, err := GetSrvMetaFromPath(c.cnt)
	if err != nil {
		return err
	}
	if !sdef.Enabled {
		return fmt.Errorf("请求的服务%s未启用", c.cnt)
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
		c.rv = c.cnt + "_result"
	}
	flowcontext.SetVarbiable(c.rv, c.GetResponseData())
	return nil
}

func (c *InnerServiceActivity) CreateResponseData(style int, data interface{}) {
	c.resultStyle = style
	c.resultData = data
}
func (c *InnerServiceActivity) GetResponseData() interface{} {
	return c.resultData
}
func (c *InnerServiceActivity) GetParam(name string) string {
	return c.params[name]
}

func (c *InnerServiceActivity) GetRequestBody() (*SRequestBody, error) {
	r := &SRequestBody{}
	err := mapstructure.Decode(c.body, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// 	{
//		"style" : "innerservice",
//		"cnt":"",
//		"params":{
//			"name":"value"
//		}
//		"rbody":{}
//	}
func CreateInnerServiceActivity(acti *activity.Activity) (activity.IActivity, error) {
	// TODO: 这个地方添加创建InnerServiceActivit的方法
	// 1、通过memRRHandler类来调用内部服务
	// 2、模仿app_test中的代码调用内部服务
	// 3、该活动执行的时候需要考虑如何获取上下文的问题，目前打算类似预定义服务的方式，将服务请求的JSON描述放到流程的定义中，使用EL表达式来替换其中的值，这似乎是一个好办法。
	// 4、似乎需要增加一个活动用来专门从请求中获取参数，然后拼接然后启动。
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
