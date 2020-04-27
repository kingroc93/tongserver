package service

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/rs/xid"
	"reflect"
	"strconv"
	"strings"
	"time"

	"tongserver.dataserver/cube"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/utils"
)

// IDSServiceHandler 数据源服务处理句柄
// IDS数据源服务支持主键查询、条件查询、排序、数据后处理、逐行行处理等
// IWriteableDataSource接口的数据源支持数据的添加删除修改
type IDSServiceHandler struct {
	SHandlerBase
}

// doBulldozer 在switch中处理推土机函数
func (c *IDSServiceHandler) doBulldozer(dataSet *datasource.DataResultSet, index int, name string, param map[string]interface{}) {
	switch name {
	case "DictMappingfunc":
		{
			ksname := param["KeyStringSourceName"].(string)
			if ksname == "" {
				return
			}
			obj := datasource.CreateIDSFromParam(datasource.IDSContainer[ksname])
			if obj == nil {
				return
			}
			ids, _ := obj.(*datasource.KeyStringSource)
			param["KeyStringSource"] = ids
			cube.DictMappingfunc(dataSet, index, param)
		}
	case "FormatDatafunc":
		cube.FormatDatafunc(dataSet, index, param)
	case "ColumnFilterFunc":
		cube.ColumnFilterFunc(dataSet, index, param)
	}
}

