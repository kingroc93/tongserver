package service

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"tongserver.dataserver/cube"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/utils"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
const (
	//返回全部数据
	SrvActionALLDATA string = "all"
	//查询动作
	SrvActionQUERY string = "query"
	//根据主键返回
	SrvActionGET string = "get"
	//返回缓存
	SrvActionCACHE string = "cache"
	//根据字段值返回
	SrvActionBYFIELD string = "byfield"
	//返回服务元数据
	SrvActionMETA string = "meta"
	//删除操作
	SrvActionDELETE string = "delete"
	//更新操作
	SrvActionUPDATE string = "update"
	//插入操作
	SrvActionINSERT string = "insert"

	//以下三个常量均为通过QueryString传入的参数名
	//针对查询自动分页中每页记录数
	RequestParamPagesize string = "_pagesize"
	//针对查询自动分页中的页索引
	RequestParamPageindex string = "_pageindex"
	//是否返回字段元数据，默认为返回
	RequestParamNofieldsinfo string = "_nofield"
	//当前请求不执行而是只返回SQL语句，仅针对IDS类型的服务有效
	RequestParamSQL string = "_sql"
	//当前请求的响应信息不直接返回
	//该参数只对query、all两个操作起作用
	RequestParamCache      string = "_cache"
	RequestParamCachebykey string = "_cachekey"
)

// IDSServiceHandler 数据源服务处理句柄
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
func (c *IDSServiceHandler) doDelete(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
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
			c.ServeJSON()
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//处理更新
func (c *IDSServiceHandler) doUpdate(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
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
			c.ServeJSON()
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 处理添加
func (c *IDSServiceHandler) doInsert(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
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
			c.ServeJSON()
		}
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
//返回所有数据
func (c *IDSServiceHandler) doAllData(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
	c.setPageParams(ids)
	resuleset, err := ids.GetAllData()
	if err != nil {
		c.createErrorResult(err.Error())
		return
	}
	if rBody == nil {
		c.setResultSet(resuleset)
		c.ServeJSON()
	} else {
		resuleset, err = c.doPostAction(c.DoBulldozer(resuleset, rBody.Bulldozer), rBody)
		if err != nil {
			c.createErrorResponseByError(err)
			return
		}
		c.setResultSet(resuleset)
		c.ServeJSON()
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 参数转换，将string的参数转换为指定的类型，针对日期类型特殊处理：
// 如前N天，lastday:1    lastday:-3
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
	v, ok := obj.(*datasource.TableDataSource)
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
	v, ok := obj.(*datasource.TableDataSource)
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
							m := make(map[string]string)
							item.Meta = &m
						}
						(*item.Meta)["CAP"] = o
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
func (c *IDSServiceHandler) doQuery(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
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
		c.createErrorResult(err.Error())
	} else {
		resuleset, err = c.doPostAction(c.DoBulldozer(resuleset, rBody.Bulldozer), rBody)
		if err != nil {
			c.createErrorResponseByError(err)
			c.ServeJSON()
		}
		c.setResultSet(resuleset)
	}

	c.ServeJSON()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回服务元数据
func (c *IDSServiceHandler) doGetMeta(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
	r := CreateRestResult(true)
	sd := make(map[string]interface{})
	r["servicedefine"] = sd
	sd["Context"] = sdef.Context
	sd["BodyType"] = sdef.BodyType
	sd["ServiceType"] = sdef.ServiceType
	sd["Namespace"] = sdef.Namespace
	sd["Enabled"] = sdef.Enabled
	sd["MsgLog"] = sdef.MsgLog
	sd["Security"] = sdef.Security
	meta := make(map[string]interface{})
	err2 := json.Unmarshal([]byte(sdef.Meta), &meta)
	if err2 == nil {
		sd["Meta"] = meta
	} else {
		sd["Meta"] = sdef.Meta
	}

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
	r["ids"] = imp

	c.Ctl.Data["json"] = r
	c.ServeJSON()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回缓存的结果数据
func (c *IDSServiceHandler) doGetCache(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
	key := c.Ctl.Input().Get(RequestParamCachebykey)
	if key == "" {
		r := CreateRestResult(false)
		r["msg"] = RequestParamCachebykey + "不得为空"
		c.Ctl.Data["json"] = r
		c.ServeJSON()
	}
	obj := utils.DataSetResultCache.Get(key)
	if obj == nil {
		r := CreateRestResult(false)
		r["msg"] = "没有找到请求的缓存信息"
		c.Ctl.Data["json"] = r
		c.ServeJSON()
		return
	}
	r, ok := obj.(RestResult)
	if !ok {
		r := CreateRestResult(false)
		r["msg"] = "缓存对象类型非法"
		c.Ctl.Data["json"] = r
		c.ServeJSON()
		return
	}
	times := r["cachetimes"].(int)
	d := r["duration"].(int)
	if times > 0 {
		times = times - 1
	}
	r["cachetimes"] = times
	if times == 0 {
		err := utils.DataSetResultCache.Delete(key)
		if err != nil {
			r["result"] = false
			r["msg"] = "删除缓存时发生错误：" + err.Error()
			c.Ctl.Data["json"] = r
			return
		}
	} else {

		err := utils.DataSetResultCache.Put(key, obj, time.Duration(d)*time.Second)
		if err != nil {
			r["result"] = false
			r["msg"] = "加入缓存时发生错误：" + err.Error()
			c.Ctl.Data["json"] = r
			return
		}
	}

	dst, ok := r["resultset"]
	if ok {
		ds, ok := dst.(*datasource.DataResultSet)
		if ok {
			resuleset, err := c.doPostAction(c.DoBulldozer(ds, rBody.Bulldozer), rBody)
			if err != nil {
				c.createErrorResponseByError(err)
				c.ServeJSON()
			}
			c.setResultSet(resuleset)
			c.ServeJSON()
			return
		}
	}
	r["result"] = true
	c.Ctl.Data["json"] = r
	c.ServeJSON()
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回当前类支持的动作类型,以及动作对应的操作函数
func (c *IDSServiceHandler) getActionMap() map[string]SerivceActionHandler {
	return map[string]SerivceActionHandler{
		SrvActionMETA:    c.doGetMeta,
		SrvActionQUERY:   c.doQuery,
		SrvActionDELETE:  c.doDelete,
		SrvActionUPDATE:  c.doUpdate,
		SrvActionINSERT:  c.doInsert,
		SrvActionALLDATA: c.doAllData,
		SrvActionGET:     c.doGetValueByKey,
		SrvActionCACHE:   c.doGetCache,
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据主键返回数据
func (c *IDSServiceHandler) doGetValueByKey(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
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
	c.ServeJSON()

}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 返回请求报文
func (c *IDSServiceHandler) getRBody() *SRequestBody {
	var rBody *SRequestBody
	if c.Ctl.Ctx.Request.Method == "POST" {
		rBody = &SRequestBody{}
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

// DoSrv 处理服务请求的入口
func (c *IDSServiceHandler) DoSrv(sdef *SDefine, inf SHandlerInterface) {
	metestr := sdef.Meta
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
	f(sdef, ids, rBody)
}
