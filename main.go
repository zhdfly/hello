package main

import (
	"fmt"
	_ "hello/routers"
	"hello/tcpserver"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"

	"github.com/astaxie/beego"
)

func main() {
	logs.SetLogger(logs.AdapterFile, `{"filename":"./log/log.log"}`)
	tcpserver.ConfigSQL()
	//tcpserver.Getvideolist()
	tcpserver.StatisDataConfig()
	tcpserver.StartModbus()
	tcpserver.TSServerConfig()
	tcpserver.GetUserInfo()
	tcpserver.AreaConfig()
	go tcpserver.StarthttpGet() //态神服务
	go tcpserver.ThreadMB()
	go tcpserver.GatewayLiteThread(10011)
	fmt.Printf("Start.........\r\n") //ceshi
	beego.BConfig.WebConfig.Session.SessionOn = true
	var FilterUser = func(ctx *context.Context) {
		key := ctx.Input.Cookie("name")
		userkey := ctx.Input.Cookie(key)
		skey := ctx.Input.Session("loginuser")
		userskey := ctx.Input.Session(skey)
		//logs.Info(key)
		//logs.Info(userkey)
		//logs.Info(skey)
		logs.Info(key, userkey, skey, userskey)
		if (key != skey || userkey != userskey) && ctx.Request.RequestURI != "/" {
			ctx.Redirect(302, "/")
			//logs.Info("Get....:", key, skey)
			//logs.Info("Get....:", userkey, userskey)
		}
	}
	beego.InsertFilter("/v1/web/*", beego.BeforeRouter, FilterUser)

	beego.Run()
}
