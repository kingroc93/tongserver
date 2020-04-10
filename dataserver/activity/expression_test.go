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
	cnt.SetParam("z", 4)
	cnt.SetParam("q", 100)
	b, err := NewExpression("q+x+y+z").doExpression(cnt)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(b)
}
