package controllers

import (
	"hello/tcpserver"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

//自定义控制器02
type AreamagController struct {
	beego.Controller
}

func (c *AreamagController) Post() {
	posttype := c.GetString("type")
	if posttype == "createarea" {
		user := c.GetSession("loginuser")
		area := c.GetString("area")
		drvs := c.GetString("drvs")
		logs.Info(drvs, area)
		rlt := tcpserver.NewArea(user.(string), area, drvs)
		c.Ctx.WriteString(rlt)
	}
	if posttype == "updataareadrv" {
		user := c.GetSession("loginuser")
		area := c.GetString("area")
		drvs := c.GetString("drvs")
		logs.Info(drvs, area)
		rlt := tcpserver.UpdataArea(user.(string), area, drvs)
		c.Ctx.WriteString(rlt)
	}
	if posttype == "deletearea" {
		user := c.GetSession("loginuser")
		area := c.GetString("area")
		logs.Info(user, area)
		rlt := tcpserver.DeleteArea(user.(string), area)
		c.Ctx.WriteString(rlt)
	}
	if posttype == "getarea" {
		user := c.GetSession("loginuser")
		logs.Info(user)
		rlt := tcpserver.GetUserArea(user)
		c.Ctx.WriteString(rlt)
	}
	if posttype == "drvbkimgok" {
		if _, err := tcpserver.UploadAreaPlaneBk(c.GetString("area"), DrvBkImgUrl); err == nil {
			c.Ctx.WriteString("OK")
		} else {
			c.Ctx.WriteString("ERR")
		}
	}
}
