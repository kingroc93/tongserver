package activity

import (
	"testing"
)

func TestStdOutActivity_Execute(t *testing.T) {
	def := make(map[string]interface{})
	def["style"] = "stdout"
	def["expressions"] = []string{"test=1+1", "v2=p1+p2"}

	cnt := NewContext()

	cnt.SetVarbiable("v1", "test var 1")
	cnt.SetVarbiable("v2", "test var 2")
	act := NewStdOutActivity(&def)
	act.Execute(cnt)
}
