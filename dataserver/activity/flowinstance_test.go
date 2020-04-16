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
	     "target": [{
			"expressions" : ["var_b=var_b+10"],
			"style" : "stdout",
			"flow" : [
			{
	           	"gate" : "to",
				"target": [{
					"style" : "stdout",
	           		"expressions" : ["var_a='next activity'"]
				}]
			}
			]
		}]
	}
	]}}`
	fl, err := NewFlowInstanceFromJSON(json)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := map[string]interface{}{
		"name": "menghui",
	}
	err = fl.Execute(&r)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestIfFlowToInstance(t *testing.T) {
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
	"gate":"ifto",
	"if":"name=='menghui'",
	"then":[
		{
			"style" : "stdout",
			"expressions" : ["var_a=\"he is tongtong's father!\""]	
		}
	],
	"else":[
		{
			"style" : "stdout",
			"expressions" : ["var_a=\"he is not tongtong's father!\""]	
		}
	]
}
]}}`
	fl, err := NewFlowInstanceFromJSON(json)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := map[string]interface{}{
		"name": "menghui2",
	}
	err = fl.Execute(&r)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}

func TestFlowLoopInstance(t *testing.T) {
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
	"flow": [{
		"gate":"loop",
		"assign":["a=1"],
		"step":["a=a+1"],
		"while": "a<10",
		"do":[{
			"style" : "stdout"
		}]
	}]
}}`
	fl, err := NewFlowInstanceFromJSON(json)
	if err != nil {
		fmt.Println(err)
		return
	}
	r := map[string]interface{}{
		"name": "menghui2",
	}
	err = fl.Execute(&r)
	if err != nil {
		fmt.Println(err)
		t.Fail()
	}
}
