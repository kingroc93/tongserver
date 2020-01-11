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

const HASHSECRET = "1@3wq,klahjaqwweq"

type SecurityController struct {
	beego.Controller
}

func (c *SecurityController) Get() {
	c.Data["logined"] = "TRUE"
	c.ServeJSON()
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
	authString := c.Ctx.Input.Header("Authorization")
	ss := strings.Split(authString, ".")
	js := utils.DecodeURLBase64(ss[0])
	if utils.GetHmacCode(js, HASHSECRET) != ss[1] {
		c.Data["json"] = service.CreateRestResult(false)
	} else {
		m := make(map[string]interface{})
		err := json.Unmarshal([]byte(js), &m)
		if err != nil {
			c.Data["json"] = service.CreateRestResult(false)
		} else {
			n := time.Now().UnixNano()

			if n-int64(m["time"].(float64)) > 60*1e9 {
				c.Data["json"] = service.CreateRestResult(false)
			} else {
				c.Data["json"] = service.CreateRestResult(true)
			}

		}
	}
	c.ServeJSON()
}

func (c *SecurityController) CreateToken() {
	uname := c.Input().Get("u")
	pwd := c.Input().Get("p")
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT USER_ID,POSITION_ID,ORG_ID,USER_NAME,USER_PASSWORD FROM JEDA_USER WHERE LOGIN_NAME=? and USER_PASSWORD=?", uname, pwd).Values(&maps)
	if err != nil {
		c.Data["json"] = service.CreateRestResult(false)
	} else {
		if len(maps) == 0 {
			c.Data["json"] = service.CreateRestResult(false)
		} else {
			r := service.CreateRestResult(true)
			r["sid"] = xid.New().String()
			r["userid"] = maps[0]["USER_ID"]
			r["time"] = time.Now().UnixNano()
			js, err := c.ConvertJson(r)
			if err != nil {
				r := service.CreateRestResult(false)
				r["msg"] = "转换data到json时发生错误," + err.Error()
			}
			c.Data["json"] = map[string]interface{}{"result": true, "token": utils.EncodeURLBase64(js) + "." + utils.GetHmacCode(js, HASHSECRET)}
		}
	}
	c.ServeJSON()
}
func (c *SecurityController) RefashToken() {

}
