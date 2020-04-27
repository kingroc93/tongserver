package activity

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseMapPath(t *testing.T) {
	p := "[1:3]"
	mv := []string{"0", "1", "2", "3", "4"}
	r, err := parsePaths(strings.Split(p, "/"), 0, mv, nil, false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r)
}
