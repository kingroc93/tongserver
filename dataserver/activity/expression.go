package activity

import (
	"github.com/antonmedv/expr"
	"tongserver.dataserver/utils"
)

// 表达式类
type Expression struct {
	define     interface{}
	expression string
}

func NewExpression(define interface{}) *Expression {
	m, ok := define.(string)
	if !ok {
		return nil
	}
	return &Expression{define: define, expression: m}
}

// 运行表达式 返回true or false
func (c *Expression) DoExpression(context IContext) (bool, error) {
	env := make(map[string]interface{})
	context.ForEachParams(func(name string, value interface{}) {
		env[name] = value
	})
	context.ForEachVariable(func(name string, value interface{}) {
		env[name] = value
	})
	program, err := expr.Compile(c.expression, expr.Env(env))
	if err != nil {
		return false, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return false, err
	}
	return utils.ConvertObj2Bool(output), nil
}
