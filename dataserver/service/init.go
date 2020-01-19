package service

import (
	"github.com/astaxie/beego"
)

// 请求返回的数据类型
type RestResult map[string]interface{}

// 创建返回数据实例
func CreateRestResult(success bool) RestResult {
	var result = make(RestResult)
	result["result"] = success
	return result
}

// 请求的通用参数
type CommonParamsType struct {
	Name   string
	Params map[string]interface{}
}

// 请求的rbody中的查询条件
type CriteriaInRBody struct {
	Field     string
	Operation string
	Value     interface{}
	Relation  string
}

// 请求报文体
type ServiceRequestBody struct {
	//新建
	Insert map[string]string
	//更新
	Update map[string]string
	//删除
	Delete string
	//操作二次确认
	OperationConfirm string
	//条件节点,针对更新、删除、查询操作
	Criteria []CriteriaInRBody
	// 排序节点，针对查询操作
	OrderBy string
	// 内连接节点，针对查询操作
	InnerJoin string
	// 聚合节点，针对查询操作
	Aggre []struct {
		Outfield  string
		Predicate string
		ColName   string
	}
	// 推土机节点，针对查询操作
	Bulldozer []*CommonParamsType
	// 后处理节点，针对查询操作
	PostAction []*CommonParamsType
}

func init() {
	// 所有服务请求的入口函数
	beego.Router("/services/?:context/?:action", &ServiceController{}, "get,post:DoSrv")

}
