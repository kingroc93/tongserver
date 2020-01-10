package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) RiverGet() {
	c.ServeJSON()
}
func (c *MainController) OrgGet() {

	c.ServeJSON()
}
