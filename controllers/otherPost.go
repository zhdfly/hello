package controllers

import (
	"encoding/json"
	"hello/tcpserver"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

type OtherHisToryPostController struct {
	beego.Controller
}

func (c *OtherHisToryPostController) Post() {
	u := c.GetString("u")
	d := c.GetString("drv")
	s := c.GetString("start")
	e := c.GetString("end")
	logs.Info(u, d)
	var user tcpserver.Usr
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM `usr` WHERE pass= ?", u).QueryRow(&user)
	if err != nil {
		rlt, _ := json.Marshal(map[string]interface{}{"Status": "错误", "Data": "参数错误"})
		c.Ctx.WriteString(string(rlt))
	} else {
		str, _ := tcpserver.GetalldotvalueR(d, s, e)
		rlt, _ := json.Marshal(map[string]interface{}{"Status": "正常", "Data": str})
		c.Ctx.WriteString(string(rlt))
	}
}

type OtherPostController struct {
	beego.Controller
}

func (c *OtherPostController) Post() {
	u := c.GetString("u")
	d := c.GetString("drv")
	logs.Info(u, d)
	var user tcpserver.Usr
	o := orm.NewOrm()
	err := o.Raw("SELECT * FROM `usr` WHERE pass= ?", u).QueryRow(&user)
	if err != nil {
		rlt, _ := json.Marshal(map[string]interface{}{"Status": "错误", "Data": "参数错误"})
		c.Ctx.WriteString(string(rlt))
	} else {
		str, _, _ := tcpserver.GetUserDrvDotInfoFromMem(user.Name, d)
		rlt, _ := json.Marshal(map[string]interface{}{"Status": "正常", "Data": str})
		c.Ctx.WriteString(string(rlt))
	}
	logs.Info(user)
	logs.Info(user.Pass)
}

// func (c *OtherPostController) Get() {
// 	u := c.GetString("u")
// 	d := c.GetString("drv")
// 	logs.Info(u, d)
// 	var user tcpserver.Usr
// 	o := orm.NewOrm()
// 	err := o.Raw("SELECT * FROM `usr` WHERE pass= ?", u).QueryRow(&user)
// 	if err != nil {
// 		rlt, _ := json.Marshal(map[string]interface{}{"Status": "错误", "Data": "参数错误"})
// 		c.Ctx.WriteString(string(rlt))
// 	} else {
// 		str, _, _ := tcpserver.GetUserDrvDotInfoFromMem(user.Name, d)
// 		rlt, _ := json.Marshal(map[string]interface{}{"Status": "正常", "Data": str})
// 		c.Ctx.WriteString(string(rlt))
// 	}
// 	logs.Info(user)
// 	logs.Info(user.Pass)
// }
