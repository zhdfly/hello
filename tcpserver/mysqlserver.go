package tcpserver

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func ConfigSQL() {
	mysqluser := beego.AppConfig.String("mysqluser")
	mysqlpass := beego.AppConfig.String("mysqlpass")
	mysqlurls := beego.AppConfig.String("mysqlurls")
	mysqldb := beego.AppConfig.String("mysqldb")
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterModel(new(Usr), new(Usrdrv), new(Dotvalue), new(Videodrv), new(Maindrv), new(Maindot), new(Userarea), new(Drvdotalarm))

	orm.RegisterDataBase("default", "mysql", mysqluser+":"+mysqlpass+"@tcp("+mysqlurls+")/"+mysqldb+"?charset=utf8&loc=Local")

	//orm.Debug = true
}
func NewAlarmNotice(drv string, dot string, msg string) string {
	o := orm.NewOrm()
	var ob Drvdotalarm
	ob.Drv = drv
	ob.Dot = dot
	ob.Status = "W"
	ob.Msg = msg
	ob.Time = time.Now().Format("2006-01-02 15:04:05")

	if created, _, err := o.ReadOrCreate(&ob, "Drv", "Dot", "Time"); err == nil {
		if created {
			NewAlarmToUser(ob)
			return "OK"
		} else {
			return "ERR"
		}
	}
	return "ERR"
}
func ActAlarmNotice(drv string, dot string, time string) string {
	o := orm.NewOrm()
	if _, err := o.Raw("UPDATE drvdotalarm SET status='O' where drv=? and dot=? and time=?", drv, dot, time).Exec(); err == nil {
		DelAlarmToUser(drv, dot, time)
		return "OK"
	}
	return "ERR"
}
func NewArea(user string, area string, drvs string) string {
	o := orm.NewOrm()
	var ob Userarea
	ob.User = user
	ob.Area = area
	ob.Drvs = drvs
	UserAreaList = append(UserAreaList, ob)
	if created, _, err := o.ReadOrCreate(&ob, "User", "Area"); err == nil {
		if created {
			return "OK"
		} else {
			return "ERR"
		}
	}
	return "ERR"
}
func UpdataArea(user string, area string, drvs string) string {
	o := orm.NewOrm()
	var ob Userarea
	ob.User = user
	ob.Area = area
	ob.Drvs = drvs

	if _, err := o.Raw("UPDATE userarea SET drvs=? where user=? and area=?", drvs, user, area).Exec(); err == nil {
		UpdataAreaList(ob)
		return "OK"
	}
	return "ERR"
}
func DeleteArea(user string, area string) string {
	o := orm.NewOrm()

	if _, err := o.Raw("DELETE FROM userarea  where user=? and area=?", user, area).Exec(); err == nil {
		DeleteAreaList(user, area)
		return "OK"
	}
	return "ERR"
}
func GetArea() []Userarea {
	o := orm.NewOrm()
	var ob []Userarea
	_, err := o.Raw("SELECT * FROM userarea").QueryRows(&ob)
	if err != nil {
		logs.Info(err)
	}
	return ob
}
func Inserttousr(usr string, pass string) string {
	o := orm.NewOrm()
	user := Usr{Name: usr, Pass: pass}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	if created, _, err := o.ReadOrCreate(&user, "Name"); err == nil {
		if created {
			CreatNewUser(usr)
			return "OK"
		} else {
			return "ERR"
		}
	}

	return "ERR"
}
func UpdataUserPass(user string, pass string) string {
	o := orm.NewOrm()
	//UPDATE `maingo`.`maindrv` SET `cmittype`='透传' WHERE `Id`=22;
	sqlstr := "UPDATE usr SET pass='" + pass + "' WHERE name='" + user + "'"
	logs.Info(sqlstr)
	_, err := o.Raw(sqlstr).Exec()
	if err == nil {
		return "OK"
	}
	return "ERR"
}
func DelUser(user string) string {
	if user == "admin" {
		return "ERR"
	}
	o := orm.NewOrm()
	//UPDATE `maingo`.`maindrv` SET `cmittype`='透传' WHERE `Id`=22;
	//删除内存中的数据
	sqlstr := "delete from usr WHERE name='" + user + "'"
	logs.Info(sqlstr)
	_, err := o.Raw(sqlstr).Exec()
	sqlstr = "delete from usrdrv WHERE usrname='" + user + "'"
	logs.Info(sqlstr)
	_, err = o.Raw(sqlstr).Exec()
	if err == nil {
		return "OK"
	}
	return "ERR"
}

