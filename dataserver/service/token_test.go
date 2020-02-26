package service

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
	"tongserver.dataserver/datasource"
)

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
func TestJeda(t *testing.T) {
	r, err := GetTokenServiceInstance().GetRoleByUserid("lvxing")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for k, _ := range r {
		fmt.Println(k)
	}
}

func TestTokenService_VerifyService(t *testing.T) {
	b := GetTokenServiceInstance().VerifyService("lvxing", "26d7e145-9d6f-434c-973d-7ef191322545", 255)
	fmt.Println(b)
	b = GetTokenServiceInstance().VerifyService("lvxing", "26d7e145-9d6f-434c-973d-7ef191322545", 255)
	fmt.Println(b)
	b = GetTokenServiceInstance().VerifyService("lvxing", "26d7e145-9d6f-434c-973d-7ef191322545", 255)
	fmt.Println(b)
}
