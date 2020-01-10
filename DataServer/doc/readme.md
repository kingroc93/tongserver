# 基于BeeGo的数据查询服务
## 基本功能
* 根据SQL语句发布Restful接口服务
    * SQL语法分析
    * 分页支持
    * 字典数据支持
    
* 将数据库表直接发布为Restful
    * 单表主键查询
    * 单表条件查询
    * 字典数据支持
    * 支持主从数据查询
    * 支持分页
    * 支持计数、求和、求平均、最大值、最小值等统计函数
    
* 基础用户、角色、组织机构、权限数据模型和数据操作接口
    * 基本数据CURD操作接口
    
* 基于JWT的安全控制接口
* 基于OAtuh2的安全控制接口
* Qrcode的二维码生成器接口
* 支持Oracle、MySQL数据库
## 基本属性
### 数据源类型
```go
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
```
### 通用数据类型
```go
const (
	Property_Datatype_INT  string = "INT"    //整数
	Property_Datatype_DOU  string = "DOUBLE" //浮点数 
	Property_Datatype_STR  string = "STRING" //字符串
	Property_Datatype_DATE string = "DATE"   //日期类型
	Property_Datatype_TIME string = "TIME"   //包含日期的时间类型
	Property_Datatype_ENUM string = "ENUM"   //枚举类型
	Property_Datatype_UNKN string = ""       //未知类型
)
```
### 属性
```go
type MyProperty struct {
	Name          string //属性名
	DataType      string //类型名in
	OutJoin       bool   //是否外联接
	Caption       string //显示名
	OutJoinDefine *OutFieldProperty
}
```
> 属性的外联接：OutJoin属性为true时，该属性为外链接属性，通过OutJoinDefine属性配置具体外链接的逻辑。
引擎在加载该属性时通过外链接配置信息获取外部数据填充该属性。外链接数据源是基础数据的一种即可。外链接配置信息包括：
>```go
>type OutFieldProperty struct {
>	Source IDataSource
>	JoinField  string
>	ValueField string
>	ValueFunc  func(record []interface{}, field []*MyProperty, Source IDataSource) interface{}
>}
>```
## 基础数据源
```go
type IDataSource interface {
	//返回数据源类型
	GetDataSourceType() DataSourceType
	//数据源初始化
	Init() error
	GetName() string
	//返回全部数据
	GetAllData() (*DataResultSet, error)
	//////返回主键信息
	//GetKeyFields() []MyProperty
	//根据主键返回数据
	QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error)
	//根据字段值返回匹配的数据
	QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error)
}
```
## 数据源
* 数据表数据源 ok

* SQL数据源 ok

* 枚举数据源 ok

  枚举数据源用于数据字典，在数据集后处理中可以作为数据字典使用

* Restful数据源

* Webservice数据源

* 静态数据源

* ~~数据源联接~~

* ~~数据源组合~~

## 数据查询
* 主键匹配查询 ok
* 字段值匹配查询 ok
* 复合条件查询 ok
* 时间序列处理
  * 针对DataTableSource，定义时间列
  * 根据分组返回最新数据
  * 返回最新的数据
  * 返回去年同期数据
  * 返回时间段数据

## 结果集后处理
* 数据分组 ok 
* 数据切面 ok
* 格式化 ok
* 列过滤 ok
* 数据字典映射 ok
* 列提取 ok
* 数据集拼接
* 数据集缓存

## 数据服务
* 基于TableDataSource类的服务  ok
* 预定义服务 
* 基于RnmuSource类的服务 ok
* 基于服务流程的服务
* 服务元数据接口 ok
* 系统管理服务
## 外部处理接口

## 安全机制

## 元数据支持

## 可视化服务

### 数据表格

### 基础统计图

```json
{
  "Criteria": [
    {
      "field": "batch_time",
      "operation": "=",
      "value": "2019-11-13",
      "relation": "and"
    }
  ],

  "PostAction":[
  	{
  		"name":"slice",
  		"params":{
  			"xfield":"item_id",
  			"yfield":["dev_id", "site_id","collect_date"],
  			"valuefield":"data_value"
  		}
  	},{
  		"name":"fieldmeta",
  		"params":{"metaurl":"idb.table.iotdata"}
  	},
  	{
  		"name":"bulldozer",
  		"params":{
  			"bulldozer":[
			  	{
			    	"name":"FormatDatafunc",
			    	"params":{
			    		"collect_date": "2006-01-02 15:04"
			    	}
			    }]}
  	}]
}
```

```json
{
  "Criteria": [
    {
      "field": "ORG_ID",
      "operation": "=",
      "value": "001013009",
      "relation": "and"
    }
  ],
  "bulldozer": [
  	{
  		"name":"DictMappingfunc",
  		"params":{
  			"outfield":        "ORG_NAME",
			"dataKeyField":    "ORG_ID",
			"KeyStringSourceName":"ORG_NAME"
  		}
  	},{
    	"name":"FormatDatafunc",
    	"params":{
    		"USER_CREATED": "2006-01-02"
    	}
    },
    {
      "name": "ColumnFilterFunc",
      "params": {
        "show": ["USER_ID", "USER_NAME", "ORG_ID", "USER_CREATED","ORG_NAME"]
      }
    }
  ]
}
```