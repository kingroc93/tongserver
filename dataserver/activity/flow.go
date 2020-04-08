package activity

import (
	"github.com/astaxie/beego/logs"
	"strings"
)

const (
	F_TO     int = 0
	F_IFTO   int = 1
	F_IFLOOP int = 2
)

type IFlow interface {
	DoFlow(flowcontext IContext)
}

type Flow struct {
	gate      string
	define    map[string]interface{}
	activitys []IActivity
}
type FlowTo struct {
	Flow
}
type FlowIfTo struct {
	Flow
}
type FlowLoop struct {
	Flow
}

func (c *FlowTo) DoFlow(flowcontext IContext) {

}

func (c *FlowTo) FlowIfTo(flowcontext IContext) {

}

func (c *FlowTo) FlowLoop(flowcontext IContext) {

}
func NewFlow(d map[string]interface{}) IFlow {
	gate := ""
	for k, v := range d {
		if k == "gate" {
			gate = strings.ToLower(v.(string))
			if gate != "to" && gate != "ifto" && gate != "loop" {
				logs.Error("gate的值非法，只能是to,ifoto,loop")
				return nil
			}
			continue
		}
		//其他类型的节点为activity

	}

	if gate == "" {
		logs.Error("传入的定义没有gate属性")
		return nil
	}
	return nil
}
