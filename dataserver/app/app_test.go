package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"tongserver.dataserver/activity"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/service"
	"tongserver.dataserver/utils"
)

var SRV_URL = "http://127.0.0.1:8081"
var TOKEN = ""

// 用于测试的RequestResponse类
type innerRRHandler struct {
	body []byte
	p    map[string]string
}

// CreateInnerRR
func CreateInnerRR(bo []byte, pa map[string]string) service.RequestResponseHandler {
	return &innerRRHandler{bo, pa}
}

// CreateResponseData
func (c *innerRRHandler) CreateResponseData(style int, data interface{}) {
	if style == service.RSP_DATA_STYLE_JSON {
		b, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		var str bytes.Buffer
		json.Indent(&str, b, "", "    ")
		fmt.Println(str.String())
	}
}

// GetParam
func (c *innerRRHandler) GetParam(name string) string {
	return c.p[name]
}

// GetRequestBody
func (c *innerRRHandler) GetRequestBody() (*service.SRequestBody, error) {
	rBody := &service.SRequestBody{}
	err := json.Unmarshal([]byte(c.body), rBody)
	if err != nil {
		return nil, fmt.Errorf("解析报文时发生错误%s", err.Error())
	}
	return rBody, nil
}

// 测试前初始化环境
func TestMain(m *testing.M) {
	err := orm.RegisterDataBase("default", "mysql", "tong:123456@tcp(127.0.0.1:3306)/idb", 30)
	datasource.DBAlias2DBTypeContainer["default"] = "mysql"
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	code := m.Run()
	os.Exit(code)
}

// 发送post请求
func postData(path string, jsondata string, header map[string]string) map[string]interface{} {
	url := SRV_URL + path //
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsondata)))
	req.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	if err != nil {
		panic(err)
	}
	rm, err := utils.ParseJSONBytes2Map(body)
	if err != nil {
		panic(err)
	}
	return rm
}

// 调用远程服务创建用户令牌，运行前先启动app.main
func TestCreateUserToken(t *testing.T) {
	jsonData := `{
	"LoginName":"lvxing",
	"Password":"123"
	}`
	rm := postData("/token/create", jsonData, nil)
	if rm["token"] == false {
		t.Fatalf("登录失败")
		return
	}
	TOKEN = rm["token"].(string)
	fmt.Println("TOKEN IS:" + TOKEN)
}

// 调用远程meta服务获取元数据，运行前先启动app.main
func TestJeda_meta(t *testing.T) {
	TestCreateUserToken(t)
	if t.Failed() {
		return
	}
	jsonData := `{}`
	postData("/services/jeda.meta/all?_repstyle=map", jsonData, map[string]string{"Authorization": TOKEN})
}

// 创建用于系统管理的默认服务
func createService(cnt string, meta map[string]interface{}) *service.SDefine {
	jsonStr, _ := json.Marshal(meta)

	ps := strings.Split(cnt, ".")
	projectid := "00000000-0000-0000-0000-000000000000"

	srv := &service.SDefine{
		ServiceId:   "JEDA_SRV_" + cnt,
		Context:     ps[1],
		BodyType:    "body",
		ServiceType: "",
		Namespace:   ps[0],
		Enabled:     true,
		MsgLog:      false,
		Security:    true,
		Meta:        string(jsonStr),
		ProjectId:   projectid}
	return srv
}

// 调用内部服务，该服务可以通过RequestResponse类实现发布
// 调用服务的all方法
func serviceCall(srvtype string, sdef *service.SDefine) {
	ReadCfg()
	CreateIDSCreator()
	f := service.SHandlerContainer[srvtype]
	rrh := CreateInnerRR([]byte("{}"), map[string]string{
		":action":   "all",
		"_repstyle": "map"})
	h := f(rrh, "lvxing")
	h.DoSrv(sdef, h)
}

// 测试用户过滤器，使用两级递归的过滤器实现
func Test_G_METAITME(t *testing.T) {
	sdef := createService("jeda.meta", map[string]interface{}{
		"ids": "default.mgr.G_META_ITEM",
		"userfilter": map[string]interface{}{
			"filterkey": "META_ID",
			"values": map[string]interface{}{
				"ids":       "default.mgr.G_META",
				"filterkey": "PROJECTID",
				"outfield":  "ID",
				"values": map[string]interface{}{
					"ids":       "default.mgr.G_USERPROJECT",
					"filterkey": "USERID",
					"outfield":  "PROJECTNAME",
					"values":    "userid",
				}}},
	})
	serviceCall(service.SrvTypeIds, sdef)
	b, err := json.Marshal(sdef)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("===== service define json =====")
	var str bytes.Buffer
	json.Indent(&str, b, "", "    ")
	fmt.Println(str.String())
}

// 测试用户过滤器,测试一级用户过滤器
func Test_G_META(t *testing.T) {
	sdef := createService("jeda.meta", map[string]interface{}{
		"ids": "default.mgr.G_META",
		"userfilter": map[string]interface{}{
			"filterkey": "PROJECTID",
			"values": map[string]interface{}{
				"ids":       "default.mgr.G_USERPROJECT",
				"filterkey": "USERID",
				"outfield":  "PROJECTNAME",
				"values":    "userid"}}})
	serviceCall(service.SrvTypeIds, sdef)
}

func TestCallInnerService(t *testing.T) {
	json := `{
	"name": "测试to flow",
	"start": {
	"params":{
		"name":{"type":"string","value":"menghui"},
		"age":{"type":"number","value":41}
	},
	"variables": {
	   "var_a": {
	     "type": "string",
	     "value": "test var"
	   },
	   "var_b": {
	     "type": "number",
	     "value": 12
	   }
	},
	"flow": [{
		"gate":"to",	
		"target":[{
			"style" : "innerservice",
			"resultvariable":"result",
			"cnt":"jeda.meta",
			"params":{
				":action":"all"
			},
			"rbody":{}
		}]
	}]
}}`
	ReadCfg()
	CreateIDSCreator()
	fl, err := activity.NewFlowInstanceFromJSON(json)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := map[string]interface{}{
		"userid": "lvxing",
	}
	err = fl.Execute(r)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
