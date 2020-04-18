package service

import "tongserver.dataserver/activity"

// "InnerService": {
// 		"style": "InnerService",
// 		"url":"${}",
// 		"params": {},
// }
// 内部服务活动
// 在流程中调用内部的服务
type InnerServiceActivity struct {
	activity.Activity
	url string
}
