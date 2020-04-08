package main

import (
	"tongserver.dataserver/app"
)

// appname = tongds
// runmode = dev
// copyrequestbody = true

// [dev]
// httpport = 8081

// jwt.token.expire = 600
// jwt.token.hashsecret = "1@3wq,klahjaqwweq"

// db.default.type = "mysql"
// db.default.ipport = "127.0.0.1:3306"
// db.default.database = "idb"
// db.default.user = "tong"
// db.default.password = "123456"
// db.default.password.encrypted = false

// redis.ip = 192.168.0.100
// redis.port = 6379
// redis.pwd = 123456

func main() {

	app.RunApp()
}
