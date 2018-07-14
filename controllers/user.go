package controllers

import (
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

//用户控制器
type MagController struct {
	beego.Controller
}

//用户GET控制器
func (this *MagController) Get() {

	usrjson, err := tcpserver.Getusrinfo()
	if err != nil {
		println(err)
	}
	this.Data["name"] = this.GetSession("loginuser")
	this.Data["usrinfo"] = usrjson
	this.TplName = "manage.html"
}

//用户POST控制器
func (this *MagController) Post() {

	var usrjson string
	var err error
	var posttype = this.GetString("type")
	//获取用户现有的设备列表
	if posttype == "userdrv" {
		usrjson, err = tcpserver.Getusrdrvinfo(this.GetString("usr"))
	}
	//获取不属于用户的设备列表
	if posttype == "alldrv" {
		usrjson, err = tcpserver.Getusrnotdrvinfo(this.GetString("usr"))
	}
	//设置用户所拥有的设备列表
	if posttype == "usrsltdrv" {
		tcpserver.Setusrdrv(this.GetString("usr"), this.GetString("drv"))
		tcpserver.HotReFrash()
		usrjson = "OK"
	}
	if err == nil {
		this.Ctx.WriteString(usrjson)
	} else {
		this.Ctx.WriteString("")
	}
}

//新增加用户控制器
type AddnewusrController struct {
	beego.Controller
}

//新增加用户控制器Post方法
func (this *AddnewusrController) Post() {
	rlt := tcpserver.Inserttousr(this.GetString("usr"), this.GetString("pas"))
	if rlt == "OK" {
		this.Ctx.WriteString("OK")
		tcpserver.HotReFrash()
	} else {
		this.Ctx.WriteString("ERR")
	}
	return
}
