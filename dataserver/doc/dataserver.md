# 说明

tongserver是一个基于配置的数据服务，主要目的就是针对数据库中的表自动生成操作数据的API，也考虑将数据表扩展为其他的数据源，比如为调用其他的API提供方便，自动生成一致的操作其他API数据的接口等。

类似之前用python写的restgang，也有所不同。restgang是基于SQL语句发布的，主要精力在于sql语句的处理，并且只提供数据查询的接口而没有数据修改的接口，权限控制也能控制到服务一级。

tongserver则将精力放在最常用的数据操作接口中，比如简单的添加、删除、修改，按照一些条件进行查询、排序、统计聚合等，也提供基于SQL语句的发布，但是不对SQL语句做过多的语法分析，也就是说你发布一个drop的SQL语句我一样执行。

简单说一下技术栈，选择golang不是因为go有多好而是因为好奇，rest框架使用beego，同样不是因为他有多好而是随手baidu了一下然后就拿来用了，至少文档是中文的而且挺全。所以beego支持的咱都可以用，但目前数据库还是用的mysql，oracle只做了一点点驱动程序的测试，其他没有继续。



# 数据源

数据源就是为调用者提供数据的东西，可以从数据库里面读取也可以从其他什么地方读取，但是要有一定的规则，不然也没法统一对外提供访问的API接口

* 数据表数据源 **√**

* SQL数据源 **√**

* ValueKey数据源 **√**

  枚举数据源用于数据字典，在数据集后处理中可以作为数据字典使用

* Restful数据源

  

## IDataSource

数据源实现一个IDataSource接口，这个接口里面有一些操作的方法。

```go
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
	QueryDataByFieldValues(fv map[string]interface{}) (*DataResultSet, error)
	GetKeyFields() []*MyProperty
	GetFields() []*MyProperty
	GetFieldByName(name string) *MyProperty
}
```

很多方法有些鸡肋，有时间在整理。实现这个接口即可以被后面的查询程序发布，数据源使用前要先调用init方法进行初始化，其实初始化就是获取一些这个数据源的基础信息，创建一个对象而已。

这个接口只提供了特别基础的功能，实用化很有差距，因此还有一堆别的接口。

```go
// IWriteableDataSource 可写的数据源接口
type IWriteableDataSource interface {
	Delete() error
	Insert(values map[string]interface{}) error
	Update(values map[string]interface{}) error
	AddCriteria(field, operation string, value interface{}) IFilterAdder
	AndCriteria(field, operation string, value interface{}) IFilterAdder
	OrCriteria(field, operation string, value interface{}) IFilterAdder
}

type IJoinedDataSource interface {
	ICriteriaDataSource
	JoinDataSource(join string, ds ICriteriaDataSource, outfield []string) IAddCriteria
}


// IAggregativeAdder 可以聚合的接口
type IAggregativeAdder interface {
	AddAggre(outfield string, aggreType *AggreType)
}
```

对于一个数据源对象，查询它实现了哪些接口就可以知道他能干什么。

## DataResultSet

一个结构体用来描述数据的结果集。

```go
// FieldDesc 返回结果时用的字段描述
type FieldDesc struct {
	FieldType string 
	Index     int //列索引
	// Meta 字段元数据，默认由PostAction中的配置信息为fieldmeta的处理程序填充
	Meta map[string]string
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
```

这里面是用二维数组来保存数据的，第一维是行，第二维是列，所以有一个Fields属性来描述各个列明对应的列索引，就是FieldDesc属性。

目前支持的数据类型包括：

```go
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
	// 字典类型
	PropertyDatetypeMap string = "MAP"
	// 数组类型
	PropertyDatetypeArray string = "ARRAY"
	// struct
	PropertyDatetypeStruct string = "STRUCT"
	// func
	PropertyDatetypeFunc string = "FUNC"
)
```

通过数据库返回的主要是：

```go
	PropertyDatatypeInt string = "INT"
	// PropertyDatatypeDou 浮点数
	PropertyDatatypeDou string = "DOUBLE"
	// PropertyDatatypeStr 字符串
	PropertyDatatypeStr string = "STRING"
	// PropertyDatatypeDate 日期
	PropertyDatatypeDate string = "DATE"
	// PropertyDatatypeTime 时间
	PropertyDatatypeTime string = "TIME"
```

这里没有考虑Blob等字段类型的情况，只考虑的最常用的数字、字符串、时间等类型。



## IQueryableTableSource

可以添加条件的数据源，由几个接口组合。

