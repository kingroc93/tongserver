package activity

import (
	"fmt"
	"reflect"
	"tongserver.dataserver/utils"
)

const (
	F_TO     int = 0
	F_IFTO   int = 1
	F_IFLOOP int = 2
)

type FlowResult uint8

const (
	FR_CONINUE FlowResult = 1
	FR_BREAK   FlowResult = 2
)

type IFlow interface {
	DoFlow(flowcontext IContext) FlowResult
}

type Flow struct {
	gate         int
	define       map[string]interface{}
	flowInstance FlowInstance
}
type FlowTo struct {
	Flow
	activitys *map[string]IActivity
}

type FlowIfTo struct {
	Flow
	expression  Expression
	thenAct     *map[string]IActivity
	elseThenAct *map[string]IActivity
}

type FlowLoop struct {
	Flow
}

func (c *FlowTo) DoFlow(flowcontext IContext) FlowResult {
	for _, v := range *c.activitys {
		v.Execute(flowcontext)
	}
	return FR_CONINUE
}

func (c *FlowIfTo) DoFlow(flowcontext IContext) FlowResult {
	return FR_CONINUE
}

func (c *FlowLoop) DoFlow(flowcontext IContext) FlowResult {
	return FR_CONINUE
}

func NewFlowTo(define *map[string]interface{}, flowInstance FlowInstance) (*FlowTo, error) {
	acts := make(map[string]IActivity)
	for k, v := range *define {
		if k == "gate" {
			//忽略gate属性
			continue
		}
		if _, ok := (*flowInstance.GlobaActivityContainer)[k]; ok {
			return nil, fmt.Errorf("Activity在flow实例中唯一")
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
	f := &FlowTo{
		Flow: Flow{
			gate: F_TO,
		},
		activitys: &acts,
	}
	return f, nil
}
func NewFlowIfTo(define map[string]interface{}, flowInstance FlowInstance) (*FlowIfTo, error) {
	f := &FlowIfTo{
		Flow: Flow{
			gate: F_IFTO,
		},
	}
	return f, nil
}
func NewFlowLoop(define map[string]interface{}, flowInstance FlowInstance) (*FlowLoop, error) {
	f := &FlowLoop{
		Flow: Flow{
			gate: F_IFLOOP,
		},
	}
	return f, nil
}

func NewFlow(d *map[string]interface{}, flowInstance FlowInstance) (IFlow, error) {
	gate, ok := (*d)["gate"]
	if !ok {
		return nil, fmt.Errorf("缺少gate属性")
	}
	sg, ok := gate.(string)
	if !ok {
		return nil, fmt.Errorf("gate属性类型必须是string")
	}

	if sg == "to" {
		return NewFlowTo(d, flowInstance)
	}
	if sg == "ifto" {
		return NewFlowIfTo(d, flowInstance)
	}
	if sg == "loop" {
		return NewFlowLoop(d, flowInstance)
	}
	return nil, nil
}
