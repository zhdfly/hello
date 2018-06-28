package controllers

import (
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
