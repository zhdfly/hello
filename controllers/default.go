package controllers

import (
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["name"] = c.GetSession("loginuser")
	beego.Info(c.GetSession("loginuser"))
	c.TplName = "index.html"
}
func (this *MainController) Post() {
	posttype := this.GetString("type")
	beego.Info(posttype)
	if posttype == "dotvalue" {
		rlt, err := tcpserver.Getdotvalue(this.GetString("drv"), this.GetString("dot"), this.GetString("start"), this.GetString("stop"))
		beego.Info(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "setwarning" {
		rlt, err := tcpserver.Setdotwarning(this.GetString("drv"), this.GetString("dot"), this.GetString("top"), this.GetString("bot"))
		beego.Info(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	//获取实时数据
	if posttype == "getreal" {
		rlt, _, err := tcpserver.GetUserDrvsFromMem(this.GetSession("loginuser"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}
