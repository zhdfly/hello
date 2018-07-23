package tcpserver

import (
	"bytes"
	"encoding/binary"
	"flag"
	"math"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	OK = iota
	ERR
	RETRYING
)

var MBDrvMap = make(map[string]int)
var MBDrvMapBak = make(map[string]int)
var MBDrv []MainDrvType
var MBDrvBak []MainDrvType

//计算MODBUS的CRC16校验
func crc16hex(str []byte, len int) [2]byte {
	var crc16 = [2]byte{0, 0}
	var lo, hi byte
	hi = 0xff
	lo = 0xff
	for i := 0; i < len; i++ {
		idx := lo ^ str[i]
		lo = (hi ^ table_h[idx])
		hi = table_l[idx]
	}
	crc16[0] = hi
	crc16[1] = lo

	return crc16
}

func GetUint16(b []byte) uint16 {
	bin_buf := bytes.NewBuffer(b)
	var x uint16
	binary.Read(bin_buf, binary.BigEndian, &x)
	return x
}
func GetUint32(b []byte) uint32 {
	bin_buf := bytes.NewBuffer(b)
	var x uint32
	binary.Read(bin_buf, binary.BigEndian, &x)
	return x
}
func GetFloat32(b []byte) float32 {
	bin_buf := bytes.NewBuffer(b)
	var x float32
	binary.Read(bin_buf, binary.LittleEndian, &x)
	return x
}
func GetFloat32ToByte(float float32) []byte {
	bits := math.Float32bits(float)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
func Getmbdotinfo() {
	MBDrv = nil
	o := orm.NewOrm()
	//获取MODBUS设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = 'MODBUS'").QueryRows(&MBDrv)
	if err != nil {
		beego.Info("ERROR:mb 001", err)
		return
	}
	//根据MODBUS设备列表信息，补充每个MODBUS设备的数据点和其他信息
	for i := 0; i < len(MBDrv); i++ {
		//根据MODBUS设备的列表信息生成索引MAP
		MBDrvMap[MBDrv[i].Drv.Name] = i + 1
		_, err = o.Raw("SELECT * FROM maindot where drvname=? order by addr asc", MBDrv[i].Drv.Name).QueryRows(&MBDrv[i].Dot)
		if err != nil {
			beego.Info("ERROR:mb 002", err)
			return
		}
		//根据数据点类型计算数据点的分类个数
		var outiolist, iniolist, outreglist, inreglist ModbuscmdType           // 临时命令参数
		var outiolistcount, iniolistcount, outreglistcount, inreglistcount int //记录当前临时命令参数中记录的地址个数
		for j := 0; j < len(MBDrv[i].Dot); j++ {
			if MBDrv[i].Dot[j].Type == inreg {
				MBDrv[i].Sensornum = MBDrv[i].Sensornum + 1
				//根据设备的点数据生成MODBUS读取指令，主要是根据单点地址生成连续读取的命令
				if inreglistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					inreglist.Cmd = inreg
					inreglist.startAddr = MBDrv[i].Dot[j].Addr
					inreglist.AddrLen = GetDataTypeLen(MBDrv[i].Dot[j].Data)
					inreglistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrv[i].Dot[j].Addr-inreglist.startAddr-inreglist.AddrLen < 20 {
						//可以把这个点和临时参数合成连续读取
						inreglist.AddrLen = (MBDrv[i].Dot[j].Addr - inreglist.startAddr) + GetDataTypeLen(MBDrv[i].Dot[j].Data)
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, inreglist) //暂缓临时命令
						//重建临时命令
						inreglist.Cmd = inreg
						inreglist.startAddr = MBDrv[i].Dot[j].Addr
						inreglist.AddrLen = GetDataTypeLen(MBDrv[i].Dot[j].Data)
					}
				}
			} else if MBDrv[i].Dot[j].Type == outreg {
				MBDrv[i].Sensornum = MBDrv[i].Sensornum + 1
				if outreglistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					outreglist.Cmd = outreg
					outreglist.startAddr = MBDrv[i].Dot[j].Addr
					outreglist.AddrLen = GetDataTypeLen(MBDrv[i].Dot[j].Data)
					outreglistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrv[i].Dot[j].Addr-outreglist.startAddr-outreglist.AddrLen < 20 {
						//可以把这个点和临时参数合成连续读取
						outreglist.AddrLen = (MBDrv[i].Dot[j].Addr - outreglist.startAddr) + GetDataTypeLen(MBDrv[i].Dot[j].Data)
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, outreglist) //暂缓临时命令
						//重建临时命令
						outreglist.Cmd = outreg
						outreglist.startAddr = MBDrv[i].Dot[j].Addr
						outreglist.AddrLen = GetDataTypeLen(MBDrv[i].Dot[j].Data)
					}
				}
			} else if MBDrv[i].Dot[j].Type == outio {
				MBDrv[i].IOnum = MBDrv[i].IOnum + 1
				if outiolistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					outiolist.Cmd = outio
					outiolist.startAddr = MBDrv[i].Dot[j].Addr
					outiolist.AddrLen = 1
					outiolistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrv[i].Dot[j].Addr-outiolist.startAddr-outiolist.AddrLen < 5 {
						//可以把这个点和临时参数合成连续读取
						outiolist.AddrLen = (MBDrv[i].Dot[j].Addr - outiolist.startAddr) + 1
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, outiolist) //暂缓临时命令
						//重建临时命令
						outiolist.Cmd = outio
						outiolist.startAddr = MBDrv[i].Dot[j].Addr
						outiolist.AddrLen = 1
					}
				}
			} else if MBDrv[i].Dot[j].Type == inio {
				MBDrv[i].IOnum = MBDrv[i].IOnum + 1
				if iniolistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					iniolist.Cmd = inio
					iniolist.startAddr = MBDrv[i].Dot[j].Addr
					iniolist.AddrLen = 1
					iniolistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrv[i].Dot[j].Addr-iniolist.startAddr-iniolist.AddrLen < 5 {
						//可以把这个点和临时参数合成连续读取
						iniolist.AddrLen = (MBDrv[i].Dot[j].Addr - iniolist.startAddr) + 1
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, iniolist) //暂缓临时命令
						//重建临时命令
						iniolist.Cmd = inio
						iniolist.startAddr = MBDrv[i].Dot[j].Addr
						iniolist.AddrLen = 1
					}
				}
			}
		}
		//所有点循环完毕之后，最后一次建立的临时命令还未暂缓到设备信息中，需要暂缓
		if inreglist.AddrLen > 0 {
			MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, inreglist)
		}
		if outreglist.AddrLen > 0 {
			MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, outreglist)
		}
		if iniolist.AddrLen > 0 {
			MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, iniolist)
		}
		if outiolist.AddrLen > 0 {
			MBDrv[i].MBCmds = append(MBDrv[i].MBCmds, outiolist)
		}
		for j := 0; j < len(MBDrv[i].MBCmds); j++ {
			MBDrv[i].MBCmds[j].Buffer[0] = uint8(MBDrv[i].Drv.Addr)
			MBDrv[i].MBCmds[j].Buffer[1] = uint8(MBDrv[i].MBCmds[j].Cmd)
			MBDrv[i].MBCmds[j].Buffer[2] = uint8(MBDrv[i].MBCmds[j].startAddr >> 8)
			MBDrv[i].MBCmds[j].Buffer[3] = uint8(MBDrv[i].MBCmds[j].startAddr)
			MBDrv[i].MBCmds[j].Buffer[4] = uint8(MBDrv[i].MBCmds[j].AddrLen >> 8)
			MBDrv[i].MBCmds[j].Buffer[5] = uint8(MBDrv[i].MBCmds[j].AddrLen)
			crc := crc16hex(MBDrv[i].MBCmds[j].Buffer[:6], 6)
			MBDrv[i].MBCmds[j].Buffer[6] = crc[1]
			MBDrv[i].MBCmds[j].Buffer[7] = crc[0]
			MBDrv[i].MBCmds[j].BufferLen = 8
		}
		beego.Info(MBDrv[i].MBCmds)
		MBDrv[i].Videonum = Getdrvvedionum(MBDrv[i].Drv.Name) //获取摄像头的个数
	}
}
func GetmbdotinfoHot() {
	MBDrvBak = nil
	o := orm.NewOrm()
	//获取态神设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = 'MODBUS'").QueryRows(&MBDrvBak)
	if err != nil {
		beego.Info("ERROR:mb 001", err)
		return
	}
	//根据态神设备列表信息，补充每个态神设备的数据点和其他信息
	for i := 0; i < len(MBDrvBak); i++ {
		//根据态神设备的列表信息生成索引MAP
		MBDrvMapBak[MBDrvBak[i].Drv.Name] = i + 1
		_, err = o.Raw("SELECT * FROM maindot where drvname=? order by addr asc", MBDrvBak[i].Drv.Name).QueryRows(&MBDrvBak[i].Dot)
		if err != nil {
			beego.Info("ERROR:mb 002", err)
			return
		}
		//根据数据点类型计算数据点的分类个数
		var outiolist, iniolist, outreglist, inreglist ModbuscmdType           // 临时命令参数
		var outiolistcount, iniolistcount, outreglistcount, inreglistcount int //记录当前临时命令参数中记录的地址个数
		for j := 0; j < len(MBDrvBak[i].Dot); j++ {
			if MBDrvBak[i].Dot[j].Type == inreg {
				MBDrvBak[i].Sensornum = MBDrvBak[i].Sensornum + 1
				//根据设备的点数据生成MODBUS读取指令，主要是根据单点地址生成连续读取的命令
				if inreglistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					inreglist.Cmd = inreg
					inreglist.startAddr = MBDrvBak[i].Dot[j].Addr
					inreglist.AddrLen = GetDataTypeLen(MBDrvBak[i].Dot[j].Data)
					inreglistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrvBak[i].Dot[j].Addr-inreglist.startAddr-inreglist.AddrLen < 20 {
						//可以把这个点和临时参数合成连续读取
						inreglist.AddrLen = (MBDrvBak[i].Dot[j].Addr - inreglist.startAddr) + GetDataTypeLen(MBDrvBak[i].Dot[j].Data)
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, inreglist) //暂缓临时命令
						//重建临时命令
						inreglist.Cmd = inreg
						inreglist.startAddr = MBDrvBak[i].Dot[j].Addr
						inreglist.AddrLen = GetDataTypeLen(MBDrvBak[i].Dot[j].Data)
					}
				}
			} else if MBDrvBak[i].Dot[j].Type == outreg {
				MBDrvBak[i].Sensornum = MBDrvBak[i].Sensornum + 1
				if outreglistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					outreglist.Cmd = outreg
					outreglist.startAddr = MBDrvBak[i].Dot[j].Addr
					outreglist.AddrLen = GetDataTypeLen(MBDrvBak[i].Dot[j].Data)
					outreglistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrvBak[i].Dot[j].Addr-outreglist.startAddr-outreglist.AddrLen < 20 {
						//可以把这个点和临时参数合成连续读取
						outreglist.AddrLen = (MBDrvBak[i].Dot[j].Addr - outreglist.startAddr) + GetDataTypeLen(MBDrvBak[i].Dot[j].Data)
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, outreglist) //暂缓临时命令
						//重建临时命令
						outreglist.Cmd = outreg
						outreglist.startAddr = MBDrvBak[i].Dot[j].Addr
						outreglist.AddrLen = GetDataTypeLen(MBDrvBak[i].Dot[j].Data)
					}
				}
			} else if MBDrvBak[i].Dot[j].Type == outio {
				MBDrvBak[i].IOnum = MBDrvBak[i].IOnum + 1
				if outiolistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					outiolist.Cmd = outio
					outiolist.startAddr = MBDrvBak[i].Dot[j].Addr
					outiolist.AddrLen = 1
					outiolistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrvBak[i].Dot[j].Addr-outiolist.startAddr-outiolist.AddrLen < 5 {
						//可以把这个点和临时参数合成连续读取
						outiolist.AddrLen = (MBDrvBak[i].Dot[j].Addr - outiolist.startAddr) + 1
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, outiolist) //暂缓临时命令
						//重建临时命令
						outiolist.Cmd = outio
						outiolist.startAddr = MBDrvBak[i].Dot[j].Addr
						outiolist.AddrLen = 1
					}
				}
			} else if MBDrvBak[i].Dot[j].Type == inio {
				MBDrvBak[i].IOnum = MBDrvBak[i].IOnum + 1
				if iniolistcount == 0 {
					//临时参数中没有地址信息，需要第一次赋值
					iniolist.Cmd = inio
					iniolist.startAddr = MBDrvBak[i].Dot[j].Addr
					iniolist.AddrLen = 1
					iniolistcount = 1
				} else {
					//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
					if MBDrvBak[i].Dot[j].Addr-iniolist.startAddr-iniolist.AddrLen < 5 {
						//可以把这个点和临时参数合成连续读取
						iniolist.AddrLen = (MBDrvBak[i].Dot[j].Addr - iniolist.startAddr) + 1
					} else {
						//当前点的地址不足以与临时点组成连续读取
						MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, iniolist) //暂缓临时命令
						//重建临时命令
						iniolist.Cmd = inio
						iniolist.startAddr = MBDrvBak[i].Dot[j].Addr
						iniolist.AddrLen = 1
					}
				}
			}
		}
		//所有点循环完毕之后，最后一次建立的临时命令还未暂缓到设备信息中，需要暂缓
		if inreglist.AddrLen > 0 {
			MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, inreglist)
		}
		if outreglist.AddrLen > 0 {
			MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, outreglist)
		}
		if iniolist.AddrLen > 0 {
			MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, iniolist)
		}
		if outiolist.AddrLen > 0 {
			MBDrvBak[i].MBCmds = append(MBDrvBak[i].MBCmds, outiolist)
		}
		for j := 0; j < len(MBDrvBak[i].MBCmds); j++ {
			//生成MODBUS命令包
			MBDrvBak[i].MBCmds[j].Buffer[0] = uint8(MBDrvBak[i].Drv.Addr)
			MBDrvBak[i].MBCmds[j].Buffer[1] = uint8(MBDrvBak[i].MBCmds[j].Cmd)
			MBDrvBak[i].MBCmds[j].Buffer[2] = uint8(MBDrvBak[i].MBCmds[j].startAddr >> 8)
			MBDrvBak[i].MBCmds[j].Buffer[3] = uint8(MBDrvBak[i].MBCmds[j].startAddr)
			MBDrvBak[i].MBCmds[j].Buffer[4] = uint8(MBDrvBak[i].MBCmds[j].AddrLen >> 8)
			MBDrvBak[i].MBCmds[j].Buffer[5] = uint8(MBDrvBak[i].MBCmds[j].AddrLen)
			crc := crc16hex(MBDrvBak[i].MBCmds[j].Buffer[:6], 6)
			MBDrvBak[i].MBCmds[j].Buffer[6] = crc[1]
			MBDrvBak[i].MBCmds[j].Buffer[7] = crc[0]
			MBDrvBak[i].MBCmds[j].BufferLen = 8
		}
		beego.Info(MBDrvBak[i].MBCmds)
		MBDrvBak[i].Videonum = Getdrvvedionum(MBDrvBak[i].Drv.Name) //获取摄像头的个数
	}
}

