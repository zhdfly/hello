package main

import (
	"fmt"
	_ "hello/routers"

	"github.com/astaxie/beego"
)

func main() {
	fmt.Printf("Start")
	fmt.Printf("Start")
	beego.Run()
}
