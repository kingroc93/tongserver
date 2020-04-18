package activity

import (
	"fmt"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/utils"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type FlowInstance struct {
	Context
	Flows
	id                     string
	name                   string
	define                 *map[string]interface{}
	GlobaActivityContainer *map[string]*IActivity
}

// 执行流程
func (c *FlowInstance) Execute(params *map[string]interface{}) error {
	if params != nil {
		for k, v := range *params {
			_, ok := (*c.varbiable)[k]
			if ok {
				(*c.varbiable)[k] = v
			} else {
				err := c.CreateVarbiable2(k, datasource.RelectType2InnerType(v), v)
				if err != nil {
					return err
				}
			}
		}
	}
	return c.ExecuteFlows(c)
}

// 根据JSON创建流程
func NewFlowInstanceFromJSON(json string) (*FlowInstance, error) {
	ma, err := utils.ParseJSONStr2Map(json)
	if err != nil {
		return nil, err
	}
	return NewFlowInstance(ma)
}

// 根据map创建流程
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
	ps := utils.GetMapFromMap(start, "params")
	flows := utils.GetArrayFromMap(start, "flow")
	if flows == nil {
		return nil, fmt.Errorf("创建流程实例失败，start节点的flow属性必须为一个数组")
	}
	gs := make(map[string]*IActivity)
	t := make(map[string]interface{})
	ty := make(map[string]string)
	inst := &FlowInstance{
		Context: Context{
			varbiable:      &t,
			varbiableTypes: &ty,
		},

		name:   n.(string),
		define: define,

		GlobaActivityContainer: &gs,
	}
	if vs != nil {
		for k, v := range *vs {
			err := inst.CreateVarbiable(k, utils.ConvertObj2Map(v))
			if err != nil {
				return nil, err
			}
		}
	}
	if ps != nil {
		for k, v := range *ps {
			err := inst.CreateVarbiable(k, utils.ConvertObj2Map(v))
			if err != nil {
				return nil, err
			}
		}
	}

	fs, err := CreateFlows(flows, inst)
	if err != nil {
		return nil, fmt.Errorf("创建NewFlowInstance实例失败，创建flows失败，%s", err.Error())
	}
	inst.flows = fs
	return inst, nil
}
