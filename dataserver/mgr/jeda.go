package mgr

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/skip2/go-qrcode"
	"strconv"
	"tongserver.dataserver/datasource"
	"tongserver.dataserver/service"
	"tongserver.dataserver/utils"
)

// JedaController 后台管理控制器
type JedaController struct {
	beego.Controller
}

// QrcodeController 生成二维码的控制器
type QrcodeController struct {
	beego.Controller
}

// reloadMetaFun 重新加载元数据的函数句柄类型
type reloadMetaFun func() error

// metaFuns 用于加载元数据的函数列表
var metaFuns = make([]reloadMetaFun, 0, 10)

// AddMetaFuns 添加加载元数据的函数句柄
func AddMetaFuns(f reloadMetaFun) {
	metaFuns = append(metaFuns, f)
}

// Get 生成二维码
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

// GetSrvs 返回所有服务
func (c *JedaController) GetSrvs() {
	c.Data["json"] = datasource.IDSContainer
	c.ServeJSON()
}

// ReloadMetaData 重新加载系统元数据
func ReloadMetaData() error {
	for _, f := range metaFuns {
		err := f()
		if err != nil {
			return err
		}
	}
	return nil
}

// ReloadMetaData 重新加载元数据
func (c *JedaController) ReloadMetaData() {
	err := ReloadMetaData()
	if err != nil {
		logs.Debug("加载系统元数据时发生错误,%v", err.Error())
		service.CreateErrorResponseByError(err, &c.Controller)
	}
	r := service.CreateRestResult(true)
	r["msg"] = "重新加载成功"
	c.ServeJSON()
}

// GetUsers 返回用户列表
func (c *JedaController) GetUsers() {

}

// GetMenu 范湖菜单信息
func (c *JedaController) GetMenu() {
	var maps []orm.Params
	o := orm.NewOrm()
	pid := c.Input().Get("pid")
	if pid == "" {
		_, err := o.Raw("SELECT MENU_ID,PARENT_MENU_ID,MENU_NAME,MENU_URL from JEDA_MENU where PARENT_MENU_ID is NULL").Values(&maps)
		if err != nil {
			service.CreateErrorResponseByError(err, &c.Controller)
			return
		}
	} else {
		_, err := o.Raw("SELECT MENU_ID,PARENT_MENU_ID,MENU_NAME,MENU_URL from JEDA_MENU where PARENT_MENU_ID=?", pid).Values(&maps)
		if err != nil {
			service.CreateErrorResponseByError(err, &c.Controller)
			return
		}
	}
	c.Data["json"] = maps
	c.ServeJSON()
}
