package service

import (
	"bytes"
	"fmt"
	"github.com/astaxie/beego/orm"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/utils"
)

var SRV_URL = "http://127.0.0.1:8081"

/*
req, _ :=
		http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/js"nq)
resp, err := client.Do(req)

*/
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

func TestCreateUserToken(t *testing.T) {
	jsonData := `{
	"LoginName":"lvxing",
	"Password":"123"
	}`
	url := SRV_URL + "/token/create"
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(jsonData)))
	req.Method = "POST"
	req.Header.Set("Content-Type", "application/json")
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
	if (*rm)["result"] == true {
		fmt.Println("OKOKOKOKOKOK")
	}
	fmt.Println("response Body:", string(body))
}
