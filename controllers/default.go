package controllers

import (
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["name"] = "beego.me"
	//c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.html"
}

type UsrController struct {
	beego.Controller
}

func (c *UsrController) Get() {
	c.Data["name"] = string(tcpserver.Buffer[0:5])
	//c.Data["Email"] = "astaxie@zhdfly.com"
	c.TplName = "index.html"
}
