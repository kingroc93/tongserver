package datasource

import (
	"fmt"
	"testing"
)

func TestSQLBuilder(t *testing.T) {
	sqlbuilder, _ := CreateSQLBuileder(DbTypeMySQL, "JEDA_USER")

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

	//CreateSQLBuileder2(dbType string, tablename string, columns []string, orderby []string, rowslimit int, rowsoffset int) (ISQLBuilder, error)
	sqld, _ := CreateSQLBuileder2(DbTypeMySQL, "JEDA_USER", []string{"USER_ID", "USER_NAME"}, nil, 0, 0)
	//sqld.AddJoin(INNER_JOIN, "JEDA_ORG", []string{"ORG_NAME"}).

	sqld.AddJoin(&PieceJoin{
		Join:      "INNER_JOIN",
		TableName: "JEDA_ORG",
		criteria: []*SQLCriteria{
			&SQLCriteria{
				PropertyName: "ORG_ID",
				Operation:    OperEq,
				Value:        &FieldNameWithTableName{Tablename: "JEDA_ORG", Fielname: "ORG_ID"},
				Complex:      CompAnd}},
		OutField: []string{"ORG_NAME"}})

	sql3, _ := sqld.CreateSelectSQL()
	fmt.Println(sql3)
}
