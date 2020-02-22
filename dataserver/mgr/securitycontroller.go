package mgr

import "C"
import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"tongserver.dataserver/utils"
)

// SecurityController 安全认证WebAPI
type SecurityController struct {
	beego.Controller
}

// VerifyToken 验证令牌是否合法的web api
func (c *SecurityController) VerifyToken() {
	_, err := GetTokenServiceInstance().VerifyToken(&c.Controller)
	if err != nil {
		r := CreateRestResult(false)
		r["msg"] = err.Error()
		c.Data["json"] = r
		c.ServeJSON()
		return
	}
	result := CreateRestResult(true)
	c.Data["json"] = result
	c.ServeJSON()
}

// checkPwd 检验密码是否一致
func (c *SecurityController) checkPwd(u string, p string) bool {
	obj := utils.JedaDataCache.Get(u)
	if obj != nil {
		if p == obj.(string) {
			return true
		}
	}
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
			c.Data["json"] = CreateRestResult(false)
			c.ServeJSON()
			return
		}
		uname = rbody.LoginName
		pwd = rbody.Password
	}

	t, err := GetTokenServiceInstance().CreateToken(uname, pwd)
	if err == nil {
		c.Data["json"] = map[string]interface{}{"result": true, "token": t}
		c.ServeJSON()
		return
	}
	r := CreateRestResult(false)
	r["msg"] = err.Error()
	c.Data["json"] = r
	c.ServeJSON()
}
