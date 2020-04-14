package activity

import (
	"fmt"
	"testing"
)

func TestNewFlowToInstance(t *testing.T) {
	json := `{
  "name": "测试to flow",
  "start": {
	"params":{
		"name":{"type":"string","value":"menghui"},
		"age":{"type":"number","value":41}
	},
    "variables": {
      "var_a": {
        "type": "string",
        "value": "test var"
      },
      "var_b": {
        "type": "number",
        "value": 12
      }
    },
  	"flow": [
	{
        "gate": "to",
        "activity1": {
			"expressions":["var_b=var_b+10"],
			"style":"stdout"
		}
	}
	]}}`
	fl, err := NewFlowInstanceFromJSON(json)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := map[string]interface{}{
		"name": "lvxing",
	}
	err = fl.Execute(&r)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
