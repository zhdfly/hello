package tcpserver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

var MBusDotList []Maindot
var MBusDrvList []Maindrv
var MBusDrv []MainDrvType

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
		fmt.Println(err)
	}
	//除了把设备信息添加到数据库中去，还需要把内存中的数据进行热更新

	//热更新数据完成
	return "ERR"
}

func Inserttodot(drv string, name string, dottype int, rw int, unit string) string {
	o := orm.NewOrm()
	var ob Maindot
	ob.Name = name
	ob.Type = dottype
	ob.Drvname = drv
	ob.Rw = rw
	ob.Unit = unit

	if created, _, err := o.ReadOrCreate(&ob, "Name", "Drvname"); err == nil {
		if created {
			return "OK"
		} else {
			return "ERR"
		}
	}
	return "ERR"
}
func Inserttodotvalue(dlist Dotvalue) string {
	o := orm.NewOrm()
	//ob := Dotvalue{Drvname: drv, Dotname: name, Value: value, Status: status, Time: time.Now().Format("2006-01-02 15:04:05")}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, err := o.Insert(&dlist); err == nil {
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
		//fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func GetMaindrvinfo() error {
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindrv").QueryRows(&MBusDrvList)
	if err == nil {
		//fmt.Println(ob)
	}
	for i := 0; i < len(MBusDrvList); i++ {
		//根据每个设备的名称选取每个设备的数据点
		_, err = o.Raw("SELECT * FROM Maindot where drvname=?", MBusDrvList[i].Name).QueryRows(&MBusDotList)
		var tmp MainDrvType
		tmp.Drv = MBusDrvList[i]
		if err == nil {
			copy(tmp.Dot, MBusDotList)
		}
		MBusDrv = append(MBusDrv, tmp)
	}
	return err
}

func Getdrvdotinfo(name string) (string, error) {
	var ob []Maindot
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindot WHERE drv = ?", name).QueryRows(&ob)
	if err == nil {
		//fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getusrdrvinfo(name string) (string, error) {
	var ob []Usrdrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM usrdrv WHERE usrname = ?", name).QueryRows(&ob)
	if err == nil {
		//fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Setusrdrv(name string, drvstring string) (string, error) {
	var ob []Usrdrv
	//usr := Usrdrv{Usrname: name}
	o := orm.NewOrm()
	o.Raw("DELETE FROM usrdrv WHERE usrname = ?", name).Exec()
	//fmt.Println(name, drvstring)
	json.Unmarshal([]byte(drvstring), &ob)
	fmt.Println(ob)
	_, err := o.InsertMulti(len(ob), ob)
	return "successNums", err
}
func Getusrnotdrvinfo(name string) (string, error) {
	var ob []Maindrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM drv where name not in (select drvname from usrdrv where usrname = ? )", name).QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Dltdrvdot(drv string, dot string) string {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM dot WHERE name = ? and drv = ?", dot, drv).Exec()
	if err == nil {
		return "OK"
	}
	return "ERR"
}

type Dotvaluertl struct {
	Data []string
	Time []string
}

func Getdotvalue(drv string, dot string, start string, stop string) (string, error) {
	var ob []string
	var obtime []string
	beego.Info(drv, dot, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT value FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, dot, start, stop).QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	_, err = o.Raw("SELECT time FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, dot, start, stop).QueryRows(&obtime)
	if err == nil {
		fmt.Println(ob)
	}
	var rlt Dotvaluertl
	rlt.Data = ob
	rlt.Time = obtime
	str, err := json.Marshal(rlt)
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
