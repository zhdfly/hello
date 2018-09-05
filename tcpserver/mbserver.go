package tcpserver

import (
	"bytes"
	"encoding/binary"
	"flag"
	"math"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/robfig/cron"
)

const (
	OK = iota
	ERR
	RETRYING
)

var MBDrvMap MainMAPstringint
var MBDrvMapBak MainMAPstringint
var MBDrv []MainDrvType
var MBDrvBak []MainDrvType
var MBCtrlCmdRltChanMap = make(map[string](chan string))

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
func CreateMbCmd(drv *MainDrvType) {
	drv.Sensornum = 0
	drv.Videonum = 0
	drv.IOnum = 0
	drv.Logicnum = 0
	drv.MBCmds = nil
	drv.MBCmdNum = 0
	//根据数据点类型计算数据点的分类个数
	var outiolist, iniolist, outreglist, inreglist ModbuscmdType           // 临时命令参数
	var outiolistcount, iniolistcount, outreglistcount, inreglistcount int //记录当前临时命令参数中记录的地址个数
	for j := 0; j < len(drv.Dot); j++ {
		if drv.Dot[j].Type == inreg {
			drv.Sensornum = drv.Sensornum + 1
			//根据设备的点数据生成MODBUS读取指令，主要是根据单点地址生成连续读取的命令
			if inreglistcount == 0 {
				//临时参数中没有地址信息，需要第一次赋值
				inreglist.Cmd = inreg
				inreglist.startAddr = drv.Dot[j].Addr
				inreglist.AddrLen = GetDataTypeLen(drv.Dot[j].Data)
				inreglistcount = 1
			} else {
				//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
				if drv.Dot[j].Addr-inreglist.startAddr-inreglist.AddrLen < 20 {
					//可以把这个点和临时参数合成连续读取
					inreglist.AddrLen = (drv.Dot[j].Addr - inreglist.startAddr) + GetDataTypeLen(drv.Dot[j].Data)
				} else {
					//当前点的地址不足以与临时点组成连续读取
					drv.MBCmds = append(drv.MBCmds, inreglist) //暂缓临时命令
					//重建临时命令
					inreglist.Cmd = inreg
					inreglist.startAddr = drv.Dot[j].Addr
					inreglist.AddrLen = GetDataTypeLen(drv.Dot[j].Data)
				}
			}
		} else if drv.Dot[j].Type == outreg {
			drv.Sensornum = drv.Sensornum + 1
			if outreglistcount == 0 {
				//临时参数中没有地址信息，需要第一次赋值
				outreglist.Cmd = outreg
				outreglist.startAddr = drv.Dot[j].Addr
				outreglist.AddrLen = GetDataTypeLen(drv.Dot[j].Data)
				outreglistcount = 1
			} else {
				//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
				if drv.Dot[j].Addr-outreglist.startAddr-outreglist.AddrLen < 20 {
					//可以把这个点和临时参数合成连续读取
					outreglist.AddrLen = (drv.Dot[j].Addr - outreglist.startAddr) + GetDataTypeLen(drv.Dot[j].Data)
				} else {
					//当前点的地址不足以与临时点组成连续读取
					drv.MBCmds = append(drv.MBCmds, outreglist) //暂缓临时命令
					//重建临时命令
					outreglist.Cmd = outreg
					outreglist.startAddr = drv.Dot[j].Addr
					outreglist.AddrLen = GetDataTypeLen(drv.Dot[j].Data)
				}
			}
		} else if drv.Dot[j].Type == outio {
			drv.IOnum = drv.IOnum + 1
			if outiolistcount == 0 {
				//临时参数中没有地址信息，需要第一次赋值
				outiolist.Cmd = outio
				outiolist.startAddr = drv.Dot[j].Addr
				outiolist.AddrLen = 1
				outiolistcount = 1
			} else {
				//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
				if drv.Dot[j].Addr-outiolist.startAddr-outiolist.AddrLen < 5 {
					//可以把这个点和临时参数合成连续读取
					outiolist.AddrLen = (drv.Dot[j].Addr - outiolist.startAddr) + 1
				} else {
					//当前点的地址不足以与临时点组成连续读取
					drv.MBCmds = append(drv.MBCmds, outiolist) //暂缓临时命令
					//重建临时命令
					outiolist.Cmd = outio
					outiolist.startAddr = drv.Dot[j].Addr
					outiolist.AddrLen = 1
				}
			}
		} else if drv.Dot[j].Type == inio {
			drv.IOnum = drv.IOnum + 1
			if iniolistcount == 0 {
				//临时参数中没有地址信息，需要第一次赋值
				iniolist.Cmd = inio
				iniolist.startAddr = drv.Dot[j].Addr
				iniolist.AddrLen = 1
				iniolistcount = 1
			} else {
				//临时参数中已存在地址信息，需要检测当前地址与参数中的地址的偏差，偏差小于50则可以设置为连续读写，否则不可以连续读取
				if drv.Dot[j].Addr-iniolist.startAddr-iniolist.AddrLen < 5 {
					//可以把这个点和临时参数合成连续读取
					iniolist.AddrLen = (drv.Dot[j].Addr - iniolist.startAddr) + 1
				} else {
					//当前点的地址不足以与临时点组成连续读取
					drv.MBCmds = append(drv.MBCmds, iniolist) //暂缓临时命令
					//重建临时命令
					iniolist.Cmd = inio
					iniolist.startAddr = drv.Dot[j].Addr
					iniolist.AddrLen = 1
				}
			}
		}
	}
	//所有点循环完毕之后，最后一次建立的临时命令还未暂缓到设备信息中，需要暂缓
	if inreglist.AddrLen > 0 {
		drv.MBCmds = append(drv.MBCmds, inreglist)
	}
	if outreglist.AddrLen > 0 {
		drv.MBCmds = append(drv.MBCmds, outreglist)
	}
	if iniolist.AddrLen > 0 {
		drv.MBCmds = append(drv.MBCmds, iniolist)
	}
	if outiolist.AddrLen > 0 {
		drv.MBCmds = append(drv.MBCmds, outiolist)
	}
	for j := 0; j < len(drv.MBCmds); j++ {
		drv.MBCmds[j].Buffer[0] = uint8(drv.Drv.Addr)
		drv.MBCmds[j].Buffer[1] = uint8(drv.MBCmds[j].Cmd)
		drv.MBCmds[j].Buffer[2] = uint8(drv.MBCmds[j].startAddr >> 8)
		drv.MBCmds[j].Buffer[3] = uint8(drv.MBCmds[j].startAddr)
		drv.MBCmds[j].Buffer[4] = uint8(drv.MBCmds[j].AddrLen >> 8)
		drv.MBCmds[j].Buffer[5] = uint8(drv.MBCmds[j].AddrLen)
		crc := crc16hex(drv.MBCmds[j].Buffer[:6], 6)
		drv.MBCmds[j].Buffer[6] = crc[1]
		drv.MBCmds[j].Buffer[7] = crc[0]
		drv.MBCmds[j].BufferLen = 8
	}
	//logs.Info(drv.MBCmds)
	drv.Videonum = Getdrvvedionum(drv.Drv.Name) //获取摄像头的个数
}
func Getmbdotinfo() {
	MBDrv = nil
	o := orm.NewOrm()
	//获取MODBUS设备的列表信息
	_, err := o.Raw("SELECT * FROM maindrv where packtype = 'MODBUS'").QueryRows(&MBDrv)
	if err != nil {
		logs.Info("ERROR:mb 001", err)
		return
	}
	//根据MODBUS设备列表信息，补充每个MODBUS设备的数据点和其他信息
	MBDrvMap.M.Lock()
	for i := 0; i < len(MBDrv); i++ {
		//根据MODBUS设备的列表信息生成索引MAP
		MBDrvMap.U[MBDrv[i].Drv.Name] = i + 1
		MBDrv[i].Drvname = MBDrv[i].Drv.Name
		_, err = o.Raw("SELECT * FROM maindot where drvname=? order by addr asc", MBDrv[i].Drv.Name).QueryRows(&MBDrv[i].Dot)
		if err != nil {
			logs.Info("ERROR:mb 002", err)
			return
		}
	}
	for i := 0; i < len(MBDrv); i++ {
		for j := 0; j < len(MBDrv[i].Dot); j++ {
			SetOneDotIntoMap(MBDrv[i].Drvname, MBDrv[i].Dot[j].Name)
		}
		CreateMbCmd(&MBDrv[i])
	}
	MBDrvMap.M.Unlock()
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
	} else {
		//开启成功
		defer l.Close()
		logs.Info("Listening on " + ":" + port)
		for {
			// 等待客户端连接
			conn, err := l.Accept()
			if err != nil {
				logs.Info("Error accepting: ", err)
				os.Exit(1)
			}
			//接收到新的客户端连接，保存连接对象到此端口的切片中
			alreadyopenport[port] = append(alreadyopenport[port], conn)
			//开启客户端接收数据线程
			go mbhandleRequest(conn, index)
		}
	}
}

