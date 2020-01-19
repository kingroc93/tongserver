package service

import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/rs/xid"
	"strings"
	"time"
	"tongserver.dataserver/utils"
)

// HASHSECRET hash算法种子
var HASHSECRET = "1@3wq,klahjaqwweq"

// TokenExpire 令牌默认过期时间60秒
var TokenExpire int64 = 60

// SecurityController 安全认证WebAPI
type SecurityController struct {
	beego.Controller
}

// VerifyToken 验证令牌是否合法，从beego控制器中获取令牌信息
func VerifyToken(c *beego.Controller) (bool, RestResult, error) {
	authString := c.Ctx.Input.Header("Authorization")
	if authString == "" {
		return false, nil, fmt.Errorf("invalid Authorization in request header")
	}
	ss := strings.Split(authString, ".")
	if len(ss) != 2 {
		return false, nil, fmt.Errorf("invalid Authorization in request header")
	}
	js := utils.DecodeURLBase64(ss[0])
	if utils.GetHmacCode(js, HASHSECRET) != ss[1] {
		return false, nil, nil
	}
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(js), &m)
	if err != nil {
		return false, nil, nil
	}
	n := time.Now().UnixNano()
	if n-int64(m["time"].(float64)) > TokenExpire*1e9 {
		return false, nil, nil
	}
	m["time"] = time.Now().UnixNano()
	return true, m, nil
}

// ConvertJSON 装换为Json格式字符串
func ConvertJSON(data interface{}, encoding ...bool) (string, error) {
	var (
		hasIndent = beego.BConfig.RunMode != beego.PROD
		content   []byte
		err       error
	)
	if hasIndent {
		content, err = json.MarshalIndent(data, "", "  ")
	} else {
		content, err = json.Marshal(data)
	}
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// VerifyToken 验证令牌是否合法的web api
func (c *SecurityController) VerifyToken() {
	r, rm, err := VerifyToken(&c.Controller)
	if err != nil {
		r := CreateRestResult(false)
		r["msg"] = err.Error()
		c.Data["json"] = r
		c.ServeJSON()
		return
	}
	result := CreateRestResult(r)
	if r {
		js, err := ConvertJSON(rm)
		if err != nil {
			r := CreateRestResult(true)
			r["msg"] = "认证成功但令牌未刷新" + err.Error()
			c.Data["json"] = r
			c.ServeJSON()
			return
		}
		result["token"] = utils.EncodeURLBase64(js) + "." + utils.GetHmacCode(js, HASHSECRET)
	}
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
	uname := c.Input().Get("u")
	pwd := c.Input().Get("p")
	if !c.checkPwd(uname, pwd) {
		c.Data["json"] = CreateRestResult(false)
	} else {
		r := CreateRestResult(true)
		r["sid"] = xid.New().String()
		r["userid"] = uname
		r["time"] = time.Now().UnixNano()
		js, err := ConvertJSON(r)
		if err != nil {
			r := CreateRestResult(false)
			r["msg"] = "转换data到json时发生错误," + err.Error()
			c.Data["json"] = r
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{"result": true, "token": utils.EncodeURLBase64(js) + "." + utils.GetHmacCode(js, HASHSECRET)}
	}
	c.ServeJSON()
}
