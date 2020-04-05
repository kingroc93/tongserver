package service

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestJeda(t *testing.T) {
	r, err := GetISevurityServiceInstance().GetRoleByUserid("lvxing")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for k, _ := range r {
		fmt.Println(k)
	}
}

func TestTokenService_VerifyService(t *testing.T) {
	b := GetISevurityServiceInstance().VerifyService("lvxing", "26d7e145-9d6f-434c-973d-7ef191322545", 255)
	fmt.Println(b)
	b = GetISevurityServiceInstance().VerifyService("lvxing", "26d7e145-9d6f-434c-973d-7ef191322545", 255)
	fmt.Println(b)
	b = GetISevurityServiceInstance().VerifyService("lvxing", "26d7e145-9d6f-434c-973d-7ef191322545", 255)
	fmt.Println(b)
}