var mbdrvportindex MainMapintints //以通讯端口作为标志的设备索引列表
type MainMapintint struct {
	U map[int]int
	M sync.Mutex
}
type MainMapintints struct {
	U map[int][]int
	C map[int]int
	B map[int]int
	M sync.Mutex
}

var mbdrvportrecvlength MainMapintint
var MBHotReFrashFalg = false //热更新标志

func StartModbus() {
	MBDrvMap.M.Lock()
	MBDrvMap.U = make(map[string]int)
	MBDrvMap.M.Unlock()
	MBDrvMapBak.M.Lock()
	MBDrvMapBak.U = make(map[string]int)
	MBDrvMapBak.M.Unlock()
	mbdrvportrecvlength.M.Lock()
	mbdrvportrecvlength.U = make(map[int]int)
	mbdrvportrecvlength.M.Unlock()
	mbdrvportindex.M.Lock()
	mbdrvportindex.U = make(map[int][]int)
	mbdrvportindex.C = make(map[int]int)
	mbdrvportindex.B = make(map[int]int)
	mbdrvportindex.M.Unlock()
	//开启MODBUS驱动初始化
	Getmbdotinfo() //获取MODBUS设备基本信息
}

var MBTimerMaster *cron.Cron

func ThreadMB() {
	mbdrvportindex.M.Lock()
	for i := 0; i < len(MBDrv); i++ {
		//建立以通信端口为标志的设备索引列表
		mbdrvportindex.U[MBDrv[i].Drv.Port] = append(mbdrvportindex.U[MBDrv[i].Drv.Port], i)
		if mbdrvportindex.C[MBDrv[i].Drv.Port] < MBDrv[i].Drv.Polltime {
			mbdrvportindex.C[MBDrv[i].Drv.Port] = MBDrv[i].Drv.Polltime
		}
		//依次开启设备的通信端口
		go MBStartTCP(strconv.Itoa(MBDrv[i].Drv.Port), MBDrv[i].Drv.Port)
	}
	MBTimerMaster = cron.New()
	for k, _ := range mbdrvportindex.U {
		//根据通信端口开启同一端口的设备轮询线程
		v := k
		mbdrvportrecvlength.M.Lock()
		if mbdrvportindex.C[v] == 0 {
			mbdrvportindex.C[v] = 60
		}
		mbdrvportindex.B[v] = 0
		str := "*/" + strconv.FormatInt(int64(mbdrvportindex.C[v]), 10) + " * * * * *"
		mbdrvportrecvlength.U[v] = 0
		mbdrvportrecvlength.M.Unlock()
		MBTimerMaster.AddFunc(str, func() { mbreqthread(v) })
	}
	mbdrvportindex.M.Unlock()
	MBTimerMaster.Start()
	//开启MODBUS设备主线程
	for {
		//监听热更新标志
		if MBHotReFrashFalg {
			//获取新的设备列表到备用区域
			MBTimerMaster.Stop()
			//以下同MODBUS初始化
			mbdrvportindex.M.Lock()
			for k, _ := range mbdrvportindex.U {
				delete(mbdrvportindex.U, k)
			}
			for i := 0; i < len(MBDrv); i++ {
				if MBDrv[i].Drv.Port != 0 {
					mbdrvportindex.U[MBDrv[i].Drv.Port] = append(mbdrvportindex.U[MBDrv[i].Drv.Port], i)
					if mbdrvportindex.C[MBDrv[i].Drv.Port] < MBDrv[i].Drv.Polltime {
						mbdrvportindex.C[MBDrv[i].Drv.Port] = MBDrv[i].Drv.Polltime
					}
					go MBStartTCP(strconv.Itoa(MBDrv[i].Drv.Port), MBDrv[i].Drv.Port)
				}
			}
			MBTimerMaster = cron.New()
			for k, _ := range mbdrvportindex.U {
				//根据通信端口开启同一端口的设备轮询线程
				v := k
				mbdrvportrecvlength.M.Lock()
				if mbdrvportindex.C[v] == 0 {
					mbdrvportindex.C[v] = 60
				}
				mbdrvportindex.B[v] = 0
				str := "*/" + strconv.FormatInt(int64(mbdrvportindex.C[v]), 10) + " * * * * *"
				mbdrvportrecvlength.U[v] = 0
				mbdrvportrecvlength.M.Unlock()
				MBTimerMaster.AddFunc(str, func() { mbreqthread(v) })
			}
			mbdrvportindex.M.Unlock()
			MBTimerMaster.Start()
			MBHotReFrashFalg = false
		}
		//简单的延时
		time.Sleep(1e9)
	}
}

