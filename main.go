package main

import (
	"fmt"
	_ "hello/routers"
	"hello/tcpserver"

	"github.com/astaxie/beego"
)

func main() {
	go tcpserver.Tcpstart("9999")
	go tcpserver.StarthttpGet()
	fmt.Printf("Start\r\n") //ceshi
	fmt.Printf("Start\r\n")
	beego.Run()
}
