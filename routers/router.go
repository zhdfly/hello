package routers

import (
	"hello/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/admin", &controllers.MainController{})
	beego.Router("/admin/drvdot", &controllers.DrvdotController{})
	beego.Router("/admin/manage", &controllers.MagController{})
	beego.Router("/admin/drvmag", &controllers.DrvmagController{})
	beego.Router("/admin/video", &controllers.VideoController{})
	beego.Router("/admin/drvshow", &controllers.DrvshowController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/usershow", &controllers.UsershowController{})
	beego.Router("/usershowdrv", &controllers.UsershowdrvController{})
	beego.Router("/admin/muxchart", &controllers.MuxChartContraller{})
}
