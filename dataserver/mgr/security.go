package mgr

import "C"
import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/rs/xid"
	"strings"
	"time"
	"tongserver.dataserver/service"
	"tongserver.dataserver/utils"
)

var HASHSECRET = "1@3wq,klahjaqwweq"
var TokenExpire int64 = 60

type SecurityController struct {
	beego.Controller
}

func (c *SecurityController) Get() {
	c.Data["logined"] = "TRUE"
	c.ServeJSON()
}

func VerifyToken(c *beego.Controller) (bool, service.RestResult) {
	authString := c.Ctx.Input.Header("Authorization")
	ss := strings.Split(authString, ".")
	js := utils.DecodeURLBase64(ss[0])
	if utils.GetHmacCode(js, HASHSECRET) != ss[1] {
		return false, nil
	} else {
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(js), &m)
		if err != nil {
			return false, nil
		} else {
			n := time.Now().UnixNano()
			if n-int64(m["time"].(float64)) > TokenExpire*1e9 {
				return false, nil
			} else {
				m["time"] = time.Now().UnixNano()
				return true, m
			}
		}
	}
}

func (c *SecurityController) ConvertJson(data interface{}, encoding ...bool) (string, error) {
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

func (c *SecurityController) VerifyToken() {
	r, rm := VerifyToken(&c.Controller)
	result := service.CreateRestResult(r)
	if r {
		js, err := c.ConvertJson(rm)
		if err != nil {
			r := service.CreateRestResult(true)
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
	} else {
		return true
	}
}

func (c *SecurityController) CreateToken() {
	uname := c.Input().Get("u")
	pwd := c.Input().Get("p")
	if !c.checkPwd(uname, pwd) {
		c.Data["json"] = service.CreateRestResult(false)
	} else {
		r := service.CreateRestResult(true)
		r["sid"] = xid.New().String()
		r["userid"] = uname
		r["time"] = time.Now().UnixNano()
		js, err := c.ConvertJson(r)
		if err != nil {
			r := service.CreateRestResult(false)
			r["msg"] = "转换data到json时发生错误," + err.Error()
			c.Data["json"] = r
			c.ServeJSON()
			return
		}
		c.Data["json"] = map[string]interface{}{"result": true, "token": utils.EncodeURLBase64(js) + "." + utils.GetHmacCode(js, HASHSECRET)}
	}
	c.ServeJSON()
}
