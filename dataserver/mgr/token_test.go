package mgr

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Before ====================")
	err := orm.RegisterDataBase("default", "mysql", "tong:123456@tcp(127.0.0.1:3306)/idb", 30)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	code := m.Run()
	fmt.Println("End ====================")
	os.Exit(code)
}
func TestJeda(t *testing.T) {

}
