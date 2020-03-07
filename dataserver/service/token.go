package service

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/rs/xid"
	"strings"
	"sync"
	"time"
	"tongserver.dataserver/datasource"
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
	VerifyService(userid string, serviceid string, rightmask int) bool
	GetRoleByUserid(userid string) (utils.StringSet, error)
}

type TokenService struct{}

var tokenService ISevurityService
var once sync.Once

// GetTokenServiceInstance 创建一个令牌服务
func GetISevurityServiceInstance() ISevurityService {
	once.Do(func() {
		tokenService = &TokenService{}
	})
	return tokenService
}

//  返回用户的角色信息
// 角色信息放入缓存，缓存key为utils.CACHE_PREFIX_SERVICEACCESS + "ROLE" + userid
func (c *TokenService) GetRoleByUserid(userid string) (utils.StringSet, error) {
	var rolemap interface{}
	rolemap = utils.JedaDataCache.Get(utils.CACHE_PREFIX_SERVICEACCESS + "ROLE" + userid)
	if rolemap == nil {
		sqld := datasource.CreateSQLDataSource("", "default", "select * from idb.JEDA_ROLE_USER where USER_ID=?")
		sqld.ParamsValues = []interface{}{userid}
		rs, err := sqld.GetAllData()
		if err != nil {
			return nil, err
		}
		if len(rs.Data) == 0 {
			logs.Info("get user role not found，user:%s", userid)
			return nil, nil
		}
		r := make(utils.StringSet)
		for _, item := range rs.Data {
			r.Put(item[rs.Fields["ROLE_ID"].Index].(string))
		}
		rolemap = r
		err = utils.JedaDataCache.Put(utils.CACHE_PREFIX_SERVICEACCESS+"ROLE"+userid, rolemap, 60*time.Second)
		if err != nil {
			logs.Error(err.Error())
		}
	}
	o, ok := rolemap.(utils.StringSet)
	if !ok {
		return nil, fmt.Errorf("GetRoleByUserid 类型转换失败rolemap.(utils.StringSet)")
	}
	return o, nil
}

// VerifyService 检验用户是否可以访问指定的服务，rightmask参数目前尚未处理
func (c *TokenService) VerifyService(userid string, serviceid string, rightmask int) bool {
	var srvmap interface{}
	srvmap = utils.JedaDataCache.Get(utils.CACHE_PREFIX_SERVICEACCESS + serviceid)
	if srvmap == nil {
		sqld := datasource.CreateSQLDataSource("", "default",
			"SELECT * FROM idb.G_USERSERVICE where SERVICEID=?")
		sqld.ParamsValues = []interface{}{serviceid}
		rs, err := sqld.GetAllData()
		if err != nil {
			logs.Error("VerifyService error:serviceid:%s\t userid:%s\t error:%s", serviceid, userid, err.Error())
			return false
		}
		serviceaccrssMap := make(map[string]utils.StringSet)
		for _, item := range rs.Data {
			o, ok := serviceaccrssMap[item[rs.Fields["SERVICEID"].Index].(string)]
			if !ok {
				o = make(utils.StringSet)
			}
			o.Put(item[rs.Fields["ROLEID"].Index].(string))
			serviceaccrssMap[item[rs.Fields["SERVICEID"].Index].(string)] = o
		}
		err = utils.JedaDataCache.Put(utils.CACHE_PREFIX_SERVICEACCESS+serviceid, serviceaccrssMap, 60*time.Second)
		if err != nil {
			logs.Error("VerifyService  utils.JedaDataCache.Put error : %s", err.Error())
		}
		srvmap = serviceaccrssMap
	}
	sm, ok := srvmap.(map[string]utils.StringSet)
	if !ok {
		logs.Error("VerifyService 类型转换错误，希望类型map[string]utils.StringSet")
		return false
	}
	v, ok := sm[serviceid]
	if !ok {
		return false
	}
	r, err := c.GetRoleByUserid(userid)
	if err != nil {
		logs.Error("VerifyService 获取用户角色信息发生错误," + err.Error())
		return false
	}
	return !r.Mix(&v).IsEmpty()
}

// checkPwd 检验密码
func (c *TokenService) checkPwd(u string, p string) bool {
	src, err := datasource.CreateIDSFromName("default.mgr.JEDA_USER")
	if err != nil {
		logs.Error("获取默认数据源时发生错误，数据源名称：default.mgr.JEDA_USER")
		return false
	}
	inf := src.(datasource.ICriteriaDataSource)
	src.(datasource.IFilterAdder).AddCriteria("LOGIN_NAME", datasource.OperEq, u).AndCriteria("USER_PASSWORD", datasource.OperEq, p)
	rs, err := inf.DoFilter()
	if err != nil {
		logs.Error("校验用户名密码时发生错误，" + err.Error())
		return false
	}
	if len(rs.Data) == 0 {
		return false
	}
	return true

}

// CreateToken 创建令牌
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
