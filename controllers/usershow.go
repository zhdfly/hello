package controllers

import (
	"fmt"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type UsershowController struct {
	beego.Controller
}

//登录页面
func (c *UsershowController) Get() {
	c.Data["name"] = c.GetSession("loginuser")
	dt, len, err := tcpserver.GetRealTimeData(c.GetSession("loginuser"))
	if err == nil {
		c.Data["tmp"] = dt
		c.Data["len"] = len
	} else {
		c.Data["tmp"] = "Error!!!"
	}

	c.TplName = "usershow.html"
}

//登录功能
func (c *UsershowController) Post() {

}

type UsershowdrvController struct {
	beego.Controller
}

//登录页面
func (c *UsershowdrvController) Get() {
	c.Data["name"] = c.GetSession("loginuser")
	u := c.GetString("user")
	d := c.GetString("drv")
	dt, _, err := tcpserver.GetDrvRealTimeData(u, d)
	if err == nil {
		c.Data["tmp"] = dt
	} else {
		c.Data["tmp"] = "Error!!!"
	}
	c.Data["video"], err = tcpserver.Getdrvvedio(d)
	c.TplName = "usershowdrv.html"
}

//登录功能
func (this *UsershowdrvController) Post() {
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
