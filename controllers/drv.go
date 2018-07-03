package controllers

import (
	"fmt"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

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

type DrvshowController struct {
	beego.Controller
}

func (c *DrvshowController) Get() {
	u := c.GetString("user")
	d := c.GetString("drv")
	dt, len, err := tcpserver.GetDrvRealTimeData(u, d)
	if err == nil {
		c.Data["tmp"] = dt
		c.Data["len"] = len
	} else {
		c.Data["tmp"] = "Error!!!"
	}
	c.Data["video"], err = tcpserver.Getdrvvedio(d)
	c.TplName = "drvshow.html"
}
func (this *DrvshowController) Post() {
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
	if posttype == "setwarning" {
		rlt, err := tcpserver.Setdotwarning(this.GetString("drv"), this.GetString("dot"), this.GetString("top"), this.GetString("bot"))
		fmt.Println(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}
