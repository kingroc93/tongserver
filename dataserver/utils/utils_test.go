package utils

import (
	"fmt"
	"testing"
	"tongserver.dataserver/utils/mapstructure"
)

type criteriaInRBody struct {
	Field     string
	Operation string
	Value     interface{}
	Relation  string
}

type sRequestBody struct {
	// Insert 新建
	Insert map[string]string
	// Update 更新
	Update map[string]string
	// Delete 删除
	Delete string
	// OperationConfirm 操作二次确认
	OperationConfirm string
	// Criteria 条件节点,针对更新、删除、查询操作
	Criteria []criteriaInRBody
}

func TestGetFieldName(t *testing.T) {
	m := make(map[string]interface{})
	m["delete"] = "asASDQWEQWESS"
	m["OperationConfirm"] = "operationConfirm"
	mm := make(map[string]string)
	m["Insert"] = mm
	mm["fa"] = "fa"
	mm["fb"] = "fb"

	mm2 := make(map[string]interface{})
	mar := make([]map[string]interface{}, 0)
	mm2["Field"] = "Field"
	mm2["Operation"] = "Operation"
	mm2["Value"] = "Value"
	mm2["Relation"] = "Relation"
	mar = append(mar, mm2)
	m["Criteria"] = mar

	sbody := &sRequestBody{}

	err := mapstructure.Decode(m, sbody)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(sbody)
}
