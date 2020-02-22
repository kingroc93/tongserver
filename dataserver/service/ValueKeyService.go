package service

import (
	"encoding/json"
	"fmt"
	"strings"
	"tongserver.dataserver/datasource"
)

// ValueKeyService 处理SrvValueKey形式的服务
type ValueKeyService struct {
	SHandlerBase
}

func (c *ValueKeyService) getActionMap() map[string]SerivceActionHandler {
	return map[string]SerivceActionHandler{
		SrvActionMETA: c.doGetMeta,
		SrvActionGET:  c.doGetValueByKey}
}

// 该服务不支持通过rBody请求数据
func (c *ValueKeyService) getRBody() *SRequestBody {
	return nil
}

func (c *ValueKeyService) doGetValueByKey(sdef *SDefine, ids datasource.IDataSource, rBody *SRequestBody) {
	fs := ids.GetKeyFields()
	params := make([]interface{}, len(fs), len(fs))
	for i, f := range fs {
		var err error
		params[i], err = c.ConvertString2Type(c.Ctl.Input().Get(f.Name), f.DataType)
		if err != nil {
			c.createErrorResponse("类型转换错误" + c.Ctl.Input().Get(f.Name) + " " + f.DataType + " err:" + err.Error())
			return
		}
	}
	resuleset, err := ids.QueryDataByKey(params...)
	if err != nil {
		c.createErrorResult(err.Error())
	} else {
		c.setResultSet(resuleset)
	}
	c.ServeJSON()
}

func (c *ValueKeyService) getServiceInterface(sdef *SDefine) (interface{}, error) {
	metestr := sdef.Meta
	meta := make(map[string]interface{})
	err2 := json.Unmarshal([]byte(metestr), &meta)
	if err2 != nil {
		return nil, fmt.Errorf("meta信息不正确,应为JSON格式")
	}
	idstr := meta["ids"].(string)

	if strings.Index(idstr, ".") == -1 {
		idstr = sdef.ProjectId + "." + idstr
	}
	obj, err := datasource.CreateIDSFromName(idstr)
	_, ok := obj.(*datasource.KeyStringSource)
	if !ok {
		return nil, fmt.Errorf("ValueKeyService只能接收KeyStringSource作为IDS属性")
	}
	return obj, err
}
