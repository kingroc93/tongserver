package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
	"strings"
	"sync"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/mgr"
	_ "tongserver.dataserver/routers"
	"tongserver.dataserver/service"
	_ "tongserver.dataserver/service"
	//	_ "github.com/mattn/go-oci8"
)

var mu sync.Mutex

// createDefaultDataIDs 注册默认数据源，这些数据源用于系统管理
func CreateDefaultDataIDs() {
	var meta map[string]interface{}
	tablesname := []string{"G_META", "G_META_ITEM", "JEDA_USER", "JEDA_ROLE", "JEDA_ROLE_USER", "JEDA_ORG", "JEDA_MENU",
		"G_USERSERVICE", "G_USERPROJECT", "G_SERVICE", "G_PROJECT", "G_IDS", "G_DATABASEURL"}
	for _, name := range tablesname {
		meta = map[string]interface{}{
			"name":      "default.mgr." + name,
			"inf":       "CreateWriteableTableDataSource",
			"dbalias":   "default",
			"tablename": name}
		datasource.IDSContainer[meta["name"].(string)] = meta
	}
}

// 注册数据源
func reloadDBUrl() error {
	ids := datasource.CreateTableDataSource("DBURL", "default", "G_DATABASEURL")
	rs, err := ids.GetAllData()
	if err != nil {
		return err
	}
	logs.Info("注册数据源")
	for _, row := range rs.Data {
		dbtype := row[rs.Fields["DBTYPE"].Index].(string)
		username := row[rs.Fields["USERNAME"].Index].(string)
		pwd := row[rs.Fields["PWD"].Index].(string)
		alias := row[rs.Fields["DBALIAS"].Index].(string)
		if dbtype == datasource.DbTypeMySQL {
			dburl := row[rs.Fields["DBURL"].Index].(string)
			logs.Info("\t%s  user:%s", dburl, username)
			dburl = strings.ReplaceAll(dburl, "{username}", username)
			dburl = strings.ReplaceAll(dburl, "{password}", pwd)
			err := orm.RegisterDataBase(alias, dbtype, dburl, 30)
			datasource.DBAlias2DBTypeContainer[alias] = dbtype
			if err != nil {
				return err
			}
		}
	}
	return nil
}
func reloadIds() error {
	ids := datasource.CreateTableDataSource("GIDS", "default", "G_IDS")
	rs, err := ids.GetAllData()
	if err != nil {
		return err
	}
	datasource.IDSContainer = make(datasource.IDSContainerType)
	CreateDefaultDataIDs()
	for _, row := range rs.Data {
		meta := make(map[string]interface{})
		err := json.Unmarshal([]byte(row[rs.Fields["META"].Index].(string)), &meta)
		if err != nil {
			logs.Error("加载数据源的时候发生错误，%s,%s", row[rs.Fields["META"].Index], err)
			continue
		}
		meta["inf"] = row[rs.Fields["INF"].Index].(string)
		meta["dbalias"] = row[rs.Fields["DBALIAS"].Index].(string)
		meta["name"] = row[rs.Fields["NAME"].Index].(string)
		meta["projectid"] = row[rs.Fields["PROJECTID"].Index].(string)
		datasource.IDSContainer[meta["projectid"].(string)+"."+meta["name"].(string)] = meta
	}
	return nil
}

func main() {
	logs.SetLogger("console")
	logs.Info("====================================================================")
	u1 := uuid.Must(uuid.NewV4(), nil)
	logs.Info(u1.String())
	dbtype := beego.AppConfig.String("db.default.type")
	username := beego.AppConfig.String("db.default.user")
	pwd := beego.AppConfig.String("db.default.password")
	if k, _ := beego.AppConfig.Bool("db.default.password.encrypted"); k {

	}

	service.HASHSECRET = beego.AppConfig.String("jwt.token.hashsecret")
	service.TokenExpire, _ = beego.AppConfig.Int64("jwt.token.expire")

	if dbtype == "mysql" {
		dburl := username + ":" + pwd + "@tcp(" + beego.AppConfig.String("db.default.ipport") + ")/" + beego.AppConfig.String("db.default.database")
		err := orm.RegisterDataBase("default", dbtype, dburl, 30)
		datasource.DBAlias2DBTypeContainer["default"] = dbtype
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	mgr.AddMetaFuns("dbalias", reloadDBUrl)
	mgr.AddMetaFuns("ids", reloadIds)
	err := mgr.ReloadMetaData()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true}))
	beego.Run()
}