```go
// ICriteriaDataSource 可以过滤的数据源接口
type ICriteriaDataSource interface {
	IDataSource
	DoFilter() (*DataResultSet, error)
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
	ICriteriaDataSource
}

```



ICriteriaDataSource接口在IDataSource基础上增加了DoFilter() (*DataResultSet, error)方法，用来根据查询条件查询数据。查询条件则通过IFilterAdder接口设定。

* AddCriteria 添加一个条件
* AndCriteria 添加一个与之前的条件为and关系的条件
* OrCriteria 添加一个与之前的条件为or关系的条件
* Orderby 添加一个排序条件

IFilterAdder接口的方法中的field参数即为数据源的列名或字段名，operation参数为操作，value参数为操作值。

* operation 支持以下操作：

  ``` go
  const (
  	// OperEq 等于
  	OperEq string = "="
  	// OperNoteq 不等于
  	OperNoteq string = "<>"
  	// OperGt 大于
  	OperGt string = ">"
  	// OperLt 小于
  	OperLt string = "<"
  	// OperGtEg 大于等于
  	OperGtEg string = ">="
  	// OperLtEg 小于等于
  	OperLtEg string = "<="
  	// OperBetween 介于--之间
  	OperBetween string = "BETWEEN"
  	// OperIn 包含
  	OperIn          string = "in"
  	OperIsNull      string = "is null"
  	OperIsNotNull   string = "is not null"
  	OperAlwaysFalse string = "alwaysfalse"
  	OperAlwaysTrue  string = "alwaystrue"
  )
  ```

  如果operation是OperBetween的话那么value必须是长度为2的数组或者切片，如果operation是OperIn的话，那么value可以是任意长度的数组或切片，也可以是一个字符串，对应SQL语法中的select in 子句。

## IJoinedDataSource

```go
type IJoinedDataSource interface {
	ICriteriaDataSource
	JoinDataSource(join string, ds ICriteriaDataSource, outfield []string) IAddCriteria
}
```



## IAggregativeAdder

```go
// IAggregativeAdder 可以聚合的接口
type IAggregativeAdder interface {
	AddAggre(outfield string, aggreType *AggreType)
}
```



## IWriteableDataSource

```go
// IWriteableDataSource 可写的数据源接口
type IWriteableDataSource interface {
	Delete() error
	Insert(values map[string]interface{}) error
	Update(values map[string]interface{}) error
	AddCriteria(field, operation string, value interface{}) IFilterAdder
	AndCriteria(field, operation string, value interface{}) IFilterAdder
	OrCriteria(field, operation string, value interface{}) IFilterAdder
}
```

需要注意的是每个方法是单独的事务，没有提供其他的事务处理机制。



# 服务





# 流程

[参考流程说明文档](./flow.md)



# MAP定位表达式

内部服务返回结果是map[string]interface{}的格式，很多外部服务返回的结果是JSON格式，可以通过json包转换为map，在流程的各个活动中可以使用map定位表达式来选取map中的元素，实现对服务返回结果的处理。


``` json
{
	"result": true,
	"resultset": {
		"Fields": {
			"ID": {
				"FieldType": "STRING",
				"Index": 0,
				"Meta": null
			},
			"METANAME": {
				"FieldType": "STRING",
				"Index": 3,
				"Meta": null
			},
			"NAMESPACE": {
				"FieldType": "STRING",
				"Index": 2,
				"Meta": null
			},
			"PROJECTID": {
				"FieldType": "STRING",
				"Index": 1,
				"Meta": null
			}
		},
		"Data": [
			["0e50d67f-5eee-41b8-8624-811f5df1dca2", "default", "jeda", "USER"],
			["bptfa6oc5r5oo5esbhhg", "jeda", "table", "meta"],
			["bq1bbgoc5r5uu2jmk9sg", "jeda", "table", "metaitem"]
		],
		"Meta": ""
	}
}
```

上述JSON转换为MAP之后

| 表达式                         | 返回值                                                       |
| ------------------------------ | ------------------------------------------------------------ |
| /result                        | true                                                         |
| /resultset/Fields/ID/FieldType | STRING                                                       |
| /resultset/Fields/ID           | map[string]interface{}{"FieldType":"STRING",<br>"Index":0,<br>"Meta":null} |
| /resultset/Data[0]             | []interface{} {"0e50d67f-5eee-41b8-8624-811f5df1dca2", "default", "jeda", "USER"} |
| /resultset/Data[0] [0]         | "0e50d67f-5eee-41b8-8624-811f5df1dca2"                       |
|                                |                                                              |
|                                |                                                              |
| /resultset/Data/@[1:2]         | Data元素的第2第3条记录组成的切片                             |
|                                |                                                              |

