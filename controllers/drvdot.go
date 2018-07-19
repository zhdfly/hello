package controllers

import (
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
		rlt, _ := tcpserver.GetDrvDotFromMem(this.GetString("drv"))
		beego.Info(rlt)
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
		rlt, err := tcpserver.GetDrvDotFromMem(this.GetString("drv"))
		beego.Info(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	//新建设备数据点
	if posttype == "creatnewdot" {
		dtype, err := this.GetInt("dtype")
		if err != nil {
			this.Ctx.WriteString("ERR:类型错误")
		}
		drw, err := this.GetInt("rw")
		if err != nil {
			this.Ctx.WriteString("ERR:读写方式错误")
		}
		daddr, err := this.GetInt("addr")
		if err != nil {
			this.Ctx.WriteString("ERR:地址错误")
		}
		ddata, err := this.GetInt("data")
		if err != nil {
			this.Ctx.WriteString("ERR:格式错误")
		}
		dtop, err := this.GetFloat("top")
		if err != nil {
			this.Ctx.WriteString("ERR:报警上限错误")
		}
		dbot, err := this.GetFloat("bot")
		if err != nil {
			this.Ctx.WriteString("ERR:报警下限错误")
		}
		dtime, err := this.GetInt("time")
		if err != nil {
			this.Ctx.WriteString("ERR:保存时间错误")
		}
		rlt := tcpserver.Inserttodot(
			this.GetString("drv"),
			this.GetString("name"),
			daddr, drw, dtype, ddata, float32(dtop), float32(dbot), dtime,
			this.GetString("unit"))
		if rlt == "OK" {
			this.Ctx.WriteString("OK")
			tcpserver.HotReFrash()
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	//删除设备数据点
	if posttype == "dltdrvdot" {
		rlt := tcpserver.Dltdrvdot(this.GetString("drv"), this.GetString("dotname"))
		this.Ctx.WriteString(rlt)
		tcpserver.HotReFrash()
	}
}