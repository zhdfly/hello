package tcpserver

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var TSDrvMap = make(map[string]int)
var TSDrvMapBak = make(map[string]int)
var TSDrv []MainDrvType
var TSDrvBak []MainDrvType

func Gettsdotinfo() {
	TSDrv = nil
	o := orm.NewOrm()
	//获取态神设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = '态神'").QueryRows(&TSDrv)
	if err != nil {
		beego.Info("ERROR:ts 001", err)
		return
	}
	//根据态神设备列表信息，补充每个态神设备的数据点和其他信息
	for i := 0; i < len(TSDrv); i++ {
		//根据态神设备的列表信息生成索引MAP
		TSDrvMap[TSDrv[i].Drv.Name] = i + 1
		_, err = o.Raw("SELECT * FROM maindot where drvname=?", TSDrv[i].Drv.Name).QueryRows(&TSDrv[i].Dot)
		if err != nil {
			beego.Info("ERROR:ts 002", err)
			return
		}
		//根据数据点类型计算数据点的分类个数
		for j := 0; j < len(TSDrv[i].Dot); j++ {
			if TSDrv[i].Dot[j].Type == inreg {
				TSDrv[i].Sensornum = TSDrv[i].Sensornum + 1
			} else {
				TSDrv[i].IOnum = TSDrv[i].IOnum + 1
			}
		}
		TSDrv[i].Videonum = Getdrvvedionum(TSDrv[i].Drv.Name) //获取摄像头的个数
		//生产态神读取数据的URL
		TSDrv[i].TSUrls = "http://211.149.159.27:5021/html5/GETTAGVAL/"
		for j := 0; j < len(TSDrv[i].Dot); j++ {
			if j == 0 {
				TSDrv[i].TSUrls = TSDrv[i].TSUrls + TSDrv[i].Drv.Name + "." + TSDrv[i].Dot[j].Name
			} else {
				TSDrv[i].TSUrls = TSDrv[i].TSUrls + "," + TSDrv[i].Drv.Name + "." + TSDrv[i].Dot[j].Name
			}
		}
		beego.Info(TSDrv[i].Drv.Name, TSDrv[i].TSUrls) //打印生成的URL
	}
}

func GettsdotinfoHot() {
	TSDrvBak = nil
	o := orm.NewOrm()
	//获取态神设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = '态神'").QueryRows(&TSDrvBak)
	if err != nil {
		beego.Info("ERROR:ts 001", err)
		return
	}
	//根据态神设备列表信息，补充每个态神设备的数据点和其他信息
	for i := 0; i < len(TSDrvBak); i++ {
		//根据态神设备的列表信息生成索引MAP
		TSDrvMap[TSDrvBak[i].Drv.Name] = i + 1
		_, err = o.Raw("SELECT * FROM maindot where drvname=?", TSDrvBak[i].Drv.Name).QueryRows(&TSDrvBak[i].Dot)
		if err != nil {
			beego.Info("ERROR:ts 002", err)
			return
		}
		//根据数据点类型计算数据点的分类个数
		for j := 0; j < len(TSDrvBak[i].Dot); j++ {
			if TSDrvBak[i].Dot[j].Type == inreg {
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
		beego.Info(TSDrvBak[i].Drv.Name, TSDrvBak[i].TSUrls) //打印生成的URL
	}
}

var TSHotReFrashFalg = false

func StarthttpGet() {
	Getindex := 0
	Gettsdotinfo()
	TSHotReFrashFalg = false
	for {
		if Getindex == len(TSDrv) {
			Getindex = 0
			for count := 0; count < 60; count++ {
				time.Sleep(5e8)
				if TSHotReFrashFalg {
					//开始启动热更新系统参数
					GettsdotinfoHot()
					copy(TSDrv, TSDrvBak)
					Getindex = 0
					TSHotReFrashFalg = false
				}
			}
		}
		if len(TSDrv[Getindex].Dot) > 0 {
			resp, err := http.Get(TSDrv[Getindex].TSUrls)
			if err != nil {
				// handle error
				log.Println("err:", err)

			} else {

				defer resp.Body.Close()

				buf := bytes.NewBuffer(make([]byte, 0, 512))

				_, _ = buf.ReadFrom(resp.Body)

				//beego.Info(len(buf.Bytes()))
				//beego.Info(length)
				//beego.Info(string(buf.Bytes()))
				strresult := string(buf.Bytes())
				strs := strings.Split(strresult, "|")
				TSDrv[Getindex].Flashtime = time.Now().Format("2006-01-02 15:04")
				for n := 0; n < len(strs); n++ {
					insertValue(Getindex, n, strs[n])
					Inserttodotvalue(TSDrv[Getindex].Dot[n])
				}
				Getindex = Getindex + 1
				//beego.Info(strs)
			}
		} else {
			Getindex = Getindex + 1
		}
	}
}
func insertValue(drvindex int, index int, data string) {
	//deindex := 0
	strd := strings.Split(data, ",")
	if len(strd) != 2 {
		return
	}
	v, err := strconv.ParseFloat(strd[1], 32)
	if err != nil {
		beego.Info("Value Cannot Convert to float32", TSDrv[drvindex].Dot[index].Name, data)
	} else {
		TSDrv[drvindex].Dot[index].Value = float32(v)
		if TSDrv[drvindex].Dot[index].Alarmtop != TSDrv[drvindex].Dot[index].Alarmbot {
			if TSDrv[drvindex].Dot[index].Value > TSDrv[drvindex].Dot[index].Alarmtop {
				TSDrv[drvindex].Dot[index].Status = "TOP"
			}
			if TSDrv[drvindex].Dot[index].Value < TSDrv[drvindex].Dot[index].Alarmbot {
				TSDrv[drvindex].Dot[index].Status = "BOT"
			}
		}
	}
}
