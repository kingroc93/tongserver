package datasource

import (
	"fmt"
	"testing"
)

func TestEnmuSource(t *testing.T) {
	datatable := CreateTableDataSource("JEDA_ORG", "default", "JEDA_ORG")
	ks := &KeyStringSource{
		DataSource: DataSource{
			Name: "ORGLIST",
		},
	}
	ks.Init()
	ks.FillDataByDataSource(datatable, "ORG_ID", "ORG_NAME")
	rs, err := ks.GetAllData()
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	printRS(rs)

}
