package routers

import (
	"hello/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/web", &controllers.MainController{})
	beego.Router("/web/drvdot", &controllers.DrvdotController{})
	beego.Router("/web/usermag", &controllers.UserMagController{})
	beego.Router("/web/drvmag", &controllers.DrvmagController{})
	beego.Router("/web/video", &controllers.VideoController{})
	beego.Router("/web/drvshow", &controllers.DrvshowController{})
	beego.Router("/login", &controllers.LoginController{})
	beego.Router("/web/muxchart", &controllers.MuxChartContraller{})
	beego.Router("/web/drvpic", &controllers.DrvPicContraller{})
	beego.Router("/channel", &controllers.OtherPostController{})
}
