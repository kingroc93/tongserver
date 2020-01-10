package datasource

import "time"

//时间序列数据表
type TMTableDataSource struct {
	TableDataSource
	//时间字段
	TMField string
}

//添加条件，最近days天
func (c TMTableDataSource) LastDays(days int) (*DataResultSet, error) {
	return nil, nil
}

//根据GroupFields分组得到每个分组的最新一条数据
func (c TMTableDataSource) LastedData(days int, GroupFields []string) (*DataResultSet, error) {
	return nil, nil
}

//返回一个时间段的数据
func (c TMTableDataSource) QueryPeriod(sDateTime time.Time, eDataTime time.Time) (*DataResultSet, error) {
	return nil, nil
}

//查询本月数据
func (c TMTableDataSource) QueryThisMonth() (*DataResultSet, error) {
	return nil, nil
}

//查询当日数据
func (c TMTableDataSource) QueryToday() (*DataResultSet, error) {
	return nil, nil
}

//查询当年数据
func (c TMTableDataSource) QueryThisYeay() (*DataResultSet, error) {
	return nil, nil
}
