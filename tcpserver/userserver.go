package tcpserver

import (
	"fmt"

	"github.com/astaxie/beego/orm"
)

var MainUser []MainUserDrv

func GetUserInfo() {
	var ob []Usr
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM usr").QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	for i := 0; i < len(ob); i++ {
		var tmp MainUserDrv
		tmp.User = ob[i].Name
		MainUser = append(MainUser, tmp)
	}
	for j := 0; j < len(MainUser); j++ {
		var drvnametmp []string
		_, err = o.Raw("SELECT drvname FROM usrdrv where usrname=?", MainUser[j].User).QueryRows(&drvnametmp)
		if err == nil {
			fmt.Println(ob)
		}
		for i := 0; i < len(drvnametmp); i++ {
			MainUser[j].Drv = append(MainUser[j].Drv, drvnametmp[i])
		}
	}
}
