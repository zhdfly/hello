package tcpserver

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego/logs"
)

//从内存数据中获取所有设备的基本信息
func GetMainDrvInfoFromMem(user interface{}) (string, error) {
	var err error
	var rlt []byte
	for i := 0; i < len(MainUser); i++ {
		if MainUser[i].User == user {
			rlt, err = json.Marshal(MainUser[i])
		}
	}
	return string(rlt), err
}

func GetDrvDotFromMem(name string) (string, error) {
	var ob []Maindot
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[name] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		ob = TSDrv[index-1].Dot
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[name]
		MBDrvMap.M.Unlock()
		if index != 0 {
			ob = MBDrv[index-1].Dot
		}
	}
	str, err := json.Marshal(ob)
	return string(str), err
}

type UserAreaDrv struct {
	Name    string
	Drvsstr []string
	Drvs    []MainDrvType
	Planebk string
}

func GetUserAreaDrvInfoFromMem(name interface{}) string {
	for TSHotReFrashFalg {
		logs.Info("正在进行热更新")
	}
	var tmpuserdrv []MainDrvType
	var tmpareadrv []UserAreaDrv

	index := 0
	indexarea := 0
	//先找到目标用户
	for index = 0; index < len(MainUser); index++ {
		if MainUser[index].User == name {
			break
		}
	}
	for indexarea = 0; indexarea < len(UserAreaList); indexarea++ {
		if UserAreaList[indexarea].User == name {
			var tmp UserAreaDrv
			tmp.Name = UserAreaList[indexarea].Area
			tmp.Drvsstr = strings.Split(UserAreaList[indexarea].Drvs, ",")
			tmp.Planebk = UserAreaList[indexarea].Planebk
			tmpareadrv = append(tmpareadrv, tmp)
		}
	}
	//logs.Info(tmpareadrv)
	if index < len(MainUser) {
		//需要索引不同的设备类型的设备列表
		TSDrvMap.M.Lock()
		for i := 0; i < len(MainUser[index].Drv); i++ {
			drvindex := TSDrvMap.U[MainUser[index].Drv[i].Name]
			if drvindex != 0 {
				tmpuserdrv = append(tmpuserdrv, TSDrv[drvindex-1])
			}
		}
		TSDrvMap.M.Unlock()
		MBDrvMap.M.Lock()
		for i := 0; i < len(MainUser[index].Drv); i++ {
			drvindex := MBDrvMap.U[MainUser[index].Drv[i].Name]
			if drvindex != 0 {
				tmpuserdrv = append(tmpuserdrv, MBDrv[drvindex-1])
			}
		}
		MBDrvMap.M.Unlock()
	}
	for j := 0; j < len(tmpareadrv); j++ {
		for k := 0; k < len(tmpareadrv[j].Drvsstr); k++ {
			for i := 0; i < len(tmpuserdrv); i++ {
				if tmpuserdrv[i].Drvname == tmpareadrv[j].Drvsstr[k] {
					tmpareadrv[j].Drvs = append(tmpareadrv[j].Drvs, tmpuserdrv[i])
				}
			}
		}
	}

	//logs.Info(tmpareadrv)
	data, _ := json.Marshal(map[string]interface{}{"User": name, "Drv": tmpuserdrv, "Area": tmpareadrv})

	return string(data)
}

type AreaDrvNum struct {
	Area string
	Num  int
}
type OnlineStatus struct {
	Status string
	Num    int
}

