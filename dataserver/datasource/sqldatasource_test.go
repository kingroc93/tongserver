package datasource

import (
	"fmt"
	"testing"
)

func TestSQLDataSource_DoFilter(t *testing.T) {
	sqld := CreateSQLDataSource("JEDA_USER", "default",
		"select b.ORG_ID,b.ORG_NAME,a.USER_ID from JEDA_USER a inner join JEDA_ORG b on a.ORG_ID=b.ORG_ID and a.USER_ID=?")
	sqld.ParamsValues = []interface{}{"lvxing"}
	//sqld.AddCriteria("USER_ID", OperEq, "lvxing")
	rs, err := sqld.DoFilter()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printRS(rs)
}

func TestSQLDataSource_GetAllData(t *testing.T) {
	sqld := &SQLDataSource{
		DBDataSource: DBDataSource{
			DataSource: DataSource{
				Name: "JEDA_USER",
				Field: []*MyProperty{
					&MyProperty{
						Name: "ORG_ID",
					},
					&MyProperty{
						Name: "ORG_NAME",
					},
					&MyProperty{
						Name: "USER_ID",
					},
				},
			},
			DBAlias: "default",
		},
		SQL: "select b.ORG_ID,b.ORG_NAME,a.USER_ID from JEDA_USER a inner join JEDA_ORG b on a.ORG_ID=b.ORG_ID"}
	sqld.Init()
	rs, err := sqld.GetAllData()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	printRS(rs)
}
