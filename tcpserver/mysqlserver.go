package tcpserver

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

type Usr struct {
	Id   int
	Name string
	Pass string
}
type Drv struct {
	Id   int
	Name string
	Port string
	Type string
	Info string
	Time string
}
type Dot struct {
	Id         int
	Name       string
	Dottype    string
	Datatype   string
	Info       string
	Val        float32
	Drv        string
	Warningtop float32
	Warningbot float32
}
type Usrdrv struct {
	Id      int
	Usrname string
	Drvname string
}
type Dotvalue struct {
	Id      int
	Drvname string
	Dotname string
	Value   float32
	Status  string
	Time    string
}
type Videodrv struct {
	Id          int
	Name        string
	Appkey      string
	Appsecret   string
	Accesstoken string
	Sn          string
	Vercode     string
	Drv         string
	Liveurl     string
}

func ConfigSQL() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(Usr), new(Drv), new(Dot), new(Usrdrv), new(Dotvalue), new(Videodrv))

	orm.RegisterDataBase("default", "mysql", "root:zhd1021@tcp(127.0.0.1:3306)/maingo?charset=utf8&loc=Local")

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

func InserttoDrv(name string, port string, types string, info string) string {
	o := orm.NewOrm()
	var ob Drv
	ob.Name = name
	ob.Port = port
	ob.Type = types
	ob.Info = info
	ob.Time = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println("1123==========", ob)
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, _, err := o.ReadOrCreate(&ob, "Name"); err == nil {
		if created {
			return "OK"
		} else {
			return "ERR"
		}
	} else {
		fmt.Println(err)
	}
	return "ERR"
}

func Inserttodot(drv string, name string, dottype string, datatype string, info string) string {
	o := orm.NewOrm()
	var ob Dot
	ob.Name = name
	ob.Dottype = dottype
	ob.Datatype = datatype
	ob.Info = info
	ob.Drv = drv

	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, _, err := o.ReadOrCreate(&ob, "Name", "Drv"); err == nil {
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
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Getdrvinfo() (string, error) {
	var ob []Drv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM drv").QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getdrvdotinfo(name string) (string, error) {
	var ob []Dot
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM dot WHERE drv = ?", name).QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getusrdrvinfo(name string) (string, error) {
	var ob []Usrdrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM usrdrv WHERE usrname = ?", name).QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Setusrdrv(name string, drvstring string) (string, error) {
	var ob []Usrdrv
	//usr := Usrdrv{Usrname: name}
	o := orm.NewOrm()
	o.Raw("DELETE FROM usrdrv WHERE usrname = ?", name).Exec()
	fmt.Println(name, drvstring)
	json.Unmarshal([]byte(drvstring), &ob)
	fmt.Println(ob)
	_, err := o.InsertMulti(len(ob), ob)
	return "successNums", err
}
func Getusrnotdrvinfo(name string) (string, error) {
	var ob []Drv
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
	t, err := strconv.ParseFloat(top, 32)
	if err != nil {
		return "ERR", err
	}
	b, err := strconv.ParseFloat(bot, 32)
	if err != nil {
		return "ERR", err
	}
	_, err = o.Raw("UPDATE dot SET warningtop = ?,warningbot = ? where name = ? and drv = ?", top, bot, dot, drv).Exec()
	if err == nil {
		SetDotWarningPara(dot, drv, float32(t), float32(b))
		return "OK", err
	}
	return "ERR", err
}
