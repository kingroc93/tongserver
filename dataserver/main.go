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
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/mgr"
	_ "tongserver.dataserver/mgr"
	_ "tongserver.dataserver/routers"
	_ "tongserver.dataserver/service"
	//	_ "github.com/mattn/go-oci8"
)

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
		if dbtype == datasource.DBType_MySQL {
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
func createDefaultDataIDs() error {
	var meta map[string]interface{}
	meta = map[string]interface{}{
		"name":      "default.mgr.G_META",
		"inf":       "CreateTableDataSource",
		"dbalias":   "default",
		"tablename": "G_META"}
	datasource.IDSContainer[meta["name"].(string)] = meta
	meta = map[string]interface{}{
		"name":      "default.mgr.G_META_ITEM",
		"inf":       "CreateTableDataSource",
		"dbalias":   "default",
		"tablename": "G_META_ITEM"}
	datasource.IDSContainer[meta["name"].(string)] = meta
	return nil
}
func reloadIds() error {
	ids := datasource.CreateTableDataSource("GIDS", "default", "G_IDS")
	rs, err := ids.GetAllData()
	if err != nil {
		return err
	}
	for _, row := range rs.Data {
		fmt.Println(row[rs.Fields["META"].Index])
		meta := make(map[string]interface{})
		err := json.Unmarshal([]byte(row[rs.Fields["META"].Index].(string)), &meta)
		if err != nil {
			logs.Error("加载数据源的时候发生错误，%s,%s", row[rs.Fields["META"].Index], err)
			continue
		}
		datasource.IDSContainer[meta["name"].(string)] = meta
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

	if dbtype == "mysql" {
		dburl := username + ":" + pwd + "@tcp(" + beego.AppConfig.String("db.default.ipport") + ")/" + beego.AppConfig.String("db.default.database")
		err := orm.RegisterDataBase("default", dbtype, dburl, 30)
		datasource.DBAlias2DBTypeContainer["default"] = dbtype
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	createDefaultDataIDs()
	mgr.AddMetaFuns(reloadDBUrl)
	mgr.AddMetaFuns(reloadIds)
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
