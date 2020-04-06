package restclient

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/http"
	"tongserver.dataserver/utils"
)

// 调动Restful服务的客户端
type RestClient struct {
	url   string
	token string
}

// 返回一个调用传入url的客户端类
func CreateRestClient(url string) *RestClient {
	return &RestClient{"url", ""}
}

func (c *RestClient) CreateToken(user string, pwd string) (string, error) {
	jsonData := `{
	"LoginName":"lvxing",
	"Password":"123"
	}`
	url := c.url + "/token/create"
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	req.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("访问 %s 时发生错误，", url, err.Error())
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	rm, err := utils.ParseJSONBytes2Map(body)
	if err != nil {
		logs.Error("访问 %s 时发生错误，", url, err.Error())
		return "", err
	}
	if (*rm)["result"] == true {
		c.token = (*rm)["token"].(string)
		return (*rm)["token"].(string), nil
	}
	return "", fmt.Errorf("创建令牌失败，服务端返回:%s", string(body))
}
