package tcpserver

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func ConfigSQL() {
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqlurls := beego.AppConfig.String("mysqlurls")
	mysqldb := beego.AppConfig.String("mysqldb")
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(Usr), new(Usrdrv), new(Dotvalue), new(Videodrv), new(Maindrv), new(Maindot))

	orm.RegisterDataBase("default", "mysql", mysqluser+":"+mysqlpass+"@tcp("+mysqlurls+")/"+mysqldb+"?charset=utf8&loc=Local")

	//orm.Debug = true
}
func Inserttousr(usr string, pass string) string {
	o := orm.NewOrm()
	user := Usr{Name: usr, Pass: pass}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, _, err := o.ReadOrCreate(&user, "Name"); err == nil {
		if created {
			return "OK"
		} else {
			return "ERR"
		}
	}
	return "ERR"
}

func InserttoMainDrv(name string, addr int, port int, types string, retrycount int, retrytime int, samplingtime int) string {
	o := orm.NewOrm()
	var ob Maindrv
	ob.Name = name
	ob.Port = port
	ob.Addr = addr
	ob.Packtype = types
	ob.Retrycount = retrycount
	ob.Retrytime = retrytime
	ob.Samplingtime = samplingtime
	ob.Time = time.Now().Format("2006-01-02 15:04:05")
	//把新增加的设备添加到数据库中
	if created, _, err := o.ReadOrCreate(&ob, "Name"); err == nil {
		if created {
			return "OK"
		} else {
			return "ERR"
		}
	} else {
		beego.Info(err)
	}
	//除了把设备信息添加到数据库中去，还需要把内存中的数据进行热更新

	//热更新数据完成
	return "ERR"
}

func Inserttodot(drv string, name string, addr int, rw int, dtype int, data int, top float32, bot float32, time int, unit string) string {
	o := orm.NewOrm()
	var err error
	var ob Maindot
	ob.Name = name
	ob.Type = dtype
	ob.Drvname = drv
	ob.Rw = rw
	ob.Unit = unit
	ob.Data = data
	ob.Addr = addr
	ob.Alarmtop = top
	ob.Alarmbot = bot
	ob.Savetime = time
	created, _, err := o.ReadOrCreate(&ob, "Name", "Drvname")
	if err == nil {
		if created {
			return "OK"
		} else {
			return "ERR"
		}
	}
	beego.Info(err)
	return "ERR"
}
func Inserttodotvalue(dlist Maindot) string {
	o := orm.NewOrm()
	var tmpdot Dotvalue
	tmpdot.Drvname = dlist.Drvname
	tmpdot.Dotname = dlist.Name
	tmpdot.Value = dlist.Value
	tmpdot.Status = dlist.Status
	tmpdot.Time = time.Now().Format("2006-01-02 15:04:05")
	//ob := Dotvalue{Drvname: drv, Dotname: name, Value: value, Status: status, Time: time.Now().Format("2006-01-02 15:04:05")}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, err := o.Insert(&tmpdot); err == nil {
		if created != 0 {
			return "OK"
		} else {
			return "ERR"
		}
	}
	return "ERR"
}
func Getusrinfo() (string, error) {
	var ob []Usr
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM usr").QueryRows(&ob)
	if err == nil {
		//beego.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

// func GetMaindrvinfo() error {
// 	o := orm.NewOrm()
// 	_, err := o.Raw("SELECT * FROM maindrv").QueryRows(&MBusDrvList)
// 	if err == nil {
// 		//beego.Info(ob)
// 	}
// 	for i := 0; i < len(MBusDrvList); i++ {
// 		//根据每个设备的名称选取每个设备的数据点
// 		_, err = o.Raw("SELECT * FROM Maindot where drvname=?", MBusDrvList[i].Name).QueryRows(&MBusDotList)
// 		var tmp MainDrvType
// 		tmp.Drv = MBusDrvList[i]
// 		if err == nil {
// 			copy(tmp.Dot, MBusDotList)
// 		}
// 		MBusDrv = append(MBusDrv, tmp)
// 	}
// 	return err
// }

func Getdrvdotinfo(name string) (string, error) {
	var ob []Maindot
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindot WHERE drv = ?", name).QueryRows(&ob)
	if err == nil {
		//beego.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getusrdrvinfo(name string) (string, error) {
	var ob []Usrdrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM usrdrv WHERE usrname = ?", name).QueryRows(&ob)
	if err == nil {
		//beego.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Setusrdrv(name string, drvstring string) (string, error) {
	var ob []Usrdrv
	//usr := Usrdrv{Usrname: name}
	o := orm.NewOrm()
	o.Raw("DELETE FROM usrdrv WHERE usrname = ?", name).Exec()
	//beego.Info(name, drvstring)
	json.Unmarshal([]byte(drvstring), &ob)
	beego.Info(ob)
	_, err := o.InsertMulti(len(ob), ob)
	return "successNums", err
}
func Getusrnotdrvinfo(name string) (string, error) {
	var ob []Maindrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindrv where name not in (select drvname from usrdrv where usrname = ? )", name).QueryRows(&ob)
	if err == nil {
		beego.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Dltdrvdot(drv string, dot string) string {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM maindot WHERE name = ? and drvname = ?", dot, drv).Exec()
	if err == nil {
		return "OK"
	}
	return "ERR"
}

func Getdotvalue(drv string, dot string, start string, stop string) (string, error) {
	var ob []string
	var obtime []string
	beego.Info(drv, dot, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT value FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, dot, start, stop).QueryRows(&ob)
	if err == nil {
		beego.Info(ob)
	}
	_, err = o.Raw("SELECT time FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, dot, start, stop).QueryRows(&obtime)
	if err == nil {
		beego.Info(ob)
	}
	var rlt Dotvaluertl
	rlt.Data = ob
	rlt.Time = obtime
	str, err := json.Marshal(rlt)
	return string(str), err
}

type Dotvaluertl struct {
	Name string
	Data []string
	Time []string
}

func Getalldotvalue(drv string, start string, stop string) (string, error) {

	var dotvalue []string
	var dottime []string
	var dotlist []Dotvaluertl
	var dotname []string
	beego.Info(drv, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM maindot where drvname=?", drv).QueryRows(&dotname)
	if err == nil {
		beego.Info(dotname)
	}
	for c, i := range dotname {
		beego.Info(c, i)
		_, err := o.Raw("SELECT time FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, i, start, stop).QueryRows(&dottime)
		if err != nil {
			beego.Info(err)
		}
		_, err = o.Raw("SELECT value FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, i, start, stop).QueryRows(&dotvalue)
		if err != nil {
			beego.Info(err)
		}
		var tmp Dotvaluertl
		tmp.Name = i
		tmp.Data = dotvalue
		tmp.Time = dottime
		dotlist = append(dotlist, tmp)
	}

	str, err := json.Marshal(dotlist)
	return string(str), err
}
func Setdotwarning(drv string, dot string, top string, bot string) (string, error) {
	o := orm.NewOrm()
	_, err := strconv.ParseFloat(top, 32)
	if err != nil {
		return "ERR", err
	}
	_, err = strconv.ParseFloat(bot, 32)
	if err != nil {
		return "ERR", err
	}
	_, err = o.Raw("UPDATE dot SET warningtop = ?,warningbot = ? where name = ? and drv = ?", top, bot, dot, drv).Exec()
	if err == nil {
		//SetDotWarningPara(dot, drv, float32(t), float32(b))
		return "OK", err
	}
	return "ERR", err
}
