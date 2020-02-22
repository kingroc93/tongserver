package mgr

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/orm"
	"github.com/rs/xid"
	"strings"
	"sync"
	"time"
	"tongserver.dataserver/utils"
)

// HASHSECRET hash算法种子
var HASHSECRET = "1@3wq,klahjaqwweq"

// TokenExpire 令牌默认过期时间60秒
var TokenExpire int64 = 60

// 令牌处理接口
type ISevurityService interface {
	VerifyToken(c *beego.Controller) (string, error)
	VerifyTokenCtx(ctx *context.Context) (string, error)
	CreateToken(userid string, pwd string) (string, error)
}
type TokenService struct{}

var tokenService ISevurityService
var once sync.Once

func GetTokenServiceInstance() ISevurityService {
	once.Do(func() {
		tokenService = &TokenService{}
	})
	return tokenService
}

func (c *TokenService) checkPwd(u string, p string) bool {
	obj := utils.JedaDataCache.Get(u)
	if obj != nil {
		if p == obj.(string) {
			return true
		}
	}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT USER_ID FROM JEDA_USER WHERE LOGIN_NAME=? and USER_PASSWORD=?", u, p).Values(&maps)
	if err != nil {
		return false
	}
	if len(maps) == 0 {
		return false
	}
	return true

}

func (c *TokenService) CreateToken(uname string, pwd string) (string, error) {
	if !c.checkPwd(uname, pwd) {
		return "", fmt.Errorf("验证失败")
	} else {
		r := make(map[string]interface{})
		r["result"] = true
		r["sid"] = xid.New().String()
		r["userid"] = uname
		r["time"] = time.Now().UnixNano()
		js, err := utils.ConvertJSON(r)
		if err != nil {
			return "", fmt.Errorf("转换data到json时发生错误," + err.Error())
		}
		return utils.EncodeURLBase64(js) + "." + utils.GetHmacCode(js, HASHSECRET), nil
	}
}

// 验证令牌并刷新令牌
// 验证成功后刷新令牌下发的时间为当前时间
func (c *TokenService) VerifyTokenCtx(ctx *context.Context) (string, error) {
	authString := ctx.Input.Header("Authorization")
	if authString == "" {
		return "", fmt.Errorf("invalid Authorization in request header")
	}
	ss := strings.Split(authString, ".")
	if len(ss) != 2 {
		return "", fmt.Errorf("invalid Authorization in request header")
	}
	js := utils.DecodeURLBase64(ss[0])
	if utils.GetHmacCode(js, HASHSECRET) != ss[1] {
		return "", fmt.Errorf("invalid Authorization in request header")
	}
	m := make(map[string]interface{})
	err := json.Unmarshal([]byte(js), &m)
	if err != nil {
		return "", fmt.Errorf("invalid Authorization in request header")
	}
	n := time.Now().UnixNano()
	if n-int64(m["time"].(float64)) > TokenExpire*1e9 {
		return "", fmt.Errorf("invalid Authorization in request header")
	}
	m["time"] = time.Now().UnixNano()
	jss, _ := utils.ConvertJSON(m)
	ctx.ResponseWriter.Header().Add("token", utils.EncodeURLBase64(jss)+"."+utils.GetHmacCode(js, HASHSECRET))
	return m["userid"].(string), nil
}

// VerifyToken 验证令牌是否合法，从beego控制器中获取令牌信息
func (c *TokenService) VerifyToken(ctl *beego.Controller) (string, error) {
	return c.VerifyTokenCtx(ctl.Ctx)
}
