package datasource

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"time"
	"tongserver.dataserver/utils"
)

// IDSContainerParam 数据源创建参数
type IDSContainerParam map[string]interface{}

// IDSContainerType 数据源容器类型
type IDSContainerType map[string]IDSContainerParam

// DBAlias2DBTypeContainer 用于保存数据连接别名和数据库类型的关系
var DBAlias2DBTypeContainer = make(map[string]string)

// IDSContainer 数据源的配置信息，从数据库中获取，由main函数加载
var IDSContainer IDSContainerType

// iDSCreator 数据源的创建函数集合，根据名字选择合适的创建函数，创建数据源
var iDSCreator = make(map[string]func(p IDSContainerParam) interface{})

func AddIdsCreator(name string, f func(p IDSContainerParam) interface{}) {
	iDSCreator[name] = f
}

// CreateIDSFromParam 根据配置参数创建数据源接口,这个配置参数是保存在数据库里面的
func CreateIDSFromParam(p IDSContainerParam) interface{} {
	if p == nil {
		return nil
	}
	inf := p["inf"].(string)
	if inf != "CreateKeyStringFromIds" {
		fu := iDSCreator[p["inf"].(string)]
		if fu == nil {
			return nil
		}
		return fu(p)
	}
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

	idsname := p["idsname"].(string)
	ids, err := CreateIDSFromName(idsname)
	if err != nil {
		logs.Error(err)
		return nil
	}
	ks.FillDataByDataSource(ids.(IDataSource), p["keyfield"].(string), p["valuefield"].(string))
	if p["cached"] == "true" {
		utils.DictDataCache.Put(p["name"].(string), ks, 5*time.Minute)
	}
	return ks
}

// RegisterIDSCreatorFun 注册数据源创建函数
func RegisterIDSCreatorFun(name string, f func(p IDSContainerParam) interface{}) {
	iDSCreator[name] = f
}

// CreateIDSFromName 根据名称返回数据源接口
func CreateIDSFromName(name string) (interface{}, error) {
	param, ok := IDSContainer[name]
	if !ok {
		return nil, fmt.Errorf("没有找到元名称为" + name + "的数据源")
	}
	obj := CreateIDSFromParam(param)
	if obj == nil {
		return nil, fmt.Errorf(name)
	}
	return obj, nil
}

// init 初始化
func init() {
}
