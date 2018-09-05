package routers

import (
	"hello/controllers"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

func init() {
	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
	beego.Router("/", &controllers.IndexController{})
	beego.Router("/v1/web", &controllers.MainController{})
	beego.Router("/v1/web/drvdot", &controllers.DrvdotController{})
	beego.Router("/v1/web/usermag", &controllers.UserMagController{})
	beego.Router("/v1/web/drvmag", &controllers.DrvmagController{})
	beego.Router("/v1/web/areamag", &controllers.AreamagController{})
	beego.Router("/v1/web/video", &controllers.VideoController{})
	beego.Router("/v1/web/drvshow", &controllers.DrvshowController{})
	beego.Router("/v1/login", &controllers.LoginController{})
	beego.Router("/v1/logout", &controllers.LogoutController{})
	beego.Router("/v1/web/muxchart", &controllers.MuxChartContraller{})
	beego.Router("/v1/web/drvpic", &controllers.DrvPicContraller{})
	beego.Router("/channel", &controllers.OtherPostController{})
	beego.Router("/history", &controllers.OtherHisToryPostController{})

	// beego.Router("/v1/login", &controllers.LoginController{})
	// beego.Router("/v1/web", &controllers.MainController{})
	// beego.Router("/v1/web/drvdot", &controllers.DrvdotController{})
	// beego.Router("/v1/web/usermag", &controllers.UserMagController{})
	// beego.Router("/v1/web/drvmag", &controllers.DrvmagController{})
	// beego.Router("/v1/web/video", &controllers.VideoController{})
	// beego.Router("/v1/web/drvshow", &controllers.DrvshowController{})
	// beego.Router("/v1/web/muxchart", &controllers.MuxChartContraller{})
	// beego.Router("/v1/web/drvpic", &controllers.DrvPicContraller{})
	// beego.Router("/v1/channel", &controllers.OtherPostController{})
}
