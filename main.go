package main

import (
	"fmt"
	_ "hello/routers"

	"github.com/astaxie/beego"
)

func main() {
	fmt.Printf("Start") //ceshi
	fmt.Printf("Start")
	beego.Run()
}
