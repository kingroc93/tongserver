package activity

import (
	"fmt"
	"github.com/antonmedv/expr"
	"strings"
	"tongserver.dataserver/utils"
)

const (
	ASSIGN_OPER  string = "="
	EXP_LEFTDIV  string = "${"
	EXP_RIGHTDIV string = "}"
)

// DoExpression
func DoExpression(expression string, context IContext) (interface{}, error) {
	return DoExpression2(expression, context.getVariableMap())
}

// 执行表达式
func DoExpression2(expression string, env map[string]interface{}) (interface{}, error) {
	program, err := expr.Compile(expression, expr.Env(env))
	if err != nil {
		return nil, err
	}
	output, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// 执行表达式并返回bool类型
func DoExpressionBool(expression string, context IContext) (bool, error) {
	return DoExpressionBool2(expression, context.getVariableMap())
}

// 运行表达式 返回true or false
func DoExpressionBool2(expression string, env map[string]interface{}) (bool, error) {
	output, err := DoExpression2(expression, env)
	if err != nil {
		return false, err
	}
	return utils.ConvertObj2Bool(output), nil
}

// 执行一堆表达式。有一个出错则返回错误信息忽略后面的表达式
// 支持赋值表达式，=号左侧为变量名，如果上下文里存在则替换，如果不存在则创建
func ExecuteExpressions(flowcontext IContext, exps []string) error {
	env := flowcontext.getVariableMap()
	vmap := make(map[string]interface{})
	for _, exp := range exps {
		v, e, ok := SplitAssignExpression(exp)
		if ok {
			vr, err := DoExpression2(e, env)
			if err != nil {
				return err
			}
			vmap[v] = vr
		} else {
			_, err := DoExpression2(exp, env)
			return err
		}
	}
	for k, v := range vmap {
		flowcontext.SetVarbiable(k, v)
	}
	return nil
}

// 分隔赋值表达式，返回值 (变量名,表达式,是否为赋值表达式)
func SplitAssignExpression(expression string) (string, string, bool) {
	exp := strings.TrimSpace(expression)
	i := strings.Index(exp, ASSIGN_OPER)
	if i == -1 {
		return "", exp, false
	}
	if i == strings.Index(exp, "==") {
		return "", exp, false
	}
	varname := exp[0:i]
	nexp := exp[i+1:]
	return varname, nexp, true
}

func ReplaceExpressionL(flowcontext IContext, elstr string) (interface{}, error) {
	si := strings.Index(elstr, EXP_LEFTDIV)
	if si == -1 {
		return elstr, nil
	}
	sj := strings.Index(elstr[si:], EXP_RIGHTDIV)
	if sj == -1 {
		return "", fmt.Errorf("EL表达式解析错误，表达式%s缺少%s", elstr, EXP_RIGHTDIV)
	}
	exp := elstr[si+2 : sj+si]
	v, err := DoExpression(exp, flowcontext)
	if err != nil {
		return nil, fmt.Errorf("EL表达式解析错误，表达式%s，%s", elstr, err.Error())
	}
	return v, nil
}

// 计算EL表达式，返回替换后的字符串
func ReplaceExpressionLStr(flowcontext IContext, elstr string) (string, error) {
	si := strings.Index(elstr, EXP_LEFTDIV)
	if si == -1 {
		return elstr, nil
	}
	sj := strings.Index(elstr[si:], EXP_RIGHTDIV)
	if sj == -1 {
		return "", fmt.Errorf("EL表达式解析错误，表达式%s缺少%s", elstr, EXP_RIGHTDIV)
	}
	exp := elstr[si+2 : sj+si]
	v, err := DoExpression(exp, flowcontext)
	if err != nil {
		return "", fmt.Errorf("EL表达式解析错误，表达式%s，%s", elstr, err.Error())
	}
	right, err := ReplaceExpressionLStr(flowcontext, elstr[si+1+sj:])
	if err != nil {
		return "", fmt.Errorf("EL表达式解析错误，表达式%s，%s", elstr, err.Error())
	}
	elstr = fmt.Sprint(elstr[:si], v, right)
	return elstr, nil
}