var connlist = make(map[string]net.Conn)
var alreadyopenport = make(map[string][]net.Conn)

func MBStartTCP(port string, index int) {
	flag.Parse()
	var l net.Listener
	var err error
	// 开启TCP监听
	l, err = net.Listen("tcp", ":"+port)
	if err != nil {
		//开启失败，一般由于端口被占用
		beego.Info("Error listening:", err)
		if len(alreadyopenport[port]) == 0 {
			beego.Info("1:")
			//return
		} else {
			beego.Info("2:")
			//return
		}
	} else {
		//开启成功
		defer l.Close()
		beego.Info("Listening on " + ":" + port)
		for {
			// 等待客户端连接
			conn, err := l.Accept()
			if err != nil {
				beego.Info("Error accepting: ", err)
				os.Exit(1)
			}
			//接收到新的客户端连接，保存连接对象到此端口的切片中
			alreadyopenport[port] = append(alreadyopenport[port], conn)
			//开启客户端接收数据线程
			go mbhandleRequest(conn, index)
		}
	}
}

var mbdrvportindex = make(map[int][]int) //以通讯端口作为标志的设备索引列表
var MBHotReFrashFalg = false             //热更新标志

func StartModbus() {
	//开启MODBUS驱动初始化
	Getmbdotinfo() //获取MODBUS设备基本信息
	for i := 0; i < len(MBDrv); i++ {
		//建立以通信端口为标志的设备索引列表
		mbdrvportindex[MBDrv[i].Drv.Port] = append(mbdrvportindex[MBDrv[i].Drv.Port], i)
		//依次开启设备的通信端口
		go MBStartTCP(strconv.Itoa(MBDrv[i].Drv.Port), MBDrv[i].Drv.Port)
	}
	for k, _ := range mbdrvportindex {
		//根据通信端口开启同一端口的设备轮询线程
		go mbreqthread(k)
	}

}
func ThreadMB() {
	//开启MODBUS设备主线程
	for {
		//监听热更新标志
		if MBHotReFrashFalg {
			//获取新的设备列表到备用区域
			GetmbdotinfoHot()
			//拷贝设备列表信息
			copy(MBDrv, MBDrvBak)
			//以下同MODBUS初始化
			for i := 0; i < len(MBDrv); i++ {
				mbdrvportindex[MBDrv[i].Drv.Port] = append(mbdrvportindex[MBDrv[i].Drv.Port], i)
				go MBStartTCP(strconv.Itoa(MBDrv[i].Drv.Port), i)
			}
			for k, _ := range mbdrvportindex {
				go mbreqthread(k)
			}
			MBHotReFrashFalg = false
		}
		//简单的延时
		time.Sleep(1e8)
	}
}

