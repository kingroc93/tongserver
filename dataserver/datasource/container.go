package datasource

import (
	"fmt"
	"time"
	"tongserver.dataserver/utils"
)

type IDSContainerParam map[string]interface{}
type IDSContainerType map[string]IDSContainerParam

var DBAlias2DBTypeContainer = make(map[string]string)

//数据源的配置信息，从数据库中获取，由main函数加载
var IDSContainer = make(IDSContainerType)

//数据源的创建函数集合，根据名字选择合适的创建函数，创建数据源
var iDSCreator = map[string]func(p IDSContainerParam) interface{}{
	"CreateTableDataSource": func(p IDSContainerParam) interface{} {
		return CreateTableDataSource(p["name"].(string), p["dbalias"].(string), p["tablename"].(string))
	},
	"CreateWriteableTableDataSource": func(p IDSContainerParam) interface{} {
		return CreateWriteableTableDataSource(p["name"].(string), p["dbalias"].(string), p["tablename"].(string))
	},
	"CreateKeyStringFromTableSource": func(p IDSContainerParam) interface{} {
		if p["cached"] == "true" {
			obj := utils.DictDataCache.Get(p["name"].(string))
			if obj != nil {
				return obj.(*KeyStringSource)
			}
		}
		ks := &KeyStringSource{
			DataSource: DataSource{
				Name: p["name"].(string),
			},
		}
		ks.Init()
		ts := CreateTableDataSource(p["name"].(string)+"_", p["dbalias"].(string), p["tablename"].(string))
		ks.FillDataByDataSource(ts, p["keyfield"].(string), p["valuefield"].(string))
		if p["cached"] == "true" {
			utils.DictDataCache.Put(p["name"].(string), ks, 5*time.Minute)
		}
		return ks
	},
}

func CreateIDSFromParam(p IDSContainerParam) interface{} {
	if p == nil {
		return nil
	}
	fu := iDSCreator[p["inf"].(string)]
	if fu == nil {
		return nil
	}
	return fu(p)
}
func RegisterIDSCreatorFun(name string, f func(p IDSContainerParam) interface{}) {
	iDSCreator[name] = f
}

func CreateIDSFromName(name string) (interface{}, error) {
	param := IDSContainer[name]
	obj := CreateIDSFromParam(param)
	if obj == nil {
		return nil, fmt.Errorf(name)
	}
	return obj, nil
}

//初始化
func init() {

	//IDSContainer["ORG_NAME"] = map[string]string{
	//	"inf":        "CreateKeyStringFromTableSource",
	//	"name":       "ORG_NAME",
	//	"dbalias":    "default",
	//	"tablename":  "JEDA_ORG",
	//	"keyfield":   "ORG_ID",
	//	"valuefield": "ORG_NAME",
	//	"cached":     "true",
	//}

}
