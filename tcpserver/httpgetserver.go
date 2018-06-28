package tcpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type ModBusDot struct {
	Dotname string
	Dottype string
	Value   string
}
type ModBusDrv struct {
	Drvname string
	Dot     []ModBusDot
}
type ModBusUserDrv struct {
	User string
	Drv  []ModBusDrv
}

var dotob []Dot
var Mbdrv []ModBusUserDrv
var MainUrl string

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
					Mbdrv[i].Drv[j].Dot = append(Mbdrv[i].Drv[j].Dot, tmpss)
				}
			}
		}
	}
}
func Creaturl() {
	MainUrl = "http://211.149.159.27:5021/html5/GETTAGVAL/"
	for i := 0; i < len(dotob); i++ {
		if i == 0 {
			MainUrl = MainUrl + dotob[i].Drv + "." + dotob[i].Name
		} else {
			MainUrl = MainUrl + "," + dotob[i].Drv + "." + dotob[i].Name
		}
	}
}
func StarthttpGet() {

	Getdotinfo()
	Creaturl()
	fmt.Println(MainUrl)
	for {
		resp, err := http.Get(MainUrl)
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
			insertValue(n, strs[n])
		}
		for i := 0; i < len(Mbdrv); i++ {
			for j := 0; j < len(Mbdrv[i].Drv); j++ {
				for l := 0; l < len(Mbdrv[i].Drv[j].Dot); l++ {
					for k := 0; k < len(dotob); k++ {
						if Mbdrv[i].Drv[j].Dot[l].Dotname == dotob[k].Name && Mbdrv[i].Drv[j].Drvname == dotob[k].Drv {
							Mbdrv[i].Drv[j].Dot[l].Value = dotob[k].Info
						}
					}
				}
			}
		}
		//fmt.Println(strs)
		time.Sleep(5e9)
	}
}
func insertValue(index int, data string) {
	//deindex := 0
	strd := strings.Split(data, ",")
	if len(strd) != 2 {
		return
	}
	dotob[index].Info = strd[1]
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
