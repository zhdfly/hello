package routers

import (
	"hello/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/addinfo", &controllers.UsrController{})
	beego.Router("/manage", &controllers.MagController{})
	beego.Router("/addnewusr", &controllers.AddnewusrController{})
	beego.Router("/drvmag", &controllers.DrvmagController{})
}
