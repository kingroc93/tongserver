package activity

import (
	"fmt"
	"tongserver.dataserver/utils"
)

const (
	F_TO     int = 0
	F_IFTO   int = 1
	F_IFLOOP int = 2
)

type FlowResult uint8

const (
	FR_ERROR   FlowResult = 0
	FR_CONINUE FlowResult = 1
	FR_BREAK   FlowResult = 2
)

type IFlow interface {
	DoFlow(flowcontext IContext) (FlowResult, error)
}

type Flows struct {
	flows []IFlow
}

func (c *Flows) ExecuteFlows(flowcontext IContext) error {
	for _, item := range c.flows {
		rs, err := item.DoFlow(flowcontext)
		if err != nil {
			return err
		}
		if rs == FR_BREAK {
			break
		}
	}
	return nil
}

type Flow struct {
	gate         int
	define       map[string]interface{}
	flowInstance FlowInstance
}
type FlowTo struct {
	Flow
	activitys *map[string]*IActivity
}

type FlowIfTo struct {
	Flow
	expression  string
	thenAct     *map[string]*IActivity
	elseThenAct *map[string]*IActivity
}

type FlowLoop struct {
	Flow
	assignExpression string
	whileExpression  string
}

func (c *Flow) executeActivitys(activitys *map[string]*IActivity, flowcontext IContext) (FlowResult, error) {
	for _, v := range *activitys {
		err := (*v).Execute(flowcontext)
		if err != nil {
			return FR_BREAK, err
		}
	}
	return FR_CONINUE, nil
}

func (c *FlowTo) DoFlow(flowcontext IContext) (FlowResult, error) {
	return c.executeActivitys(c.activitys, flowcontext)
}

func (c *FlowIfTo) DoFlow(flowcontext IContext) (FlowResult, error) {
	r, err := DoExpressionBool(c.expression, flowcontext)
	if err != nil {
		return FR_ERROR, fmt.Errorf("执行表达式 %s 发生错误,%s", c.expression, err.Error())
	}
	if r {
		return c.executeActivitys(c.thenAct, flowcontext)
	} else {
		if c.elseThenAct != nil {
			return c.executeActivitys(c.elseThenAct, flowcontext)
		}
	}
	return FR_CONINUE, nil

}

func (c *FlowLoop) DoFlow(flowcontext IContext) (FlowResult, error) {
	return FR_CONINUE, nil
}

// {
//        "gate": "to",
//        "activity1": {
//			"expressions":["var_b=var_b+10"],
//			"style":"stdout"
//		}
//	}
func NewFlowTo(define *map[string]interface{}, flowInstance *FlowInstance) (*FlowTo, error) {
	acts, err := CreateActivitys(define, flowInstance, []string{"gate"})
	if err != nil {
		return nil, err
	}
	f := &FlowTo{
		Flow: Flow{
			gate: F_TO,
		},
		activitys: acts,
	}
	return f, nil
}

//{
//    "gate":"ifto",
//    "if":"表达式",
//    "then":{},
//    "else":{}
//}
func NewFlowIfTo(define *map[string]interface{}, flowInstance *FlowInstance) (*FlowIfTo, error) {
	f := &FlowIfTo{
		Flow: Flow{
			gate: F_IFTO,
		},
	}
	exp, ok := (*define)["if"]
	if !ok {
		return nil, fmt.Errorf("创建ifto失败，没有if属性")
	}
	f.expression = exp.(string)
	then := utils.GetMapFromMap(define, "then")
	if then == nil {
		return nil, fmt.Errorf("创建ifto失败，没有then属性")
	}
	am, err := CreateActivitys(then, flowInstance, nil)
	if err != nil {
		return nil, fmt.Errorf("创建ifto失败，then创建失败，%s", err.Error())
	}
	f.thenAct = am
	els := utils.GetMapFromMap(define, "else")
	if els != nil {
		f.elseThenAct, err = CreateActivitys(els, flowInstance, nil)
		if err != nil {
			return nil, fmt.Errorf("创建ifto失败，else 创建失败，%s", err.Error())
		}
	}
	return f, nil
}
func NewFlowLoop(define *map[string]interface{}, flowInstance *FlowInstance) (*FlowLoop, error) {
	f := &FlowLoop{
		Flow: Flow{
			gate: F_IFLOOP,
		},
	}
	return f, nil
}

func NewFlow(d *map[string]interface{}, flowInstance *FlowInstance) (IFlow, error) {
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
func CreateFlows(flows []interface{}, inst *FlowInstance) ([]IFlow, error) {
	if len(flows) > 0 {
		iflows := make([]IFlow, len(flows), len(flows))
		for index, item := range flows {
			inf := utils.ConvertObj2Map(item)
			if inf == nil {
				return nil, fmt.Errorf("创建流程实例失败，start节点的flow属性必须为一个对象数组")
			}
			f, err := NewFlow(inf, inst)
			if err != nil {
				return nil, err
			}
			iflows[index] = f
		}
		return iflows, nil
	}
	return nil, nil
}
