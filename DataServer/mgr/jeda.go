package mgr

import (
	"awesome/datasource"
	"awesome/utils"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/skip2/go-qrcode"
	"strconv"
	"awesome/service"
)

type JedaController struct {
	beego.Controller
}
type QrcodeController struct {
	beego.Controller
}

type reloadMetaFun func () error

//用于加载元数据的函数列表
var metaFuns = make([]reloadMetaFun,0,10)

func AddMetaFuns(f reloadMetaFun){
	metaFuns= append(metaFuns, f)
}

func (c *QrcodeController) Get() {
	cnt := c.Input().Get("c")
	bs := c.Input().Get("t")
	size, errw := strconv.Atoi(c.Input().Get("s"))
	if errw != nil {
		size = 256
	}
	c.Ctx.Output.ContentType("png")
	var png []byte
	if bs == "64" {
		cnt = utils.DecodeURLBase64(cnt)
	}
	png, _ = qrcode.Encode(cnt, qrcode.Medium, size)
	c.Ctx.Output.Body(png)
}
func (c *JedaController) GetSrvs() {
	c.Data["json"] = datasource.IDSContainer
	c.ServeJSON()
}

func (c *JedaController) ReloadMetaData(){
	for _,f :=range metaFuns{
		err:=f()
		if err!=nil {
			logs.Debug("加载系统元数据时发生错误,%v %v",err.Error(),f)
			service.CreateErrorResponseByError(err,&c.Controller)
		}
	}
	r:=service.CreateRestResult(true)
	r["msg"]="重新加载成功"
	c.ServeJSON()
}

func (c *JedaController) GetUsers() {

}

func (c *JedaController) GetMenu() {
	var maps []orm.Params
	o := orm.NewOrm()
	pid := c.Input().Get("pid")
	if pid == "" {
		o.Raw("SELECT MENU_ID,PARENT_MENU_ID,MENU_NAME,MENU_URL from JEDA_MENU where PARENT_MENU_ID is NULL").Values(&maps)
	} else {
		o.Raw("SELECT MENU_ID,PARENT_MENU_ID,MENU_NAME,MENU_URL from JEDA_MENU where PARENT_MENU_ID=?", pid).Values(&maps)
	}
	c.Data["json"] = maps
	c.ServeJSON()
}