func (c *IDSServiceHandler) getCache(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) (*utils.RestResult, error) {
	r, err := c.SHandlerBase.getCache(sdef, ids, rBody)
	if r != nil {
		dst, ok := (*r)["resultset"]
		if ok {
			ds, ok := dst.(*datasource.DataResultSet)
			if ok {
				resuleset, err := c.doPostAction(c.DoBulldozer(ds, rBody.Bulldozer), rBody)
				if err != nil {
					return nil, err
				}
				(*r)["resultset"] = resuleset
			}
		}
	}
	return r, err
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回缓存的结果数据
func (c *IDSServiceHandler) doGetCache(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	r, err := c.getCache(sdef, ids, rBody)
	if r != nil {
		(*r)["result"] = true
	}
	if err != nil {
		(*r)["msg"] = err.Error()
	}
	c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
}

// DoBulldozer 处理推土机函数
func (c *IDSServiceHandler) DoBulldozer(dataSet *datasource.DataResultSet, bulldozer []*CommonParamsType) *datasource.DataResultSet {
	if bulldozer == nil {
		return dataSet
	}
	if len(bulldozer) == 0 {
		return dataSet
	}
	L := len(dataSet.Data)
	for i := 0; i < L; i++ {
		for _, v := range bulldozer {
			c.doBulldozer(dataSet, i, v.Name, v.Params)
		}
	}
	return dataSet
}

// 判断ids是否为DataSource.IWriteableDataSource接口,判断当前请求是否为post,如果是IWriteableDataSource接口则返回
// IWriteableDataSource接口实例,否则返回nil
func (c *IDSServiceHandler) checkWriteableInf(ids interface{}) datasource.IWriteableDataSource {
	inf, ok := ids.(datasource.IWriteableDataSource)
	if !ok {
		c.createErrorResponse("请求的服务没有实现DataSource.IWriteableDataSource接口")
		return nil
	}
	return inf
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//从报文中提取字段值并进行转换
func (c *IDSServiceHandler) getVauleMapFromStringMap(svalue map[string]string, ids datasource.IDataSource) map[string]interface{} {
	values := make(map[string]interface{})
	for k, v := range svalue {
		fs := ids.GetFieldByName(k)
		if fs == nil {
			c.createErrorResponse("Insert节点中描述的字段" + k + "不存在")
			return nil
		}
		fv, err := c.ConvertString2Type(v, fs.DataType)
		if err != nil {
			c.createErrorResponse("字段值类型转换失败，字段：" + k + ",值：" + v + "，预期类型：" + fs.DataType)
			return nil
		}
		values[k] = fv
	}
	return values
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理删除
func (c *IDSServiceHandler) doDelete(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	if inf := c.checkWriteableInf(ids); inf != nil {
		if rBody.Delete != "true" {
			c.createErrorResponse("报文Delete节点的值必须为true")
			return
		}
		if len(rBody.Criteria) == 0 {
			if rBody.OperationConfirm != "delete" {
				c.createErrorResponse("删除操作，但是报文中没有条件节点，此时OperationConfirm节点的值必须为delete")
				return
			}
		}
		if err := c.fillCriteriaFromRbody(ids, rBody); err != nil {
			c.createErrorResponse(err.Error())
			return
		}
		if err := inf.Delete(); err != nil {
			c.createErrorResponse(err.Error())
		} else {
			r := utils.CreateRestResult(true)
			r["msg"] = "处理成功"
			c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//处理更新
func (c *IDSServiceHandler) doUpdate(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	if inf := c.checkWriteableInf(ids); inf != nil {
		if rBody.Update == nil {
			c.createErrorResponse("报文没有update节点")
			return
		}
		if len(rBody.Criteria) == 0 {
			if rBody.OperationConfirm != "update" {
				c.createErrorResponse("更新操作，但是报文中没有条件节点，此时OperationConfirm节点的值必须为update")
				return
			}
		}
		values := c.getVauleMapFromStringMap(rBody.Update, ids)
		if values == nil {
			return
		}
		if err := c.fillCriteriaFromRbody(ids, rBody); err != nil {
			c.createErrorResponse(err.Error())
			return
		}

		if err := inf.Update(values); err != nil {
			c.createErrorResponse(err.Error())
		} else {
			r := utils.CreateRestResult(true)
			r["msg"] = "处理成功"
			c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理添加
func (c *IDSServiceHandler) doInsert(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	if inf := c.checkWriteableInf(ids); inf != nil {
		if rBody.Insert == nil {
			c.createErrorResponse("报文没有insert节点")
			return
		}
		values := c.getVauleMapFromStringMap(rBody.Insert, ids)
		for k, v := range values {
			if v == "newguid()" {
				values[k] = xid.New().String()
			}
		}
		if values == nil {
			return
		}
		err := inf.Insert(values)
		if err != nil {
			c.createErrorResponse(err.Error())
		} else {
			r := utils.CreateRestResult(true)
			r["msg"] = "处理成功"
			c.RRHandler.CreateResponseData(RSP_DATA_STYLE_JSON, r)
		}
	}
}

// 			"values":{
// 				"outfield": "PROJECTNAME"
// 				"ids": "default.mgr.G_USERPROJECT",
//				"filterkey": "PROJECTID",
//				"values":"userid"
// 			}
func (c *IDSServiceHandler) getUserFilterValues(values interface{}) []string {
	if values == nil {
		logs.Error("userfilter节点的values子节点为nil")
		return nil
	}
	switch reflect.TypeOf(values).Kind() {
	case reflect.String:
		if strings.ToLower(values.(string)) == "userid" {
			return []string{c.CurrentUserId}
		} else {
			logs.Error("userfilter节点的values子节点如果是string类型，但值不为userid，该值会被忽略")
			return nil
		}
	case reflect.Map:
		umap := values.(map[string]interface{})
		outfield := umap["outfield"].(string)
		ids := umap["ids"].(string)
		filterkey := umap["filterkey"].(string)
		values := c.getUserFilterValues(umap["values"])
		if len(values) == 0 {
			logs.Error("getUserFilterValues：values 节点定义的数据为空")
			return nil
		}
		obj := datasource.CreateIDSFromParam(datasource.IDSContainer[ids])
		if obj == nil {
			logs.Error("doUserFilter：服务元数据中定义的userfilter节点，中引用的id为" + ids + "的数据源不存在")
			return nil
		}
		dids, ok := obj.(datasource.IQueryableTableSource)
		if !ok {
			logs.Error("doUserFilter：服务元数据中定义的userfilter节点，中引用的id为" + ids + "的数据源没有实现IDataSource接口")
			return nil
		}
		var rs *datasource.DataResultSet
		var err error
		if len(values) == 1 {
			rs, err = dids.QueryDataByFieldValues(map[string]interface{}{filterkey: values[0]})
		} else {
			dids.AddCriteria(filterkey, datasource.OperIn, values)
			rs, err = dids.DoFilter()
		}
		if err != nil {
			logs.Error("doUserFilter：服务元数据中定义的userfilter节点，中引用的id为" + ids + "的数据源在查询数据时发生错误：" + err.Error())
			return nil
		}
		os := make([]string, len(rs.Data), len(rs.Data))
		for index, item := range rs.Data {
			os[index] = item[rs.Fields[outfield].Index].(string)
		}
		return os
	default:
		logs.Error("userfilter节点的values子节点必须是string类型或者map类型")
		return nil
	}
	return nil
}

// 处理用户过滤器，添加用户过滤器的服务，用户查询只返回当前用户的信息
// 根据当前的用户信息对数据进行筛选，通过直接在rBody添加相应的查询条件实现
// 在服务定义元数据中配置过滤的目标列，以及与用户信息的对照操作
// 操作为in或者=，为=时条件为目标字段值等于当前用户id
// 操作为in时，定义目标字段的值包含在idsname定义的数据源中根据userfield等于当前用户id，该数据源必须实现ICriteriaDataSource和IFilterAdder接口
// {
// 		"ids": "default.mgr.G_META",
// 		"userfilter": {
// 			"filterkey": "PROJECTID",
// 			"values":{
// 				"outfield": "PROJECTNAME"
// 				"ids": "default.mgr.G_USERPROJECT",
//				"values": {
//		 			"filterkey": "PROJECTID",
//					"userfield": "USERID",
//				}
// 			}
// 		}
// }
func (c *IDSServiceHandler) doUserFilter(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody **SRequestBody) (bool, error) {
	us, ok := meta["userfilter"]
	if !ok {
		//没有找到userfilter节点，直接返回
		return false, nil
	}

	//一旦找到userfilter节点，则不符合节点要求的数据就不能返回
	if c.CurrentUserId == "" {
		logs.Error("doUserFilter：当前调用者用户id为空")
		return false, fmt.Errorf("doUserFilter：当前调用者用户id为空")
	}
	_, ok = ids.(datasource.ICriteriaDataSource)
	if !ok {
		logs.Info("doUserFilter：服务元数据中定义的userfilter节点，但是服务的数据源没有实现ICriteriaDataSource接口")
		return false, fmt.Errorf("doUserFilter：服务元数据中定义的userfilter节点，但是服务的数据源没有实现ICriteriaDataSource接口")
	}

	umap, ok := us.(map[string]interface{})
	if !ok {
		return false, fmt.Errorf("doUserFilter：userfilter节点格式不正确")
	}
	dfieldname := ""

	if umap["filterkey"] != nil {
		dfieldname = umap["filterkey"].(string)
	}

	values := c.getUserFilterValues(umap["values"])

	if *rBody == nil {
		*rBody = &SRequestBody{}
		t := *rBody
		t.Criteria = make([]CriteriaInRBody, 0, 1)
	}
	t := *rBody

	if len(values) == 0 {
		logs.Error("doUserFilter：定义的values节点返回的数据为空")
		t.Criteria = append(t.Criteria, CriteriaInRBody{
			Field:     dfieldname,
			Operation: datasource.OperAlwaysFalse,
			Value:     "",
			Relation:  datasource.CompAnd})
		return true, nil
	}
	if len(values) == 1 {
		t.Criteria = append(t.Criteria, CriteriaInRBody{
			Field:     dfieldname,
			Operation: datasource.OperEq,
			//OperEq操作值处理values节点返回的第一个值
			Value:    values[0], //c.CurrentUserId,
			Relation: datasource.CompAnd})
		return true, nil
	}
	if len(values) > 1 {
		t.Criteria = append(t.Criteria, CriteriaInRBody{
			Field:     dfieldname,
			Operation: datasource.OperIn,
			Value:     values,
			Relation:  datasource.CompAnd})
		return true, nil
	}
	return false, fmt.Errorf("doUserFilter：values 节点定义的数据为空")
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//返回所有数据
func (c *IDSServiceHandler) doAllData(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	var resuleset *datasource.DataResultSet
	var err error
	evool, err := c.doUserFilter(sdef, meta, ids, &rBody)
	if err != nil {
		c.createErrorResponse(err.Error())
		return
	}
	if evool {
		c.doQuery(sdef, meta, ids, rBody)
		return
	} else {
		resuleset, err = ids.GetAllData()
	}
	c.setPageParams(ids)
	if err != nil {
		c.createErrorResponse(err.Error())
		return
	}
	if rBody == nil {
		c.setResultSet(resuleset)

	} else {
		resuleset, err = c.doPostAction(c.DoBulldozer(resuleset, rBody.Bulldozer), rBody)
		if err != nil {
			c.createErrorResponse(err.Error())
			return
		}
		c.setResultSet(resuleset)

	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 参数转换，将string的参数转换为指定的类型，针对日期类型特殊处理：
// 如前N天，lastday:1    lastday:-3
// addmonth
// addyear
// now
// thismonth
// thisyear
// today
func (c *IDSServiceHandler) convertParamValues(value string, datatype string) (interface{}, error) {
	//特殊处理日期类型
	if datatype == datasource.PropertyDatatypeTime || datatype == datasource.PropertyDatatypeDate {
		switch {
		case strings.HasPrefix(value, "addday"):
			{
				//前N天，lastday:1    lastday:-3
				ss := strings.Split(value, ":")
				var days = 0
				if len(ss) == 1 {
					days = -1
				} else {
					days, _ = strconv.Atoi(ss[1])
				}
				if days == 0 {
					days = -1
				}
				pd, _ := time.ParseDuration(strconv.Itoa(days*24) + "h")
				return c.convertParamValues(time.Now().Add(pd).Format("2006-01-02 15:04:05"), datatype)

			}
		case strings.HasPrefix(value, "addmonth"):
			{
				ss := strings.Split(value, ":")
				var ms = 0
				if len(ss) == 1 {
					ms = -1
				} else {
					ms, _ = strconv.Atoi(ss[1])
				}
				if ms == 0 {
					ms = -1
				}
				return c.convertParamValues(time.Now().AddDate(0, ms, 0).Format("2006-01-02 15:04:05"), datatype)
			}
		case strings.HasPrefix(value, "addyear"):
			{
				ss := strings.Split(value, ":")
				var ms = 0
				if len(ss) == 1 {
					ms = -1
				} else {
					ms, _ = strconv.Atoi(ss[1])
				}
				if ms == 0 {
					ms = -1
				}
				return c.convertParamValues(time.Now().AddDate(0, 0, ms).Format("2006-01-02 15:04:05"), datatype)
			}
		case value == "now":
			{
				//预定义当前时刻
				return c.convertParamValues(time.Now().Format("2006-01-02 15:04:05"), datatype)
			}
		case strings.HasPrefix(value, "thismonth"):
			{
				//预定义当月，thismonth后跟时间，如thismonth,08:00:00
				ss := strings.Split(value, ",")
				n := time.Now()
				timeStr := time.Date(n.Year(), n.Month(), 1, 0, 0, 0, 0, n.Location()).Format("2006-01-02")
				if len(ss) == 1 {
					timeStr += timeStr + " 00:00:00"
				} else {
					timeStr += timeStr + " " + ss[1]
				}
				return c.convertParamValues(timeStr, datatype)
			}
		case strings.HasPrefix(value, "thisyear"):
			{
				ss := strings.Split(value, ",")
				n := time.Now()
				timeStr := time.Date(n.Year(), 1, 1, 0, 0, 0, 0, n.Location()).Format("2006-01-02")
				if len(ss) == 1 {
					timeStr += timeStr + " 00:00:00"
				} else {
					timeStr += timeStr + " " + ss[1]
				}
				return c.convertParamValues(timeStr, datatype)
			}
		case strings.HasPrefix(value, "today"):
			{
				//预定义当日时刻，today后跟时间，如today,08:00:00
				ss := strings.Split(value, ",")
				timeStr := time.Now().Format("2006-01-02")
				if len(ss) == 1 {
					timeStr += timeStr + " 00:00:00"
				} else {
					timeStr += timeStr + " " + ss[1]
				}
				return c.convertParamValues(timeStr, datatype)
			}
		}
	}
	pv, err := c.ConvertString2Type(value, datatype)
	if err != nil {
		return nil, fmt.Errorf("类型转换错误 " + value + " " + datatype + " " + err.Error())
	}
	return pv, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 添加一个查询条件
func (c *IDSServiceHandler) addOneCriteria(v *CriteriaInRBody, ids datasource.IDataSource) error {
	f := ids.GetFieldByName(v.Field)
	if f == nil {
		return fmt.Errorf("没有找到Criteria中定义的字段名" + v.Field)
	}
	var pv interface{}
	switch reflect.TypeOf(v.Value).Kind() {
	case reflect.Slice, reflect.Array:
		{
			s := reflect.ValueOf(v.Value)
			pvs := make([]interface{}, s.Len(), s.Len())
			for i := 0; i < s.Len(); i++ {
				var e error
				pvs[i], e = c.convertParamValues(s.Index(i).Interface().(string), f.DataType)
				if e != nil {
					return e
				}
			}
			pv = pvs
		}
	default:
		{
			p, err := c.convertParamValues(v.Value.(string), f.DataType)
			if err != nil {
				return err
			}
			pv = p
		}
	}

	fc, _ := ids.(datasource.IFilterAdder)
	if strings.ToUpper(v.Relation) == "AND" {
		fc.AndCriteria(v.Field, v.Operation, pv)
	}
	if strings.ToUpper(v.Relation) == "OR" {
		fc.OrCriteria(v.Field, v.Operation, pv)
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//根据请求的报文填充Criteria,ids必须实现DataSource.IFilterAdder接口
func (c *IDSServiceHandler) fillCriteriaFromRbody(ids datasource.IDataSource, rBody *SRequestBody) error {
	_, okfc := ids.(datasource.IFilterAdder)
	//处理条件
	if !okfc {
		return fmt.Errorf("请求的服务没有实现IFilterAdder接口,不能处理Criteria节点")
	}
	for _, v := range rBody.Criteria {
		err := c.addOneCriteria(&v, ids)
		if err != nil {
			return err
		}
	}
	return nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据元数据ID返回元数据结果集
func (c *IDSServiceHandler) getMetaData(metaid string) (*datasource.DataResultSet, error) {
	obj, err := datasource.CreateIDSFromName("default.mgr.G_META_ITEM")
	if err != nil {
		return nil, err
	}
	v, ok := obj.(datasource.IQueryableTableSource)
	if !ok {
		return nil, fmt.Errorf("获取默认数据源default.mgr.G_META_ITEM出错")
	}
	v.AddCriteria("META_ID", datasource.OperEq, metaid)
	r, err := v.DoFilter()
	if err != nil {
		return nil, err
	}
	return r, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//根据project和namespace返回在数据库中定义的该服务返回数据的的元数据ID
func (c *IDSServiceHandler) getMetaID(project string, namespace string, metaname string) (string, error) {
	obj, err := datasource.CreateIDSFromName("default.mgr.G_META")
	if err != nil {
		return "", err
	}
	v, ok := obj.(datasource.IQueryableTableSource)
	if !ok {
		return "", fmt.Errorf("获取默认数据源default.mgr.G_META出错")
	}
	v.AddCriteria("NAMESPACE", datasource.OperEq, namespace)
	v.AndCriteria("METANAME", datasource.OperEq, metaname)
	if project != "" {
		v.AndCriteria("PROJECTID", datasource.OperEq, project)
	}
	rs, err := v.DoFilter()
	if err != nil {
		return "", err
	}
	metaid := rs.Data[0][rs.Fields["ID"].Index].(string)
	return metaid, nil
}

//数据后处理
func (c *IDSServiceHandler) doPostAction(dataSet *datasource.DataResultSet, rBody *SRequestBody) (*datasource.DataResultSet, error) {
	if len(rBody.PostAction) == 0 {
		return dataSet, nil
	}
	var rdataset = dataSet
	for _, item := range rBody.PostAction {
		switch item.Name {
		//根据字段分组
		case "fieldgroup":
			{
				field := item.Params["field"].(string)
				rdataset = cube.GroupField(rdataset, field)
			}
		//行列转换
		case "row2column":
			{
				fields := strings.Split(item.Params["fields"].(string), ",")
				rdataset = cube.Row2Column(rdataset, fields...)
			}
		//添加字段的元数据
		case "fieldmeta":
			{
				url := item.Params["metaurl"].(string)
				ss := strings.Split(url, ".")
				project := ""
				namespace := ""
				metaname := ""
				if len(ss) == 3 {
					project = ss[0]
					namespace = ss[1]
					metaname = ss[2]
				} else if len(ss) == 2 {
					namespace = ss[0]
					metaname = ss[1]
				} else {
					return nil, fmt.Errorf("处理PostAction发生错误metaurl非法")
				}
				metaid, err := c.getMetaID(project, namespace, metaname)
				if err != nil {
					return nil, err
				}
				metaset, err := c.getMetaData(metaid)
				if err != nil {
					return nil, err
				}
				mf := map[string]string{}
				for _, item := range metaset.Data {
					mf[item[metaset.Fields["NAME"].Index].(string)] = item[metaset.Fields["VALUE"].Index].(string)
				}
				for k, item := range rdataset.Fields {
					o, ok := mf[k]
					if ok {
						if item.Meta == nil {
							item.Meta = make(map[string]string)
						}
						item.Meta["CAP"] = o
					}
				}
			}
		//数据转置
		case "slice":
			{
				//rss := Slice(rs, "ITEM_ID", []string{"DEV_ID", "SITE_ID","COLLECT_DATE"}, "DATA_VALUE")
				yfi := item.Params["yfield"].([]interface{})
				yf := make([]string, len(yfi), len(yfi))
				for i, v := range yfi {
					yf[i] = v.(string)
				}
				rdataset = cube.Slice(rdataset, item.Params["xfield"].(string), yf, item.Params["valuefield"].(string))
			}
		//按行处理
		case "bulldozer":
			{
				ps := item.Params["bulldozer"].([]interface{})
				tmpb := make([]*CommonParamsType, len(ps), len(ps))
				for i, v := range ps {
					fmt.Println(reflect.TypeOf(v))
					p := &CommonParamsType{
						Name:   v.(map[string]interface{})["name"].(string),
						Params: v.(map[string]interface{})["params"].(map[string]interface{}),
					}
					tmpb[i] = p
				}
				rdataset = c.DoBulldozer(rdataset, tmpb)
			}
		}
	}
	return rdataset, nil
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理查询报文
func (c *IDSServiceHandler) doQuery(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	if rBody == nil {
		c.createErrorResponse("query操作必须POST方式提交rbody信息")
		return
	}
	c.setPageParams(ids)
	fids, ok := ids.(datasource.ICriteriaDataSource)
	if !ok {
		c.createErrorResponse("请求的服务没有实现ICriteriaDataSource接口,不能处理Query请求")
		return
	}
	if len(rBody.Criteria) != 0 {
		err := c.fillCriteriaFromRbody(ids, rBody)
		if err != nil {
			c.createErrorResponse(err.Error())
			return
		}
	}
	fc, okfc := ids.(datasource.IFilterAdder)
	if len(rBody.OrderBy) != 0 {
		//处理排序
		if !okfc {
			c.createErrorResponse("请求的服务没有实现IFilterAdder接口,不能处理Criteria节点")
			return
		}
		os := strings.Split(rBody.OrderBy, ",")
		for _, ov := range os {
			orders := strings.Split(strings.Trim(ov, " "), " ")
			if len(orders) == 2 {
				fc.Orderby(orders[0], orders[1])
			}
		}
	}
	if len(rBody.Aggre) != 0 {
		//处理聚合
		ag, okag := ids.(datasource.IAggregativeAdder)
		if !okag {
			c.createErrorResponse("请求的服务没有实现IAggregativeAdder,不能处理Aggre节点")
			return
		}
		for _, agg := range rBody.Aggre {
			/*	AggCount int = 1
				AggSum   int = 2
				AggAvg   int = 3
				AggMax   int = 4
				AggMin   int = 5  */
			pred := strings.ToUpper(agg.Predicate)
			p := 0
			switch pred {
			case "COUNT":
				p = datasource.AggCount
			case "SUM":
				p = datasource.AggSum
			case "AVG":
				p = datasource.AggAvg
			case "MAX":
				p = datasource.AggMax
			case "MIN":
				p = datasource.AggMin
			}
			ag.AddAggre(agg.Outfield, &datasource.AggreType{
				Predicate: p,
				ColName:   agg.ColName,
			})
		}
	}
	resuleset, err := fids.DoFilter()
	if err != nil {
		c.createErrorResponse(err.Error())
	} else {
		resuleset, err = c.doPostAction(c.DoBulldozer(resuleset, rBody.Bulldozer), rBody)
		if err != nil {
			c.createErrorResponse(err.Error())

		}
		c.setResultSet(resuleset)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回当前类支持的动作类型,以及动作对应的操作函数
// 该方法由 func (c *IDSServiceHandler) DoSrv(sdef *SDefine, inf SHandlerInterface)方法调用
func (c *IDSServiceHandler) getActionMap() map[string]SerivceActionHandler {
	r := c.SHandlerBase.getActionMap()
	r[SrvActionCACHE] = c.doGetCache
	r[SrvActionQUERY] = c.doQuery
	r[SrvActionDELETE] = c.doDelete
	r[SrvActionUPDATE] = c.doUpdate
	r[SrvActionINSERT] = c.doInsert
	r[SrvActionALLDATA] = c.doAllData
	r[SrvActionGET] = c.doGetValueByKey
	return r
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据主键返回数据
func (c *IDSServiceHandler) doGetValueByKey(sdef *SDefine, meta map[string]interface{}, ids datasource.IDataSource, rBody *SRequestBody) {
	fs := ids.GetKeyFields()
	params := make([]interface{}, len(fs), len(fs))
	for i, f := range fs {
		var err error
		params[i], err = c.ConvertString2Type(c.RRHandler.GetParam(f.Name), f.DataType)
		if err != nil {
			c.createErrorResponse("类型转换错误" + c.RRHandler.GetParam(f.Name) + " " + f.DataType + " err:" + err.Error())
			return
		}
	}
	resuleset, err := ids.QueryDataByKey(params...)
	if err != nil {
		c.createErrorResponse(err.Error())
	} else {
		c.setResultSet(c.DoBulldozer(resuleset, rBody.Bulldozer))
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回请求报文
func (c *IDSServiceHandler) getRBody() *SRequestBody {
	rBody, err := c.RRHandler.GetRequestBody()
	if err != nil {
		c.createErrorResponse("解析报文时发生错误" + err.Error())
		return nil
	}
	return rBody
}