func GetUserMainMsg(name interface{}) string {
	var tmpalarm []Drvdotalarm
	var tmpareanum []AreaDrvNum
	var onoffline []OnlineStatus
	var tmpuserdrv []MainDrvType
	onlineNum := 0
	offlinenum := 0
	for index := 0; index < len(MainUser); index++ {
		if MainUser[index].User == name {
			tmpalarm = MainUser[index].Alarm
			TSDrvMap.M.Lock()
			for i := 0; i < len(MainUser[index].Drv); i++ {
				drvindex := TSDrvMap.U[MainUser[index].Drv[i].Name]
				if drvindex != 0 {
					if TSDrv[drvindex-1].Online {
						onlineNum++
					} else {
						offlinenum++
					}
					if TSDrv[drvindex-1].Drv.Point == 1 {
						tmpuserdrv = append(tmpuserdrv, TSDrv[drvindex-1])
					}
				}
			}
			TSDrvMap.M.Unlock()
			MBDrvMap.M.Lock()
			for i := 0; i < len(MainUser[index].Drv); i++ {
				drvindex := MBDrvMap.U[MainUser[index].Drv[i].Name]
				if drvindex != 0 {
					if MBDrv[drvindex-1].Online {
						onlineNum++
					} else {
						offlinenum++
					}
					if MBDrv[drvindex-1].Drv.Point == 1 {
						tmpuserdrv = append(tmpuserdrv, MBDrv[drvindex-1])
					}
				}
			}
			MBDrvMap.M.Unlock()
		}
	}
	for indexarea := 0; indexarea < len(UserAreaList); indexarea++ {
		if UserAreaList[indexarea].User == name {
			var tmp AreaDrvNum
			tmp.Num = len(strings.Split(UserAreaList[indexarea].Drvs, ","))
			tmp.Area = UserAreaList[indexarea].Area
			tmpareanum = append(tmpareanum, tmp)
		}
	}
	var online OnlineStatus
	online.Status = "在线"
	online.Num = onlineNum
	onoffline = append(onoffline, online)
	var offline OnlineStatus
	offline.Status = "离线"
	offline.Num = offlinenum
	onoffline = append(onoffline, offline)
	data, _ := json.Marshal(map[string]interface{}{"User": name, "Alarm": tmpalarm, "AreaNum": tmpareanum, "OnLine": onoffline, "Point": tmpuserdrv})
	return string(data)
}
func GetUserDrvsFromMem(name interface{}) (string, int, error) {
	for TSHotReFrashFalg {
		logs.Info("正在进行热更新")
	}
	var tmpuserdrv []MainDrvType

	index := 0
	for index = 0; index < len(MainUser); index++ {
		if MainUser[index].User == name {
			break
		}
	}
	if index < len(MainUser) {
		//需要索引不同的设备类型的设备列表
		TSDrvMap.M.Lock()
		for i := 0; i < len(MainUser[index].Drv); i++ {
			drvindex := TSDrvMap.U[MainUser[index].Drv[i].Name]
			if drvindex != 0 {
				tmpuserdrv = append(tmpuserdrv, TSDrv[drvindex-1])
			}
		}
		TSDrvMap.M.Unlock()
		MBDrvMap.M.Lock()
		for i := 0; i < len(MainUser[index].Drv); i++ {
			drvindex := MBDrvMap.U[MainUser[index].Drv[i].Name]
			if drvindex != 0 {
				tmpuserdrv = append(tmpuserdrv, MBDrv[drvindex-1])
			}
		}
		MBDrvMap.M.Unlock()
	}
	data, err := json.Marshal(map[string]interface{}{"User": name, "Drv": tmpuserdrv})

	return string(data), len(data), err
}

func GetUserDrvInfoFromMem(user interface{}, drv string) (string, int, error) {
	for TSHotReFrashFalg {
		logs.Info("正在进行热更新")
	}
	var tmpuserdrv MainDrvType
	var ishavedrv = false
	index := 0
	for index = 0; index < len(MainUser); index++ {
		if MainUser[index].User == user {
			break
		}
	}

	if index < len(MainUser) {
		for i := 0; i < len(MainUser[index].Drv); i++ {
			if MainUser[index].Drv[i].Name == drv {
				ishavedrv = true
			}
		}
		if ishavedrv {
			//需要索引不同的设备类型的设备列表
			TSDrvMap.M.Lock()
			drvindex := TSDrvMap.U[drv]
			TSDrvMap.M.Unlock()
			if drvindex != 0 {
				tmpuserdrv = TSDrv[drvindex-1]
			} else {
				MBDrvMap.M.Lock()
				drvindex = MBDrvMap.U[drv]
				MBDrvMap.M.Unlock()
				if drvindex != 0 {
					tmpuserdrv = MBDrv[drvindex-1]
				}
			}
		}
	}

	data, err := json.Marshal(tmpuserdrv)

	return string(data), len(data), err
}

