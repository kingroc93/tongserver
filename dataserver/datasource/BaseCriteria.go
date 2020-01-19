package datasource

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 基本查询条件
type BaseCriteria struct {
	filter    []*TDFilter
	orderlist []string
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加排序字段
func (this *BaseCriteria) Orderby(field string, dir string) IFilterAdder {
	this.orderlist = append(this.orderlist, field+" "+dir)
	return this
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 内部调用的添加一个查询条件
func (this *BaseCriteria) addCriteria(propertyname, operation, complex string, value interface{}) IFilterAdder {
	this.filter = append(this.filter, &TDFilter{
		PropertyName: propertyname,
		Operation:    operation,
		Value:        value,
		Complex:      complex,
	})
	return this
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加一个查询条件
func (this *BaseCriteria) AddCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if this.filter == nil {
		this.filter = make([]*TDFilter, 0, 10)
	}
	return this.addCriteria(propertyname, operation, CompNone, value)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 清空查询条件
func (this *BaseCriteria) ClearCriteria() {
	this.filter = nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加一个与条件
func (this *BaseCriteria) AndCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if this.filter == nil {
		this.AddCriteria(propertyname, operation, value)
		return this
	}
	return this.addCriteria(propertyname, operation, CompAnd, value)
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加一个或条件
func (this *BaseCriteria) OrCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if this.filter == nil {
		this.AddCriteria(propertyname, operation, value)
		return this
	}
	return this.addCriteria(propertyname, operation, CompOr, value)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// TableDataSource的条件,在基本条件上增加了聚合条件
type TableDataSourceCriteria struct {
	BaseCriteria
	//OutFields []string
	aggre map[string]*AggreType //key为聚合后返回的字段名
}

// 添加一个聚合条件
func (tc *TableDataSourceCriteria) AddAggre(outfield string, aggreType *AggreType) {
	if tc.aggre == nil {
		tc.aggre = make(map[string]*AggreType)
	}
	tc.aggre[outfield] = aggreType
}
