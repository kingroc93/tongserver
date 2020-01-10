package service

import (
	"github.com/astaxie/beego"
)


func CreateErrorResponse(msg string,Ctl *beego.Controller){
	r := CreateRestResult(false)
	r["msg"] = msg
	Ctl.Data["json"] = r
	Ctl.ServeJSON()
}
func  CreateErrorResponseByError(err error,Ctl *beego.Controller) {
	r := CreateRestResult(false)
	r["msg"] = err.Error()
	Ctl.Data["json"] = r
	Ctl.ServeJSON()
}
func CreateErrorResult(msg string,Ctl *beego.Controller) {
	r := CreateRestResult(false)
	r["msg"] = msg
	Ctl.Data["json"] = r
	Ctl.ServeJSON()
}