func GetUserDrvDotInfoFromMem(user interface{}, drv string) (string, int, error) {
	for TSHotReFrashFalg {
		logs.Info("正在进行热更新")
	}
	var tmpuserdrv MainDrvType
	var ishavedrv = false
	index := 0
	for index = 0; index < len(MainUser); index++ {
		if MainUser[index].User == user {
			break
		}
	}

	if index < len(MainUser) {
		for i := 0; i < len(MainUser[index].Drv); i++ {
			if MainUser[index].Drv[i].Name == drv {
				ishavedrv = true
			}
		}
		if ishavedrv {
			//需要索引不同的设备类型的设备列表
			TSDrvMap.M.Lock()
			drvindex := TSDrvMap.U[drv]
			TSDrvMap.M.Unlock()
			if drvindex != 0 {
				tmpuserdrv = TSDrv[drvindex-1]
			} else {
				MBDrvMap.M.Lock()
				drvindex = MBDrvMap.U[drv]
				MBDrvMap.M.Unlock()
				if drvindex != 0 {
					tmpuserdrv = MBDrv[drvindex-1]
				}
			}
		}
	}

	data, err := json.Marshal(tmpuserdrv.Dot)

	return string(data), len(data), err
}
func UpdataDrvDataInfo(name string) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[name] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		CreateTsUrl(&TSDrv[index-1])
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[name]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			CreateMbCmd(&MBDrv[index-1])
		}
	}
}
func UpdataDotToMem(drv string, name string, ob Maindot) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		for i := 0; i < len(TSDrv[index-1].Dot); i++ {
			if TSDrv[index-1].Dot[i].Name == name {
				TSDrv[index-1].Dot = append(TSDrv[index-1].Dot[:i], TSDrv[index-1].Dot[i+1:]...)
				TSDrv[index-1].Dot = append(TSDrv[index-1].Dot, ob)
			}
		}
		CreateTsUrl(&TSDrv[index-1])
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			for i := 0; i < len(MBDrv[index-1].Dot); i++ {
				if MBDrv[index-1].Dot[i].Name == name {
					MBDrv[index-1].Dot = append(MBDrv[index-1].Dot[:i], MBDrv[index-1].Dot[i+1:]...)
					MBDrv[index-1].Dot = append(MBDrv[index-1].Dot, ob)
				}
			}
			CreateMbCmd(&MBDrv[index-1])
		}
	}
}
func DelDrvDot(name string, dot string) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[name] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		for i := 0; i < len(TSDrv[index-1].Dot); i++ {
			if TSDrv[index-1].Dot[i].Name == dot {
				TSDrv[index-1].Dot = append(TSDrv[index-1].Dot[:i], TSDrv[index-1].Dot[i+1:]...)
			}
		}
		CreateTsUrl(&TSDrv[index-1])
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[name]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			for i := 0; i < len(MBDrv[index-1].Dot); i++ {
				if MBDrv[index-1].Dot[i].Name == dot {
					MBDrv[index-1].Dot = append(MBDrv[index-1].Dot[:i], MBDrv[index-1].Dot[i+1:]...)
				}
			}
			CreateMbCmd(&MBDrv[index-1])
		}
	}
}
func SetDrvDotAlarm(drv string, dot string, top float32, bot float32) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		for i := 0; i < len(TSDrv[index-1].Dot); i++ {
			if dot == TSDrv[index-1].Dot[i].Name {
				TSDrv[index-1].Dot[i].Alarmtop = top
				TSDrv[index-1].Dot[i].Alarmbot = bot
			}
		}
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			for i := 0; i < len(MBDrv[index-1].Dot); i++ {
				if dot == MBDrv[index-1].Dot[i].Name {
					MBDrv[index-1].Dot[i].Alarmtop = top
					MBDrv[index-1].Dot[i].Alarmbot = bot
				}
			}
		}
	}
}
func CreatNewDot(drv string, dot Maindot) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		TSDrv[index-1].Dot = append(TSDrv[index-1].Dot, dot)
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			MBDrv[index-1].Dot = append(MBDrv[index-1].Dot, dot)
		}
	}
}
func CreatNewDrv(user interface{}, drvname string, ob Maindrv) {
	//为用户动态追加新的设备
	for i := 0; i < len(MainUser); i++ {
		if MainUser[i].User == user {
			MainUser[i].Drv = append(MainUser[i].Drv, ob)
			logs.Info(MainUser[i])
		}
	}
	var drvtmp MainDrvType
	drvtmp.Drvname = drvname
	drvtmp.Drv = ob
	if ob.Packtype == "MODBUS" {
		MBDrv = append(MBDrv, drvtmp)
		MBDrvMap.M.Lock()
		MBDrvCount := len(MBDrvMap.U)
		MBDrvMap.U[drvname] = MBDrvCount + 1
		MBDrvMap.M.Unlock()
		MBHotReFrashFalg = true
		logs.Info(MBDrv, MBDrvMap)
	} else if ob.Packtype == "态神" {
		TSDrv = append(TSDrv, drvtmp)
		TSDrvMap.M.Lock()
		TSDrvCount := len(TSDrvMap.U)
		TSDrvMap.U[drvname] = TSDrvCount + 1
		TSDrvMap.M.Unlock()
		TSHotReFrashFalg = true
		logs.Info(TSDrv, TSDrvMap)
	}
}
func DeleteDrvFromMem(drv string) {
	index := 0
	//
	//不能直接删除设备列表中的对象，否则会对MAP索引造成影响，
	//可以清空设备列表中对象的DOT列表，这样就不会再进行数据的读取
	//
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		TSDrv[index-1].Dot = nil
		TSHotReFrashFalg = true
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			MBDrv[index-1].Dot = nil
			MBDrv[index-1].MBCmds = nil
			MBDrv[index-1].MBCmdNum = 0
			MBDrv[index-1].Drv.Port = 0
			//热更新一下
			MBHotReFrashFalg = true
		}
	}
	//还需要对用户设备列表进行删除
	for i := 0; i < len(MainUser); i++ {
		for j := 0; j < len(MainUser[i].Drv); j++ {
			if MainUser[i].Drv[j].Name == drv {
				MainUser[i].Drv = append(MainUser[i].Drv[:j], MainUser[i].Drv[j+1:]...)
			}
		}
	}
	//需要从区域列表中删除

	for i := 0; i < len(UserAreaList); i++ {
		drvstr := ""
		if len(UserAreaList[i].Drvs) > 0 {
			drvlist := strings.Split(UserAreaList[i].Drvs, ",")
			for j := 0; j < len(drvlist); j++ {
				if drvlist[j] == drv {
					for k := 0; k < len(drvlist); k++ {
						if k != j {
							if drvstr == "" {
								drvstr += drvlist[k]
							} else {
								drvstr += "," + drvlist[k]
							}
						}
					}
					UserAreaList[i].Drvs = drvstr
					UpdataArea(UserAreaList[i].User, UserAreaList[i].Area, UserAreaList[i].Drvs)
				}
			}
		}

	}
	//统计信息中的MAP中的信息不需要删除
}
func CreatNewUser(username string) {
	//追加新用户
	var usertmp MainUserDrv
	usertmp.User = username
	// //为用户增加设备信息 admin用户下拥有所有设备的信息
	MainUser = append(MainUser, usertmp)
	MainUserMAP.M.Lock()
	MainUserMAP.U[username] = len(MainUserMAP.U)
	MainUserMAP.M.Unlock()
}
func DeleteOldUser(username string) {
	MainUserMAP.M.Lock()
	index := MainUserMAP.U[username]
	MainUserMAP.M.Unlock()
	if index != 0 {
		MainUser[index-1].User = ""
		MainUser[index-1].Drv = nil
	}
}
func AddUserSelectDrv(user string, drvlist []Usrdrv) {

	//为用户增加设备信息 admin用户下拥有所有设备的信息
	MainUserMAP.M.Lock()
	userindex := MainUserMAP.U[user]
	adminindex := MainUserMAP.U["admin"]
	MainUserMAP.M.Unlock()
	for i := 0; i < len(drvlist); i++ {
		for j := 0; j < len(MainUser[adminindex-1].Drv); j++ {
			if drvlist[i].Drvname == MainUser[adminindex-1].Drv[j].Name {
				var drvtmp Maindrv
				drvtmp = MainUser[adminindex-1].Drv[j]
				MainUser[userindex].Drv = append(MainUser[userindex].Drv, drvtmp)
			}
		}
	}
}

