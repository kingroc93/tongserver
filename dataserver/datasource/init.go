package datasource

// TDFilter 过滤条件
type TDFilter SQLCriteria

// IDataSource 数据源接口
//数据源接口不是线程安全的，因此每一个Web请求都需要创建独立的数据源类
type IDataSource interface {
	//返回数据源类型
	GetDataSourceType() DSType
	//数据源初始化
	Init() error
	GetName() string
	//返回全部数据
	GetAllData() (*DataResultSet, error)
	SetRowsLimit(limit int)
	SetRowsOffset(offset int)
	//根据主键返回数据
	QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error)
	//根据字段值返回匹配的数据
	QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error)
	GetKeyFields() []*MyProperty
	GetFields() []*MyProperty
	GetFieldByName(name string) *MyProperty
}

// IWriteableDataSource 可写的数据源接口
type IWriteableDataSource interface {
	Delete() error
	Insert(values map[string]interface{}) error
	Update(values map[string]interface{}) error
	AddCriteria(field, operation string, value interface{}) IFilterAdder
	AndCriteria(field, operation string, value interface{}) IFilterAdder
	OrCriteria(field, operation string, value interface{}) IFilterAdder
}

// ICriteriaDataSource 可以过滤的数据源接口
type ICriteriaDataSource interface {
	IDataSource
	DoFilter() (*DataResultSet, error)
}

// 支持内连接和外链接的数据源
type IJoinedDataSource interface {
	ICriteriaDataSource
	JoinDataSource(join string, ds ICriteriaDataSource, outfield []string) IAddCriteria
}

// IFilterAdder 过滤条件接口
type IFilterAdder interface {
	AddCriteria(field, operation string, value interface{}) IFilterAdder
	AndCriteria(field, operation string, value interface{}) IFilterAdder
	OrCriteria(field, operation string, value interface{}) IFilterAdder
	Orderby(field string, dir string) IFilterAdder
}

type IQueryableTableSource interface {
	IFilterAdder
	IJoinedDataSource
}

// IAggregativeAdder 可以聚合的接口
type IAggregativeAdder interface {
	AddAggre(outfield string, aggreType *AggreType)
}

// DSType 数据源类型
type DSType int8

const (
	// DataSourceTypeSQL SQL数据源
	DataSourceTypeSQL DSType = 1
	// DataSourceTypeSqltable 数据库表数据源
	DataSourceTypeSqltable DSType = 2
	// DataSourceTypeRest REST服务数据源
	DataSourceTypeRest DSType = 3
	// DataSourceTypeEnmu 枚举数据源
	DataSourceTypeEnmu DSType = 4
	// DataSourceTypeInner 联接数据源
	DataSourceTypeInner DSType = 5
)
const (
	// DbTypeMySQL MySQL数据库类型
	DbTypeMySQL string = "mysql"
	// DbTypeOracle Oracle数据库类型
	DbTypeOracle string = "oracle"
)

// GetDataSourceTypeStr 根据数据源类型返回数据源类型的String表达
func GetDataSourceTypeStr(t DSType) string {
	switch t {
	case DataSourceTypeSQL:
		return "SQL"
	case DataSourceTypeSqltable:
		return "SQLTable"
	case DataSourceTypeRest:
		return "RESTService"
	case DataSourceTypeEnmu:
		return "ENMU"
	case DataSourceTypeInner:
		return "INNER"
	}
	return "UNKNOW"
}

const (
	// PropertyDatatypeInt 整数
	PropertyDatatypeInt string = "INT"
	// PropertyDatatypeDou 浮点数
	PropertyDatatypeDou string = "DOUBLE"
	// PropertyDatatypeStr 字符串
	PropertyDatatypeStr string = "STRING"
	// PropertyDatatypeDate 日期
	PropertyDatatypeDate string = "DATE"
	// PropertyDatatypeTime 时间
	PropertyDatatypeTime string = "TIME"
	// PropertyDatatypeEnum 枚举类型
	PropertyDatatypeEnum string = "ENUM"
	// PropertyDatatypeDs 数据集，表示该数据数值是一个数据集
	PropertyDatatypeDs string = "DATASET"
	// PropertyDatatypeUnkn 未知类型
	PropertyDatatypeUnkn string = ""
)

// OutFieldProperty 联接时使用的输出字段
type OutFieldProperty struct {
	Source     IDataSource `json:"-"`
	JoinField  string
	ValueField string
	ValueFunc  func(record []interface{}, field []*MyProperty, Source IDataSource) interface{} `json:"-"`
}

