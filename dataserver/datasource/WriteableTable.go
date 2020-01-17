package datasource

type WriteableTableSource struct {
	TableDataSource
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *WriteableTableSource) Delete() error {
	sqlb, err := c.createSQLBuilder()
	if err != nil {
		return err
	}
	sqlb.ClearCriteria()
	for _, item := range c.filter {
		sqlb.AddCriteria(item.PropertyName, item.Operation, item.Complex, item.Value)
	}
	sql, p := sqlb.CreateDeleteSQL()
	_, err2 := c.openedDB.Exec(sql, p...)
	return err2
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *WriteableTableSource) Insert(values map[string]interface{}) error {
	sqlb, err := CreateSQLBuileder(DBAlias2DBTypeContainer[c.DBAlias], c.TableName)
	if err != nil {
		return err
	}
	sql, ps := sqlb.CreateInsertSQLByMap(values)
	_, err2 := c.openedDB.Exec(sql, ps...)
	return err2
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
func (c *WriteableTableSource) Update(values map[string]interface{}) error {
	sqlb, err := c.createSQLBuilder()
	if err != nil {
		return err
	}
	sqlb.ClearCriteria()
	for _, item := range c.filter {
		sqlb.AddCriteria(item.PropertyName, item.Operation, item.Complex, item.Value)
	}
	sql, ps := sqlb.CreateUpdateSQL(values)
	_, err2 := c.openedDB.Exec(sql, ps...)
	return err2
}
