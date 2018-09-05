package tcpserver

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/robfig/cron"
)

var TSDrvMap MainMAPstringint
var TSDrvMapBak MainMAPstringint
var TSDrv []MainDrvType
var TSDrvBak []MainDrvType

func TSServerConfig() {
	TSDrvMap.M.Lock()
	TSDrvMap.U = make(map[string]int)
	TSDrvMap.M.Unlock()
	TSDrvMapBak.M.Lock()
	TSDrvMapBak.U = make(map[string]int)
	TSDrvMapBak.M.Unlock()
	Gettsdotinfo()
	TSHotReFrashFalg = false

}
func CreateTsUrl(drv *MainDrvType) {
	//根据数据点类型计算数据点的分类个数
	drv.Sensornum = 0
	drv.Videonum = 0
	drv.IOnum = 0
	drv.Logicnum = 0
	for j := 0; j < len(drv.Dot); j++ {

		if drv.Dot[j].Type == inreg || drv.Dot[j].Type == outreg {
			drv.Sensornum = drv.Sensornum + 1
		} else {
			drv.IOnum = drv.IOnum + 1
		}
	}
	drv.Videonum = Getdrvvedionum(drv.Drv.Name) //获取摄像头的个数
	//生产态神读取数据的URL
	drv.TSUrls = "http://211.149.159.27:5021/html5/GETTAGVAL/"
	for j := 0; j < len(drv.Dot); j++ {
		if j == 0 {
			drv.TSUrls = drv.TSUrls + drv.Drv.Name + "." + drv.Dot[j].Name
		} else {
			drv.TSUrls = drv.TSUrls + "," + drv.Drv.Name + "." + drv.Dot[j].Name
		}
	}
	logs.Info(drv.Drv.Name, drv.TSUrls) //打印生成的URL
}
func Gettsdotinfo() {
	TSDrv = nil
	o := orm.NewOrm()
	//获取态神设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = '态神'").QueryRows(&TSDrv)
	if err != nil {
		logs.Info("ERROR:ts 001", err)
		return
	}
	//根据态神设备列表信息，补充每个态神设备的数据点和其他信息
	TSDrvMap.M.Lock()
	for i := 0; i < len(TSDrv); i++ {
		//根据态神设备的列表信息生成索引MAP
		TSDrvMap.U[TSDrv[i].Drv.Name] = i + 1
		TSDrv[i].Drvname = TSDrv[i].Drv.Name
		_, err = o.Raw("SELECT * FROM maindot where drvname=?", TSDrv[i].Drv.Name).QueryRows(&TSDrv[i].Dot)
		if err != nil {
			logs.Info("ERROR:ts 002", err)
			return
		}
	}
	TSDrvMap.M.Unlock()

	for i := 0; i < len(TSDrv); i++ {
		for j := 0; j < len(TSDrv[i].Dot); j++ {
			SetOneDotIntoMap(TSDrv[i].Drvname, TSDrv[i].Dot[j].Name)
		}
		CreateTsUrl(&TSDrv[i])

	}

}

func GettsdotinfoHot() {
	TSDrvBak = nil
	o := orm.NewOrm()
	//获取态神设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = '态神'").QueryRows(&TSDrvBak)
	if err != nil {
		logs.Info("ERROR:ts 001", err)
		return
	}
	TSDrvMapBak.M.Lock()
	//根据态神设备列表信息，补充每个态神设备的数据点和其他信息
	for i := 0; i < len(TSDrvBak); i++ {
		//根据态神设备的列表信息生成索引MAP

		TSDrvMapBak.U[TSDrvBak[i].Drv.Name] = i + 1

		TSDrvBak[i].Drvname = TSDrvBak[i].Drv.Name
		_, err = o.Raw("SELECT * FROM maindot where drvname=?", TSDrvBak[i].Drv.Name).QueryRows(&TSDrvBak[i].Dot)
		if err != nil {
			logs.Info("ERROR:ts 002", err)
			return
		}
		//根据数据点类型计算数据点的分类个数
		for j := 0; j < len(TSDrvBak[i].Dot); j++ {
			if TSDrvBak[i].Dot[j].Type == inreg || TSDrvBak[i].Dot[j].Type == outreg {
				TSDrvBak[i].Sensornum = TSDrvBak[i].Sensornum + 1
			} else {
				TSDrvBak[i].IOnum = TSDrvBak[i].IOnum + 1
			}
		}
		TSDrvBak[i].Videonum = Getdrvvedionum(TSDrvBak[i].Drv.Name) //获取摄像头的个数
		//生产态神读取数据的URL
		TSDrvBak[i].TSUrls = "http://211.149.159.27:5021/html5/GETTAGVAL/"
		for j := 0; j < len(TSDrvBak[i].Dot); j++ {
			if j == 0 {
				TSDrvBak[i].TSUrls = TSDrvBak[i].TSUrls + TSDrvBak[i].Drv.Name + "." + TSDrvBak[i].Dot[j].Name
			} else {
				TSDrvBak[i].TSUrls = TSDrvBak[i].TSUrls + "," + TSDrvBak[i].Drv.Name + "." + TSDrvBak[i].Dot[j].Name
			}
		}
		logs.Info(TSDrvBak[i].Drv.Name, TSDrvBak[i].TSUrls) //打印生成的URL
	}
	TSDrvMapBak.M.Unlock()
}

