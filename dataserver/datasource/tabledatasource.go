package datasource

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
	"tongserver.dataserver/utils"
)

// TableDataSource 数据库表数据源,封装单表操作
type TableDataSource struct {
	DBDataSource
	TableName string
}

// fillColumn 填充列信息
func (c *TableDataSource) fillColumn() error {
	k := c.DBAlias + "_" + c.TableName + "_Column"
	v := utils.DataSourceCache.Get(k)
	if v != nil {
		c.Field = v.([]*MyProperty)
		return nil
	}

	sqlb, err := CreateSQLBuileder(DBAlias2DBTypeContainer[c.DBAlias], c.TableName)
	if err != nil {
		return err
	}
	rs, err := c.querySQLData(sqlb.CreateGetColsSQL())
	if err != nil {
		return err
	}
	c.Field = make([]*MyProperty, len(rs.Data))
	for i, item := range rs.Data {
		c.Field[i] = &MyProperty{
			Name:     item[0].(string),
			DataType: ConvertMySQLType2CommonType(strings.ToUpper(item[1].(string))),
		}
	}
	return utils.DataSourceCache.Put(k, c.Field, 10*time.Minute)
}

// fillKeyFields 填充主键信息
func (c *TableDataSource) fillKeyFields() error {
	k := c.DBAlias + "_" + c.TableName + "_KeyField"
	v := utils.DataSourceCache.Get(k)
	if v != nil {
		c.KeyField = v.([]*MyProperty)
		return nil
	}
	sqlb, err := CreateSQLBuileder(DBAlias2DBTypeContainer[c.DBAlias], c.TableName)
	if err != nil {
		return err
	}

	rs, err := c.querySQLData(sqlb.CreateKeyFieldsSQL())
	if err != nil {
		return err
	}
	c.KeyField = make([]*MyProperty, len(rs.Data))
	for i, item := range rs.Data {
		c.KeyField[i] = &MyProperty{
			Name:     item[0].(string),
			DataType: ConvertMySQLType2CommonType(strings.ToUpper(item[1].(string))),
		}
	}

	return utils.DataSourceCache.Put(k, c.KeyField, 10*time.Minute)
}

// createSQLBuilder 创建SQL构造器
func (c *TableDataSource) createSQLBuilder() (ISQLBuilder, error) {
	return CreateSQLBuileder2(DBAlias2DBTypeContainer[c.DBAlias], c.TableName, c.convertPropertys2Cols(c.Field), c.orderlist, c.RowsLimit, c.RowsOffset)
}

// GetDataSourceType 返回数据源类型
func (c *TableDataSource) GetDataSourceType() DSType {
	return DataSourceTypeSqltable
}

// Init 初始化
func (c *TableDataSource) Init() error {
	if c.TableName == "" {
		return fmt.Errorf("tableName is nil")
	}
	if c.DBAlias == "" {
		return fmt.Errorf("Name is nil")
	}
	var err error
	c.openedDB, err = orm.GetDB(c.DBAlias)
	if err != nil {
		return err
	}
	c.palesql = true
	defer func() { c.palesql = false }()

	err = c.fillKeyFields()
	if err != nil {
		return err
	}
	if c.AutoFillFields {
		return c.fillColumn()
	}
	if len(c.Field) == 0 {
		logs.Warn("AutoFillFields属性为false并且Fields属性长度为0,查询语句会转换为*")
	}
	return nil
}

// QueryDataByFieldValues 根据字段值返回数据
func (c *TableDataSource) QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error) {
	c.ClearCriteria()
	for pname, value := range *fv {
		c.AndCriteria(pname, OperEq, value)
	}
	return c.DoFilter()
}

// QueryDataByKey 根据主键返回数据
func (c *TableDataSource) QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error) {
	if len(keyvalues) == 0 {
		return nil, fmt.Errorf("key values is none")
	}
	c.ClearCriteria()
	for i, v := range keyvalues {
		c.AndCriteria(c.KeyField[i].Name, OperEq, v)
	}

	return c.DoFilter()
}

// GetAllData 返回全部数据
func (c *TableDataSource) GetAllData() (*DataResultSet, error) {
	sqlstr, err := c.createSQLBuilder()
	if err != nil {
		return nil, err
	}

	sql, ps := sqlstr.CreateSelectSQL()
	return c.querySQLData(sql, ps...)
}

// DoFilter 根据查询条件返回数据
func (c *TableDataSource) DoFilter() (*DataResultSet, error) {
	sqlb, err := c.createSQLBuilder()
	if err != nil {
		return nil, err
	}
	sqlb.ClearCriteria()
	for _, item := range c.filter {
		sqlb.AddCriteria(item.PropertyName, item.Operation, item.Complex, item.Value)
	}
	for k, item := range c.aggre {
		sqlb.AddAggre(k, item)
	}
	sqlstr, param := sqlb.CreateSelectSQL()
	return c.querySQLData(sqlstr, param...)
}

//
//
//////////////////////////////////////////////////////////////////////////////////////////////////////////