// MyProperty 对象属性
type MyProperty struct {
	// Name属性名
	Name string
	// DataType 类型名
	DataType string
	// OutJoin 是否为联接字段
	OutJoin bool
	// Caption 显示名
	Caption string
	// OutJoinDefine 联接定义
	OutJoinDefine *OutFieldProperty
	// Hidden 是否隐藏,该属性目前没有处理
	Hidden bool
}

// DataSource 数据源
type DataSource struct {
	dtype    DSType
	Name     string
	KeyField []*MyProperty
	Field    []*MyProperty
}

// FieldDesc 返回结果时用的字段描述
type FieldDesc struct {
	FieldType string
	Index     int
	// Meta 字段元数据，默认由PostAction中的配置信息为fieldmeta的处理程序填充
	Meta *map[string]string
}

// FieldDescType 字段描述类型
type FieldDescType map[string]*FieldDesc

// DataResultSet 返回的结果集
type DataResultSet struct {
	// Fields 字段列表
	Fields FieldDescType
	// Data 二维表数据
	Data [][]interface{}
	// Meta ResultSet的元数据
	Meta string
}

// Copy 复制FieldDescType
func (f FieldDescType) Copy() FieldDescType {
	r := make(FieldDescType)
	for k, v := range f {
		r[k] = v
	}
	return r
}

// ConvertMySQLType2CommonType 类型表达形式转换
func ConvertMySQLType2CommonType(t string) string {
	switch t {
	case "VARCHAR":
		return PropertyDatatypeStr
	case "NVARCHAR":
		return PropertyDatatypeStr
	case "CHAR":
		return PropertyDatatypeStr
	case "INT":
		return PropertyDatatypeInt
	case "TINYINT":
		return PropertyDatatypeInt
	case "SMALLINT":
		return PropertyDatatypeInt
	case "MEDIUMINT":
		return PropertyDatatypeInt
	case "INTEGER":
		return PropertyDatatypeInt
	case "BIGINT":
		return PropertyDatatypeInt
	case "FLOAT":
		return PropertyDatatypeDou
	case "DOUBLE":
		return PropertyDatatypeDou
	case "TIMESTAMP":
		return PropertyDatatypeTime
	case "DATE":
		return PropertyDatatypeDate
	case "DATETIME":
		return PropertyDatatypeTime
	case "TIME":
		return PropertyDatatypeTime
	}
	return PropertyDatatypeUnkn
}

// SourceCompare 对比连个数据源是否是一个数据源
func SourceCompare(dsa, dsb *DBDataSource) bool {
	if dsa == nil && dsb == nil {
		return false
	}
	return dsa.DBAlias == dsb.DBAlias
}

// CreateWriteableTableDataSource 创建可写的数据表数据源
func CreateWriteableTableDataSource(name, dbAlias, tablename string, fields ...string) *WriteableTableSource {
	ids := &WriteableTableSource{
		TableDataSource{
			DBDataSource: DBDataSource{
				DataSource: DataSource{
					Name: name,
				},
				DBAlias:        dbAlias,
				AutoFillFields: true,
			},
			TableName: tablename,
		}}
	ids.Init()
	return ids
}

// CreateSQLDataSource 创建SQL数据源sql语句必须是查询语句
func CreateSQLDataSource(name, dbAlias, sql string, fields ...string) *SQLDataSource {
	var fs []*MyProperty = nil
	if len(fields) != 0 {
		fs = make([]*MyProperty, len(fields), len(fields))
		for i, e := range fields {
			fs[i] = &MyProperty{Name: e}
		}
	}
	sqld := &SQLDataSource{
		DBDataSource: DBDataSource{
			DataSource: DataSource{
				Name:  name,
				Field: fs},
			DBAlias:        dbAlias,
			AutoFillFields: false,
		},
		SQL: sql}
	sqld.Init()
	return sqld
}

// CreateTableDataSource 创建数据表数据源
func CreateTableDataSource(name, dbAlias, tablename string, fields ...string) *TableDataSource {

	ids := &TableDataSource{
		DBDataSource: DBDataSource{
			DataSource: DataSource{
				Name: name,
			},
			DBAlias:        dbAlias,
			AutoFillFields: true,
		},
		TableName: tablename,
	}
	ids.Init()
	return ids
}

// CreateFieldNonType 根据数数据组创建一堆字段属性,不设定属性类型
func CreateFieldNonType(fields ...string) []*MyProperty {
	temp := make([]*MyProperty, len(fields), len(fields))
	for i, v := range fields {
		temp[i] = &MyProperty{
			Name: v,
		}
	}
	return temp
}
