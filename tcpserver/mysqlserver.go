package tcpserver

import (
	"encoding/json"
	"fmt"

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
	Info string
	Time string
}
type Dot struct {
	Id       int
	Name     string
	Dottype  string
	Datatype string
	Info     string
	Drv      string
}
type UsrDrv struct {
	Id      int
	Usrname string
	Drvname string
}

func ConfigSQL() {
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(Usr), new(Drv), new(Dot), new(UsrDrv))

	orm.RegisterDataBase("default", "mysql", "root:zhd1021@tcp(127.0.0.1:3306)/maingo?charset=utf8&loc=Local")

	orm.Debug = true
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

func Inserttodot(drv string, name string, dottype string, datatype string, info string) string {
	o := orm.NewOrm()
	var ob Dot
	ob.Name = name
	ob.Dottype = dottype
	ob.Datatype = datatype
	ob.Info = info
	ob.Drv = drv
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, _, err := o.ReadOrCreate(&ob, "Name"); err == nil {
		if created {
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
	_, err := o.Raw("SELECT name FROM drv").QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getusrdrvinfo(name string) (string, error) {
	var ob []UsrDrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM usrdrv WHERE usrname = ?", name).QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
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
