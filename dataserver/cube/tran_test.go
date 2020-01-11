package cube

import (
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strconv"
	"testing"
	"time"
	"tongserver.dataserver/datasource"
)

func TestMain(m *testing.M) {
	fmt.Println("Before ====================")
	err := orm.RegisterDataBase("default", "mysql", "tong:123456@tcp(127.0.0.1:3306)/idb", 30)
	err = orm.RegisterDataBase("pest", "mysql", "tong:123456@tcp(127.0.0.1:3306)/pest", 30)
	datasource.DBAlias2DBTypeContainer["default"] = "mysql"
	datasource.DBAlias2DBTypeContainer["pest"] = "mysql"
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	code := m.Run()
	fmt.Println("End ====================")
	os.Exit(code)
}
func printRS(rs *datasource.DataResultSet) {
	fmt.Println("=======================================================================")
	for k, v := range rs.Fields {
		fmt.Printf("%s:%s:%s\n", k, v.FieldType, strconv.Itoa(v.Index))
	}
	fmt.Println("=======================================================================")
	for _, row := range rs.Data {
		for i, item := range row {
			fmt.Printf("%d-%s\t", i, item)
		}
		fmt.Println()
	}
}

func TestDictMapping(t *testing.T) {
	users := datasource.CreateTableDataSource("JEDA_USER", "default", "JEDA_USER")
	org := datasource.CreateTableDataSource("JEDA_ORG", "default", "JEDA_ORG")
	ks := &datasource.KeyStringSource{
		DataSource: datasource.DataSource{
			Name: "ORGLIST",
		},
	}
	ks.Init()
	ks.FillDataByDataSource(org, "ORG_ID", "ORG_NAME")
	//func DictMappingfunc(dataset *DataSource.DataResultSet,index int,params map[string]interface{}){
	data, _ := users.GetAllData()
	L := len(data.Data)

	for i := 0; i < L; i++ {
		DictMappingfunc(data, i, map[string]interface{}{
			"outfield":        "ORG_NAME",
			"dataKeyField":    "ORG_ID",
			"KeyStringSource": ks,
		})
		FormatDatafunc(data, i, map[string]interface{}{
			"USER_CREATED": "2006-01-02",
		})
		ColumnFilterFunc(data, i, map[string]interface{}{
			"show": []interface{}{"USER_ID", "USER_NAME", "ORG_ID", "ORG_NAME", "USER_CREATED"},
		})

	}
	data = Row2Colume(data, "USER_ID", "USER_NAME", "ORG_NAME")
	printRS(data)
}
func TestGroupField(t *testing.T) {
	datatable := datasource.CreateTableDataSource("ST_RIVER_R", "default", "ST_RIVER_R")
	datatable.RowsLimit = 1000
	datatable.Orderby("tm", "desc")
	rs, err := datatable.DoFilter()
	if err != nil {
		t.Fatalf("获取数据时发生错误，%s", err.Error())
	}
	rss := GroupField(rs, "STCD")
	for _, item := range rss.Data {
		fmt.Println("*****************************")
		fmt.Println(item[0])
		fmt.Println("*****************************")
		printRS(item[1].(*datasource.DataResultSet))
	}
}
func TestCompositeDataSource(t *testing.T) {
	datatable := datasource.CreateTableDataSource("iot_data_bas", "pest", "iot_data_bas")
	datatable.Field = datasource.CreateFieldNonType("DEV_ID", "ITEM_ID", "SITE_ID", "DATA_VALUE", "COLLECT_DATE")
	timestamp1, _ := time.Parse("2006-01-02", "2019-11-13")
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02", "2019-11-13", loc)

	fmt.Println(theTime)
	fmt.Println(timestamp1)

	datatable.AddCriteria("batch_time", datasource.OPER_EQ, timestamp1)
	rs, _ := datatable.DoFilter()
	//rs,_=monitemtable.GetAllData()
	printRS(rs)

	fmt.Println("===============================================================================================")
	fmt.Println("===============================================================================================")

	rss := Slice(rs, "ITEM_ID", []string{"DEV_ID", "SITE_ID", "COLLECT_DATE"}, "DATA_VALUE")
	printRS(rss)

}