//设置MODBUS的IO输出
//参数：被操作的设备，被操作的数据点，目标值
func SetMBioValue(drv string, dot string, value int) {
	//从设备MAP中找到此设备的索引
	index := MBDrvMap[drv]
	//判断设备名称是否合法 只有再MAP找不到对应的值的时候才会为0，建立MAP时默认最小为1
	if index > 0 {
		//遍历设备中的数据点
		for _, c := range MBDrv[index-1].Dot {
			if c.Name == dot {
				//找到数据点
				//建立命令数据包
				MBDrv[index-1].MBStatus = mbready
				var tmp ModbuscmdType
				tmp.startAddr = c.Addr
				tmp.AddrLen = 1
				tmp.Buffer[0] = uint8(MBDrv[index-1].Drv.Addr)
				tmp.Buffer[1] = 0x05
				tmp.Buffer[2] = uint8(c.Addr >> 8)
				tmp.Buffer[3] = uint8(c.Addr)
				tmp.Buffer[4] = uint8(value >> 8)
				tmp.Buffer[5] = uint8(value)
				crc := crc16hex(tmp.Buffer[:6], 6)
				tmp.Buffer[6] = crc[1]
				tmp.Buffer[7] = crc[0]
				//数据包需要值定数据包类型为Ctr（被控数据包类型）
				tmp.Type = Ctr
				tmp.BufferLen = 8
				//插入数据包
				MBDrv[index-1].MBCmds = append(MBDrv[index-1].MBCmds, tmp)
			}
		}
	}
}

