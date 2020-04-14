package activity

import (
	"github.com/antonmedv/expr"
	"strings"
	"tongserver.dataserver/utils"
)

const (
	ASSIGN_OPER string = "="
)

func DoExpression(expression string, context IContext) (interface{}, error) {
	return DoExpression2(expression, *context.getVariableMap())
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

func DoExpressionBool(expression string, context IContext) (bool, error) {
	return DoExpressionBool2(expression, *context.getVariableMap())
}

// 运行表达式 返回true or false
func DoExpressionBool2(expression string, env map[string]interface{}) (bool, error) {
	output, err := DoExpression2(expression, env)
	if err != nil {
		return false, err
	}
	return utils.ConvertObj2Bool(output), nil
}

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
