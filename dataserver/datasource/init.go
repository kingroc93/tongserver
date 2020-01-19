package datasource

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
type TDFilter SQLCriteria

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//数据源接口
//数据源接口不是线程安全的，因此每一个Web请求都需要创建独立的数据源类
type IDataSource interface {
	//返回数据源类型
	GetDataSourceType() DataSourceType
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

//可写的数据源接口
type IWriteableDataSource interface {
	Delete() error
	Insert(values map[string]interface{}) error
	Update(values map[string]interface{}) error
	AddCriteria(field, operation string, value interface{}) IFilterAdder
	AndCriteria(field, operation string, value interface{}) IFilterAdder
	OrCriteria(field, operation string, value interface{}) IFilterAdder
}

//可以过滤的数据源接口
type ICriteriaDataSource interface {
	IDataSource
	DoFilter() (*DataResultSet, error)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//过滤条件接口
type IFilterAdder interface {
	AddCriteria(field, operation string, value interface{}) IFilterAdder
	AndCriteria(field, operation string, value interface{}) IFilterAdder
	OrCriteria(field, operation string, value interface{}) IFilterAdder
	Orderby(field string, dir string) IFilterAdder
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 可以聚合的接口
type IAggregativeAdder interface {
	AddAggre(outfield string, aggreType *AggreType)
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 数据源类型
type DataSourceType int8

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	// SQL数据源
	DataSourceType_SQL DataSourceType = 1
	// 数据库表数据源
	DataSourceType_SQLTABLE DataSourceType = 2
	// REST服务数据源
	DataSourceType_REST DataSourceType = 3
	// 枚举数据源
	DataSourceType_ENMU DataSourceType = 4
	// 联接数据源
	DataSourceType_INNER DataSourceType = 5
)
const (
	// MySQL数据库类型
	DBType_MySQL string = "mysql"
	// Oracle数据库类型
	DBType_Oracle string = "oracle"
)

// 根据数据源类型返回数据源类型的String表达
func GetDataSourceTypeStr(t DataSourceType) string {
	switch t {
	case DataSourceType_SQL:
		return "SQL"
	case DataSourceType_SQLTABLE:
		return "SQLTable"
	case DataSourceType_REST:
		return "RESTService"
	case DataSourceType_ENMU:
		return "ENMU"
	case DataSourceType_INNER:
		return "INNER"
	}
	return "UNKNOW"
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//通用类型

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	//整数
	Property_Datatype_INT string = "INT"
	//浮点数
	Property_Datatype_DOU string = "DOUBLE"
	//字符串
	Property_Datatype_STR string = "STRING"
	//日期
	Property_Datatype_DATE string = "DATE"
	//时间
	Property_Datatype_TIME string = "TIME"
	//枚举类型
	Property_Datatype_ENUM string = "ENUM"
	//数据集，表示该数据数值是一个数据集
	Property_Datatype_DS string = "DATASET"
	//未知类型
	Property_Datatype_UNKN string = ""
)

// 联接时使用的输出字段
type OutFieldProperty struct {
	Source     IDataSource `json:"-"`
	JoinField  string
	ValueField string
	ValueFunc  func(record []interface{}, field []*MyProperty, Source IDataSource) interface{} `json:"-"`
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//对象属性
type MyProperty struct {
	//属性名
	Name string
	//类型名in
	DataType string
	//是否为联接字段
	OutJoin bool
	//显示名
	Caption string
	//联接定义
	OutJoinDefine *OutFieldProperty
	//是否隐藏,该属性目前没有处理
	Hidden bool
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//数据源
type DataSource struct {
	dtype    DataSourceType
	Name     string
	KeyField []*MyProperty
	Field    []*MyProperty
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回结果时用的字段描述
type FieldDesc struct {
	FieldType string
	Index     int
	//字段元数据，默认由PostAction中的配置信息为fieldmeta的处理程序填充
	Meta *map[string]string
}

// 字段描述类型
type FieldDescType map[string]*FieldDesc

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回的结果集
type DataResultSet struct {
	// 字段列表
	Fields FieldDescType
	// 二维表数据
	Data [][]interface{}
	//ResultSet的元数据
	Meta string
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 复制FieldDescType
func (f FieldDescType) Copy() FieldDescType {
	r := make(FieldDescType)
	for k, v := range f {
		r[k] = v
	}
	return r
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 类型表达形式转换
func ConvertMySQLType2CommonType(t string) string {
	switch t {
	case "VARCHAR":
		return Property_Datatype_STR
	case "NVARCHAR":
		return Property_Datatype_STR
	case "CHAR":
		return Property_Datatype_STR
	case "INT":
		return Property_Datatype_INT
	case "TINYINT":
		return Property_Datatype_INT
	case "SMALLINT":
		return Property_Datatype_INT
	case "MEDIUMINT":
		return Property_Datatype_INT
	case "INTEGER":
		return Property_Datatype_INT
	case "BIGINT":
		return Property_Datatype_INT
	case "FLOAT":
		return Property_Datatype_DOU
	case "DOUBLE":
		return Property_Datatype_DOU
	case "TIMESTAMP":
		return Property_Datatype_TIME
	case "DATE":
		return Property_Datatype_DATE
	case "DATETIME":
		return Property_Datatype_TIME
	case "TIME":
		return Property_Datatype_TIME
	}
	return Property_Datatype_UNKN
}

// 对比连个数据源是否是一个数据源
func DataSourceCompare(dsa, dsb *DBDataSource) bool {
	if dsa == nil && dsb == nil {
		return false
	}
	return dsa.DBAlias == dsb.DBAlias
}

// 创建可写的数据表数据源
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

// 创建数据表数据源
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

// 根据数数据组创建一堆字段属性,不设定属性类型
func CreateFieldNonType(fields ...string) []*MyProperty {
	temp := make([]*MyProperty, len(fields), len(fields))
	for i, v := range fields {
		temp[i] = &MyProperty{
			Name: v,
		}
	}
	return temp
}