var TSHotReFrashFalg = false

func TSDrvThread(drv *MainDrvType) {
	if len(drv.Dot) > 0 {
		resp, err := http.Get(drv.TSUrls)
		if err != nil {
			// handle error
			log.Println("err:", err)

		} else {

			defer resp.Body.Close()

			buf := bytes.NewBuffer(make([]byte, 0, 512))

			_, _ = buf.ReadFrom(resp.Body)

			//logs.Info(len(buf.Bytes()))
			//logs.Info(length)
			//logs.Info(string(buf.Bytes()))
			strresult := string(buf.Bytes())
			strs := strings.Split(strresult, "|")
			drv.Flashtime = time.Now().Format("2006-01-02 15:04")
			drv.Online = true
			for n := 0; n < len(strs); n++ {
				insertValue(drv, n, strs[n])
				//Inserttodotvalue(TSDrv[Getindex].Dot[n])
			}
			TimerSaveDrvDotValue(drv.Drvname, 0)
			//logs.Info(strs)
		}
	}
}

var TSTimerMaster *cron.Cron

func InsertNewTimerThreat(drv *MainDrvType) {
	str := "*/" + strconv.FormatInt(int64(drv.Drv.Polltime), 10) + " * * * * *"
	v := drv
	logs.Info(v)
	TSTimerMaster.AddFunc(str, func() { TSDrvThread(v) })
}
func StarthttpGet() {
	TSTimerMaster := cron.New()
	for i := 0; i < len(TSDrv); i++ { //strconv.FormatInt(int64(TSDrv[i].Drv.Samplingtime), 10)
		str := "*/" + strconv.FormatInt(int64(TSDrv[i].Drv.Polltime), 10) + " * * * * *"
		v := &TSDrv[i]
		TSTimerMaster.AddFunc(str, func() { TSDrvThread(v) })
	}
	TSTimerMaster.Start()
	for {
		if TSHotReFrashFalg {
			//开始启动热更新系统参数
			TSTimerMaster.Stop()
			TSTimerMaster = cron.New()
			for i := 0; i < len(TSDrv); i++ { //strconv.FormatInt(int64(TSDrv[i].Drv.Samplingtime), 10)
				str := "*/" + strconv.FormatInt(int64(TSDrv[i].Drv.Polltime), 10) + " * * * * *"
				v := &TSDrv[i]
				TSTimerMaster.AddFunc(str, func() { TSDrvThread(v) })
			}
			TSTimerMaster.Start()
			TSHotReFrashFalg = false
		} else {
			time.Sleep(1e9)
		}
	}
}
func insertValue(drv *MainDrvType, index int, data string) {
	//deindex := 0
	strd := strings.Split(data, ",")
	if len(strd) != 2 {
		return
	}
	v, err := strconv.ParseFloat(strd[1], 32)
	if err != nil {
		logs.Info("Value Cannot Convert to float32", drv.Dot[index].Name, data)
	} else {
		drv.Dot[index].Value = float32(v)
		SetOneDotValueToMap(drv.Drvname, drv.Dot[index].Name, drv.Dot[index].Value)
		if drv.Dot[index].Alarmtop > drv.Dot[index].Alarmbot {
			if drv.Dot[index].Value > drv.Dot[index].Alarmtop {
				drv.Dot[index].Status = "TOP"
			}
			if drv.Dot[index].Value < drv.Dot[index].Alarmbot {
				drv.Dot[index].Status = "BOT"
			}
		}
	}
}
