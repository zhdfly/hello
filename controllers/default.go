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
	c.Data["name"] = "beego.me"
	dt, len, err := tcpserver.GetRealTimeData("*")
	if err == nil {
		c.Data["tmp"] = dt
		c.Data["len"] = len
	} else {
		c.Data["tmp"] = "Error!!!"
	}

	c.TplName = "index.html"
}

type UsrController struct {
	beego.Controller
}

func (c *UsrController) Get() {
	c.Data["name"] = "beego.me"
	c.Data["result"] = "NO"
	//c.Data["Email"] = "astaxie@zhdfly.com"
	c.TplName = "addinfo.html"
}

//自定义控制器02
type AddnewdotinfoController struct {
	beego.Controller
}

//实现Post方法
func (this *AddnewdotinfoController) Post() {
	rlt := tcpserver.Inserttodot(this.GetString("drv"), this.GetString("name"), this.GetString("dottype"), this.GetString("datatype"), this.GetString("info"))
	if rlt == "OK" {
		this.Ctx.WriteString("OK")
	} else {
		this.Ctx.WriteString("ERR")
	}
}

//自定义控制器02
type MagController struct {
	beego.Controller
}

//实现Post方法
func (this *MagController) Get() {

	usrjson, err := tcpserver.Getusrinfo()
	if err != nil {
		println(err)
	}
	this.Data["name"] = "beego.me"
	this.Data["usrinfo"] = usrjson
	this.TplName = "manage.html"
}

//实现Post方法
func (this *MagController) Post() {

	var usrjson string
	var err error
	var posttype = this.GetString("type")
	if posttype == "userdrv" {
		usrjson, err = tcpserver.Getusrdrvinfo(this.GetString("usr"))
	}
	if posttype == "alldrv" {
		usrjson, err = tcpserver.Getusrnotdrvinfo(this.GetString("usr"))
	}
	if err == nil {
		this.Ctx.WriteString(usrjson)
	} else {
		this.Ctx.WriteString("")
	}
}

//自定义控制器02
type AddnewusrController struct {
	beego.Controller
}

//实现Post方法
func (this *AddnewusrController) Post() {
	rlt := tcpserver.Inserttousr(this.GetString("usr"), this.GetString("pas"))
	if rlt == "OK" {
		this.Ctx.WriteString("OK")
	} else {
		this.Ctx.WriteString("ERR")
	}
	return
}

//自定义控制器02
type DrvmagController struct {
	beego.Controller
}

//实现Post方法
func (this *DrvmagController) Get() {
	usrjson, err := tcpserver.Getusrdrvinfo(this.GetString("usr"))
	if err != nil {
		println(err)
	}
	fmt.Println(usrjson)
	this.Data["name"] = "beego.me"
	this.Data["usrdrvinfo"] = usrjson
	this.TplName = "drvmag.html"
}
