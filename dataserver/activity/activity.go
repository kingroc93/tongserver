package activity

import (
	"fmt"
	"reflect"
	"strings"
	"tongserver.dataserver/utils"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type IActivity interface {
	Execute(flowcontext IContext) error
}

// 基础活动
// 基础活动包括一个类型属性、唯一活动名称属性，还可以附加一组表达式
type Activity struct {
	Flows
	Style  int
	Name   string
	Exp    []string
	define *map[string]interface{}
}

func (c *Activity) executeExp(flowcontext IContext) error {
	env := *flowcontext.getVariableMap()
	vmap := make(map[string]interface{})
	for _, exp := range c.Exp {
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

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 脚本活动
type ScriptActivity struct {
	Activity
	Script string
}

func (c *ScriptActivity) Execute(flowcontext IContext) {

}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 控制台输出活动
type StdOutActivity struct {
	Activity
}

// 创建一个控制台输出活动
func NewStdOutActivity(def *map[string]interface{}) *StdOutActivity {
	act := &StdOutActivity{Activity{define: def}}
	m := (*def)["expressions"]
	if m != nil {
		dd := m.([]interface{})
		act.Exp = make([]string, len(dd), len(dd))
		for index, d := range dd {
			act.Exp[index] = d.(string)
		}
	}

	return act
}

// 创建一组活动
func CreateActivitys(define *map[string]interface{}, flowInstance *FlowInstance, igname []string) (*map[string]IActivity, error) {
	acts := make(map[string]IActivity)
actLoop:
	for k, v := range *define {
		if igname != nil {
			for _, ig := range igname {
				if k == ig {
					continue actLoop
				}
			}
		}
		if reflect.TypeOf(v).Kind() == reflect.String {
			//值是字符串时表示直接引用之前声明过得activity
			act, ok := (*flowInstance.GlobaActivityContainer)[k]
			if !ok {
				return nil, fmt.Errorf("没有找到已经声明的activity，%s", k)
			}
			acts[k] = act
		} else {
			m := utils.ConvertObj2Map(v)
			if m == nil {
				return nil, fmt.Errorf("创建flow失败，toflow中除了gate属性外其他属性应该为对象")
			}
			act, err := CreateActivity(m)
			if err != nil {
				return nil, fmt.Errorf("创建flow失败，%s", err)
			}
			acts[k] = act
			(*flowInstance.GlobaActivityContainer)[k] = act
		}
	}
	return &acts, nil
}

// 创建活动
func CreateActivity(def *map[string]interface{}) (IActivity, error) {
	style, ok := (*def)["style"]
	if !ok {
		return nil, fmt.Errorf("创建Activiti失败，缺少style属性")
	}

	sStyle := strings.ToLower(style.(string))

	if style == "normal" {

	}
	if sStyle == "stdout" {
		return NewStdOutActivity(def), nil
	}
	if sStyle == "innerservice" {
		return nil, nil
	}
	if sStyle == "message" {
		return nil, nil
	}
	if sStyle == "process" {
		return nil, nil
	}
	return nil, nil
}

// StdOutActivity 执行方法，目前是将所有变量输出到控制台，其实没什么用
func (c *StdOutActivity) Execute(flowcontext IContext) error {
	err := c.executeExp(flowcontext)
	if err != nil {
		return err
	}
	flowcontext.ForEachVariable(func(name string, value interface{}) {
		fmt.Printf("%s : %s \n", name, value)
	})
	return c.ExecuteFlows(flowcontext)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
