package mgr

import (
	"github.com/astaxie/beego"
)

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
