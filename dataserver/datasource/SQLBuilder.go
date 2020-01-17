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
	OPER_EQ      string = "="
	OPER_NOTEQ   string = "<>"
	OPER_GT      string = ">"
	OPER_LT      string = "<"
	OPER_GT_EG   string = ">="
	OPER_LT_EG   string = "<="
	OPER_BETWEEN string = "BETWEEN"
	OPER_IN      string = "in"
)

const (
	COMP_AND  string = "and"
	COMP_OR   string = "or"
	COMP_NOT  string = "not"
	COMP_NONE string = ""
)

type SQLCriteria struct {
	PropertyName string
	Operation    string
	Value        interface{}
	Complex      string
}
type AggreType struct {
	Predicate int    //谓词
	ColName   string //字段名
}

const (
	AGG_COUNT int = 1
	AGG_SUM   int = 2
	AGG_AVG   int = 3
	AGG_MAX   int = 4
	AGG_MIN   int = 5
)

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

type MySQLSQLBuileder struct {
	SQLBuilder
}

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
func CreateSQLBuileder(dbType string, tablename string) (ISQLBuilder, error) {
	switch dbType {
	case DBType_MySQL:
		return &MySQLSQLBuileder{
			SQLBuilder: SQLBuilder{
				tableName: tablename}}, nil
	}
	return nil, fmt.Errorf("不支持的数据库类型" + dbType)
}

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

func (c *MySQLSQLBuileder) CreateGetColsSQL() string {
	if c.objectTable == "" {
		return "SELECT column_name,data_type FROM information_schema.columns WHERE table_name='" + c.tableName + "'"
	} else {
		return ""
	}
}

func (c *MySQLSQLBuileder) ClearCriteria() {
	c.criteria = nil
}

func (tc *MySQLSQLBuileder) AddAggre(outfield string, aggreType *AggreType) {
	if tc.aggre == nil {
		tc.aggre = make(map[string]*AggreType)
	}
	tc.aggre[outfield] = aggreType
}

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

func (c *MySQLSQLBuileder) createWhereSubStr() (string, []interface{}) {
	var sqlwhere string
	param := make([]interface{}, 0, len(c.criteria))
	for i, cr := range c.criteria {
		var exp string
		switch cr.Operation {
		case OPER_BETWEEN:
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
		case OPER_IN:
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
			if cr.Complex == COMP_AND || cr.Complex == COMP_OR {
				sqlwhere = fmt.Sprint(sqlwhere, " ", cr.Complex, " ", exp)
			}
		} else {
			sqlwhere = fmt.Sprint(sqlwhere, " ", exp)
		}
	}
	//sql += " WHERE " + sqlwhere
	return " WHERE " + sqlwhere, param
}
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
			case AGG_COUNT:
				p = "COUNT("
			case AGG_AVG:
				p = "AVG("
			case AGG_MAX:
				p = "MAX("
			case AGG_MIN:
				p = "MIN("
			case AGG_SUM:
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
