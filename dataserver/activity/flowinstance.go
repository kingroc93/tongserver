package activity

import (
	"fmt"
	"tongserver.dataserver/utils"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FlowInstance struct {
	Context
	id                     string
	name                   string
	define                 *map[string]interface{}
	flows                  []IFlow
	GlobaActivityContainer *map[string]IActivity
}

func (c *FlowInstance) Execute(params *map[string]interface{}) error {
	if len(*params) != len(*(c.params)) {
		return fmt.Errorf("参数错误")
	}
	for k, v := range *params {
		(*c.params)[k] = v
	}
	for _, item := range c.flows {
		item.DoFlow(c)
	}
	return nil
}

func NewFlowInstance(define *map[string]interface{}) (*FlowInstance, error) {
	n, ok := (*define)["name"]
	if !ok {
		return nil, fmt.Errorf("创建流程实例失败，没有name属性")
	}

	//开始节点
	start := utils.GetMapFromMap(define, "start")
	if start == nil {
		return nil, fmt.Errorf("创建流程实例失败，没有start属性应该是一个对象")
	}
	vs := utils.GetMapFromMap(start, "variables")
	if vs == nil {
		return nil, fmt.Errorf("创建流程实例失败，没有variables属性")
	}
	ps := utils.GetMapFromMap(start, "param")
	if ps == nil {
		return nil, fmt.Errorf("创建流程实例失败，没有param属性")
	}

	flows := utils.GetArrayFromMap(start, "flow")
	if flows == nil {
		return nil, fmt.Errorf("创建流程实例失败，start节点的flow属性必须为一个数组")
	}

	gs := make(map[string]IActivity)

	inst := &FlowInstance{
		Context: Context{
			params:    ps,
			varbiable: vs,
		},
		name:                   n.(string),
		define:                 define,
		flows:                  make([]IFlow, len(flows), len(flows)),
		GlobaActivityContainer: &gs,
	}

	for index, item := range flows {
		inf := utils.ConvertObj2Map(item)
		if inf == nil {
			return nil, fmt.Errorf("创建流程实例失败，start节点的flow属性必须为一个对象数组")
		}
		f, err := NewFlow(inf, inst)
		if err == nil {
			return nil, err
		}
		inst.flows[index] = f
	}
	return inst, nil
}
