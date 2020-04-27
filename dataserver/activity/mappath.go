package activity

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"tongserver.dataserver/utils"
)

func parsePaths(ts []string, index int, m interface{}, context IContext, enabledEL bool) (interface{}, error) {
	var (
		name string
		err  error
	)
	if index >= len(ts) {
		return m, nil
	}
	name = ts[index]
	if enabledEL {
		// 计算EL表达式，获取路径名称
		name, err = ReplaceExpressionLStr(context, name)
		if err != nil {
			return nil, err
		}
	}
	ty := reflect.TypeOf(m)
	ns := name[:1]
	if ns == "@" {

	}
	if ns == "[" {
		if ty.Kind() != reflect.Array && ty.Kind() != reflect.Slice {
			return nil, fmt.Errorf("mappath解析错误，[]运算符只能用于数组或切片")
		}
		mapvalue := reflect.ValueOf(m)

		ind := strings.Index(name, "]")
		if ind == -1 {
			return nil, fmt.Errorf("mappath解析错误，语法错误，%s", name)
		}
		sub := name[1:ind] // sub是[]内部的字符串
		if strings.Index(sub, ":") == -1 {
			// 没有找到":"
			i, err := utils.String(sub).Int()
			if err != nil {
				return nil, fmt.Errorf("mappath解析错误，字符串 %s 不能转换为数字，%s", sub, name)
			}
			if i > mapvalue.Len() {
				return nil, fmt.Errorf("mappath解析，[%s]越界，数组长度为%s", strconv.Itoa(i), strconv.Itoa(mapvalue.Len()))
			}
			return mapvalue.Index(i).Interface(), nil
		} else {
			// 找到":"  处理切片语法
			si := strings.Split(sub, ":")
			if len(si) != 2 {
				return nil, fmt.Errorf("mappath解析错误，语法错误,%s", sub)
			}
			start := 0
			end := mapvalue.Len()
			if si[0] != "" {
				start, err = utils.String(si[0]).Int()
				if err != nil {
					return nil, fmt.Errorf("mappath解析错误，字符串 %s 不能转换为数字，%s", si[0], name)
				}
			}
			if si[0] != "" {
				end, err = utils.String(si[1]).Int()
				if err != nil {
					return nil, fmt.Errorf("mappath解析错误，字符串 %s 不能转换为数字，%s", si[1], name)
				}
			}
			if start >= end {
				return nil, fmt.Errorf("mappath解析错误，%s，起始值大于等于终止值", sub)
			}
			L := end - start + 1
			result := reflect.MakeSlice(ty, L, L)
			for i := 0; i < L; i++ {
				result.Index(i).Set(mapvalue.Index(start + i))
			}
			return result.Interface(), nil
		}

	}
	return m, nil
}

func ParseMapPath(m map[string]interface{}, ps string, context IContext, enabledEL bool) (interface{}, error) {
	ps = strings.TrimSpace(ps)
	if ps == "" {
		return nil, nil
	}
	if m == nil {
		return nil, nil
	}
	if ps[:1] != "/" {
		return nil, fmt.Errorf("mappath表达式必须以/开头")
	}
	ps = ps[1:]
	ts := strings.Split(ps, "/")
	return parsePaths(ts, 0, m, context, enabledEL)
}
