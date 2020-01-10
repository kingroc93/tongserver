package cube

import (
	"awesome/datasource"
	"fmt"
	"time"
)

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//数据集后处理句柄，返回新的结果集
type DataSetBulldozerOperator func(dataset *datasource.DataResultSet, params map[string]interface{})

//数据集后处理行句柄，不返回新的结果集
type HandlerFunc func(dataset *datasource.DataResultSet, index int, params map[string]interface{})

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//格式化数据集中的数据，格式化以后的数据类型均为string，FieldDesc中的数据类型不变
func FormatDatafunc(dataset *datasource.DataResultSet, index int, params map[string]interface{}) {
	for k, v := range params {
		if dataset.Fields[k] == nil {
			continue
		}
		obj := dataset.Data[index][dataset.Fields[k].Index]
		switch obj.(type) {
		case time.Time:
			dataset.Data[index][dataset.Fields[k].Index] = obj.(time.Time).Format(v.(string))
		default:
			dataset.Data[index][dataset.Fields[k].Index] = fmt.Sprintf(v.(string), obj)
		}

	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 根据隐藏的列返回显示的列
func CreateShowFieldList(dataset *datasource.DataResultSet, hidden []string) []string {
	fd := make([]string, 1, 1)
	for k, _ := range dataset.Fields {
		inhidden := false
		for _, i := range hidden {
			if i == k {
				inhidden = true
				break
			}
		}
		if !inhidden {
			fd = append(fd, k)
		}
	}
	return fd
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 列过滤,该函数必须在最后一个，否则会造成错误
func ColumnFilterFunc(dataset *datasource.DataResultSet, index int, params map[string]interface{}) {
	if len(dataset.Data) == 0 || index >= len(dataset.Data) {
		return
	}
	if params["show"] != nil {
		d := make([]interface{}, 0, len(dataset.Data[0]))
		for _, v := range params["show"].([]interface{}) {
			d = append(d, dataset.Data[index][dataset.Fields[v.(string)].Index])
		}
		dataset.Data[index] = d
		if index == len(dataset.Data)-1 {
			fd := make(map[string]*datasource.FieldDesc)
			for i, k := range params["show"].([]interface{}) {
				fd[k.(string)] = &datasource.FieldDesc{
					Index:     i,
					FieldType: dataset.Fields[k.(string)].FieldType,
				}
			}
			dataset.Fields = fd
		}
		return
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//在原结果集的基础上将字典的映射值添加到结果集中，返回新的结果集
func DictMappingfunc(dataset *datasource.DataResultSet, index int, params map[string]interface{}) {
	if len(dataset.Data) == 0 || index >= len(dataset.Data) {
		return
	}
	nfield := params["outfield"].(string)
	dataKeyField := params["dataKeyField"].(string)
	ks := params["KeyStringSource"].(*datasource.KeyStringSource)
	if nfield == "" || dataKeyField == "" || ks == nil {
		return
	}
	f, ok := dataset.Fields[dataKeyField]
	if !ok {
		return
	}
	v := dataset.Data[index][f.Index]
	dataset.Data[index] = append(dataset.Data[index], ks.GetDataByKey(fmt.Sprint(v)))

	dataset.Fields[nfield] = &datasource.FieldDesc{
		Index: len(dataset.Data[index]) - 1,
	}
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//数据集转置
func Slice(dataset *datasource.DataResultSet, xfieldname string, yfieldnames []string, valuefields ...string) *datasource.DataResultSet {
	fieldlist := make(map[interface{}]int)
	temp := make(map[interface{}]([][]interface{}))
	xindex := 0
	for _, row := range dataset.Data {
		y := row[dataset.Fields[yfieldnames[0]].Index]
		if y == nil {
			continue
		}
		list := temp[y]
		if list == nil {
			list = make([][]interface{}, 0, 10)
		}
		list = append(list, row)
		temp[y] = list
		x := row[dataset.Fields[xfieldname].Index]
		if _, ok := fieldlist[x]; !ok {
			fieldlist[x] = xindex
			xindex++
		}
	}
	table := make([][]interface{}, len(temp), len(temp))
	yindex := 0
	vL := len(valuefields)
	yfLen := len(yfieldnames)
	for _, v := range temp {
		table[yindex] = make([]interface{}, xindex+yfLen, xindex+yfLen)
		for _, item := range v {
			for i, fv := range yfieldnames {
				table[yindex][i] = item[dataset.Fields[fv].Index]
			}
			if vL == 0 {
				table[yindex][yfLen+fieldlist[item[dataset.Fields[xfieldname].Index]]] = item
			} else if vL == 1 {
				table[yindex][yfLen+fieldlist[item[dataset.Fields[xfieldname].Index]]] = item[dataset.Fields[valuefields[0]].Index]
			} else {
				vt := make([]interface{}, vL, vL)
				for i, fs := range valuefields {
					vt[i] = item[dataset.Fields[fs].Index]
				}
				table[yindex][yfLen+fieldlist[item[dataset.Fields[xfieldname].Index]]] = vt
			}
		}
		yindex++
	}
	var result = &datasource.DataResultSet{}
	result.Data = table
	result.Fields = make(map[string]*datasource.FieldDesc)
	//result.Fields[yfieldname] = &DataSource.FieldDesc{Index: 0}
	for i, fv := range yfieldnames {
		result.Fields[fv] = &datasource.FieldDesc{Index: i}
	}

	for k, v := range fieldlist {

		result.Fields[fmt.Sprintf("F_%v", k)] = &datasource.FieldDesc{
			Index: v + len(yfieldnames),
		}
	}
	return result
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
//数据集分组
func GroupField(data *datasource.DataResultSet, fieldname string) map[interface{}](*datasource.DataResultSet) {
	temp := make(map[interface{}](*datasource.DataResultSet))
	for _, row := range data.Data {
		value := row[data.Fields[fieldname].Index]
		if value == nil {
			continue
		}
		list := temp[value]
		if list == nil {
			list = &datasource.DataResultSet{}
			list.Fields=data.Fields.Copy()
		}
		list.Data = append(list.Data, row)
		temp[value] = list
	}
	return temp
}

//func GroupField(data *DataSource.DataResultSet, fieldname string) map[interface{}]([][]interface{}) {
//	temp := make(map[interface{}][][]interface{})
//	for _, row := range data.Data {
//		value := row[data.Fields[fieldname].Index]
//		if value == nil {
//			continue
//		}
//		list := temp[value]
//		if list == nil {
//			list = make([][]interface{}, 0, 10)
//		}
//		list = append(list, row)
//		temp[value] = list
//	}
//	return temp
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// 提取列
func Row2Colume(dataset *datasource.DataResultSet, fieldsname ...string) *datasource.DataResultSet {
	if len(fieldsname) == 0 {
		return dataset
	}
	var result = &datasource.DataResultSet{}
	result.Fields = make(map[string]*datasource.FieldDesc)
	for i, v := range fieldsname {
		result.Fields[v] = &datasource.FieldDesc{
			FieldType: v,
			Index:     i,
		}
	}
	//Data   [][]interface{}
	result.Data = make([][]interface{}, len(fieldsname), len(fieldsname))
	for i, _ := range result.Data {
		result.Data[i] = make([]interface{}, len(dataset.Data), len(dataset.Data))
	}
	for j, item := range dataset.Data {
		for i, v := range fieldsname {
			fd, ok := dataset.Fields[v]
			if !ok {
				continue
			}
			result.Data[i][j] = item[fd.Index]
		}
	}
	return result
}
