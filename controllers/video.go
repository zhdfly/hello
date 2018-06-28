package controllers

import (
	"github.com/astaxie/beego"
)

type VideoController struct {
	beego.Controller
}

//登录页面
func (c *VideoController) Get() {
	c.Data["name"] = c.GetSession("loginuser")
	c.TplName = "videodrv.html"
}

//登录功能
func (c *VideoController) Post() {

}
