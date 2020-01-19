package datasource

const (
	// InnerJoin 内连接
	InnerJoin string = "INNER"
	// LeftJoin 左链接
	LeftJoin string = "LEFT"
	// RightJoin 右链接
	RightJoin string = "RIGHT"
	// FullJoin 全链接
	FullJoin string = "FULL"
)

// CompositeDataSourceItem 复合数据项定义
type CompositeDataSourceItem struct {
	BaseCriteria
	Source     IDataSource
	OutField   []string
	JoinMethod string
}

// CompositeTableDataSource 复合数据表数据源
type CompositeTableDataSource struct {
	//BaseCriteria
	DataSource
	TableName        string
	joinedDataSource []*CompositeDataSourceItem
}

// GetDataSourceType 返回数据源类型
func (c *CompositeTableDataSource) GetDataSourceType() DSType {
	return DataSourceTypeInner
}

// GetAllData 返回全部数据
func (c *CompositeTableDataSource) GetAllData() (*DataResultSet, error) {
	panic("implement me")
}

// QueryDataByKey 根据主键返回数据
func (c *CompositeTableDataSource) QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error) {
	panic("implement me")
}

// QueryDataByFieldValues 根据字段值返回数据
func (c *CompositeTableDataSource) QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error) {
	panic("implement me")
}

// Init 初始化
func (c *CompositeTableDataSource) Init() error {
	c.joinedDataSource = make([]*CompositeDataSourceItem, 0, 10)
	return nil
}

// JoinDataSource 联接其他数据源
func (c *CompositeTableDataSource) JoinDataSource(tableSource *TableDataSource, OutField []string, JoinMethod string) *CompositeDataSourceItem {
	item := &CompositeDataSourceItem{
		Source:     tableSource,
		OutField:   OutField,
		JoinMethod: JoinMethod,
	}
	c.joinedDataSource = append(c.joinedDataSource, item)
	return item
}
