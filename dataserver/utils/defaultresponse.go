package utils

import (
	"github.com/astaxie/beego"
)

// RestResult 请求返回的数据类型
type RestResult map[string]interface{}

// CreateRestResult 创建返回数据实例
func CreateRestResult(success bool) RestResult {
	var result = make(RestResult)
	result["result"] = success
	return result
}

// CreateErrorResponse 创建错误信息的响应
func CreateErrorResponse(msg string, Ctl *beego.Controller) {
	r := CreateRestResult(false)
	r["msg"] = msg
	Ctl.Data["json"] = r
	Ctl.ServeJSON()
}

// CreateErrorResponseByError 根据Error类创建错误信息
func CreateErrorResponseByError(err error, Ctl *beego.Controller) {
	r := CreateRestResult(false)
	r["msg"] = err.Error()
	Ctl.Data["json"] = r
	Ctl.ServeJSON()
}
