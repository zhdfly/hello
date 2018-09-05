package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"hello/tcpserver"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/axgle/mahonia"
)

type LoginController struct {
	beego.Controller
}

//登录页面
// func (c *LoginController) Get() {
// 	c.TplName = "login.html"
// }
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	return result
}

//登录功能
func (c *LoginController) Post() {
	name := c.GetString("name")
	pwd := c.GetString("pwd")
	islogin := "OK"
	logs.Info(name, pwd)
	var user tcpserver.Usr
	o := orm.NewOrm()
	_ = o.Raw("SELECT * FROM `usr` WHERE name= ?", name).QueryRow(&user)
	logs.Info(user)
	logs.Info(user.Pass)
	if user.Pass != pwd || len(user.Pass) == 0 {
		islogin = "ERR"
		c.Ctx.WriteString(islogin)

	} else {
		se := md5.Sum([]byte(name + time.Now().Format("2006-01-02 15:04:05")))
		c.SetSession("loginuser", name)
		c.SetSession(name, hex.EncodeToString(se[:]))
		logs.Info(c.CruSession)
		c.Ctx.SetCookie("name", name, 36000, "/")
		c.Ctx.SetCookie(name, hex.EncodeToString(se[:]), 36000, "/")
		c.Ctx.WriteString(islogin)
	}
}

//退出
type LogoutController struct {
	beego.Controller
}

//登录退出功能
func (c *LogoutController) Post() {
	v := c.GetSession("loginuser")
	if v != nil {
		//删除指定的session
		logs.Info(c.CruSession)
		c.DelSession(c.GetSession("loginuser"))
		c.DelSession("loginuser")

		logs.Info("当前的session:")
		logs.Info(c.CruSession)
	}
	c.Ctx.WriteString("OK")
}
