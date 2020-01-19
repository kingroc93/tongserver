package datasource

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

var mu sync.Mutex

const (
	// 等于
	OperEq string = "="
	// 不等于
	OperNoteq string = "<>"
	// 大于
	OperGt string = ">"
	// 小于
	OperLt string = "<"
	// 大于等于
	OperGtEg string = ">="
	// 小于等于
	OperLtEg string = "<="
	// 介于--之间
	OperBetween string = "BETWEEN"
	// 包含
	OperIn string = "in"
)

const (
	// 与
	CompAnd string = "and"
	// 或
	CompOr string = "or"
	// 非
	CompNot string = "not"
	// 未知
	CompNone string = ""
)

// SQL查询条件
type SQLCriteria struct {
	PropertyName string
	Operation    string
	Value        interface{}
	Complex      string
}

// 聚合类型
type AggreType struct {
	//谓词
	Predicate int
	//字段名
	ColName string
}

const (
	// 计数
	AggCount int = 1
	// 求和
	AggSum int = 2
	// 求算术平均
	AggAvg int = 3
	// 最大值
	AggMax int = 4
	// 最小值
	AggMin int = 5
)

// SQL构造器接口
type ISQLBuilder interface {
	AddCriteria(field, operation, complex string, value interface{}) ISQLBuilder
	CreateSelectSQL() (string, []interface{})
	CreateInsertSQLByMap(fieldvalues map[string]interface{}) (string, []interface{})
	CreateDeleteSQL() (string, []interface{})
	CreateUpdateSQL(fieldvalues map[string]interface{}) (string, []interface{})
	CreateKeyFieldsSQL() string
	CreateGetColsSQL() string
	ClearCriteria()
	AddAggre(outfield string, aggreType *AggreType)
}

// SQL 构造器类
type SQLBuilder struct {
	//表名
	tableName string
	//抽象表，ObjectTable为一个SQL语句，返回一个结果集，这个结果集作为查询的表参与Select语句，相当于select * from (ObjectTable) as tableName
	objectTable string
	//字段名
	columns []string
	//排序字段
	orderBy    []string
	criteria   []SQLCriteria
	rowsLimit  int
	rowsOffset int
	aggre      map[string]*AggreType
}

// MySQl的SQL构造器
type MySQLSQLBuileder struct {
	SQLBuilder
}

// 创建SQL构造器
func CreateSQLBuileder2(dbType string, tablename string, columns []string, orderby []string, rowslimit int, rowsoffset int) (ISQLBuilder, error) {
	switch dbType {
	case DBType_MySQL:
		return &MySQLSQLBuileder{
			SQLBuilder: SQLBuilder{
				tableName:  tablename,
				columns:    columns,
				orderBy:    orderby,
				rowsLimit:  rowslimit,
				rowsOffset: rowsoffset}}, nil
	}
	return nil, fmt.Errorf("不支持的数据库类型" + dbType)
}

// 创建SQL构造器
func CreateSQLBuileder(dbType string, tablename string) (ISQLBuilder, error) {
	switch dbType {
	case DBType_MySQL:
		return &MySQLSQLBuileder{
			SQLBuilder: SQLBuilder{
				tableName: tablename}}, nil
	}
	return nil, fmt.Errorf("不支持的数据库类型" + dbType)
}

// 返回查询数据库表主键信息的SQL语句
func (c *MySQLSQLBuileder) CreateKeyFieldsSQL() string {
	if c.objectTable == "" {
		sqlstr := "SELECT a.column_name,b.data_type FROM INFORMATION_SCHEMA.`KEY_COLUMN_USAGE` a" +
			" inner join information_schema.columns b on a.table_name=b.table_name and a.column_name=b.column_name " +
			" WHERE a.table_name='" + c.tableName + "' AND a.constraint_name='PRIMARY'"
		return sqlstr
	} else {
		return ""
	}
}

// 返回获取数据库表全部字段的SQL语句
func (c *MySQLSQLBuileder) CreateGetColsSQL() string {
	if c.objectTable == "" {
		return "SELECT column_name,data_type FROM information_schema.columns WHERE table_name='" + c.tableName + "'"
	} else {
		return ""
	}
}

// 清楚查询条件
func (c *MySQLSQLBuileder) ClearCriteria() {
	c.criteria = nil
}

// 添加聚合
func (tc *MySQLSQLBuileder) AddAggre(outfield string, aggreType *AggreType) {
	if tc.aggre == nil {
		tc.aggre = make(map[string]*AggreType)
	}
	tc.aggre[outfield] = aggreType
}

// 删除条件
func (c *MySQLSQLBuileder) AddCriteria(field, operation, complex string, value interface{}) ISQLBuilder {
	mu.Lock()
	if c.criteria == nil {
		c.criteria = make([]SQLCriteria, 0, 10)
	}
	mu.Unlock()
	c.criteria = append(c.criteria, SQLCriteria{
		PropertyName: field,
		Operation:    operation,
		Value:        value,
		Complex:      complex,
	})
	return c
}

