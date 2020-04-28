package activity

import (
	"fmt"
	"strings"
	"testing"
)

func TestParseMapPath(t *testing.T) {
	p := "[1:2]"
	mv := []string{"0", "1", "2", "3", "4"}
	r, err := parsePaths(strings.Split(p, "/"), 0, mv, nil, false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r)

	m := map[string]interface{}{
		"a": 11,
		"b": map[string]interface{}{
			"aa": 111,
			"bb": 222,
			"cc": []string{"0", "1", "2", "3", "4"},
		}}
	p = "/b/cc/[1]"
	r, err = ParseMapPath(m, p, nil, false)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r)
}
