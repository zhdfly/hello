package tcpserver

import (
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func HotReFrash() {
	TSHotReFrashFalg = true
	MBHotReFrashFalg = true
}

var MainUser []MainUserDrv
var MainUserMAP MainMAPstringint //= make(map[string]int)

func GetUserInfo() {
	MainUserMAP.M.Lock()
	MainUserMAP.U = make(map[string]int)
	MainUserMAP.M.Unlock()
	var ob []Usr
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM usr").QueryRows(&ob)
	if err != nil {
		logs.Info(err)
	}
	for i := 0; i < len(ob); i++ {
		var tmp MainUserDrv
		tmp.User = ob[i].Name
		MainUser = append(MainUser, tmp)
		MainUserMAP.M.Lock()
		MainUserMAP.U[ob[i].Name] = i + 1
		MainUserMAP.M.Unlock()
	}
	for j := 0; j < len(MainUser); j++ {
		var drvnametmp []string
		_, err = o.Raw("SELECT drvname FROM usrdrv where usrname=?", MainUser[j].User).QueryRows(&drvnametmp)
		if err != nil {
			logs.Info(err)
		}
		for i := 0; i < len(drvnametmp); i++ {
			//MainUser[j].Drv = append(MainUser[j].Drv, drvnametmp[i])
			var drvtmp []Maindrv
			_, err = o.Raw("SELECT * FROM maindrv where name=?", drvnametmp[i]).QueryRows(&drvtmp)
			if err != nil {
				logs.Info(err)
			}
			for _, tmp := range drvtmp {
				MainUser[j].Drv = append(MainUser[j].Drv, tmp)
			}
		}
	}
}
