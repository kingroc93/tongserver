package service

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestRequestBody(t *testing.T) {
	jsonstr := `{
  "Criteria": [
    {
      "field": "filedname",
      "operation": "=|>|<|!=|>=|<=",
      "value": "",
      "relation": "and"
    }
  ],
  "OrderBy": "TM desc",
  "InnerJoin": "",
  "Aggre": [
    {
      "Outfield": "MAX_TM",
      "Predicate": "COUNT",
      "ColName": "TM"
    }
  ],
  "bulldozer": [
    {
      "name": "bull1",
      "params": {
        "aa": 1,
        "bb": 2
      }
    }
  ]
}`
	var s SRequestBody
	json.Unmarshal([]byte(jsonstr), &s)
	fmt.Println(s)

}
