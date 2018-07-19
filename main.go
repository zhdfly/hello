package main

import (
	"fmt"
	_ "hello/routers"
	"hello/tcpserver"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
)

func main() {
	tcpserver.ConfigSQL()
	tcpserver.Getvideolist()
	tcpserver.StartModbus()
	tcpserver.GetUserInfo()
	go tcpserver.Tcpstart("9999")
	go tcpserver.StarthttpGet() //态神服务
	go tcpserver.ThreadMB()
	fmt.Printf("Start\r\n") //ceshi
	fmt.Printf("Start\r\n")
	beego.BConfig.WebConfig.Session.SessionOn = true
	var FilterUser = func(ctx *context.Context) {
		key := ctx.Input.Cookie("name")
		userkey := ctx.Input.Cookie(key)
		skey := ctx.Input.Session("loginuser")
		userskey := ctx.Input.Session(skey)
		beego.Info(key)
		beego.Info(userkey)
		beego.Info(skey)
		beego.Info(userskey)
		if (key != skey || userkey != userskey) && ctx.Request.RequestURI != "/login" {
			ctx.Redirect(302, "/login")
			beego.Info("Get....:", key, skey)
			beego.Info("Get....:", userkey, userskey)
		}
	}
	var FilterAdmin = func(ctx *context.Context) {
		key := ctx.Input.Cookie("name")
		userkey := ctx.Input.Cookie(key)
		skey := ctx.Input.Session("loginuser")
		userskey := ctx.Input.Session(skey)
		if key == skey && userkey == userskey {
			ctx.Redirect(302, "/web")
		} else {
			ctx.Redirect(302, "/login")
		}
	}
	beego.InsertFilter("/web/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/", beego.BeforeRouter, FilterAdmin)

	beego.Run()
}
