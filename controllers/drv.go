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

//设备管理GET控制器
func (this *DrvmagController) Get() {
	//获取设备信息  这个设备信息本来是通过读取数据库得到的，需要改成从内存中读取数据得到
	usrjson, err := tcpserver.GetMainDrvInfoFromMem()
	if err != nil {
		println(err)
	}
	fmt.Println(usrjson)
	this.Data["name"] = this.GetSession("loginuser")
	this.Data["usrdrvinfo"] = usrjson
	this.TplName = "drvmag.html"
}

//设备管理Post控制器
func (this *DrvmagController) Post() {
	var posttype = this.GetString("type")
	beego.Info(posttype)
	//新建设备 修改为新定义的结构
	if posttype == "creatnewdrv" {
		name := this.GetString("name")
		addr, _ := this.GetInt("addr")
		port, _ := this.GetInt("port")
		packtype := this.GetString("packtype")
		rcount, _ := this.GetInt("rcount")
		rtime, _ := this.GetInt("rtime")
		spltime, _ := this.GetInt("spltime")
		rlt := tcpserver.InserttoMainDrv(
			name,
			addr,
			port,
			packtype,
			rcount,
			rtime,
			spltime)
		if rlt == "OK" {
			this.Ctx.WriteString("OK")
		} else {
			this.Ctx.WriteString("ERR")
		}
	}

	//获取设备数据点
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
	u := c.GetSession("loginuser")
	d := c.GetString("drv")
	c.Data["drv"] = d
	c.Data["name"] = u
	c.Data["video"], _ = tcpserver.Getdrvvedio(d)
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
	//获取数据点的实时信息
	if posttype == "getdotinfo" {
		u := this.GetString("user")
		d := this.GetString("drv")
		beego.Info(u, d)
		dt, _, err := tcpserver.GetDrvRealTimeData(u, d)
		if err == nil {
			this.Ctx.WriteString(dt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}

type MuxChartContraller struct {
	beego.Controller
}

func (c *MuxChartContraller) Get() {
	u := c.GetString("user")
	d := c.GetString("drv")
	dt, len, err := tcpserver.GetDrvRealTimeData(u, d)
	if err == nil {
		c.Data["tmp"] = dt
		c.Data["len"] = len
	} else {
		c.Data["tmp"] = "Error!!!"
	}
	c.Data["drv"] = d
	c.Data["video"], err = tcpserver.Getdrvvedio(d)
	c.TplName = "muxchart.html"
}
func (this *MuxChartContraller) Post() {
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
