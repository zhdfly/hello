package tcpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ModBusDot struct {
	Dotname       string
	Dottype       string
	Dotwarningtop float32
	Dotwarningbot float32
	Dotstatus     string
	Value         float32
}
type ModBusDrv struct {
	Drvname   string
	Sensornum int
	IOnum     int
	Logicnum  int
	Videonum  int
	Flashtime string
	Dot       []ModBusDot
}
type ModBusUserDrv struct {
	User string
	Drv  []ModBusDrv
}
type DrvBaseInfo struct {
}

var dotob []Dot
var Mbdrv []ModBusUserDrv
var MainUrl []string
var Urlpacknum int

func GetRealTimeData(name interface{}) (string, int, error) {
	index := 0
	if name == "*" {
		data, err := json.Marshal(Mbdrv)
		//fmt.Println(data)
		return string(data), len(data), err
	}

	for index = 0; index < len(Mbdrv); index++ {
		if Mbdrv[index].User == name {
			break
		}
	}
	var data []byte
	var err error
	if index < len(Mbdrv) {
		data, err = json.Marshal(Mbdrv[index])
	}
	//fmt.Println(data)
	return string(data), len(data), err
}
func GetDrvRealTimeData(user string, drv string) (string, int, error) {
	index := 0
	drvindex := 0
	for index = 0; index < len(Mbdrv); index++ {
		if Mbdrv[index].User == user {
			for drvindex = 0; drvindex < len(Mbdrv[index].Drv); drvindex++ {
				if Mbdrv[index].Drv[drvindex].Drvname == drv {
					break
				}
			}
			break
		}
	}
	var data []byte
	var err error
	if index < len(Mbdrv) && drvindex < len(Mbdrv[index].Drv) {
		data, err = json.Marshal(Mbdrv[index].Drv[drvindex])
		fmt.Println(err)
	}
	fmt.Println(index, drvindex, user, drv, data)
	return string(data), len(data), err
}
func Getdotinfo() {
	Mbdrv = nil
	var ob []Usr

	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM dot").QueryRows(&dotob)
	if err == nil {
		fmt.Println(dotob)
	}
	_, err = o.Raw("SELECT name FROM usr").QueryRows(&ob)
	if err == nil {
		fmt.Println(ob)
	}
	for i := 0; i < len(ob); i++ {
		var tmp ModBusUserDrv
		tmp.User = ob[i].Name
		Mbdrv = append(Mbdrv, tmp)
	}
	for j := 0; j < len(Mbdrv); j++ {
		var drvnametmp []string
		_, err = o.Raw("SELECT drvname FROM usrdrv where usrname=?", Mbdrv[j].User).QueryRows(&drvnametmp)
		if err == nil {
			fmt.Println(ob)
		}
		for i := 0; i < len(drvnametmp); i++ {
			var tmp ModBusDrv
			tmp.Drvname = drvnametmp[i]
			Mbdrv[j].Drv = append(Mbdrv[j].Drv, tmp)
		}
	}
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drv); j++ {
			for k := 0; k < len(dotob); k++ {
				if Mbdrv[i].Drv[j].Drvname == dotob[k].Drv {
					var tmpss ModBusDot
					tmpss.Dotname = dotob[k].Name
					tmpss.Dottype = dotob[k].Datatype
					tmpss.Dotwarningtop = dotob[k].Warningtop
					tmpss.Dotwarningbot = dotob[k].Warningbot
					if tmpss.Dottype == "数值类" {
						Mbdrv[i].Drv[j].Sensornum = Mbdrv[i].Drv[j].Sensornum + 1
					} else {
						Mbdrv[i].Drv[j].IOnum = Mbdrv[i].Drv[j].IOnum + 1
					}
					Mbdrv[i].Drv[j].Dot = append(Mbdrv[i].Drv[j].Dot, tmpss)
				}
			}
		}
	}
}
func Creaturl() {
	MainUrl = nil
	var dotnum = len(dotob)
	var Urlpacknum = dotnum/10 + 1
	for p := 0; p < Urlpacknum; p++ {
		url := "http://211.149.159.27:5021/html5/GETTAGVAL/"
		if (dotnum - p*10) >= 10 {
			for i := 0; i < 10; i++ {
				if i == 0 {
					url = url + dotob[p*10+i].Drv + "." + dotob[p*10+i].Name
				} else {
					url = url + "," + dotob[p*10+i].Drv + "." + dotob[p*10+i].Name
				}
			}
		} else {
			for i := 0; i < (dotnum - p*10); i++ {
				if i == 0 {
					url = url + dotob[p*10+i].Drv + "." + dotob[p*10+i].Name
				} else {
					url = url + "," + dotob[p*10+i].Drv + "." + dotob[p*10+i].Name
				}
			}
		}
		MainUrl = append(MainUrl, url)
	}
	fmt.Println(MainUrl)
}
func StarthttpGet() {
	Getindex := 0
	Getdotinfo()
	Creaturl()
	fmt.Println(MainUrl)
	for {
		if Getindex == len(MainUrl) {
			Getindex = 0
			time.Sleep(3e10)
		}
		resp, err := http.Get(MainUrl[Getindex])
		if err != nil {
			// handle error
			log.Println(err)
			return
		}

		defer resp.Body.Close()

		buf := bytes.NewBuffer(make([]byte, 0, 512))

		_, _ = buf.ReadFrom(resp.Body)

		//fmt.Println(len(buf.Bytes()))
		//fmt.Println(length)
		//fmt.Println(string(buf.Bytes()))
		strresult := string(buf.Bytes())
		strs := strings.Split(strresult, "|")
		for n := 0; n < len(strs); n++ {
			insertValue(n+Getindex*10, strs[n])
		}
		for n := Getindex * 10; n < len(dotob); n++ {
			//Inserttodotvalue(dotob[n].Drv, dotob[n].Name, dotob[n].Val, "OK")
			var tmpdot Dotvalue
			tmpdot.Drvname = dotob[n].Drv
			tmpdot.Dotname = dotob[n].Name
			tmpdot.Value = dotob[n].Val
			tmpdot.Status = "OK"
			if dotob[n].Warningtop != dotob[n].Warningbot && dotob[n].Val > dotob[n].Warningtop {
				tmpdot.Status = "TOP"
			}
			if dotob[n].Warningtop != dotob[n].Warningbot && dotob[n].Val < dotob[n].Warningbot {
				tmpdot.Status = "BOT"
			}

			tmpdot.Time = time.Now().Format("2006-01-02 15:04:05")
			Inserttodotvalue(tmpdot)
		}
		for i := 0; i < len(Mbdrv); i++ {
			for j := 0; j < len(Mbdrv[i].Drv); j++ {
				for l := 0; l < len(Mbdrv[i].Drv[j].Dot); l++ {
					for k := 0; k < len(dotob); k++ {
						if Mbdrv[i].Drv[j].Dot[l].Dotname == dotob[k].Name && Mbdrv[i].Drv[j].Drvname == dotob[k].Drv {
							Mbdrv[i].Drv[j].Dot[l].Value = dotob[k].Val
							if Mbdrv[i].Drv[j].Dot[l].Dotwarningtop != Mbdrv[i].Drv[j].Dot[l].Dotwarningbot && Mbdrv[i].Drv[j].Dot[l].Value > Mbdrv[i].Drv[j].Dot[l].Dotwarningtop {
								Mbdrv[i].Drv[j].Dot[l].Dotstatus = "TOP"
							}
							if Mbdrv[i].Drv[j].Dot[l].Dotwarningtop != Mbdrv[i].Drv[j].Dot[l].Dotwarningbot && Mbdrv[i].Drv[j].Dot[l].Value < Mbdrv[i].Drv[j].Dot[l].Dotwarningbot {
								Mbdrv[i].Drv[j].Dot[l].Dotstatus = "BOT"
							}
						}
					}
				}
			}
		}
		Getindex = Getindex + 1
		//fmt.Println(strs)

	}
}
func insertValue(index int, data string) {
	//deindex := 0
	strd := strings.Split(data, ",")
	if len(strd) != 2 {
		return
	}
	v, err := strconv.ParseFloat(strd[1], 32)
	if err != nil {
		dotob[index].Val = 0
	} else {
		dotob[index].Val = float32(v)
	}
	// for i := 0; i < len(dotob); i++ {
	// 	if deindex == index {
	// 		//if Mbdrv[i].Drv[j].Dottype == strd[0] {
	// 		dotob[i].Info = strd[1]
	// 		//}
	// 		//break
	// 		deindex++
	// 	}
	// }
}