//设置MODBUS的IO输出
//参数：被操作的设备，被操作的数据点，目标值
func SetMBioValue(drv string, dot string, value int) {
	//从设备MAP中找到此设备的索引
	MBDrvMap.M.Lock()
	index := MBDrvMap.U[drv]
	MBDrvMap.M.Unlock()
	//判断设备名称是否合法 只有再MAP找不到对应的值的时候才会为0，建立MAP时默认最小为1
	if index > 0 {
		//遍历设备中的数据点
		for di, c := range MBDrv[index-1].Dot {
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
				go mbctrlthread(index-1, MBDrv[index-1].Drv.Port, di, float32(value))
			}
		}
	}
}

//设置MODBUS寄存器
//参数：被操作的设备，被操作的数据点，目标值
func SetMBregValue(drv string, dot string, value float32) {
	//从设备MAP中找到此设备的索引
	MBDrvMap.M.Lock()
	index := MBDrvMap.U[drv]
	MBDrvMap.M.Unlock()
	len := 0
	//先格式化目标值为INT型，float型无法执行>>（位移）操作
	intvalue := int(value)
	//判断设备名称是否合法 只有再MAP找不到对应的值的时候才会为0，建立MAP时默认最小为1
	if index > 0 {
		for di, c := range MBDrv[index-1].Dot {
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
				go mbctrlthread(index-1, MBDrv[index-1].Drv.Port, di, value)
			}
		}
	}
}
func mbctrlthread(i int, port int, dotindex int, value float32) {
	timecount := 0

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
			mbdrvportrecvlength.M.Lock()
			mbdrvportrecvlength.U[port] = 0
			mbdrvportrecvlength.M.Unlock()
			_, err := con.Write(MBDrv[i].MBCmds[MBDrv[i].MBCmdNum].Buffer[:MBDrv[i].MBCmds[MBDrv[i].MBCmdNum].BufferLen])
			if err != nil {
				logs.Info(err)
				//发送失败，一般是连接失效，需要删除这个连接
				alreadyopenport[strconv.Itoa(port)] = append(alreadyopenport[strconv.Itoa(port)][:index], alreadyopenport[strconv.Itoa(port)][index+1:]...)
			}
			//改变设备状态
			MBDrv[i].MBStatus = mbbusy
		}
	}

	for {
		time.Sleep(5e8)
		timecount++
		if MBDrv[i].MBStatus == mbready {
			if MBDrv[i].Dot[dotindex].Type > 2 {
				MBDrv[i].Dot[dotindex].Value = value
			} else {
				if value == 0 {
					MBDrv[i].Dot[dotindex].Value = 0
				} else {
					MBDrv[i].Dot[dotindex].Value = 1
				}
			}
			break
		}
		if timecount > 10 {
			MBCtrlCmdRltChanMap[MBDrv[i].Drvname] <- "ERR"
		}
	}
}
func mblitethread(i int, port int) {
	timecount := 0
	//logs.Info("启动单")
	for {
		//检测设备的状态
		//logs.Info("设备状态", MBDrv[i])
		if MBDrv[i].MBStatus == mbready {
			//设备处于空闲状态
			timecount = 0
			if len(alreadyopenport[strconv.Itoa(port)]) == 0 {
				MBDrv[i].MBStatus = mbtimeout
				MBDrv[i].Online = false
				logs.Info("此设备连接丢失", MBDrv[i].Drvname)
				break
			}
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
					mbdrvportrecvlength.M.Lock()
					//logs.Info(MBDrv[i].Drvname, "S", mbdrvportrecvlength.U[port])
					mbdrvportrecvlength.U[port] = 0
					//logs.Info(MBDrv[i].Drvname, "E", MBDrv[i].MBCmdNum)
					mbdrvportrecvlength.M.Unlock()
					_, err := con.Write(MBDrv[i].MBCmds[MBDrv[i].MBCmdNum].Buffer[:MBDrv[i].MBCmds[MBDrv[i].MBCmdNum].BufferLen])
					if err != nil {
						logs.Info(err)
						//发送失败，一般是连接失效，需要删除这个连接
						alreadyopenport[strconv.Itoa(port)] = append(alreadyopenport[strconv.Itoa(port)][:index], alreadyopenport[strconv.Itoa(port)][index+1:]...)
					}
					//改变设备状态
					MBDrv[i].MBStatus = mbbusy

				} else {
					//logs.Info("设备里面没有数据点")
					return
				}
			}
		} else if MBDrv[i].MBStatus == mbbusy {
			//判断设备在忙碌状态下的时间，是否超过了采样时间，超过了就认为设备通讯超时
			if timecount >= 50 {
				MBDrv[i].MBStatus = mbtimeout
				MBDrv[i].Online = false
				logs.Info("读取设备信息超时", MBDrv[i].Drvname)
				break
			}
		} else if MBDrv[i].MBStatus == mbtimeout {
			//设备在上一次轮询中超时，则需要跳过一次，等下一次再次轮询
			MBDrv[i].MBStatus = mbready
			break
		} else if MBDrv[i].MBStatus == mbwite {
			// 	//只有所有的命令包发送一边之后才会进入等待状态，等待经过一个采样时间之后，也就是采样周期之后，再次开始
			MBDrv[i].MBStatus = mbready
			break
		}
		time.Sleep(1e8)
		timecount++
	}
	//logs.Info(port)
}

