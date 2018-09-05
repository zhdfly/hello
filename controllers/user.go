package controllers

import (
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

//用户控制器
type UserMagController struct {
	beego.Controller
}

//用户GET控制器
// func (this *UserMagController) Get() {

// 	usrjson, err := tcpserver.Getusrinfo()
// 	if err != nil {
// 		println(err)
// 	}
// 	this.Data["name"] = this.GetSession("loginuser")
// 	this.Data["usrinfo"] = usrjson
// 	this.TplName = "usermag.html"
// }

//用户POST控制器
func (this *UserMagController) Post() {

	var usrjson string
	var err error
	var posttype = this.GetString("type")
	if posttype == "userlist" {
		usrjson, err = tcpserver.Getusrinfo()
	}
	//获取用户现有的设备列表
	if posttype == "userdrv" {
		usrjson, err = tcpserver.Getusrdrvinfo(this.GetString("user"))
	}
	//获取不属于用户的设备列表
	if posttype == "alldrv" {
		usrjson, err = tcpserver.Getusralldrvinfo(this.GetString("user"))
	}
	//设置用户所拥有的设备列表
	if posttype == "usrsltdrv" {
		tcpserver.Setusrdrv(this.GetString("user"), this.GetString("drv"))
		//清空用户下的所有设备信息，并重新添加，只涉及到设备的基本信息
		//tcpserver.HotReFrash()
		usrjson = "OK"
	}
	if posttype == "createuser" {
		usrjson = tcpserver.Inserttousr(this.GetString("name"), this.GetString("pass"))
		//return
	}
	if posttype == "deluser" {
		//只需要删除与用户有关，与设备无关
		usrjson = tcpserver.DelUser(this.GetString("name"))
	}
	if posttype == "changeuserpass" {
		usrjson = tcpserver.UpdataUserPass(this.GetString("name"), this.GetString("pass"))
	}
	if err == nil {
		this.Ctx.WriteString(usrjson)
	} else {
		this.Ctx.WriteString("ERR")
	}
}

//新增加用户控制器
type AddnewusrController struct {
	beego.Controller
}

//新增加用户控制器Post方法
//增加新用户
func (this *AddnewusrController) Post() {
	rlt := tcpserver.Inserttousr(this.GetString("usr"), this.GetString("pas"))
	if rlt == "OK" {
		this.Ctx.WriteString("OK")
		//添加一个用户对象即可
		//tcpserver.HotReFrash()
	} else {
		this.Ctx.WriteString("ERR")
	}
	return
}
