package controllers

import (
	"encoding/json"
	"hello/tcpserver"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OtherPostController struct {
	beego.Controller
}

func (c *OtherPostController) Post() {
	u := c.GetString("u")
	d := c.GetString("drv")
	beego.Info(u, d)
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
	beego.Info(user)
	beego.Info(user.Pass)
}
func (c *OtherPostController) Get() {
	u := c.GetString("u")
	d := c.GetString("drv")
	beego.Info(u, d)
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
	beego.Info(user)
	beego.Info(user.Pass)
}
