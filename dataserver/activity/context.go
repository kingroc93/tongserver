package activity

import (
	"fmt"
	"strings"
	"tongserver.dataserver/datasource"
)

const (
	ACT_INNERSRV int = 1
	ACT_MESSAGE  int = 2
	ACT_STDOUT   int = 3
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type ForEachFun func(name string, value interface{})

type IContext interface {
	GetVarbiableByName(name string) interface{}
	SetVarbiable(name string, value interface{})
	CreateVarbiable(name string, def map[string]interface{}) error
	GetVarbiableNames() []string
	ForEachVariable(fun ForEachFun)
	getVariableMap() map[string]interface{}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
type Context struct {
	varbiable      map[string]interface{}
	varbiableTypes map[string]string
}

// getVariableMap
func (c *Context) getVariableMap() map[string]interface{} {
	return c.varbiable
}

func (c *Context) CreateVarbiable2(name string, ty string, value interface{}) error {
	_, ok := c.varbiable[name]
	if ok {
		return fmt.Errorf("变量名" + name + "已经存在")
	}
	if !datasource.ValidPropertyType(ty) {
		return fmt.Errorf("不支持的数据类型" + ty)
	}
	c.varbiable[name] = value
	c.varbiableTypes[name] = ty
	return nil
}

// CreateVarbiable
func (c *Context) CreateVarbiable(name string, def map[string]interface{}) error {
	_, ok := c.varbiable[name]
	if ok {
		return fmt.Errorf("name属性定义的变量名在当前上下文中已经存在")
	}
	t, ok := def["type"]
	ty := datasource.PropertyDatatypeStr
	if ok {
		ty = strings.ToUpper(t.(string))
	}
	v, ok := def["value"]
	if !ok {
		c.varbiable[name] = datasource.DefaultTypeValue(ty)
	} else {
		c.varbiable[name] = v
	}
	c.varbiableTypes[name] = ty
	return nil
}

// GetVarbiableByName
func (c *Context) GetVarbiableByName(name string) interface{} {
	if c.varbiable == nil {
		return nil
	}
	return c.varbiable[name]
}

// SetVarbiable
func (c *Context) SetVarbiable(name string, value interface{}) {
	c.varbiable[name] = value
}

// GetVarbiableNames
func (c *Context) GetVarbiableNames() []string {
	keys := make([]string, 0, len(c.varbiable))
	for k := range c.varbiable {
		keys = append(keys, k)
	}
	return keys
}

// ForEachVariable
func (c *Context) ForEachVariable(fun ForEachFun) {
	for k, v := range c.varbiable {
		fun(k, v)
	}
}

// NewContext
func NewContext() *Context {
	v := make(map[string]interface{})
	return &Context{varbiable: v}
}
