package controllers

import (
	"fmt"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

type DrvdotController struct {
	beego.Controller
}

//设备数据点GET控制器
func (this *DrvdotController) Get() {
	posttype := this.GetString("type")
	if posttype == "drvdot" {
		rlt, _ := tcpserver.Getdrvdotinfo(this.GetString("drv"))
		fmt.Println(rlt)
		this.Data["name"] = this.GetSession("loginuser")
		this.Data["drv"] = this.GetString("drv")
		this.Data["usrdrvinfo"] = rlt
		this.TplName = "drvdot.html"
	}
}

//设备数据点POST控制器
func (this *DrvdotController) Post() {
	posttype := this.GetString("type")
	//获取设备数据点
	if posttype == "drvdot" {
		rlt, err := tcpserver.Getdrvdotinfo(this.GetString("drv"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	//新建设备数据点
	// if posttype == "creatnewdot" {
	// 	rlt := tcpserver.Inserttodot(this.GetString("drv"), this.GetString("name"), this.GetString("dtype"), this.GetString("datatype"), this.GetString("info"))
	// 	if rlt == "OK" {
	// 		this.Ctx.WriteString("OK")
	// 		tcpserver.HotReFrash()
	// 	} else {
	// 		this.Ctx.WriteString("ERR")
	// 	}
	// }
	//删除设备数据点
	if posttype == "dltdrvdot" {
		rlt := tcpserver.Dltdrvdot(this.GetString("drv"), this.GetString("dotname"))
		this.Ctx.WriteString(rlt)
		tcpserver.HotReFrash()
	}
}
