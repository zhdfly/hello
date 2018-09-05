package controllers

import (
	"fmt"
	"hello/tcpserver"
	"path"
	"strconv"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

var DrvBkImgUrl string

//自定义控制器02
type DrvmagController struct {
	beego.Controller
}

//设备管理GET控制器
// func (this *DrvmagController) Get() {
// 	//获取设备信息  这个设备信息本来是通过读取数据库得到的，需要改成从内存中读取数据得到
// 	usrjson, err := tcpserver.GetMainDrvInfoFromMem(this.GetSession("loginuser"))
// 	if err != nil {
// 		println(err)
// 	}
// 	logs.Info(usrjson)
// 	this.Data["name"] = this.GetSession("loginuser")
// 	this.Data["usrdrvinfo"] = usrjson
// 	this.TplName = "drvmag.html"
// }

//设备管理Post控制器
func (this *DrvmagController) Post() {
	u := this.GetSession("loginuser")
	var posttype = this.GetString("type")
	logs.Info(posttype)
	//新建设备 修改为新定义的结构
	if posttype == "creatnewdrv" {
		name := this.GetString("name")
		addr, _ := this.GetInt("addr")
		port, _ := this.GetInt("port")
		packtype := this.GetString("packtype")
		cmittype := this.GetString("cmittype")
		idcode, _ := this.GetInt("idcode")
		polltime, _ := this.GetInt("polltime")
		rlt := tcpserver.InserttoMainDrv(
			u,
			name,
			addr,
			port,
			packtype,
			cmittype,
			idcode,
			polltime)
		if rlt == "OK" {
			this.Ctx.WriteString("OK")
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	if posttype == "deletedrv" {
		_, err := tcpserver.DeleteDrv(this.GetString("drv"))
		if err == nil {
			this.Ctx.WriteString("OK")
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	if posttype == "pointdrv" {
		_, err := tcpserver.PointDrv(this.GetString("drv"))
		if err == nil {
			this.Ctx.WriteString("OK")
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	if posttype == "unpointdrv" {
		_, err := tcpserver.UnPointDrv(this.GetString("drv"))
		if err == nil {
			this.Ctx.WriteString("OK")
		} else {
			this.Ctx.WriteString("ERR")
		}
	}
	if posttype == "updataxy" {
		x, _ := this.GetInt("planex")
		y, _ := this.GetInt("planey")
		_, err := tcpserver.UpdataPlaneDrv(this.GetString("drv"), x, y)
		if err == nil {
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
	if posttype == "drvlist" {
		rlt, err := tcpserver.GetMainDrvInfoFromMem(this.GetSession("loginuser"))
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

// func (c *DrvshowController) Get() {
// 	u := c.GetSession("loginuser")
// 	d := c.GetString("drv")
// 	c.Data["drv"] = d
// 	c.Data["name"] = u
// 	c.Data["video"], _ = tcpserver.Getdrvvedio(d)
// 	c.TplName = "drvshow.html"
// }
func (this *DrvshowController) Post() {
	posttype := this.GetString("type")
	if posttype == "dotvalue" {
		rlt, err := tcpserver.Getdotvalue(this.GetString("drv"), this.GetString("dot"), this.GetString("start"), this.GetString("stop"))
		logs.Info(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "setwarning" {
		dtop, err := this.GetFloat("top")
		if err != nil {
			this.Ctx.WriteString("ERR")
		}
		dbot, err := this.GetFloat("bot")
		if err != nil {
			this.Ctx.WriteString("ERR")
		}
		rlt, err := tcpserver.Setdotwarning(this.GetString("drv"), this.GetString("dot"), float32(dtop), float32(dbot))
		logs.Info(rlt)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	//获取数据点的实时信息
	if posttype == "getdotinfo" {
		u := this.GetSession("loginuser")
		d := this.GetString("drv")
		dt, _, err := tcpserver.GetUserDrvInfoFromMem(u, d)
		if err == nil {
			this.Ctx.WriteString(dt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "setiovalue" {
		d := this.GetString("drv")
		dot := this.GetString("dot")
		v := this.GetString("value")
		logs.Info(d, dot, v)
		//新建MAP通道
		var chtmp = make(chan string)
		tcpserver.MBCtrlCmdRltChanMap[d] = chtmp
		if v == "0" {
			tcpserver.SetMBioValue(d, dot, 0x0000)
		} else {
			tcpserver.SetMBioValue(d, dot, 0xff00)
		}
		logs.Info(tcpserver.MBCtrlCmdRltChanMap)
		//删除MAP通道
		setrlt := <-tcpserver.MBCtrlCmdRltChanMap[d]
		logs.Info(setrlt)
		delete(tcpserver.MBCtrlCmdRltChanMap, d)
		this.Ctx.WriteString(setrlt)
	}
	if posttype == "setregvalue" {
		d := this.GetString("drv")
		dot := this.GetString("dot")
		v := this.GetString("value")
		logs.Info(d, dot, v)
		i, err := strconv.ParseFloat(v, 32)
		if err != nil {
			this.Ctx.WriteString("ERR")
		} else {
			var chtmp = make(chan string)
			tcpserver.MBCtrlCmdRltChanMap[d] = chtmp
			tcpserver.SetMBregValue(d, dot, float32(i))
			logs.Info(tcpserver.MBCtrlCmdRltChanMap)
			setrlt := <-tcpserver.MBCtrlCmdRltChanMap[d]
			delete(tcpserver.MBCtrlCmdRltChanMap, d)
			logs.Info(setrlt)
			this.Ctx.WriteString(setrlt)
		}
	}
}

type MuxChartContraller struct {
	beego.Controller
}

// func (c *MuxChartContraller) Get() {
// 	u := c.GetString("user")
// 	d := c.GetString("drv")
// 	c.Data["name"] = u
// 	dt, len, err := tcpserver.GetUserDrvInfoFromMem(u, d)
// 	if err == nil {
// 		c.Data["tmp"] = dt
// 		c.Data["len"] = len
// 	} else {
// 		c.Data["tmp"] = "Error!!!"
// 	}
// 	c.Data["drv"] = d
// 	c.Data["video"], err = tcpserver.Getdrvvedio(d)
// 	c.TplName = "muxchart.html"
// }
func (this *MuxChartContraller) Post() {
	posttype := this.GetString("type")
	if posttype == "dotvalue" {
		rlt, err := tcpserver.Getalldotvalue(this.GetString("drv"), this.GetString("start"), this.GetString("stop"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "dotvaluer" {
		rlt, err := tcpserver.GetalldotvalueR(this.GetString("drv"), this.GetString("start"), this.GetString("stop"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "comdrvvaluer" {
		drvs := this.GetString("drvs")
		s := this.GetString("start")
		e := this.GetString("stop")
		logs.Info(drvs, s, e)
		rlt, err := tcpserver.GetdrvsalldotvalueR(strings.Split(drvs, ","), s, e)
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "dotavgvalue" {
		rlt, err := tcpserver.GetDotAvgValues(this.GetString("drv"), this.GetString("start"), this.GetString("stop"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
	if posttype == "dotdailyvalue" {
		rlt, err := tcpserver.GetDotDailyValues(this.GetString("drv"), this.GetString("start"), this.GetString("stop"))
		if err == nil {
			this.Ctx.WriteString(rlt)
		} else {
			this.Ctx.WriteString("")
		}
	}
}

type DrvPicContraller struct {
	beego.Controller
}

// func (c *DrvPicContraller) Get() {
// 	u := c.GetSession("loginuser")
// 	d := c.GetString("drv")
// 	c.Data["name"] = u
// 	c.Data["drv"] = d
// 	logs.Info(u, d)
// 	c.TplName = "drvpic.html"
// }
func (this *DrvPicContraller) Post() {
	//image，这是一个key值，对应的是html中input type-‘file’的name属性值
	f, h, _ := this.GetFile("file")
	//得到文件的名称
	fileName := h.Filename
	arr := strings.Split(fileName, ":")
	if len(arr) > 1 {
		index := len(arr) - 1
		fileName = arr[index]
	}
	fmt.Println("文件名称:")
	fmt.Println(fileName)
	//关闭上传的文件，不然的话会出现临时文件不能清除的情况
	f.Close()
	this.SaveToFile("file", path.Join("static/img", fileName))
	DrvBkImgUrl = "/static/img/" + fileName
	this.Ctx.WriteString("/static/img/" + fileName)
}
