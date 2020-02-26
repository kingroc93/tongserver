package mgr

import "C"
import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"tongserver.dataserver/service"
	"tongserver.dataserver/utils"
)

type ControllerWithVerify struct {
}

// SecurityController 安全认证WebAPI
type SecurityController struct {
	beego.Controller
	ControllerWithVerify
}

func (c *ControllerWithVerify) Verifty(ctl *beego.Controller) bool {
	_, err := service.GetTokenServiceInstance().VerifyToken(ctl)
	if err != nil {
		r := utils.CreateRestResult(false)
		r["msg"] = err.Error()
		ctl.Data["json"] = r
		ctl.ServeJSON()
		return false
	}
	return true
}

// VerifyToken 验证令牌是否合法的web api
func (c *SecurityController) VerifyToken() {
	if !c.Verifty(&c.Controller) {
		return
	}
	result := utils.CreateRestResult(true)
	c.Data["json"] = result
	c.ServeJSON()
}

// checkPwd 检验密码是否一致
func (c *SecurityController) checkPwd(u string, p string) bool {

	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT USER_ID,POSITION_ID,ORG_ID,USER_NAME,USER_PASSWORD FROM JEDA_USER WHERE LOGIN_NAME=? and USER_PASSWORD=?", u, p).Values(&maps)
	if err != nil {
		return false
	}
	if len(maps) == 0 {
		return false
	}
	return true

}

// CreateToken 创建令牌
func (c *SecurityController) CreateToken() {
	uname := c.Input().Get("LoginName")
	pwd := c.Input().Get("Password")

	if uname == "" && pwd == "" {
		//尝试通过RequestBody获取
		rbody := &struct {
			LoginName string
			Password  string
		}{}
		err := json.Unmarshal([]byte(c.Ctx.Input.RequestBody), rbody)
		if err != nil {
			c.Data["json"] = utils.CreateRestResult(false)
			c.ServeJSON()
			return
		}
		uname = rbody.LoginName
		pwd = rbody.Password
	}

	t, err := service.GetTokenServiceInstance().CreateToken(uname, pwd)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"result": true, "token": t}
		c.ServeJSON()
		return
	}
	r := utils.CreateRestResult(false)
	r["msg"] = err.Error()
	c.Data["json"] = r
	c.ServeJSON()
}
