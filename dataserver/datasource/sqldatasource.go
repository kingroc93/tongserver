package datasource

//
//import (
//	"fmt"
//	"github.com/astaxie/beego/orm"
//)
//
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
//type SQLDataSource struct {
//	DBDataSource
//	SQL          string
//	Params       []MyProperty
//	ParamsValues []interface{}
//}
//
///////////////////////////////////////////////////////////////////////////////////////////////////////////////
//// SQLDataSource
////返回数据源类型
//func (c *SQLDataSource) GetDataSourceType() DataSourceType {
//	return DataSourceType_SQL
//}
//
////数据源初始化
//func (c *SQLDataSource) Init() error {
//	if c.SQL == "" {
//		return fmt.Errorf("SQL is nil")
//	}
//	if c.DBAlias == "" {
//		return fmt.Errorf("Name is nil")
//	}
//	var err error
//	c.openedDB, err = orm.GetDB(c.DBAlias)
//	if err != nil {
//		return err
//	}
//	if c.AutoFillFields {
//		return c.fillFields()
//	}
//
//	return nil
//}
//func (c *SQLDataSource) fillFields() error {
//	sqlb := &MySQLSQLBuileder{
//		ObjectTable: c.SQL,
//		tableName:   c.Name,
//		Columns:     []string{"*"},
//		RowsLimit:   1,
//		RowsOffset:  0,
//	}
//	sqlstr, _ := sqlb.CreateSelectSQL()
//	rs, err := c.querySQLData(sqlstr, c.ParamsValues...)
//	if err != nil {
//		return err
//	}
//	fds := rs.Fields
//	c.Field = make([]*MyProperty, len(fds))
//	for k, v := range fds {
//		c.Field[v.Index] = &MyProperty{
//			Name:     k,
//			DataType: ConvertMySQLType2CommonType(v.FieldType),
//		}
//	}
//	return nil
//}
//func (c *SQLDataSource) GetName() string {
//	return c.Name
//}
//
//func (c *SQLDataSource) QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error) {
//	c.ClearCriteria()
//	for pname, value := range *fv {
//		c.AndCriteria(pname, OperEq, value)
//	}
//	return c.DoFilter()
//}
//
//func (c *SQLDataSource) QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error) {
//	if len(keyvalues) == 0 {
//		return nil, fmt.Errorf("key values is none!")
//	}
//	c.ClearCriteria()
//	for i, v := range keyvalues {
//		c.AndCriteria(c.KeyField[i].Name, OperEq, v)
//	}
//
//	return c.DoFilter()
//}
//
//func (c *SQLDataSource) createSQLBuilder() *MySQLSQLBuileder {
//	sqlb := &MySQLSQLBuileder{
//		ObjectTable: c.SQL,
//		tableName:   c.Name,
//		Columns:     c.convertPropertys2Cols(c.Field),
//		OrderBy:     c.orderlist,
//		RowsLimit:   c.RowsLimit,
//		RowsOffset:  c.RowsOffset,
//	}
//	return sqlb
//}
//
////返回全部数据
//func (c *SQLDataSource) GetAllData() (*DataResultSet, error) {
//	sqlstr, _ := c.createSQLBuilder().CreateSelectSQL()
//	return c.querySQLData(sqlstr, c.ParamsValues...)
//}
//
//func (c *SQLDataSource) DoFilter() (*DataResultSet, error) {
//	sqlb := c.createSQLBuilder()
//	sqlb.ClearCriteria()
//	for _, item := range c.filter {
//		sqlb.AddCriteria(item.PropertyName, item.Operation, item.Complex, item.Value)
//	}
//	for k, item := range c.aggre {
//		sqlb.AddAggre(k, item)
//	}
//	sqlstr, param := sqlb.CreateSelectSQL()
//	p := append(c.ParamsValues, param...)
//	return c.querySQLData(sqlstr, p...)
//}
//
