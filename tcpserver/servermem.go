package tcpserver

import (
	"encoding/json"

	"github.com/astaxie/beego"
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
	index = TSDrvMap[name] //从态神的设备列表中查找
	//需要索引不同的设备类型的设备列表
	if index != 0 {
		//再态神的列表中查找到了设备
		ob = TSDrv[index-1].Dot
	} else {
		index = MBDrvMap[name]
		if index != 0 {
			ob = MBDrv[index-1].Dot
		}
	}
	str, err := json.Marshal(ob)
	return string(str), err
}
func GetUserDrvsFromMem(name interface{}) (string, int, error) {
	for TSHotReFrashFalg {
		beego.Info("正在进行热更新")
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
		for i := 0; i < len(MainUser[index].Drv); i++ {
			drvindex := TSDrvMap[MainUser[index].Drv[i].Name]
			if drvindex != 0 {
				tmpuserdrv = append(tmpuserdrv, TSDrv[drvindex-1])
			}
		}
		for i := 0; i < len(MainUser[index].Drv); i++ {
			drvindex := MBDrvMap[MainUser[index].Drv[i].Name]
			if drvindex != 0 {
				tmpuserdrv = append(tmpuserdrv, MBDrv[drvindex-1])
			}
		}
	}
	data, err := json.Marshal(map[string]interface{}{"User": name, "Drv": tmpuserdrv})

	return string(data), len(data), err
}
func GetUserDrvInfoFromMem(user interface{}, drv string) (string, int, error) {
	for TSHotReFrashFalg {
		beego.Info("正在进行热更新")
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
			drvindex := TSDrvMap[drv]
			if drvindex != 0 {
				tmpuserdrv = TSDrv[drvindex-1]
			} else {
				drvindex = MBDrvMap[drv]
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
		beego.Info("正在进行热更新")
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
			drvindex := TSDrvMap[drv]
			if drvindex != 0 {
				tmpuserdrv = TSDrv[drvindex-1]
			} else {
				drvindex = MBDrvMap[drv]
				if drvindex != 0 {
					tmpuserdrv = MBDrv[drvindex-1]
				}
			}
		}
	}

	data, err := json.Marshal(tmpuserdrv.Dot)

	return string(data), len(data), err
}
