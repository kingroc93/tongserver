package activity

import (
	"github.com/antonmedv/expr"
	"tongserver.dataserver/utils"
)

// 表达式类
// 使用github.com/antonmedv/expr作为表达式处理引擎
type Expression struct {
	expression string
}

// 执行表达式
func (c *Expression) doExpression(context IContext) (interface{}, error) {
	env := make(map[string]interface{})
	context.ForEachParams(func(name string, value interface{}) {
		env[name] = value
	})
	context.ForEachVariable(func(name string, value interface{}) {
		env[name] = value
	})
	program, err := expr.Compile(c.expression, expr.Env(env))
	if err != nil {
		return nil, err
	}

	output, err := expr.Run(program, env)
	if err != nil {
		return nil, err
	}
	return output, nil
}

// 运行表达式 返回true or false
func (c *Expression) DoExpressionBool(context IContext) (bool, error) {
	output, err := c.doExpression(context)
	if err != nil {
		return false, err
	}
	return utils.ConvertObj2Bool(output), nil
}

// 赋值表达式
type AssignExpress struct {
	Expression
	AssignName string
}

// 执行赋值表达式
func (c *AssignExpress) DoExpression(context IContext) (interface{}, error) {
	out, err := c.doExpression(context)
	if err != nil {
		return nil, err
	}
	context.SetVarbiable(c.AssignName, out)
	return out, err
}

// 创建一个赋值表达式
// 为了方便从定义信息中获取数据，这里的参数为interface{}类型
// 实际define为map[string]string类型，包含assign属性和exp属性
// exp属性是string类型定义一个表达式
// assign属性是string类型，定义表达式运行成功后将值写入的变量名。
// 该变量名可以不存在。
func NewAssignExpress(define interface{}) *AssignExpress {
	n, ok := define.(map[string]string)
	if ok {
		return &AssignExpress{
			Expression: Expression{
				expression: n["exp"],
			},
			AssignName: n["assign"]}
	}
	return nil
}

// 创建一个表达式，为了方便从定义信息中获取数据，这里的参数为interface{}类型，实际必须是string类型
func NewExpression(define interface{}) *Expression {
	m, ok := define.(string)
	if ok {
		return &Expression{expression: m}
	}
	return nil
}
