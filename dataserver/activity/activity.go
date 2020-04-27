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
	GetDef() map[string]interface{}
}

// 基础活动
// 基础活动包括一个类型属性、唯一活动名称属性，还可以附加一组表达式
type Activity struct {
	Flows
	Style  string
	Exp    []string
	define map[string]interface{}
}

func (c *Activity) ExecuteExp(flowcontext IContext) error {
	return ExecuteExpressions(flowcontext, c.Exp)
}

func (c *Activity) GetDef() map[string]interface{} {
	return c.define
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
func NewStdOutActivity(acti *Activity) (IActivity, error) {
	act := &StdOutActivity{Activity: *acti}
	return act, nil
}

// 创建一组活动
func CreateActivitys(define []interface{}, flowInstance *FlowInstance) ([]*IActivity, error) {
	acts := make([]*IActivity, 0, len(define))
	for _, v := range define {
		if reflect.TypeOf(v).Kind() == reflect.String {
			//值是字符串时表示直接引用之前声明过得activity
			act, ok := flowInstance.GlobaActivityContainer[v.(string)]
			if !ok {
				return nil, fmt.Errorf("没有找到已经声明的activity，%s", v)
			}
			acts = append(acts, act)
		}
		if reflect.TypeOf(v).Kind() == reflect.Map {
			defm := v.(map[string]interface{})
			actname, hasname := defm["name"]
			if hasname {
				_, ok := flowInstance.GlobaActivityContainer[actname.(string)]
				if ok {
					return nil, fmt.Errorf("创建flow失败，Actitity名称不唯一")
				}

			}

			m := utils.ConvertObj2Map(v)
			if m == nil {
				return nil, fmt.Errorf("创建flow失败，toflow中除了gate属性外其他属性应该为对象")
			}
			var actp IActivity = nil
			if hasname {
				flowInstance.GlobaActivityContainer[actname.(string)] = &actp
			}
			act, err := CreateActivity(m, flowInstance)
			if err != nil {
				if hasname {
					delete(flowInstance.GlobaActivityContainer, actname.(string))
				}
				return nil, fmt.Errorf("创建flow失败，%s", err)
			}
			acts = append(acts, &act)
			if hasname {
				*(flowInstance.GlobaActivityContainer)[actname.(string)] = act
			}

		}
	}
	return acts, nil
}

// 创建活动
func CreateActivity(def map[string]interface{}, inst *FlowInstance) (IActivity, error) {
	style, ok := (def)["style"]
	if !ok {
		return nil, fmt.Errorf("创建Activiti失败，缺少style属性")
	}

	sStyle := strings.ToLower(style.(string))

	acti := &Activity{define: def}
	acti.Style = sStyle
	////////////////////////////////////////////////////
	m := (def)["expressions"]
	if m != nil {
		dd := m.([]interface{})
		acti.Exp = make([]string, len(dd), len(dd))
		for index, d := range dd {
			acti.Exp[index] = d.(string)
		}
	}
	////////////////////////////////////////////////////
	flows := utils.GetArrayFromMap(def, "flow")
	fs, err := CreateFlows(flows, inst)
	if err != nil {
		return nil, fmt.Errorf("创建Activity失败，%s", err.Error())
	}
	acti.flows = fs

	f, ok := acitvityCreatorFunContainer[sStyle]
	if !ok {
		return nil, fmt.Errorf("没有找到style属性为%s的构造器", sStyle)
	}
	return f(acti)
}

// StdOutActivity 执行方法，目前是将所有变量输出到控制台，其实没什么用
func (c *StdOutActivity) Execute(flowcontext IContext) error {
	err := c.ExecuteExp(flowcontext)
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

// 用于创建flow的构造器
type AcitvityCreatorFun func(acti *Activity) (IActivity, error)

// flow的构造器的容器
var acitvityCreatorFunContainer = make(map[string]AcitvityCreatorFun)

func RegisterAcitvityCreator(gateName string, f AcitvityCreatorFun) {
	acitvityCreatorFunContainer[gateName] = f
}
