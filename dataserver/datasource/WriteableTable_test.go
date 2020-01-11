package datasource

import (
	"fmt"
	"testing"
)

func TestCreateInsertSQLService(t *testing.T) {
	ids := &WriteableTableSource{
		TableDataSource{
			DBDataSource: DBDataSource{
				DataSource: DataSource{
					Name: "G_SERCVICE",
				},
				DBAlias:        "default",
				AutoFillFields: true,
			},
			TableName: "G_SERCVICE",
		}}
	ids.Init()
	err := ids.Insert(map[string]interface{}{
		"ID":        "05442082-76d9-41da-b563-a19914131993",
		"DBTYPE":    "mysql",
		"DBURL":     "{username}:{password}@tcp(127.0.0.1:3306)/idb",
		"USERNAME":  "tong",
		"PWD":       "123456",
		"PROJECTID": "",
		"DBALIAS":   "idb"})
	err = ids.Insert(map[string]interface{}{
		"ID":        "f903de9b-9a96-4014-a991-cb01e7d96318",
		"DBTYPE":    "mysql",
		"DBURL":     "{username}:{password}@tcp(127.0.0.1:3306)/pest",
		"USERNAME":  "tong",
		"PWD":       "123456",
		"PROJECTID": "",
		"DBALIAS":   "pest"})

	fmt.Println(err)
}

func TestCreateInsertSQL(t *testing.T) {
	ids := &WriteableTableSource{
		TableDataSource{
			DBDataSource: DBDataSource{
				DataSource: DataSource{
					Name: "G_DATABASEURL",
				},
				DBAlias:        "default",
				AutoFillFields: true,
			},
			TableName: "G_DATABASEURL",
		}}
	ids.Init()
	err := ids.Insert(map[string]interface{}{
		"ID":        "05442082-76d9-41da-b563-a19914131993",
		"DBTYPE":    "mysql",
		"DBURL":     "{username}:{password}@tcp(127.0.0.1:3306)/idb",
		"USERNAME":  "tong",
		"PWD":       "123456",
		"PROJECTID": "",
		"DBALIAS":   "idb"})
	err = ids.Insert(map[string]interface{}{
		"ID":        "f903de9b-9a96-4014-a991-cb01e7d96318",
		"DBTYPE":    "mysql",
		"DBURL":     "{username}:{password}@tcp(127.0.0.1:3306)/pest",
		"USERNAME":  "tong",
		"PWD":       "123456",
		"PROJECTID": "",
		"DBALIAS":   "pest"})

	fmt.Println(err)
}
