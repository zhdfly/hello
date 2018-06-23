package tcpserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ModBusDrv struct {
	Drvname string
	Dotname string
	Dottype string
	Value   string
}
type ModBusUserDrv struct {
	User string
	Drv  []ModBusDrv
}

var Mbdrv []ModBusUserDrv

func GetRealTimeData(name string) (string, int, error) {
	index := 0
	if name == "*" {
		data, err := json.Marshal(Mbdrv)
		fmt.Println(data)
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
func Getdotinfo() {

}
func StarthttpGet() {
	var node ModBusUserDrv
	node.User = "admin"
	var drvnode ModBusDrv
	drvnode.Drvname = "MODBUS"
	drvnode.Dotname = "空温"
	drvnode.Dottype = "4"
	node.Drv = append(node.Drv, drvnode)
	node.Drv = append(node.Drv, drvnode)
	node.Drv = append(node.Drv, drvnode)
	var nodes ModBusUserDrv
	nodes.User = "admins"
	var drvnodes ModBusDrv
	drvnodes.Drvname = "modbus2"
	drvnodes.Dotname = "空湿1"
	drvnodes.Dottype = "4"
	nodes.Drv = append(nodes.Drv, drvnodes)
	Mbdrv = append(Mbdrv, node)
	Mbdrv = append(Mbdrv, nodes)
	mainurl := "http://211.149.159.27:5021/html5/GETTAGVAL/"
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drv); j++ {
			if i == 0 && j == 0 {
				mainurl = mainurl + Mbdrv[i].Drv[j].Drvname + "." + Mbdrv[i].Drv[j].Dotname
			} else {
				mainurl = mainurl + "," + Mbdrv[i].Drv[j].Drvname + "." + Mbdrv[i].Drv[j].Dotname
			}
		}
	}
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
		strresult := string(buf.Bytes())
		strs := strings.Split(strresult, "|")
		for n := 0; n < len(strs); n++ {
			insertValue(n, strs[n])
		}
		//fmt.Println(strs)
		time.Sleep(1e9)
	}
}
func insertValue(index int, data string) {
	deindex := 0
	strd := strings.Split(data, ",")
	if len(strd) != 2 {
		return
	}
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drv); j++ {
			if deindex == index {
				if Mbdrv[i].Drv[j].Dottype == strd[0] {
					Mbdrv[i].Drv[j].Value = strd[1]
				}
				//break
			}
			deindex++
		}
	}
}
