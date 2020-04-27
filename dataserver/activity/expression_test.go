package activity

import (
	"fmt"
	"github.com/antonmedv/expr"
	"testing"
)

func TestExpression(t *testing.T) {
	env := map[string]interface{}{
		"greet":   "Hello, %v!",
		"names":   []string{"world", "you"},
		"sprintf": fmt.Sprintf,
	}

	code := `sprintf(greet, names[0])`

	program, err := expr.Compile(code, expr.Env(env))
	if err != nil {
		panic(err)
	}

	output, err := expr.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(output)
}

func TestExpression_DoExpression(t *testing.T) {
	cnt := NewContext()
	cnt.SetVarbiable("x", 1)
	cnt.SetVarbiable("y", 2)

	b, err := DoExpression("q+x+y+z", cnt)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(b)
}

// 测试EL表达式
func TestReplaceExpressionL(t *testing.T) {
	cnt := NewContext()
	cnt.SetVarbiable("x", 1)
	cnt.SetVarbiable("y", 2)
	s, err := ReplaceExpressionLStr(cnt, "x=${x+y},${y-x}")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(s)
}

func TestSplitAssignExpression(t *testing.T) {
	fmt.Println(SplitAssignExpression("a =1+1"))
	fmt.Println(SplitAssignExpression(" a = 1+1"))
	fmt.Println(SplitAssignExpression("a ==1+1"))
}
