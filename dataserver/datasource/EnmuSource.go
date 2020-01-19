package datasource

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

// KeyStringSource Key-String类型的数据源
type KeyStringSource struct {
	DataSource
	fields   FieldDescType
	valueMap map[string]string
}

// SetRowsLimit 该数据源此方法无意义
func (c *KeyStringSource) SetRowsLimit(limit int) {

}

// SetRowsOffset 该数据源此方法无意义
func (c *KeyStringSource) SetRowsOffset(offset int) {

}

// GetKeyFields 返回键字段
func (c *KeyStringSource) GetKeyFields() []*MyProperty {
	return c.KeyField
}

// GetDataSourceType 返回数据源类型
func (c *KeyStringSource) GetDataSourceType() DSType {
	return DataSourceTypeEnmu
}

// FillDataByDataSource 通过其他数据源填充
func (c *KeyStringSource) FillDataByDataSource(source IDataSource, keyfield string, valuefield string) {
	ds, err := source.GetAllData()
	if err != nil {
		logs.Error("获取数据时出错！" + err.Error())
		return
	}
	for _, r := range ds.Data {

		c.valueMap[fmt.Sprint(r[ds.Fields[keyfield].Index])] = fmt.Sprint(r[ds.Fields[valuefield].Index])
	}
}

// Init 数据源初始化
func (c *KeyStringSource) Init() error {
	c.Field = []*MyProperty{
		{
			Name:          "KEY",
			DataType:      PropertyDatatypeStr,
			OutJoin:       false,
			Caption:       "KEY",
			OutJoinDefine: nil,
		},
		{
			Name:          "VALUE",
			DataType:      PropertyDatatypeStr,
			OutJoin:       false,
			Caption:       "KEY",
			OutJoinDefine: nil,
		},
	}
	c.KeyField = []*MyProperty{
		{
			Name:     "KEY",
			DataType: PropertyDatatypeStr,
		},
	}
	c.fields = make(FieldDescType)
	c.fields["KEY"] = &FieldDesc{
		Index:     0,
		FieldType: PropertyDatatypeStr,
	}
	c.fields["VALUE"] = &FieldDesc{
		Index:     1,
		FieldType: PropertyDatatypeStr,
	}
	c.valueMap = make(map[string]string)
	return nil
}

// GetAllData 返回所有数据
func (c *KeyStringSource) GetAllData() (*DataResultSet, error) {
	var result = &DataResultSet{}
	result.Fields = c.fields
	result.Data = make([][]interface{}, len(c.valueMap), len(c.valueMap))
	i := 0
	for k, v := range c.valueMap {
		item := make([]interface{}, 2, 2)
		item[0] = k
		item[1] = v
		result.Data[i] = item
		i++
	}
	return result, nil
}

// SetValueMap 设定key-string数据
func (c *KeyStringSource) SetValueMap(v map[string]string) {
	c.valueMap = v
}

// GetValueMap 返回key-string数据
func (c *KeyStringSource) GetValueMap() map[string]string {
	return c.valueMap
}

// GetDataByKey 根据键值返回数据
func (c *KeyStringSource) GetDataByKey(key string) string {
	return c.valueMap[key]
}

// QueryDataByKey 根据键值返回数据
func (c *KeyStringSource) QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error) {
	var result = &DataResultSet{}
	result.Fields = c.fields
	result.Data = make([][]interface{}, 1, 1)
	result.Data[0] = make([]interface{}, 2, 2)
	result.Data[0][0] = keyvalues[0]
	result.Data[0][1] = c.valueMap[keyvalues[0].(string)]
	return result, nil
}

// QueryDataByFieldValues 该数据源此方法无意义
func (c *KeyStringSource) QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error) {
	return nil, fmt.Errorf("Use QueryDataByKey")
}