// 创建查询Where语句
func (c *MySQLSQLBuileder) createWhereSubStr() (string, []interface{}) {
	var sqlwhere string
	param := make([]interface{}, 0, len(c.criteria))
	for i, cr := range c.criteria {
		var exp string
		switch cr.Operation {
		case OperBetween:
			{
				switch reflect.TypeOf(cr.Value).Kind() {
				case reflect.Slice, reflect.Array:
					s := reflect.ValueOf(cr.Value)
					if s.Len() != 2 {
						panic("the BETWEEN operation in SQLBuilder the params must be array or slice, and length must be 2")
					}
					param = append(param, s.Index(0).Interface(), s.Index(1).Interface())
					exp = fmt.Sprint(c.tableName, ".", cr.PropertyName, " BETWEEN ? and ?")
				default:
					{
						panic("the BETWEEN operation in SQLBuilder the params must be array or slice, and length must be 2")
					}
				}
			}
		case OperIn:
			{
				switch reflect.TypeOf(cr.Value).Kind() {
				case reflect.Slice, reflect.Array:
					s := reflect.ValueOf(cr.Value)
					ins := ""
					for si := 0; si < s.Len(); si++ {
						ins = ins + "?,"
						param = append(param, s.Index(si).Interface())
					}
					ins = strings.TrimRight(ins, ",")
					exp = fmt.Sprint(c.tableName, ".", cr.PropertyName, " in (", ins, ")")
				default:
					{
						exp = fmt.Sprint(c.tableName, ".", cr.PropertyName, " in (?)")
						param = append(param, cr.Value)
					}
				}
			}
		default:
			{
				exp = fmt.Sprint(c.tableName, ".", cr.PropertyName, cr.Operation, "?")
				param = append(param, cr.Value)
			}
		}
		if i != 0 {
			if cr.Complex == CompAnd || cr.Complex == CompOr {
				sqlwhere = fmt.Sprint(sqlwhere, " ", cr.Complex, " ", exp)
			}
		} else {
			sqlwhere = fmt.Sprint(sqlwhere, " ", exp)
		}
	}
	//sql += " WHERE " + sqlwhere
	return " WHERE " + sqlwhere, param
}

// 创建删除数据的SQL语句
func (c *MySQLSQLBuileder) CreateDeleteSQL() (string, []interface{}) {
	sql := "DELETE FROM " + c.tableName
	if c.criteria != nil {
		where, ps := c.createWhereSubStr()
		sql += where
		return sql, ps
	} else {
		return sql, nil
	}
}

//
func (c *MySQLSQLBuileder) CreateUpdateSQL(fieldvalues map[string]interface{}) (string, []interface{}) {
	sql := "UPDATE " + c.tableName + " SET "
	params := make([]interface{}, len(fieldvalues), len(fieldvalues))
	i := 0
	for k, v := range fieldvalues {
		if i != 0 {
			sql += ","
		}
		sql += k + "=?"
		params[i] = v
		i++
	}
	if c.criteria != nil {
		where, ps := c.createWhereSubStr()
		sql += where
		params = append(params, ps...)
	}
	return sql, params
}

func (c *MySQLSQLBuileder) CreateInsertSQLByMap(fieldvalues map[string]interface{}) (string, []interface{}) {
	params := make([]interface{}, len(fieldvalues), len(fieldvalues))
	sql := "INSERT INTO " + c.tableName + " ("
	ps := ""
	i := 0
	for k, v := range fieldvalues {
		if i != 0 {
			sql += ","
			ps += ","
		}
		ps += "?"
		sql += k
		params[i] = v
		i++
	}
	sql = sql + ") VALUES (" + ps + ")"
	return sql, params
}
func (c *MySQLSQLBuileder) CreateSelectSQL() (string, []interface{}) {
	if c.objectTable != "" && len(c.criteria) == 0 && c.rowsLimit == 0 && c.rowsOffset == 0 && len(c.orderBy) == 0 && len(c.aggre) == 0 && (len(c.columns) == 0 || c.columns[0] == "*") {
		return c.objectTable, nil
	}
	var sql = "SELECT "
	var param []interface{}
	param = nil
	groupFields := make([]string, 0, 10)
	cols := c.columns
	if len(c.aggre) != 0 {
		//计算 group by子句中的字段列表
		if len(c.columns) != 0 {
			cols = make([]string, 0, 10)
			for _, col := range c.columns {
				if strings.Trim(col, " ") != "*" {
					cols = append(cols, c.tableName+"."+col)
					groupFields = append(groupFields, col)
				}
			}
		}
		//将聚合函数添加到选择字段列表
		for field, aggre := range c.aggre {
			var p string
			switch aggre.Predicate {
			case AggCount:
				p = "COUNT("
			case AggAvg:
				p = "AVG("
			case AggMax:
				p = "MAX("
			case AggMin:
				p = "MIN("
			case AggSum:
				p = "SUM("
			}
			p += c.tableName + "." + aggre.ColName + ") as " + field
			cols = append(cols, p)
		}
	}
	if len(cols) == 0 {
		//cols长度为0，选择*
		sql += c.tableName + ".* "
	} else {
		//生成选择的字段列表
		for i, fs := range cols {
			if i != 0 {
				sql += ","
			}
			sql += fs
		}
	}
	if c.objectTable == "" {
		sql += " FROM " + c.tableName
	} else {
		sql += " FROM (" + c.objectTable + ") as " + c.tableName
	}
	if c.criteria != nil {
		where, ps := c.createWhereSubStr()
		sql += where
		param = append(param, ps...)
	}

	if len(groupFields) != 0 {
		var grs string
		for index, gr := range groupFields {
			if index != 0 {
				grs = fmt.Sprint(",", grs)
			}
			grs = fmt.Sprint(grs, c.tableName+"."+gr)
		}
		sql += " GROUP BY " + grs
	}

	if len(c.orderBy) != 0 {
		sql += " ORDER BY "
		for i, o := range c.orderBy {
			if i != 0 {
				sql += ","
			}
			sql += o
		}
	}

	if c.rowsLimit != 0 {
		sql += " LIMIT " + strconv.Itoa(c.rowsOffset) + "," + strconv.Itoa(c.rowsLimit)
	}

	return sql, param
}

////
////
/////////////////////////////////////////////////////////////////////////////////////////////////////////////

type SQLInnerJoin struct {
}
