package controllers

import (
	"fmt"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type UsershowController struct {
	beego.Controller
}

//非管理员权限用户登录
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

//非管理员权限用户登录POST
func (c *UsershowController) Post() {
	posttype := c.GetString("type")
	//获取实时数据
	if posttype == "getreal" {
		rlt, _, err := tcpserver.GetRealTimeData(c.GetSession("loginuser"))
		fmt.Println(rlt)
		if err == nil {
			c.Ctx.WriteString(rlt)
		} else {
			c.Ctx.WriteString("")
		}
	}
}

type UsershowdrvController struct {
	beego.Controller
}

//非管理员权限用户获取设备详细信息
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

//非管理员权限用户获取设备详细信息POST
func (this *UsershowdrvController) Post() {
	posttype := this.GetString("type")
	//获取数据点的历史信息
	if posttype == "dotvalue" {
		rlt, err := tcpserver.Getdotvalue(this.GetString("drv"), this.GetString("dot"), this.GetString("start"), this.GetString("stop"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	//获取数据点的实时信息
	if posttype == "getdotinfo" {
		u := this.GetString("user")
		d := this.GetString("drv")
		beego.Info(u, d)
		dt, _, err := tcpserver.GetDrvRealTimeData(u, d)
		if err == nil {
			this.Ctx.WriteString(dt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	//设置数据点的报警阈值
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
