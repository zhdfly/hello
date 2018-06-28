package controllers

import (
	"GOGP/library/common"
	"fmt"
	"hello/tcpserver"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type LoginController struct {
	beego.Controller
}

//登录页面
func (c *LoginController) Get() {
	c.TplName = "login.html"
}

//登录功能
func (c *LoginController) Post() {
	name := c.GetString("name")
	pwd := c.GetString("pwd")
	islogin := "OK"
	fmt.Println(name, pwd)
	var user tcpserver.Usr
	o := orm.NewOrm()
	_ = o.Raw("SELECT * FROM `usr` WHERE name= ?", name).QueryRow(&user)
	beego.Info(user)
	fmt.Println(user.Pass)
	if user.Pass != pwd {
		islogin = "ERR"
		c.Ctx.WriteString(islogin)

	} else {
		if name != "admin" {
			islogin = "YES"
		}
		c.SetSession("loginuser", name)
		c.SetSession(name, common.Md5String(name+common.GetString(time.Now().Unix())))
		c.Ctx.SetCookie("name", name, 1800, "/")
		c.Ctx.SetCookie(name, name, 1800, "/")
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
	islogin := false
	if v != nil {
		//删除指定的session
		c.DelSession("loginuser")
		//销毁全部的session
		c.DestroySession()
		islogin = true

		fmt.Println("当前的session:")
		fmt.Println(c.CruSession)
	}
	c.Data["json"] = map[string]interface{}{"islogin": islogin}
	c.ServeJSON()
}
