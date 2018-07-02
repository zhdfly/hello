package controllers

import (
	"fmt"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type UsrController struct {
	beego.Controller
}

func (this *UsrController) Get() {
	posttype := this.GetString("type")
	if posttype == "drvdot" {
		rlt, _ := tcpserver.Getdrvdotinfo(this.GetString("drv"))
		fmt.Println(rlt)
		this.Data["name"] = this.GetSession("loginuser")
		this.Data["drv"] = this.GetString("drv")
		this.Data["usrdrvinfo"] = rlt
		this.TplName = "addinfo.html"
	}
}
func (this *UsrController) Post() {
	posttype := this.GetString("type")
	if posttype == "drvdot" {
		rlt, err := tcpserver.Getdrvdotinfo(this.GetString("drv"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "dltdrvdot" {
		rlt := tcpserver.Dltdrvdot(this.GetString("drv"), this.GetString("dotname"))
		this.Ctx.WriteString(rlt)
		tcpserver.Dltdot(this.GetString("dotname"), this.GetString("drv"))
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
	if posttype == "usrsltdrv" {
		tcpserver.Setusrdrv(this.GetString("usr"), this.GetString("drv"))
		tcpserver.Getdotinfo()
		tcpserver.Creaturl()
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
		tcpserver.Getdotinfo()
	} else {
		this.Ctx.WriteString("ERR")
	}
	return
}
