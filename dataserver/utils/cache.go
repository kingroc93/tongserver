package utils

import "github.com/astaxie/beego/cache"

//用于保存数据集的配置信息的缓存
var DataSourceCache, _ = cache.NewCache("memory", `{"interval":600}`)

//用于保存数据字典信息
var DictDataCache, _ = cache.NewCache("memory", `{"interval":600}`)
