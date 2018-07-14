package modbus

import "github.com/astaxie/beego"
import "hello/tcpserver"

const (
	OK = iota
	ERR
	RETRYING
)
const (
	outio   = 1
	inio    = 2
	holdreg = 3
	inreg   = 4
)

func Initmodbus() {
	beego.Info("开始初始化MODBUS驱动")
	//读取所有的MODBUS设备
	tcpserver.GetMaindrvinfo()
	beego.Info(tcpserver.MBusDrv)

}
