package datasource

type BaseCriteria struct {
	filter    []*TDFilter
	orderlist []string
}

func (this *BaseCriteria) Orderby(field string, dir string) IFilterAdder {
	this.orderlist = append(this.orderlist, field+" "+dir)
	return this
}

func (this *BaseCriteria) addCriteria(propertyname, operation, complex string, value interface{}) IFilterAdder {
	this.filter = append(this.filter, &TDFilter{
		PropertyName: propertyname,
		Operation:    operation,
		Value:        value,
		Complex:      complex,
	})
	return this
}

func (this *BaseCriteria) AddCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if this.filter == nil {
		this.filter = make([]*TDFilter, 0, 10)
	}
	return this.addCriteria(propertyname, operation, COMP_NONE, value)
}

func (this *BaseCriteria) ClearCriteria() {
	this.filter = nil
}

func (this *BaseCriteria) AndCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if this.filter == nil {
		this.AddCriteria(propertyname, operation, value)
		return this
	}
	return this.addCriteria(propertyname, operation, COMP_AND, value)
}

func (this *BaseCriteria) OrCriteria(propertyname, operation string, value interface{}) IFilterAdder {
	if this.filter == nil {
		this.AddCriteria(propertyname, operation, value)
		return this
	}
	return this.addCriteria(propertyname, operation, COMP_OR, value)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

type TableDataSourceCriteria struct {
	BaseCriteria
	//OutFields []string
	aggre map[string]*AggreType //key为聚合后返回的字段名
}

func (tc *TableDataSourceCriteria) AddAggre(outfield string, aggreType *AggreType) {
	if tc.aggre == nil {
		tc.aggre = make(map[string]*AggreType)
	}
	tc.aggre[outfield] = aggreType
}