// CREATE TABLE `maingo`.`modbus1_dot` (
// 	`Id` int(11) NOT NULL AUTO_INCREMENT,
// 	`dotname` varchar(20) NULL DEFAULT NULL,
// 	`time` datetime NULL DEFAULT NULL,
// 	`空温` float(11,4) NOT NULL DEFAULT 0,
// 	  `空温1` float(11,4) NOT NULL DEFAULT 0,
// 	  `空温2` float(11,4) NOT NULL DEFAULT 0,
// 	  `空温3` float(11,4) NOT NULL DEFAULT 0,
// 	PRIMARY KEY (`Id`)
//   )  DEFAULT CHARSET=utf8  COLLATE=utf8_general_ci;

func InserttoMainDrv(user interface{}, name string, addr int, port int, types string, cmittype string, idcode int, polltime int) string {
	o := orm.NewOrm()
	var ob Maindrv
	ob.Name = name
	ob.Port = port
	ob.Addr = addr
	ob.Packtype = types
	ob.Cmittype = cmittype
	ob.Idcode = idcode
	ob.Polltime = polltime
	ob.Time = time.Now().Format("2006-01-02 15:04:05")
	//把新增加的设备添加到数据库中
	if created, _, err := o.ReadOrCreate(&ob, "Name"); err == nil {
		if created {
			//为新设备添加新的历史数据表
			_, err = o.Raw("CREATE TABLE `maingo`.`" + name + "_dot` (`Id` int(11) NOT NULL AUTO_INCREMENT,`time` datetime NULL DEFAULT NULL,PRIMARY KEY (`Id`))  DEFAULT CHARSET=utf8  COLLATE=utf8_general_ci;").Exec()
			_, err = o.Raw("CREATE TABLE `maingo`.`" + name + "_dotvalueday` (`Id` int(11) NOT NULL AUTO_INCREMENT,`dotname` varchar(20) NULL DEFAULT NULL,`high` float(11,4) NOT NULL DEFAULT 0,`low` float(11,4) NOT NULL DEFAULT 0,`time` datetime NOT NULL DEFAULT '0000-00-00 00:00:00',PRIMARY KEY (`Id`))  DEFAULT CHARSET=utf8  COLLATE=utf8_general_ci;").Exec()
			_, err = o.Raw("CREATE TABLE `maingo`.`" + name + "_dotavgvaluehour` (`Id` int(11) NOT NULL AUTO_INCREMENT,`dotname` varchar(20) NULL DEFAULT NULL,`value` float(11,4) NOT NULL DEFAULT 0,`time` datetime NULL DEFAULT NULL,PRIMARY KEY (`Id`))  DEFAULT CHARSET=utf8  COLLATE=utf8_general_ci;").Exec()
			var userdrv Usrdrv
			userdrv.Drvname = name
			userdrv.Usrname = "admin"
			_, _, err = o.ReadOrCreate(&userdrv, "Id")
			CreatNewDrv(user, name, ob)
			if err == nil {
				return "OK"
			}
		} else {
			return "ERR"
		}
	} else {
		logs.Info(err)
	}
	//除了把设备信息添加到数据库中去，还需要把内存中的数据进行热更新

	//热更新数据完成
	return "ERR"
}
func UpdataMainDrv(user interface{}, name string, addr int, port int, types string, cmittype string, idcode int, polltime int) string {
	o := orm.NewOrm()
	var ob Maindrv
	ob.Name = name
	ob.Port = port
	ob.Addr = addr
	ob.Packtype = types
	ob.Cmittype = cmittype
	ob.Idcode = idcode
	ob.Polltime = polltime
	ob.Time = time.Now().Format("2006-01-02 15:04:05")
	//把新增加的设备添加到数据库中
	if _, err := o.Update(&ob); err == nil {
		CreatNewDrv(user, name, ob)
		return "OK"
	} else {
		logs.Info(err)
	}
	//除了把设备信息添加到数据库中去，还需要把内存中的数据进行热更新

	//热更新数据完成
	return "ERR"
}
func PointDrv(drv string) (string, error) {
	o := orm.NewOrm()
	var err error
	if _, err = o.Raw("UPDATE maindrv SET point=1 where name=?", drv).Exec(); err == nil {
		PointDrvMem(drv)
		return "OK", err
	}
	return "ERR", err
}
func UnPointDrv(drv string) (string, error) {
	o := orm.NewOrm()
	var err error
	if _, err = o.Raw("UPDATE maindrv SET point=0 where name=?", drv).Exec(); err == nil {
		UnPointDrvMem(drv)
		return "OK", err
	}
	return "ERR", err
}
func UploadAreaPlaneBk(area string, bk string) (string, error) {
	o := orm.NewOrm()
	var err error
	if _, err = o.Raw("UPDATE userarea SET planebk=? where area=?", bk, area).Exec(); err == nil {
		UploadAreaPlaneBkMem(area, bk)
		return "OK", err
	}
	return "ERR", err
}
func GetDrvAlarm(drv string) string {
	var ob []Drvdotalarm
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM drvdotalarm where drv=?", drv).QueryRows(&ob)
	if err == nil {
		//logs.Info(ob)
	}
	str, _ := json.Marshal(ob)
	return string(str)
}
func UpdataPlaneDrv(drv string, x int, y int) (string, error) {
	o := orm.NewOrm()
	var err error
	if _, err = o.Raw("UPDATE maindrv SET x=?, y=? where name=?", x, y, drv).Exec(); err == nil {
		UpdataPlaneDrvMem(drv, x, y)
		return "OK", err
	}
	logs.Info(err)
	return "ERR", err
}
func DeleteDrv(drv string) (string, error) {
	//删除设备，需要删除设备的所有衍生表，用户设备记录，主设备表中记录，设备的数据点
	//先从内存中删除，删除数据库过程中出错
	DeleteDrvFromMem(drv)
	o := orm.NewOrm()
	sql := "DROP TABLE `maingo`.`" + drv + "_dot`,`maingo`.`" + drv + "_dotavgvaluehour`,`maingo`.`" + drv + "_dotvalueday`" //删除衍生表
	logs.Info(sql)
	_, err := o.Raw(sql).Exec()
	if err != nil {
		logs.Info(err)
		return "ERR", err
	}
	sql = "delete  from maingo.maindrv where name='" + drv + "'" //删除设备表中数据
	logs.Info(sql)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Info(err)
		return "ERR", err
	}
	sql = "delete  from maindot where drvname='" + drv + "'" //删除数据点表中的数据
	logs.Info(sql)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Info(err)
		return "ERR", err
	}
	sql = "delete  from usrdrv where drvname='" + drv + "'" //删除用户设备表中的数据
	logs.Info(sql)
	_, err = o.Raw(sql).Exec()
	if err != nil {
		logs.Info(err)
		return "ERR", err
	}
	DeleteDrvFromMem(drv)
	logs.Info(drv, err)
	return "OK", err
}

