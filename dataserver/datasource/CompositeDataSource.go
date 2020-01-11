package datasource

/////////////////////////////////////////////////////////////////////////////////////////////////////////////

const (
	//内连接
	INNER_JOIN string = "INNER"
	//左链接
	LEFT_JOIN string = "LEFT"
	//右链接
	RIGHT_JOIN string = "RIGHT"
	//全链接
	FULL_JOIN string = "FULL"
)

//李磊  63202576 522
type CompositeDataSourceItem struct {
	BaseCriteria
	Source     IDataSource
	OutField   []string
	JoinMethod string
}

type CompositeTableDataSource struct {
	//BaseCriteria
	DataSource
	TableName        string
	joinedDataSource []*CompositeDataSourceItem
}

func (c *CompositeTableDataSource) GetDataSourceType() DataSourceType {
	return DataSourceType_INNER
}

func (c *CompositeTableDataSource) GetAllData() (*DataResultSet, error) {
	panic("implement me")
}

func (c *CompositeTableDataSource) QueryDataByKey(keyvalues ...interface{}) (*DataResultSet, error) {
	panic("implement me")
}

func (c *CompositeTableDataSource) QueryDataByFieldValues(fv *map[string]interface{}) (*DataResultSet, error) {
	panic("implement me")
}

func (c *CompositeTableDataSource) Init() error {
	c.joinedDataSource = make([]*CompositeDataSourceItem, 0, 10)
	return nil
}
func (c *CompositeTableDataSource) JoinDataSource(tableSource *TableDataSource, OutField []string, JoinMethod string) *CompositeDataSourceItem {
	item := &CompositeDataSourceItem{
		Source:     tableSource,
		OutField:   OutField,
		JoinMethod: JoinMethod,
	}
	c.joinedDataSource = append(c.joinedDataSource, item)
	return item
}