//设置MODBUS寄存器
//参数：被操作的设备，被操作的数据点，目标值
func SetMBregValue(drv string, dot string, value float32) {
	//从设备MAP中找到此设备的索引
	index := MBDrvMap[drv]
	len := 0
	//先格式化目标值为INT型，float型无法执行>>（位移）操作
	intvalue := int(value)
	//判断设备名称是否合法 只有再MAP找不到对应的值的时候才会为0，建立MAP时默认最小为1
	if index > 0 {
		for _, c := range MBDrv[index-1].Dot {
			if c.Name == dot {
				//找到目标数据点
				MBDrv[index-1].MBStatus = mbready
				var tmp ModbuscmdType
				tmp.startAddr = c.Addr
				tmp.AddrLen = 1
				tmp.Buffer[len] = uint8(MBDrv[index-1].Drv.Addr)
				len++
				//这里需要判断数据点的数据格式，dword和float均占用四个字节，所以发送时的功能码是0x10，byte,word是两个字节，功能码是0x06
				if c.Data == databyte || c.Data == dataword {
					tmp.Buffer[len] = 0x06
					len++
					tmp.Buffer[len] = uint8(c.Addr >> 8)
					len++
					tmp.Buffer[len] = uint8(c.Addr)
					len++
					tmp.Buffer[len] = uint8(intvalue >> 8)
					len++
					tmp.Buffer[len] = uint8(intvalue)
					len++
				} else if c.Data == datadword {
					tmp.Buffer[len] = 0x10
					len++
					tmp.Buffer[len] = uint8(c.Addr >> 8)
					len++
					tmp.Buffer[len] = uint8(c.Addr)
					len++
					tmp.Buffer[len] = 0x00
					len++
					tmp.Buffer[len] = 0x02
					len++
					tmp.Buffer[len] = 0x04
					len++
					tmp.Buffer[len] = uint8(intvalue >> 24)
					len++
					tmp.Buffer[len] = uint8(intvalue >> 16)
					len++
					tmp.Buffer[len] = uint8(intvalue >> 8)
					len++
					tmp.Buffer[len] = uint8(intvalue)
					len++
				} else if c.Data == datafloat {
					//float值不需要用到开头格式化为INT的数值，直接使用就可以，但需要变为byte[]数组
					fb := GetFloat32ToByte(value)
					tmp.Buffer[len] = 0x10
					len++
					tmp.Buffer[len] = uint8(c.Addr >> 8)
					len++
					tmp.Buffer[len] = uint8(c.Addr)
					len++
					tmp.Buffer[len] = 0x00
					len++
					tmp.Buffer[len] = 0x02
					len++
					tmp.Buffer[len] = 0x04
					len++
					tmp.Buffer[len] = fb[0]
					len++
					tmp.Buffer[len] = fb[1]
					len++
					tmp.Buffer[len] = fb[3]
					len++
					tmp.Buffer[len] = fb[4]
					len++
				}
				crc := crc16hex(tmp.Buffer[:len], len)
				tmp.Buffer[len] = crc[1]
				len++
				tmp.Buffer[len] = crc[0]
				len++
				tmp.Type = Ctr
				tmp.BufferLen = len
				MBDrv[index-1].MBCmds = append(MBDrv[index-1].MBCmds, tmp)
			}
		}
	}
}

