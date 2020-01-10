package service

import (
	"awesome/datasource"
	"awesome/cube"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	//返回全部数据
	SrvAction_ALLDATA string = "all"
	//查询动作
	SrvAction_QUERY string = "query"
	//根据主键返回
	SrvAction_GET string = "get"
	//根据字段值返回
	SrvAction_BYFIELD string = "byfield"
	//返回服务元数据
	SrvAction_META string = "meta"
	//删除操作
	SrvAction_DELETE string = "delete"
	//更新操作
	SrvAction_UPDATE string = "update"
	//插入操作
	SrvAction_INSERT string = "insert"

	//以下三个常量均为通过QueryString传入的参数名
	//针对查询自动分页中每页记录数
	REQUEST_PARAM_PAGESIZE string = "_pagesize"
	//针对查询自动分页中的页索引
	REQUEST_PARAM_PAGEINDEX string = "_pageindex"
	//是否返回字段元数据，默认为返回
	REQUEST_PARAM_NOFIELDSINFO string = "_nofield"
	//当前请求不执行而是只返回SQL语句，仅针对IDS类型的服务有效
	REQUEST_PARAM_SQL string = "_sql"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
type IDSServiceHandler struct {
	ServiceHandlerBase
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 在switch中处理推土机函数
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理推土机函数
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

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 判断ids是否为DataSource.IWriteableDataSource接口,判断当前请求是否为post,如果是IWriteableDataSource接口则返回
// IWriteableDataSource接口实例,否则返回nil
func (c *IDSServiceHandler) checkMethodAndWriteableInf(ids interface{}) datasource.IWriteableDataSource {
	if c.Ctl.Ctx.Input.Method() != "POST" {
		c.createErrorResponse("Query动作必须发起POST请求")
		return nil
	}
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
func (c *IDSServiceHandler) doDelete(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	if inf := c.checkMethodAndWriteableInf(ids); inf != nil {
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
			c.createErrorResponseByError(err)
			return
		}
		if err := inf.Delete(); err != nil {
			c.createErrorResponseByError(err)
		} else {
			c.setResult("处理成功")
			c.ServeJson()
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//处理更新
func (c *IDSServiceHandler) doUpdate(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	if inf := c.checkMethodAndWriteableInf(ids); inf != nil {
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
			c.createErrorResponseByError(err)
			return
		}

		if err := inf.Update(values); err != nil {
			c.createErrorResponseByError(err)
		} else {
			c.setResult("处理成功")
			c.ServeJson()
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理添加
func (c *IDSServiceHandler) doInsert(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	if inf := c.checkMethodAndWriteableInf(ids); inf != nil {
		if rBody.Insert == nil {
			c.createErrorResponse("报文没有insert节点")
			return
		}
		values := c.getVauleMapFromStringMap(rBody.Insert, ids)
		if values == nil {
			return
		}
		err := inf.Insert(values)
		if err != nil {
			c.createErrorResponseByError(err)
		} else {
			c.setResult("处理成功")
			c.ServeJson()
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//返回所有数据
func (c *IDSServiceHandler) doAllData(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	c.setPageParams(ids)
	resuleset, err := ids.GetAllData()
	if err != nil {
		c.createErrorResult(err.Error())
		return
	}
	if rBody == nil {
		c.setResultSet(resuleset)
		c.ServeJson()
	} else {
		resuleset, err = c.doPostAction(c.DoBulldozer(resuleset, rBody.Bulldozer), rBody)
		if err != nil {
			c.createErrorResponseByError(err)
			return
		}
		c.setResultSet(resuleset)
		c.ServeJson()
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//根据请求的报文填充Criteria,ids必须实现DataSource.IFilterAdder接口
func (c *IDSServiceHandler) fillCriteriaFromRbody(ids datasource.IDataSource, rBody *ServiceRequestBody) error {
	fc, okfc := ids.(datasource.IFilterAdder)
	//处理条件
	if !okfc {
		return fmt.Errorf("请求的服务没有实现IFilterAdder接口,不能处理Criteria节点")
	}
	for i, v := range rBody.Criteria {
		f := ids.GetFieldByName(v.Field)
		if f == nil {
			return fmt.Errorf("没有找到Criteria中定义的字段名" + v.Field)
		}
		pv, err := c.ConvertString2Type(v.Value, f.DataType)
		if err != nil {
			return fmt.Errorf("类型转换错误 " + v.Value + " " + f.DataType + " " + err.Error())
		}
		if i == 0 {
			fc.AddCriteria(v.Field, v.Operation, pv)
		} else {
			if strings.ToUpper(v.Relation) == "AND" {
				fc.AndCriteria(v.Field, v.Operation, pv)
			}
			if strings.ToUpper(v.Relation) == "OR" {
				fc.OrCriteria(v.Field, v.Operation, pv)
			}
		}
	}
	return nil
}
func (c *IDSServiceHandler) getMetaData(metaid string) (*datasource.DataResultSet, error) {
	obj, err := datasource.CreateIDSFromName("default.mgr.G_META_ITEM")
	if err != nil {
		return nil, err
	}
	v, ok := obj.(*datasource.TableDataSource)
	if !ok {
		return nil, fmt.Errorf("获取默认数据源default.mgr.G_META_ITEM出错")
	}
	v.AddCriteria("META_ID", datasource.OPER_EQ, metaid)
	r, err := v.DoFilter()
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (c *IDSServiceHandler) getMetaID(project string, namespace string, metaname string) (string, error) {
	obj, err := datasource.CreateIDSFromName("default.mgr.G_META")
	if err != nil {
		return "", err
	}
	v, ok := obj.(*datasource.TableDataSource)
	if !ok {
		return "", fmt.Errorf("获取默认数据源default.mgr.G_META出错")
	}
	v.AddCriteria("NAMESPACE", datasource.OPER_EQ, namespace)
	v.AndCriteria("METANAME", datasource.OPER_EQ, metaname)
	if project != "" {
		v.AndCriteria("PROJECTID", datasource.OPER_EQ, project)
	}
	rs, err := v.DoFilter()
	if err != nil {
		return "", err
	}
	metaid := rs.Data[0][rs.Fields["ID"].Index].(string)
	return metaid, nil
}
func (c *IDSServiceHandler) doPostAction(dataSet *datasource.DataResultSet, rBody *ServiceRequestBody) (*datasource.DataResultSet, error) {
	if len(rBody.PostAction) == 0 {
		return dataSet, nil
	}
	var rdataset = dataSet
	for _, item := range rBody.PostAction {
		switch item.Name {
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
							m := make(map[string]string)
							item.Meta = &m
						}
						(*item.Meta)["CAP"] = o
					}
				}
			}
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
//处理查询报文
func (c *IDSServiceHandler) doQuery(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	if rBody == nil {
		c.doAllData(ids, rBody)
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
			c.createErrorResponseByError(err)
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
			/*	AGG_COUNT int = 1
				AGG_SUM   int = 2
				AGG_AVG   int = 3
				AGG_MAX   int = 4
				AGG_MIN   int = 5*/
			pred := strings.ToUpper(agg.Predicate)
			p := 0
			switch pred {
			case "COUNT":
				p = datasource.AGG_COUNT
			case "SUM":
				p = datasource.AGG_SUM
			case "AVG":
				p = datasource.AGG_AVG
			case "MAX":
				p = datasource.AGG_MAX
			case "MIN":
				p = datasource.AGG_MIN
			}
			ag.AddAggre(agg.Outfield, &datasource.AggreType{
				Predicate: p,
				ColName:   agg.ColName,
			})
		}
	}
	resuleset, err := fids.DoFilter()
	if err != nil {
		c.createErrorResult(err.Error())
	} else {
		resuleset, err = c.doPostAction(c.DoBulldozer(resuleset, rBody.Bulldozer), rBody)
		if err != nil {
			c.createErrorResponseByError(err)
			c.ServeJson()
		}
		c.setResultSet(resuleset)
	}

	c.ServeJson()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//返回服务元数据
func (c *IDSServiceHandler) doGetMeta(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	r := CreateRestResult(true)
	r["name"] = ids.GetName()
	r["type"] = datasource.GetDataSourceTypeStr(ids.GetDataSourceType())
	r["keyfields"] = ids.GetKeyFields()
	r["fields"] = ids.GetFields()
	imp := []string{"IDataSource"}
	if _, ok := ids.(datasource.ICriteriaDataSource); ok {
		imp = append(imp, "ICriteriaDataSource")
	}
	if _, ok := ids.(datasource.IFilterAdder); ok {
		imp = append(imp, "IFilterAdder")
	}
	if _, ok := ids.(datasource.IAggregativeAdder); ok {
		imp = append(imp, "IAggregativeAdder")
	}
	if _, ok := ids.(datasource.IWriteableDataSource); ok {
		imp = append(imp, "IWriteableDataSource")
	}
	r["implemention"] = imp
	c.Ctl.Data["json"] = r
	c.ServeJson()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回当前类支持的动作类型,以及动作对应的操作函数
func (c *IDSServiceHandler) getActionMap() map[string]SerivceActionHandler {
	return map[string]SerivceActionHandler{
		SrvAction_META:    c.doGetMeta,
		SrvAction_QUERY:   c.doQuery,
		SrvAction_DELETE:  c.doDelete,
		SrvAction_UPDATE:  c.doUpdate,
		SrvAction_INSERT:  c.doInsert,
		SrvAction_ALLDATA: c.doAllData,
		SrvAction_GET:     c.doGetValueByKey,
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据主键返回数据
func (c *IDSServiceHandler) doGetValueByKey(ids datasource.IDataSource, rBody *ServiceRequestBody) {
	fs := ids.GetKeyFields()
	params := make([]interface{}, len(fs), len(fs))
	for i, f := range fs {
		var err error
		params[i], err = c.ConvertString2Type(c.Ctl.Input().Get(f.Name), f.DataType)
		if err != nil {
			c.createErrorResponse("类型转换错误" + c.Ctl.Input().Get(f.Name) + " " + f.DataType + " err:" + err.Error())
			return
		}
	}
	resuleset, err := ids.QueryDataByKey(params...)
	if err != nil {
		c.createErrorResult(err.Error())
	} else {
		c.setResultSet(c.DoBulldozer(resuleset, rBody.Bulldozer))
	}
	c.ServeJson()

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回请求报文
func (c *IDSServiceHandler) getRBody() *ServiceRequestBody {
	var rBody *ServiceRequestBody
	if c.Ctl.Ctx.Request.Method == "POST" {
		rBody = &ServiceRequestBody{}
		err := json.Unmarshal([]byte(c.Ctl.Ctx.Input.RequestBody), rBody)
		if err != nil {
			c.createErrorResponse("解析报文时发生错误" + err.Error())
		}
	} else {
		rBody = nil
	}
	return rBody
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据元数据返回处理服务的接口
func (c *IDSServiceHandler) getServiceInterface(metestr string) (interface{}, error) {
	meta := make(map[string]interface{})
	err2 := json.Unmarshal([]byte(metestr), &meta)
	if err2 != nil {
		return nil, fmt.Errorf("meta信息不正确,应为JSON格式")
	}
	return datasource.CreateIDSFromName(meta["ids"].(string))
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理服务请求的入口
func (c *IDSServiceHandler) DoSrv(metestr string, inf ServiceHandlerInterface) {

	//////////////////////////////////////////////////////////////////////////
	//调用传入的接口中的方法实现下面的功能,因为需要通过不同的接口实现来实现不同的行为
	obj, err := inf.getServiceInterface(metestr)
	if err != nil {
		c.createErrorResponseByError(err)
		return
	}
	rBody := inf.getRBody()
	//////////////////////////////////////////////////////////////////////////

	ids, ok := obj.(datasource.IDataSource)
	if !ok {
		c.createErrorResponse("请求的服务没有实现IDataSource接口")
		return
	}
	act := c.Ctl.Ctx.Input.Param(":action")
	amap := inf.getActionMap()
	f, ok := amap[act]
	if !ok {
		c.createErrorResponse("请求的动作当前服务没有实现")
		return
	}
	f(ids, rBody)
}