func UpdataAreaList(ob Userarea) {
	for i := 0; i < len(UserAreaList); i++ {
		if UserAreaList[i].Area == ob.Area && UserAreaList[i].User == ob.User {
			UserAreaList[i].Drvs = ob.Drvs
		}
	}
}
func DeleteAreaList(user string, area string) {
	for i := 0; i < len(UserAreaList); i++ {
		if UserAreaList[i].Area == area && UserAreaList[i].User == user {
			UserAreaList = append(UserAreaList[:i], UserAreaList[i+1:]...)
		}
	}
}
func NewAlarmToUser(ob Drvdotalarm) {
	for i := 0; i < len(MainUser); i++ {
		for j := 0; j < len(MainUser[i].Drv); j++ {
			if MainUser[i].Drv[j].Name == ob.Drv {
				MainUser[i].Alarm = append(MainUser[i].Alarm, ob)
				break
			}
		}
	}
}
func DelAlarmToUser(drv string, dot string, time string) {
	for i := 0; i < len(MainUser); i++ {
		for j := 0; j < len(MainUser[i].Alarm); j++ {
			if MainUser[i].Alarm[j].Drv == drv && MainUser[i].Alarm[j].Dot == dot && MainUser[i].Alarm[j].Time == time {
				MainUser[i].Alarm = append(MainUser[i].Alarm[:j], MainUser[i].Alarm[j+1:]...)
				break
			}
		}
	}
}
func PointDrvMem(drv string) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		TSDrv[index-1].Drv.Point = 1
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			MBDrv[index-1].Drv.Point = 1
		}
	}
}
func UnPointDrvMem(drv string) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		TSDrv[index-1].Drv.Point = 0
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			MBDrv[index-1].Drv.Point = 0
		}
	}
}
func UploadAreaPlaneBkMem(area string, bk string) {
	for i := 0; i < len(UserAreaList); i++ {
		if UserAreaList[i].Area == area {
			UserAreaList[i].Planebk = bk
		}
	}
}
func UpdataPlaneDrvMem(drv string, x int, y int) {
	index := 0
	//按个查找不同得分设备列表
	TSDrvMap.M.Lock()
	index = TSDrvMap.U[drv] //从态神的设备列表中查找
	TSDrvMap.M.Unlock()
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		TSDrv[index-1].Drv.X = x
		TSDrv[index-1].Drv.Y = y
	} else {
		MBDrvMap.M.Lock()
		index = MBDrvMap.U[drv]
		MBDrvMap.M.Unlock()
		if index != 0 {
			//在MODBUS的列表中查找到了设备
			MBDrv[index-1].Drv.X = x
			MBDrv[index-1].Drv.Y = y
		}
	}
}
