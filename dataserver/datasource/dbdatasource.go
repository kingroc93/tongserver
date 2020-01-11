package datasource

import (
	"database/sql"
	"fmt"
	"github.com/astaxie/beego/logs"
	"tongserver.dataserver/utils"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
type DBDataSource struct {
	DataSource
	TableDataSourceCriteria

	DBAlias        string
	RowsLimit      int
	RowsOffset     int
	AutoFillFields bool

	openedDB *sql.DB `json:"-"`
	palesql  bool
}

/////////////////////////////////////////////////////////////////////////////////////////////////////////////
// DataSource
func (c *DataSource) Init() {
	panic("")
}
func (c *DataSource) convertPropertys2Cols(ps []*MyProperty) []string {
	L := 0
	for _, v := range ps {
		if !v.OutJoin {
			L++
		}
	}
	result := make([]string, L, L)
	L = 0
	for _, v := range ps {
		if !v.OutJoin {
			result[L] = v.Name
			L++
		}
	}
	return result
}
func (c *DataSource) GetFields() []*MyProperty {
	return c.Field
}
func (c *DataSource) GetDataSourceType() DataSourceType {
	panic("")
}
func (c *DataSource) GetName() string {
	return c.Name
}
func (c *DataSource) GetKeyFieldByName(name string) *MyProperty {
	for _, v := range c.KeyField {
		if v.Name == name {
			return v
		}
	}
	return nil
}
func (c *DataSource) GetFieldByName(name string) *MyProperty {
	for _, v := range c.Field {
		if v.Name == name {
			return v
		}
	}
	return nil
}
func (c *DBDataSource) SetRowsLimit(limit int) {
	c.RowsLimit = limit
}
func (c *DBDataSource) SetRowsOffset(offset int) {
	c.RowsOffset = offset
}

func (c *DBDataSource) GetKeyFields() []*MyProperty {
	return c.KeyField
}
func (c *DBDataSource) convertData(value interface{}, fieldType string) interface{} {
	var str utils.String
	switch v := value.(type) {
	case []byte:
		str = utils.String(v)
	case string:
		str = utils.String(v)
	default:
		str = utils.String(fmt.Sprintf("%v", v))
	}

	var item interface{}
	switch fieldType {
	case "VARCHAR", "NVARCHAR", "CHAR", "SQLT_AFC", "SQLT_CHR", "SQLT_VCS":
		item = str.String()
	case "INT", "MEDIUMINT", "INTEGER", "SQLT_INT":
		item, _ = str.Int32()
	case "TINYINT":
		item, _ = str.Int8()
	case "SMALLINT":
		item, _ = str.Int16()
	case "BIGINT":
		item, _ = str.Int64()
	case "FLOAT", "SQLT_FLT", "SQLT_BFLOAT":
		item, _ = str.Float32()
	case "DOUBLE", "SQLT_BDOUBLE":
		item, _ = str.Float64()
	case "TIMESTAMP", "SQLT_TIMESTAMP", "SQLT_TIMESTAMP_TZ", "SQLT_TIMESTAMP_LTZ":
		item, _ = str.DateTime()
	case "DATE", "SQLT_DAT":
		item, _ = str.Date()
	case "DATETIME":
		item, _ = str.DateTime()
	case "TIME":
		item, _ = str.DateTime()
	default:
		item = str.String()
	}
	return item
}
func (c *DBDataSource) getRecordByRef(refs []interface{}, cols []string, colsTypes *FieldDescType) ([]interface{}, []*MyProperty) {
	if c.Field == nil || len(c.Field) == 0 || c.palesql {
		item := make([]interface{}, len(cols), len(cols))
		for i, fieldname := range cols {
			item[i] = c.convertData(*refs[i].(*interface{}), (*colsTypes)[fieldname].FieldType)
		}
		return item, nil
	} else {
		item := make([]interface{}, len(c.Field), len(c.Field))
		Oj := make([]*MyProperty, 0, len(c.Field))
		for i, v := range c.Field {
			if !v.OutJoin {
				item[i] = c.convertData(*refs[(*colsTypes)[v.Name].Index].(*interface{}), (*colsTypes)[v.Name].FieldType)
			} else {
				Oj = append(Oj, v)
			}
		}
		return item, Oj
	}
}
func (c *DBDataSource) querySQLData(sqlstr string, params ...interface{}) (*DataResultSet, error) {
	var err error
	logs.Debug(sqlstr)
	if c.openedDB == nil {
		return nil, fmt.Errorf("OpenedDB is nil")
	}
	rs, err := c.openedDB.Query(sqlstr, params...) //获取所有数据

	if err != nil {
		return nil, err
	}
	defer rs.Close()
	cols, err := rs.Columns()
	if err != nil {
		return nil, err
	}
	colsTypes, err := rs.ColumnTypes()
	if err != nil {
		return nil, err
	}
	var result = &DataResultSet{}
	fm := make(FieldDescType)
	for i, item := range cols {
		fm[item] = &FieldDesc{
			FieldType: ConvertMySQLType2CommonType(colsTypes[i].DatabaseTypeName()),
			Index:     i,
		}
	}
	refs := make([]interface{}, len(cols))
	for i := range refs {
		var ref interface{}
		refs[i] = &ref
	}
	result.Fields = make(FieldDescType)
	if c.palesql || c.Field == nil {
		result.Fields = fm
	} else {
		for index, item := range c.Field {
			var typ string
			if fm[item.Name] != nil {
				typ = fm[item.Name].FieldType
			} else {
				typ = item.DataType
			}
			result.Fields[item.Name] = &FieldDesc{
				FieldType: typ,
				Index:     index,
			}
		}
	}
	datas := make([][]interface{}, 0, 100)
	for rs.Next() {
		err := rs.Scan(refs...)
		if err != nil {
			return nil, err
		}
		item, Ofs := c.getRecordByRef(refs, cols, &fm)
		if Ofs != nil {
			//存在通过Join加载其他数据源的字段
			for _, f := range Ofs {
				if f.OutJoinDefine == nil {
					continue
				}
				kv := item[result.Fields[f.OutJoinDefine.JoinField].Index]
				rfs, err := f.OutJoinDefine.Source.QueryDataByKey(kv)
				if err != nil {
					logs.Error(err)
					continue
				}
				if len(rfs.Data) == 0 {
					continue
				}
				if rfs.Fields[f.OutJoinDefine.ValueField] == nil {
					logs.Error("ValueField错误没有找到字段" + f.OutJoinDefine.ValueField)
					continue
				}
				item[result.Fields[f.Name].Index] = rfs.Data[0][rfs.Fields[f.OutJoinDefine.ValueField].Index]
			}
		}
		datas = append(datas, item)
	}
	result.Data = datas

	return result, nil
}

//
//
/////////////////////////////////////////////////////////////////////////////////////////////////////////////
