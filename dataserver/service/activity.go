package service

import (
	"tongserver.dataserver/activity"
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
	cnt string
}

// Execute
func (InnerServiceActivity) Execute(flowcontext activity.IContext) error {
	//QueryServiceFromDB(cnt string, ns string, context string) (*SDefine, error)
	return nil
}

// memRRHandler
type memRRHandler struct {
	body       []byte
	p          map[string]string
	resultData interface{}
}

// GetResult
func (c *memRRHandler) GetResult() interface{} {
	return c.resultData
}

// CreateResponseData
func (c *memRRHandler) CreateResponseData(style int, data interface{}) {
	c.resultData = data
}

// GetParam
func (c *memRRHandler) GetParam(name string) string {
	return (c.p)[name]
}

// GetRequestBody
func (c *memRRHandler) GetRequestBody() []byte {
	return c.body
}

func CreateInnerServiceActivity(acti *activity.Activity) (activity.IActivity, error) {
	return nil, nil
}
