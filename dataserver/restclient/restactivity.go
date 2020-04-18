package restclient

import "tongserver.dataserver/activity"

// 调用restful服务接口的活动
type RestClientActivity struct {
	activity.Activity
	// 服务完整的URL
	url string
}
