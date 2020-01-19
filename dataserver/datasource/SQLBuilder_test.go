package datasource

import (
	"fmt"
	"testing"
)

func TestSQLBuilder(t *testing.T) {
	sqlbuilder, _ := CreateSQLBuileder(DBType_MySQL, "JEDA_USER")

	sql, ps := sqlbuilder.CreateInsertSQLByMap(map[string]interface{}{
		"USER_ID":       "112123",
		"USER_PASSWORD": "123",
		"ORG_ID":        13001})
	fmt.Println("======== INSERT SQL =========")

	fmt.Println(sql)
	fmt.Println(ps)
	fmt.Println("======== UPDATE SQL =========")
	sqlbuilder.AddCriteria("USER_ID", OperEq, CompAnd, "112123")
	sql2, ps2 := sqlbuilder.CreateUpdateSQL(map[string]interface{}{
		"USER_ID":       "112123",
		"USER_PASSWORD": "123",
		"ORG_ID":        13001})
	fmt.Println(sql2)
	fmt.Println(ps2)
}
