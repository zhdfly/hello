package controllers

import (
	"github.com/astaxie/beego"
)

type VideoController struct {
	beego.Controller
}

//登录页面
// func (c *VideoController) Get() {
// 	c.Data["name"] = c.GetSession("loginuser")
// 	c.Data["user"] = c.GetSession("loginuser")

// 	c.TplName = "videodrv.html"
// }

//登录功能
func (c *VideoController) Post() {
	posttype := c.GetString("type")
	if posttype == "newvideo" {
		//优先检测是否在当前的设备列表中存在APPKEY和APPACCESS相同的设备，如果存在则不用重新获取TOKEN
		//如不存在，则需要先获取TOKEN

	}
}