// ALTER TABLE `maingo`.`modbus_dot`
//   ADD COLUMN `空湿2` float(11,3) NOT NULL DEFAULT 0;
func Inserttodot(drv string, name string, addr int, rw int, dtype int, data int, top float32, bot float32, vtop float32, vbot float32, unit string, bynum float32) string {
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
	ob.Savetime = 30
	ob.Valuetop = vtop
	ob.Valuebot = vbot
	ob.Bynum = bynum
	created, _, err := o.ReadOrCreate(&ob, "Name", "Drvname")
	if err == nil {
		if created {
			//
			CreatNewDot(drv, ob)
			//增加设备数据点，需要在统计MAP中增加相应的点
			SetOneDotIntoMap(drv, name)
			//需要更新关联设备的数据信息，此处需要判断设备类型
			UpdataDrvDataInfo(drv)
			//需要在设备数据表中增加新建的数据点的字段
			_, err = o.Raw("ALTER TABLE `maingo`.`" + drv + "_dot` ADD COLUMN `" + drv + "_" + name + "` float(11,3) NOT NULL DEFAULT 0;").Exec()
			return "OK"
		} else {
			return "ERR"
		}
	}
	logs.Info(err)
	return "ERR"
}
func Updatadot(drv string, name string, addr int, rw int, dtype int, data int, top float32, bot float32, vtop float32, vbot float32, unit string, bynum float32) string {
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
	ob.Savetime = 30
	ob.Valuetop = vtop
	ob.Valuebot = vbot
	ob.Bynum = bynum
	_, err = o.Raw("UPDATE maindot SET type=?,rw=?,data=?,unit=?,addr=?,valuetop=?,valuebot=?,bynum=?,alarmtop = ?,alarmbot = ? where name = ? and drvname = ?", dtype, rw, data, unit, addr, vtop, vbot, bynum, top, bot, name, drv).Exec()
	if err == nil {
		UpdataDotToMem(drv, name, ob)
		return "OK"
	}
	logs.Info(err)
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
func TimerSaveDrvDotValue(drv string, drvtype int) string {
	o := orm.NewOrm()
	var sqlstr = "INSERT INTO `maingo`.`" + drv + "_dot` SET "
	if drvtype == 0 {
		//态神设备
		TSDrvMap.M.Lock()
		index := TSDrvMap.U[drv]
		TSDrvMap.M.Unlock()
		if index > 0 {
			for _, c := range TSDrv[index-1].Dot {
				var dstr = ""
				dstr += drv + "_" + c.Name
				dstr += "="
				dstr += strconv.FormatFloat(float64(c.Value), 'E', -1, 32)
				dstr += ","
				sqlstr += dstr
			}
			sqlstr += "time='" + time.Now().Format("2006-01-02 15:04:05") + "'"
		}
	} else if drvtype == 1 {
		//态神设备
		MBDrvMap.M.Lock()
		index := MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index > 0 {
			for _, c := range MBDrv[index-1].Dot {
				var dstr = ""
				dstr += drv + "_" + c.Name
				dstr += "="
				dstr += strconv.FormatFloat(float64(c.Value), 'E', -1, 32)
				dstr += ","
				sqlstr += dstr
			}
			sqlstr += "time='" + time.Now().Format("2006-01-02 15:04:05") + "'"
		}
	}
	//`dotname`='12',`time`='1899-12-30 02:00:00',`数据1`=0.000;
	//ob := Dotvalue{Drvname: drv, Dotname: name, Value: value, Status: status, Time: time.Now().Format("2006-01-02 15:04:05")}
	// 三个返回参数依次为：是否新创建的，对象 Id 值，错误
	logs.Info(sqlstr)
	if _, err := o.Raw(sqlstr).Exec(); err == nil {
		return "OK"
	} else {
		return "ERR"

	}
	return "ERR"
}
func Getusrinfo() (string, error) {
	var ob []Usr
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM usr").QueryRows(&ob)
	if err == nil {
		//logs.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getdrvdotinfo(name string) (string, error) {
	var ob []Maindot
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindot WHERE drv = ?", name).QueryRows(&ob)
	if err == nil {
		//logs.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

func Getusrdrvinfo(name string) (string, error) {
	var ob []Usrdrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM usrdrv WHERE usrname = ?", name).QueryRows(&ob)
	if err == nil {
		//logs.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Setusrdrv(name string, drvstring string) (string, error) {
	var ob []Usrdrv
	//usr := Usrdrv{Usrname: name}
	o := orm.NewOrm()
	o.Raw("DELETE FROM usrdrv WHERE usrname = ?", name).Exec()
	logs.Info(name, drvstring)
	json.Unmarshal([]byte(drvstring), &ob)
	logs.Info(ob)
	_, err := o.InsertMulti(len(ob), ob)
	AddUserSelectDrv(name, ob)
	return "OK", err
}
func Getusrnotdrvinfo(name string) (string, error) {
	var ob []Maindrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindrv where name not in (select drvname from usrdrv where usrname = ? )", name).QueryRows(&ob)
	if err == nil {
		logs.Info(ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Getusralldrvinfo(name string) (string, error) {
	var ob []Maindrv
	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM maindrv").QueryRows(&ob)
	if err == nil {
		logs.Info(name, ob)
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func Dltdrvdot(drv string, dot string) string {
	o := orm.NewOrm()
	_, err := o.Raw("DELETE FROM maindot WHERE name = ? and drvname = ?", dot, drv).Exec()
	if err == nil {
		//ALTER TABLE `maingo`.`modbus_dot` DROP COLUMN 光照
		//删除设备的数据点
		DelDrvDot(drv, dot)
		_, err = o.Raw("ALTER TABLE `maingo`.`" + drv + "_dot` DROP COLUMN `" + drv + "_" + dot + "`").Exec()
		_, err = o.Raw("DELETE FROM " + drv + "_dotavgvaluehour WHERE dotname='" + dot + "'").Exec()
		_, err = o.Raw("DELETE FROM " + drv + "_dotvalueday WHERE dotname='" + dot + "'").Exec()
		return "OK"
	}
	return "ERR"
}

func Getdotvalue(drv string, dot string, start string, stop string) (string, error) {
	var ob []string
	var obtime []string
	logs.Info(drv, dot, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT value FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, dot, start, stop).QueryRows(&ob)
	if err == nil {
		logs.Info(ob)
	}
	_, err = o.Raw("SELECT time FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, dot, start, stop).QueryRows(&obtime)
	if err == nil {
		logs.Info(ob)
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
	logs.Info(drv, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM maindot where drvname=?", drv).QueryRows(&dotname)
	if err == nil {
		logs.Info(dotname)
	}
	for c, i := range dotname {
		logs.Info(c, i)
		logs.Info("读取时间")
		_, err := o.Raw("SELECT time FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, i, start, stop).QueryRows(&dottime)
		if err != nil {
			logs.Info(err)
		}
		logs.Info("读取数据")
		_, err = o.Raw("SELECT value FROM dotvalue where drvname=? and dotname=? and time >= ? and time <= ?", drv, i, start, stop).QueryRows(&dotvalue)
		if err != nil {
			logs.Info(err)
		}
		logs.Info("建立列表")
		var tmp Dotvaluertl
		tmp.Name = i
		tmp.Data = dotvalue
		tmp.Time = dottime
		dotlist = append(dotlist, tmp)
	}
	logs.Info("start")
	var valuemaps = make(map[int]map[string]string)
	if len(dotlist) > 0 {
		numlen := len(dotlist[0].Data)
		for i := 0; i < len(dotlist); i++ {
			if numlen > len(dotlist[i].Data) {
				numlen = len(dotlist[i].Data)
			}
		}
		for j := 0; j < numlen; j++ {
			var s = make(map[string]string)
			for i := 0; i < len(dotlist); i++ {
				s[dotlist[i].Name] = dotlist[i].Data[j]
			}
			s["Time"] = dotlist[0].Time[j]
			valuemaps[j] = s
		}
	}
	logs.Info("stop")
	str, err := json.Marshal(dotlist)
	return string(str), err
}

func GetalldotvalueR(drv string, start string, stop string) (string, error) {

	var dotname []string
	var dotvalue []orm.Params
	logs.Info(drv, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM maindot where drvname=?", drv).QueryRows(&dotname)
	rlt, err := o.Raw("SELECT * FROM "+drv+"_dot where time >= ? and time <= ?", start, stop).Values(&dotvalue)
	if err == nil {
		logs.Info(rlt)
	} else {
		logs.Info(err)
	}
	for j := 0; j < len(dotname); j++ {
		dotname[j] = drv + "_" + dotname[j]
	}
	str, err := json.Marshal(map[string]interface{}{"Name": drv, "Dots": dotname, "Values": dotvalue})

	return string(str), err
}
func GetdrvsalldotvalueR(drv []string, start string, stop string) (string, error) {
	var allvalue []orm.Params
	var alldotname []string
	for i := 0; i < len(drv); i++ {
		var dotname []string
		var dotvalue []orm.Params
		logs.Info(drv[i], start, stop)
		o := orm.NewOrm()
		_, err := o.Raw("SELECT name FROM maindot where drvname=?", drv[i]).QueryRows(&dotname)
		rlt, err := o.Raw("SELECT * FROM "+drv[i]+"_dot where time >= ? and time <= ?", start, stop).Values(&dotvalue)
		if err == nil {
			logs.Info(rlt)
		} else {
			logs.Info(err)
		}
		for j := 0; j < len(dotname); j++ {
			dotname[j] = drv[i] + "_" + dotname[j]
			alldotname = append(alldotname, dotname[j])
		}
		for j := 0; j < len(dotvalue); j++ {
			allvalue = append(allvalue, dotvalue[j])
		}
	}
	logs.Info(allvalue)
	str, err := json.Marshal(map[string]interface{}{"Dots": alldotname, "Values": allvalue})

	return string(str), err
}
func GetDotAvgValues(drv string, start string, stop string) (string, error) {

	var dotname []string
	var dotvalue []orm.Params
	logs.Info(drv, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM maindot where drvname=?", drv).QueryRows(&dotname)
	rlt, err := o.Raw("SELECT * FROM "+drv+"_dotavgvaluehour where time >= ? and time <= ?", start, stop).Values(&dotvalue)
	if err == nil {
		logs.Info(rlt)
	} else {
		logs.Info(err)
	}
	str, err := json.Marshal(map[string]interface{}{"Name": drv, "Dots": dotname, "Values": dotvalue})
	logs.Info(string(str))
	return string(str), err
}
func GetDotDailyValues(drv string, start string, stop string) (string, error) {

	var dotname []string
	var dotvalue []orm.Params
	logs.Info(drv, start, stop)
	o := orm.NewOrm()
	_, err := o.Raw("SELECT name FROM maindot where drvname=?", drv).QueryRows(&dotname)
	rlt, err := o.Raw("SELECT * FROM "+drv+"_dotvalueday where time >= ? and time <= ?", start, stop).Values(&dotvalue)
	if err == nil {
		logs.Info(rlt)
	} else {
		logs.Info(err)
	}
	str, err := json.Marshal(map[string]interface{}{"Name": drv, "Dots": dotname, "Values": dotvalue})
	logs.Info(string(str))
	return string(str), err
}
func Setdotwarning(drv string, dot string, top float32, bot float32) (string, error) {
	o := orm.NewOrm()
	_, err := o.Raw("UPDATE maindot SET alarmtop = ?,alarmbot = ? where name = ? and drvname = ?", top, bot, dot, drv).Exec()
	if err == nil {
		SetDrvDotAlarm(drv, dot, top, bot)

		return "OK", err
	}
	return "ERR", err
}

func SaveDailyValue(drvname string, dotname string, valuehigh float32, valuelow float32) string {
	o := orm.NewOrm()
	var err error
	var sqlstr = "INSERT INTO `maingo`.`" + drvname + "_dotvalueday` SET `dotname`='" + dotname + "',`high`=" + strconv.FormatFloat(float64(valuehigh), 'E', -1, 32) + ",`low`=" + strconv.FormatFloat(float64(valuelow), 'E', -1, 32) + ",`time`='" + time.Now().Add(-1e10).Format("2006-01-02") + "';"

	_, err = o.Raw(sqlstr).Exec()
	if err == nil {
		return "OK"
	}
	logs.Info(err)
	return "ERR"
}
func SaveHourlyValue(drvname string, dotname string, value float32) string {
	o := orm.NewOrm()
	var err error
	var sqlstr = "INSERT INTO `maingo`.`" + drvname + "_dotavgvaluehour` SET `dotname`='" + dotname + "',`value`=" + strconv.FormatFloat(float64(value), 'E', -1, 32) + ",`time`='" + time.Now().Add(-1e10).Format("2006-01-02 15") + "';"

	_, err = o.Raw(sqlstr).Exec()
	if err == nil {
		return "OK"
	}
	logs.Info(err)
	return "ERR"
}
