package controllers

import (
	"hello/tcpserver"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type IndexController struct {
	beego.Controller
}

func (c *IndexController) Get() {
	c.TplName = "index.html"
}

type MainController struct {
	beego.Controller
}

func (this *MainController) Post() {
	posttype := this.GetString("type")
	logs.Info(posttype)
	//获取实时数据
	if posttype == "getreal" {
		rlt, _, err := tcpserver.GetUserDrvsFromMem(this.GetSession("loginuser"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "getareadrvlist" {
		rlt := tcpserver.GetUserAreaDrvInfoFromMem(this.GetSession("loginuser"))
		this.Ctx.WriteString(rlt)
	}
	if posttype == "getusermainmsg" {
		rlt := tcpserver.GetUserMainMsg(this.GetSession("loginuser"))
		this.Ctx.WriteString(rlt)
	}
	if posttype == "getdrvalarm" {
		rlt := tcpserver.GetDrvAlarm(this.GetString("drv"))
		this.Ctx.WriteString(rlt)
	}
	if posttype == "actuseralarm" {
		drv := this.GetString("drv")
		dot := this.GetString("dot")
		time := this.GetString("time")
		rlt := tcpserver.ActAlarmNotice(drv, dot, time)
		this.Ctx.WriteString(rlt)
	}
}
