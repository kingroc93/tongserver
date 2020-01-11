package utils

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	s := GetHmacCode("aaasadfasdfasdfasdfasdfaswqer13r1rewrqwfsdafasdfqewrqweqfasdvasd", "menghui")
	fmt.Println(s)
}
