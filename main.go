package main

import (
	"fmt"
	"hello/modbus"
	_ "hello/routers"
	"hello/tcpserver"

	"github.com/astaxie/beego/context"

	"github.com/astaxie/beego"
)

func main() {
	tcpserver.ConfigSQL()
	tcpserver.Getvideolist()
	modbus.Initmodbus()
	tcpserver.GetUserInfo()
	go tcpserver.Tcpstart("9999")
	go tcpserver.StarthttpGet() //态神服务
	fmt.Printf("Start\r\n")     //ceshi
	fmt.Printf("Start\r\n")
	beego.BConfig.WebConfig.Session.SessionOn = true
	var FilterUser = func(ctx *context.Context) {
		key := ctx.Input.Cookie("name")
		userkey := ctx.Input.Cookie(key)
		skey := ctx.Input.Session("loginuser")
		userskey := ctx.Input.Session(skey)
		if key != skey && userkey != userskey && ctx.Request.RequestURI != "/login" {
			ctx.Redirect(302, "/login")
			fmt.Println("Get....:", key, skey)
		}
	}
	var FilterAdmin = func(ctx *context.Context) {
		key := ctx.Input.Cookie("name")
		userkey := ctx.Input.Cookie(key)
		skey := ctx.Input.Session("loginuser")
		userskey := ctx.Input.Session(skey)
		fmt.Println("Get....:", key, skey)
		if (key != skey && userkey != userskey) || key != "admin" {
			ctx.Redirect(302, "/login")

		}
	}
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
	beego.InsertFilter("/admin/*", beego.BeforeRouter, FilterAdmin)

	beego.Run()
}
