package tcpserver

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func HotReFrash() {
	TSHotReFrashFalg = true
	MBHotReFrashFalg = true
}

var MainUser []MainUserDrv
var MainUserMAP = make(map[string]int)

func GetUserInfo() {
	var ob []Usr
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM usr").QueryRows(&ob)
	if err != nil {
		beego.Info(err)
	}
	for i := 0; i < len(ob); i++ {
		var tmp MainUserDrv
		tmp.User = ob[i].Name
		MainUser = append(MainUser, tmp)
		MainUserMAP[ob[i].Name] = i + 1
	}
	for j := 0; j < len(MainUser); j++ {
		var drvnametmp []string
		_, err = o.Raw("SELECT drvname FROM usrdrv where usrname=?", MainUser[j].User).QueryRows(&drvnametmp)
		if err != nil {
			beego.Info(err)
		}
		for i := 0; i < len(drvnametmp); i++ {
			//MainUser[j].Drv = append(MainUser[j].Drv, drvnametmp[i])
			var drvtmp []Maindrv
			_, err = o.Raw("SELECT * FROM maindrv where name=?", drvnametmp[i]).QueryRows(&drvtmp)
			if err != nil {
				beego.Info(err)
			}
			for _, tmp := range drvtmp {
				MainUser[j].Drv = append(MainUser[j].Drv, tmp)
			}
		}
	}
}
