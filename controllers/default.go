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
	c.Data["name"] = c.GetSession("loginuser")
	fmt.Println(c.GetSession("loginuser"))
	dt, len, err := tcpserver.GetRealTimeData(c.GetSession("loginuser"))
	if err == nil {
		c.Data["tmp"] = dt
		c.Data["len"] = len
	} else {
		c.Data["tmp"] = "Error!!!"
	}

	c.TplName = "index.html"
}
func (this *MainController) Post() {
	posttype := this.GetString("type")
	if posttype == "dotvalue" {
		rlt, err := tcpserver.Getdotvalue(this.GetString("drv"), this.GetString("dot"), this.GetString("start"), this.GetString("stop"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}

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

//自定义控制器02
type DrvmagController struct {
	beego.Controller
}

//实现Post方法
func (this *DrvmagController) Get() {
	usrjson, err := tcpserver.Getdrvinfo()
	if err != nil {
		println(err)
	}
	fmt.Println(usrjson)
	this.Data["name"] = "beego.me"
	this.Data["usrdrvinfo"] = usrjson
	this.TplName = "drvmag.html"
}

//实现Post方法
func (this *DrvmagController) Post() {
	var posttype = this.GetString("type")
	beego.Info(posttype)
	if posttype == "creatnewdrv" {
		rlt := tcpserver.InserttoDrv(this.GetString("name"), this.GetString("port"), this.GetString("types"), this.GetString("info"))
		if rlt == "OK" {
			this.Ctx.WriteString("OK")
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	if posttype == "creatnewdot" {
		rlt := tcpserver.Inserttodot(this.GetString("drv"), this.GetString("name"), this.GetString("dtype"), this.GetString("datatype"), this.GetString("info"))
		if rlt == "OK" {
			this.Ctx.WriteString("OK")
			tcpserver.Addnewdot(this.GetString("name"), this.GetString("datatype"), this.GetString("dtype"), this.GetString("drv"))
			tcpserver.Creaturl()
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	if posttype == "drvdot" {
		rlt, err := tcpserver.Getdrvdotinfo(this.GetString("drv"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}
