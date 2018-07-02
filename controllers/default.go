package controllers

import (
	"fmt"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["name"] = c.GetSession("loginuser")
	fmt.Println(c.GetSession("loginuser"))
	dt, len, err := tcpserver.GetRealTimeData(c.GetSession("loginuser"))
	if err == nil {
		c.Data["tmp"] = dt
		c.Data["len"] = len
	} else {
		c.Data["tmp"] = "Error!!!"
	}

	c.TplName = "index.html"
}
func (this *MainController) Post() {
	posttype := this.GetString("type")
	if posttype == "dotvalue" {
		rlt, err := tcpserver.Getdotvalue(this.GetString("drv"), this.GetString("dot"), this.GetString("start"), this.GetString("stop"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "setwarning" {
		rlt, err := tcpserver.Setdotwarning(this.GetString("drv"), this.GetString("dot"), this.GetString("top"), this.GetString("bot"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}
