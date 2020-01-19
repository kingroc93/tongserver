package datasource

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	//_ "github.com/mattn/go-oci8"
	"os"
	"testing"
)

var ids *TableDataSource

func TestMain(m *testing.M) {
	fmt.Println("Before ====================")
	err := orm.RegisterDataBase("default", "mysql", "tong:123456@tcp(127.0.0.1:3306)/idb", 30)
	err = orm.RegisterDataBase("pest", "mysql", "tong:123456@tcp(127.0.0.1:3306)/pest", 30)
	DBAlias2DBTypeContainer["default"] = "mysql"
	DBAlias2DBTypeContainer["pest"] = "mysql"
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	code := m.Run()
	fmt.Println("End ====================")
	os.Exit(code)
}

func createIDS() *TableDataSource {
	if ids == nil {
		ids = &TableDataSource{
			DBDataSource: DBDataSource{
				DataSource: DataSource{
					Name: "JEDA_USER",
					Field: []*MyProperty{
						&MyProperty{
							Name: "ORG_ID",
						},
						&MyProperty{
							Name:    "ORG_NAME",
							OutJoin: true,
							OutJoinDefine: &OutFieldProperty{
								Source:     createIDSOrg(),
								ValueField: "ORG_NAME",
								JoinField:  "ORG_ID",
								ValueFunc:  nil,
							},
						},
						&MyProperty{
							Name: "USER_ID",
						},
					},
				},
				DBAlias: "default",
			},
			TableName: "JEDA_USER",
		}
		ids.Init()
	}
	return ids
}
func createIDSOrg() *TableDataSource {
	idso := &TableDataSource{
		DBDataSource: DBDataSource{
			DataSource: DataSource{
				Name: "JEDA_ORG",
				Field: []*MyProperty{
					&MyProperty{
						Name: "ORG_ID",
					},
					&MyProperty{
						Name: "ORG_NAME",
					},
				},
			},
			DBAlias: "default",
		},
		TableName: "JEDA_ORG",
	}
	idso.Init()
	return idso
}

func createIDSRiver() *TableDataSource {
	if ids == nil {
		ids = &TableDataSource{
			DBDataSource: DBDataSource{
				DataSource: DataSource{
					Name: "default",
				},
				DBAlias: "default",
			},
			TableName: "ST_RIVER_R",
		}
		ids.Init()
	}
	return ids
}

func printRS(rs *DataResultSet) {
	fmt.Println("=======================================================================")
	for k, v := range rs.Fields {
		fmt.Printf("%s:%s:%s\n", k, v.FieldType, strconv.Itoa(v.Index))
	}
	fmt.Println("=======================================================================")
	for _, row := range rs.Data {
		for _, item := range row {
			fmt.Printf("%s\t", item)
		}
		fmt.Println()
	}
}

func TestRiverIDS(t *testing.T) {
	ids := createIDSRiver()
	ids.RowsLimit = 1000
	ids.AddCriteria("Z", OperLt, 5.00)
	ids.AddAggre("CNT", &AggreType{
		Predicate: AggCount,
		ColName:   "Z",
	})
	rs, _ := ids.DoFilter()
	printRS(rs)
}

func TestTableDataSourceGetAllData(t *testing.T) {
	ids := createIDS()
	rs, _ := ids.GetAllData()
	printRS(rs)
	ids.RowsLimit = 10
	rs, _ = ids.GetAllData()

	fmt.Println("===============================================================================================")
	printRS(rs)

	data, err := json.Marshal(ids)
	fmt.Println(err)
	fmt.Println(string(data))
}

func TestAddCriteria(t *testing.T) {
	ids := createIDS()
	//	var inf interface{}
	//	inf = ids
	ids.AddCriteria("ORG_ID", OperEq, "001031")
	rs, err := ids.DoFilter()
	if err != nil {
		fmt.Print(err)
	}
	printRS(rs)
}
