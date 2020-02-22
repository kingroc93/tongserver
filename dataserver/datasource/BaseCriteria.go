package datasource

// BaseCriteria 基本查询条件
type BaseCriteria struct {
	filter    []*TDFilter
	orderlist []string
}

// Orderby 添加排序字段
func (c *BaseCriteria) Orderby(field string, dir string) IFilterAdder {
	c.orderlist = append(c.orderlist, field+" "+dir)
	return c
}

// addCriteria 内部调用的添加一个查询条件
func (c *BaseCriteria) addCriteria(propertyname, operation, complex string, value interface{}) IFilterAdder {
	c.filter = append(c.filter, &TDFilter{
		PropertyName: propertyname,
		Operation:    operation,
		Value:        value,
		Complex:      complex,
	})

	return c
}

// AddCriteria 添加一个查询条件
func (c *BaseCriteria) AddCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if c.filter == nil {
		c.filter = make([]*TDFilter, 0, 10)
	}
	return c.addCriteria(propertyname, operation, CompNone, value)
}

// ClearCriteria 清空查询条件
func (c *BaseCriteria) ClearCriteria() {
	c.filter = nil
}

// AndCriteria 添加一个与条件
func (c *BaseCriteria) AndCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if c.filter == nil {
		c.AddCriteria(propertyname, operation, value)
		return c
	}
	return c.addCriteria(propertyname, operation, CompAnd, value)
}

// OrCriteria 添加一个或条件
func (c *BaseCriteria) OrCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if c.filter == nil {
		c.AddCriteria(propertyname, operation, value)
		return c
	}
	return c.addCriteria(propertyname, operation, CompOr, value)
}

// TableDataSourceCriteria TableDataSource的条件,在基本条件上增加了聚合条件
type TableDataSourceCriteria struct {
	BaseCriteria
	//OutFields []string
	aggre map[string]*AggreType //key为聚合后返回的字段名
}

// AddAggre 添加一个聚合条件
func (c *TableDataSourceCriteria) AddAggre(outfield string, aggreType *AggreType) {
	if c.aggre == nil {
		c.aggre = make(map[string]*AggreType)
	}
	c.aggre[outfield] = aggreType
}
