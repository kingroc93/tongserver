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
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/service"
	"tongserver.dataserver/utils"
)

var SRV_URL = "http://127.0.0.1:8081"
var TOKEN = ""

//type RequestResponseHandler interface {
//	CreateResponseData(style int, data interface{})
//	GetParam(name string) string
//	GetRequestBody() []byte
//}

type innerRRHandler struct {
	body []byte
	p    *map[string]string
}

func CreateInnerRR(bo []byte, pa *map[string]string) service.RequestResponseHandler {
	return &innerRRHandler{bo, pa}
}

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

func (c *innerRRHandler) GetParam(name string) string {
	return (*c.p)[name]
}

func (c *innerRRHandler) GetRequestBody() []byte {
	return c.body
}

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
func postData(path string, jsondata string, header *map[string]string) *map[string]interface{} {
	url := SRV_URL + path //
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsondata)))
	req.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
	if header != nil {
		for key, value := range *header {
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

func TestCreateUserToken(t *testing.T) {
	jsonData := `{
	"LoginName":"lvxing",
	"Password":"123"
	}`
	rm := postData("/token/create", jsonData, nil)
	if (*rm)["token"] == false {
		t.Fatalf("登录失败")
		return
	}
	TOKEN = (*rm)["token"].(string)
	fmt.Println("TOKEN IS:" + TOKEN)
}

func TestJeda_meta(t *testing.T) {
	TestCreateUserToken(t)
	if t.Failed() {
		return
	}
	jsonData := `{}`
	postData("/services/jeda.meta/all?_repstyle=map", jsonData, &map[string]string{"Authorization": TOKEN})
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

func serviceCall(sdef *service.SDefine) {
	ReadCfg()
	CreateIDSCreator()
	f := service.SHandlerContainer[service.SrvTypeIds]
	rrh := CreateInnerRR([]byte("{}"), &map[string]string{
		":action":   "all",
		"_repstyle": "map"})
	h := f(rrh, "lvxing")
	h.DoSrv(sdef, h)
}

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
	serviceCall(sdef)
}

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
	serviceCall(sdef)
}