//通信端口的轮询线程
func mbreqthread(port int) {
	timecount := 0
	for {
		for _, i := range mbdrvportindex[port] {
			//遍历端口下的所有设备的索引
			//检测设备的状态
			if MBDrv[i].MBStatus == mbready {
				//设备处于空闲状态
				for index := 0; index < len(alreadyopenport[strconv.Itoa(port)]); index++ {
					//遍历设备通讯端口下的所有的有效连接
					//所有的有效连接都需要发送一下，因为不知道设备具体存在那个连接中
					con := alreadyopenport[strconv.Itoa(port)][index]
					//获取其中一个有效连接
					//判断设备中是否有命令包列表
					if len(MBDrv[i].MBCmds) > 0 {
						for j := 0; j < len(MBDrv[i].MBCmds); j++ {
							if MBDrv[i].MBCmds[j].Type == Ctr {
								//判断是否存在控制指令 Ctr
								//存在的话，优先发送
								MBDrv[i].MBCmdNum = j
								break
							}
						}
						//发送数据包
						_, err := con.Write(MBDrv[i].MBCmds[MBDrv[i].MBCmdNum].Buffer[:MBDrv[i].MBCmds[MBDrv[i].MBCmdNum].BufferLen])
						if err != nil {
							beego.Info(err)
							//发送失败，一般是连接失效，需要删除这个连接
							alreadyopenport[strconv.Itoa(port)] = append(alreadyopenport[strconv.Itoa(port)][:index], alreadyopenport[strconv.Itoa(port)][index+1:]...)
						}
						//改变设备状态
						MBDrv[i].MBStatus = mbbusy
					}
				}
			} else if MBDrv[i].MBStatus == mbbusy {
				//判断设备在忙碌状态下的时间，是否超过了采样时间，超过了就认为设备通讯超时
				if timecount >= MBDrv[i].Drv.Samplingtime && timecount%MBDrv[i].Drv.Samplingtime == 0 {
					MBDrv[i].MBStatus = mbtimeout
				}
			} else if MBDrv[i].MBStatus == mbtimeout {
				//设备超时之后需要等待10秒再次发送尝试
				if timecount >= MBDrv[i].Drv.Samplingtime && timecount%MBDrv[i].Drv.Samplingtime == 10 {
					MBDrv[i].MBStatus = mbready
				}
			} else if MBDrv[i].MBStatus == mbwite {
				//只有所有的命令包发送一边之后才会进入等待状态，等待经过一个采样时间之后，也就是采样周期之后，再次开始
				if timecount%MBDrv[i].Drv.Samplingtime == 0 {
					MBDrv[i].MBStatus = mbready
				}
			}
			time.Sleep(1e8)
		}
		timecount++
	}
	beego.Info("Done!")
}
func mbhandleRequest(conn net.Conn, port int) {
	var MBRecv [10240]byte
	ipStr := conn.RemoteAddr().String()
	defer func() {
		beego.Info("Disconnected :" + ipStr)
		conn.Close()
	}()
	for {
		//等待数据接收
		n, err := conn.Read(MBRecv[0:])
		if err != nil {
			return
		}
		if n > 0 {
			MBRecv[n] = 0
			//注意：下面操作并未考虑TCP数据的粘包问题，因为每个设备有且只有一个连接且此线程为单独设备独占，再加上MODBUS为主从模式，基本不会出现粘包问题
			//接收到数据包，判断是否是合法的数据，判断长度或者特殊字符
			if MBRecv[2] == byte(n-5) || MBRecv[1] == setio || MBRecv[1] == setreg || MBRecv[1] == setregs {
				//获取CRC校验值
				crc := crc16hex(MBRecv[0:n-2], n-2)
				if crc[0] == MBRecv[n-1] && crc[1] == MBRecv[n-2] {
					//CRC校验通过
					for _, index := range mbdrvportindex[port] {
						//遍历通信端口下的所有设备
						if MBRecv[0] == byte(MBDrv[index].Drv.Addr) {
							//判断数据包的地址是否和设备地址一致
							//提取当前发送的请求包中的起始地址
							packstartaddr := MBDrv[index].MBCmds[MBDrv[index].MBCmdNum].startAddr
							//提取当前发送的请求包中的结束地址
							packstopaddr := packstartaddr + MBDrv[index].MBCmds[MBDrv[index].MBCmdNum].AddrLen
							MBDrv[index].MBCmdNum++
							for i := 0; i < len(MBDrv[index].Dot); i++ {
								//遍历设备的所有数据点
								//判断数据点的地址是否在起始地址和结束地址之间
								if MBDrv[index].Dot[i].Addr >= packstartaddr && MBDrv[index].Dot[i].Addr < packstopaddr {
									//判断数据点的类型是否和返回数据包中的匹配
									if (MBRecv[1] == inreg || MBRecv[1] == outreg) && MBDrv[index].Dot[i].Type == int(MBRecv[1]) {
										//寄存器数据
										//寄存器数据，每个寄存器占用两个字节
										packindex := (MBDrv[index].Dot[i].Addr - packstartaddr) * 2
										//根据数据点数据格式提取数据
										if MBDrv[index].Dot[i].Data == databyte {
											MBDrv[index].Dot[i].Value = float32(MBRecv[packindex+3]) * MBDrv[index].Dot[i].Bynum
											Inserttodotvalue(MBDrv[index].Dot[i])
										} else if MBDrv[index].Dot[i].Data == dataword {
											MBDrv[index].Dot[i].Value = float32(GetUint16(MBRecv[packindex+3:packindex+5])) * MBDrv[index].Dot[i].Bynum
											Inserttodotvalue(MBDrv[index].Dot[i])
										} else if MBDrv[index].Dot[i].Data == datadword {
											MBDrv[index].Dot[i].Value = float32(GetUint32(MBRecv[packindex+3:packindex+7])) * MBDrv[index].Dot[i].Bynum
											Inserttodotvalue(MBDrv[index].Dot[i])
										} else if MBDrv[index].Dot[i].Data == datafloat {
											MBDrv[index].Dot[i].Value = GetFloat32(MBRecv[packindex+3:packindex+7]) * MBDrv[index].Dot[i].Bynum
											Inserttodotvalue(MBDrv[index].Dot[i])
										}
									} else if (MBRecv[1] == inio || MBRecv[1] == outio) && MBDrv[index].Dot[i].Type == int(MBRecv[1]) {
										//IO点数据
										//IO点的数据返回是按位来返回的，每8个点占用一个字节，所以需要通过移位来提取数据
										packindex := MBDrv[index].Dot[i].Addr - packstartaddr
										packnum := packindex / 8
										valueindex := uint8(packindex % 8)
										if (MBRecv[3+packnum]>>valueindex)&0x01 == 1 {
											MBDrv[index].Dot[i].Value = 1
											Inserttodotvalue(MBDrv[index].Dot[i])
										} else {
											MBDrv[index].Dot[i].Value = 0
											Inserttodotvalue(MBDrv[index].Dot[i])
										}
									} else if MBRecv[1] == setio {
										//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
										for ci, c := range MBDrv[index].MBCmds {
											if c.Type == Ctr {
												MBDrv[index].MBCmdNum = 0
												MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
											}
										}
									} else if MBRecv[1] == setreg {
										//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
										for ci, c := range MBDrv[index].MBCmds {
											if c.Type == Ctr {
												MBDrv[index].MBCmdNum = 0
												MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
											}
										}
									} else if MBRecv[1] == setregs {
										//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
										for ci, c := range MBDrv[index].MBCmds {
											if c.Type == Ctr {
												MBDrv[index].MBCmdNum = 0
												MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
											}
										}
									}
								}
							}
							MBDrv[index].Flashtime = time.Now().Format("2006-01-02 15:04")
							if MBDrv[index].MBCmdNum >= len(MBDrv[index].MBCmds) {
								MBDrv[index].MBCmdNum = 0
								MBDrv[index].MBStatus = mbwite
							} else {
								MBDrv[index].MBStatus = mbready
							}
						}
					}
				}
			}
		}
	}
	beego.Info("Done!")
}