//通信端口的轮询线程
func mbreqthread(port int) {
	//logs.Info(port)
	mbdrvportindex.M.Lock()
	if mbdrvportindex.B[port] == 0 {
		mbdrvportindex.B[port] = 1
		for _, i := range mbdrvportindex.U[port] {
			mbdrvportindex.M.Unlock()
			//遍历端口下的所有设备的索引
			v := i
			p := port
			//logs.Info("遍历端口列表", v, p)
			mblitethread(v, p)
			mbdrvportindex.M.Lock()
		}
		mbdrvportindex.B[port] = 0
	} else {
		logs.Info("上一个线程未结束", port)
	}
	mbdrvportindex.M.Unlock()

}
func mbhandleRequest(conn net.Conn, port int) {
	var MBRecv [1024]byte
	var Ubuffer [1024]byte
	ipStr := conn.RemoteAddr().String()
	defer func() {
		logs.Info("Disconnected :" + ipStr)
		conn.Close()
	}()
	for {
		//等待数据接收

		n, err := conn.Read(Ubuffer[0:])
		mbdrvportrecvlength.M.Lock()
		bufferlength := mbdrvportrecvlength.U[port]
		mbdrvportrecvlength.M.Unlock()
		for i := 0; i < n; i++ {
			MBRecv[bufferlength+i] = Ubuffer[i]
		}
		logs.Info(bufferlength, port, n, MBRecv[:bufferlength], Ubuffer[:n])
		n += bufferlength
		if err != nil {
			return
		}
		if n > 0 {
			MBRecv[n] = 0
			//logs.Info(n, MBRecv[:n])
			//注意：下面操作并未考虑TCP数据的粘包问题，因为每个设备有且只有一个连接且此线程为单独设备独占，再加上MODBUS为主从模式，基本不会出现粘包问题
			//接收到数据包，判断是否是合法的数据，判断长度或者特殊字符
			if MBRecv[2] == byte(n-5) {
				mbdrvportrecvlength.M.Lock()
				mbdrvportrecvlength.U[port] = 0
				mbdrvportrecvlength.M.Unlock()
				//获取CRC校验值
				crc := crc16hex(MBRecv[0:n-2], n-2)
				if crc[0] == MBRecv[n-1] && crc[1] == MBRecv[n-2] {
					//CRC校验通过
					mbdrvportindex.M.Lock()
					for _, index := range mbdrvportindex.U[port] {
						//遍历通信端口下的所有设备
						if MBRecv[0] == byte(MBDrv[index].Drv.Addr) && len(MBDrv[index].MBCmds) > 0 && MBDrv[index].MBStatus == mbbusy {
							//判断数据包的地址是否和设备地址一致
							//提取当前发送的请求包中的起始地址
							//logs.Info("错误", index)
							//logs.Info("错误", MBDrv)
							//logs.Info("错误", MBDrv[index])
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
											MBDrv[index].Dot[i].Value = float32(MBRecv[packindex+3+1]) * MBDrv[index].Dot[i].Bynum
											//Inserttodotvalue(MBDrv[index].Dot[i])
										} else if MBDrv[index].Dot[i].Data == dataword {
											MBDrv[index].Dot[i].Value = float32(GetUint16(MBRecv[packindex+3:packindex+5])) * MBDrv[index].Dot[i].Bynum
											//Inserttodotvalue(MBDrv[index].Dot[i])
										} else if MBDrv[index].Dot[i].Data == datadword {
											MBDrv[index].Dot[i].Value = float32(GetUint32(MBRecv[packindex+3:packindex+7])) * MBDrv[index].Dot[i].Bynum
											//Inserttodotvalue(MBDrv[index].Dot[i])
										} else if MBDrv[index].Dot[i].Data == datafloat {
											MBDrv[index].Dot[i].Value = GetFloat32(MBRecv[packindex+3:packindex+7]) * MBDrv[index].Dot[i].Bynum
											//Inserttodotvalue(MBDrv[index].Dot[i])
										}
										if MBDrv[index].Dot[i].Alarmtop > MBDrv[index].Dot[i].Alarmbot {
											if MBDrv[index].Dot[i].Value > MBDrv[index].Dot[i].Alarmtop {
												if MBDrv[index].Dot[i].Status != "TOP" {
													MBDrv[index].Dot[i].Status = "TOP"
													//触发上限报警
													NewAlarmNotice(MBDrv[index].Drvname, MBDrv[index].Dot[i].Name, "TOP")
												}
											} else if MBDrv[index].Dot[i].Value < MBDrv[index].Dot[i].Alarmbot {
												if MBDrv[index].Dot[i].Status != "BOT" {
													MBDrv[index].Dot[i].Status = "BOT"
													//触发下限报警
													NewAlarmNotice(MBDrv[index].Drvname, MBDrv[index].Dot[i].Name, "BOT")
												}
											} else {
												if MBDrv[index].Dot[i].Status == "BOT" || MBDrv[index].Dot[i].Status == "TOP" {
													MBDrv[index].Dot[i].Status = "OK"
													//报警解除
													NewAlarmNotice(MBDrv[index].Drvname, MBDrv[index].Dot[i].Name, "OK")
												}
											}
										}
										SetOneDotValueToMap(MBDrv[index].Drvname, MBDrv[index].Dot[i].Name, MBDrv[index].Dot[i].Value)
									} else if (MBRecv[1] == inio || MBRecv[1] == outio) && MBDrv[index].Dot[i].Type == int(MBRecv[1]) {
										//IO点数据
										//IO点的数据返回是按位来返回的，每8个点占用一个字节，所以需要通过移位来提取数据
										packindex := MBDrv[index].Dot[i].Addr - packstartaddr
										packnum := packindex / 8
										valueindex := uint8(packindex % 8)
										if (MBRecv[3+packnum]>>valueindex)&0x01 == 1 {
											MBDrv[index].Dot[i].Value = 1
											//Inserttodotvalue(MBDrv[index].Dot[i])
										} else {
											MBDrv[index].Dot[i].Value = 0
											//Inserttodotvalue(MBDrv[index].Dot[i])
										}
									} else if MBRecv[1] == setio {
										//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
										for ci, c := range MBDrv[index].MBCmds {
											if c.Type == Ctr {
												MBDrv[index].MBCmdNum = 0
												logs.Info(MBDrv[index].Drvname)
												MBCtrlCmdRltChanMap[MBDrv[index].Drvname] <- "OK"
												MBDrv[index].MBStatus = mbready
												MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
											}
										}
									} else if MBRecv[1] == setreg {
										//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
										for ci, c := range MBDrv[index].MBCmds {
											if c.Type == Ctr {
												MBDrv[index].MBCmdNum = 0
												MBCtrlCmdRltChanMap[MBDrv[index].Drvname] <- "OK"
												MBDrv[index].MBStatus = mbready
												MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
											}
										}
									} else if MBRecv[1] == setregs {
										//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
										for ci, c := range MBDrv[index].MBCmds {
											if c.Type == Ctr {
												MBDrv[index].MBCmdNum = 0
												MBCtrlCmdRltChanMap[MBDrv[index].Drvname] <- "OK"
												MBDrv[index].MBStatus = mbready
												MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
											}
										}
									}
								}
							}
							MBDrv[index].Flashtime = time.Now().Format("2006-01-02 15:04")
							MBDrv[index].Online = true
							if MBDrv[index].MBCmdNum >= len(MBDrv[index].MBCmds) {
								MBDrv[index].MBCmdNum = 0
								MBDrv[index].MBStatus = mbwite
								TimerSaveDrvDotValue(MBDrv[index].Drvname, 1)
							} else {
								MBDrv[index].MBStatus = mbready
							}
						}
					}
					mbdrvportindex.M.Unlock()
				} else {
					logs.Info("crc err")
				}
			} else {

				if MBRecv[1] == setio || MBRecv[1] == setreg || MBRecv[1] == setregs {
					//获取CRC校验值
					mbdrvportrecvlength.M.Lock()
					mbdrvportrecvlength.U[port] = 0
					mbdrvportrecvlength.M.Unlock()
					crc := crc16hex(MBRecv[0:n-2], n-2)
					if crc[0] == MBRecv[n-1] && crc[1] == MBRecv[n-2] {
						//CRC校验通过
						mbdrvportindex.M.Lock()
						for _, index := range mbdrvportindex.U[port] {
							//遍历通信端口下的所有设备
							if MBRecv[0] == byte(MBDrv[index].Drv.Addr) && MBDrv[index].MBStatus == mbbusy {
								//判断数据包的地址是否和设备地址一致
								//提取当前发送的请求包中的起始地址
								packstartaddr := MBDrv[index].MBCmds[MBDrv[index].MBCmdNum].startAddr
								//提取当前发送的请求包中的结束地址
								packstopaddr := packstartaddr + MBDrv[index].MBCmds[MBDrv[index].MBCmdNum].AddrLen
								for i := 0; i < len(MBDrv[index].Dot); i++ {
									//遍历设备的所有数据点
									//判断数据点的地址是否在起始地址和结束地址之间
									if MBDrv[index].Dot[i].Addr >= packstartaddr && MBDrv[index].Dot[i].Addr < packstopaddr {
										//判断数据点的类型是否和返回数据包中的匹配
										if MBRecv[1] == setio {
											//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
											for ci, c := range MBDrv[index].MBCmds {
												if c.Type == Ctr {
													MBDrv[index].MBCmdNum = 0
													logs.Info(MBDrv[index].Drvname)
													MBCtrlCmdRltChanMap[MBDrv[index].Drvname] <- "OK"
													MBDrv[index].MBStatus = mbready
													MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
												}
											}
										} else if MBRecv[1] == setreg {
											//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
											for ci, c := range MBDrv[index].MBCmds {
												if c.Type == Ctr {
													MBDrv[index].MBCmdNum = 0
													MBCtrlCmdRltChanMap[MBDrv[index].Drvname] <- "OK"
													MBDrv[index].MBStatus = mbready
													MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
												}
											}
										} else if MBRecv[1] == setregs {
											//设置Ctr数据包的返回，确认状态，直接删除Ctr数据包，重置数据包的计数值
											for ci, c := range MBDrv[index].MBCmds {
												if c.Type == Ctr {
													MBDrv[index].MBCmdNum = 0
													MBCtrlCmdRltChanMap[MBDrv[index].Drvname] <- "OK"
													MBDrv[index].MBStatus = mbready
													MBDrv[index].MBCmds = append(MBDrv[index].MBCmds[:ci], MBDrv[index].MBCmds[ci+1:]...)
												}
											}
										}
									}
								}
							}
						}
						mbdrvportindex.M.Unlock()
					} else {
						logs.Info("crc err")
					}
				} else {
					mbdrvportrecvlength.M.Lock()
					mbdrvportrecvlength.U[port] = n
					mbdrvportrecvlength.M.Unlock()
					logs.Info("length err:", n)
				}

			}
		}
	}
	logs.Info("Done!")
}
