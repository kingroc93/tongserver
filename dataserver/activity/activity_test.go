package activity

import (
	"testing"
)

func TestStdOutActivity_Execute(t *testing.T) {
	cnt := NewContext()
	cnt.SetParam("p1", "test param1")
	cnt.SetParam("p2", "test param2")
	cnt.SetVarbiable("v1", "test var 1")
	cnt.SetVarbiable("v2", "test var 2")
	act := NewStdOutActivity(nil)
	act.Execute(cnt)
}
