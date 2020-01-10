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
type DataSourceType int8

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	DataSourceType_SQL      DataSourceType = 1
	DataSourceType_SQLTABLE DataSourceType = 2
	DataSourceType_REST     DataSourceType = 3
	DataSourceType_ENMU     DataSourceType = 4
	DataSourceType_INNER    DataSourceType = 5
)
const (
	DBType_MySQL  string = "mysql"
	DBType_Oracle string = "oracle"
)

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
	Property_Datatype_INT  string = "INT"
	Property_Datatype_DOU  string = "DOUBLE"
	Property_Datatype_STR  string = "STRING"
	Property_Datatype_DATE string = "DATE"
	Property_Datatype_TIME string = "TIME"
	Property_Datatype_ENUM string = "ENUM"
	Property_Datatype_UNKN string = ""
)

type OutFieldProperty struct {
	Source     IDataSource `json:"-"`
	JoinField  string
	ValueField string
	ValueFunc  func(record []interface{}, field []*MyProperty, Source IDataSource) interface{} `json:"-"`
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//对象属性
type MyProperty struct {
	Name          string //属性名
	DataType      string //类型名in
	OutJoin       bool
	Caption       string //显示名
	OutJoinDefine *OutFieldProperty
	Hidden        bool
}

type DBField struct {
	Name     string
	DataType string
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
type FieldDesc struct {
	FieldType string
	Index     int
	//字段元数据，默认由PostAction中的配置信息为fieldmeta的处理程序填充
	Meta      *map[string]string
}

type FieldDescType  map[string]*FieldDesc

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DataResultSet struct {
	Fields FieldDescType
	Data   [][]interface{}
	//ResultSet的元数据
	Meta   string
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
type paramAddInterface interface {
	AddParamValue(obj interface{}) paramAddInterface
}


func (f FieldDescType) Copy() FieldDescType{
	r:=make(FieldDescType)
	for k,v :=range f{
		r[k]=v
	}
	return r
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
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

func DataSourceCompare(dsa, dsb *DBDataSource) bool {
	if dsa == nil && dsb == nil {
		return false
	}
	return dsa.DBAlias == dsb.DBAlias
}
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

func CreateFieldNonType(fields ...string) []*MyProperty {
	temp := make([]*MyProperty, len(fields), len(fields))
	for i, v := range fields {
		temp[i] = &MyProperty{
			Name: v,
		}
	}
	return temp
}
