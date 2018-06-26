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

type ModBusDrv struct {
	Drvname string
	Dotname string
	Dottype string
	Value   string
}
type ModBusUserDrv struct {
	User    string
	Drvname []string
	Drv     []ModBusDrv
}

var dotob []Dot
var DrvDot []ModBusDrv
var Mbdrv []ModBusUserDrv

func GetRealTimeData(name string) (string, int, error) {
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
	data, err := json.Marshal(Mbdrv[index])
	//fmt.Println(data)
	return string(data), len(data), err
}
func Getdotinfo() string {
	Mbdrv = nil
	var ob []Usr

	o := orm.NewOrm()
	_, err := o.Raw("SELECT * FROM dot").QueryRows(&dotob)
	if err == nil {
		fmt.Println(dotob)
	}
	for i := 0; i < len(dotob); i++ {
		var tmp ModBusDrv
		tmp.Dotname = dotob[i].Name
		tmp.Drvname = dotob[i].Drv
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
			Mbdrv[j].Drvname = append(Mbdrv[j].Drvname, drvnametmp[i])
		}
	}
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drvname); j++ {
			for k := 0; k < len(dotob); k++ {
				if Mbdrv[i].Drvname[j] == dotob[k].Drv {
					var tmpss ModBusDrv
					tmpss.Drvname = dotob[k].Drv
					tmpss.Dotname = dotob[k].Name
					tmpss.Dottype = dotob[k].Dottype
					Mbdrv[i].Drv = append(Mbdrv[i].Drv, tmpss)
				}
			}
		}
	}
	mainurl := "http://211.149.159.27:5021/html5/GETTAGVAL/"
	// for i := 0; i < len(Mbdrv); i++ {
	// 	for j := 0; j < len(Mbdrv[i].Drv); j++ {
	// 		if i == 0 && j == 0 {
	// 			mainurl = mainurl + Mbdrv[i].Drv[j].Drvname + "." + Mbdrv[i].Drv[j].Dotname
	// 		} else {
	// 			mainurl = mainurl + "," + Mbdrv[i].Drv[j].Drvname + "." + Mbdrv[i].Drv[j].Dotname
	// 		}
	// 	}
	// }
	for i := 0; i < len(dotob); i++ {
		if i == 0 {
			mainurl = mainurl + dotob[i].Drv + "." + dotob[i].Name
		} else {
			mainurl = mainurl + "," + dotob[i].Drv + "." + dotob[i].Name
		}
	}
	return mainurl
}
func StarthttpGet() {

	mainurl := Getdotinfo()
	fmt.Println(mainurl)
	for {
		resp, err := http.Get(mainurl)
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
				for k := 0; k < len(dotob); k++ {
					if Mbdrv[i].Drv[j].Dotname == dotob[k].Name && Mbdrv[i].Drv[j].Drvname == dotob[k].Drv {
						Mbdrv[i].Drv[j].Value = dotob[k].Info
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
