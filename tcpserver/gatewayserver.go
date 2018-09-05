package tcpserver

import (
	"net"

	"github.com/astaxie/beego"
)

var SensorName = [8]string{"空气温度", "空气湿度", "土壤温度", "土壤湿度", "光照强度", "二氧化碳", "", ""}
var SensorUnit = [8]string{"℃", "%", "℃", "%", "KLux", "PPM", "", ""}

func gatewaycrc16hex(str []byte, len int) [2]byte {
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
func chkError(err error) {
	if err != nil {
		beego.Info(err)
	}
}

type Gatewaynodechild struct {
	Name  string
	Unit  string
	Bynum float32
	Value float32
}
type Getawaynode struct {
	Id           int
	Addr         int
	Noderssi     int
	Datatype     int
	Nodetype     int
	Nodetypeplus []byte
	Node         map[string]Gatewaynodechild
	Gatewayrssi  int
	Gatewatdb    int
	Time         string
	Online       bool
}
type Gatewaymain struct {
	Name   string
	Port   int
	Node   map[int]Getawaynode
	Online bool
	Time   string
}

var GateWayMain Gatewaymain

func clientHandle(conn *net.UDPConn) {
	//defer conn.Close()
	buf := make([]byte, 20480)
	//读取数据
	//注意这里返回三个参数
	//第二个是udpaddr
	//下面向客户端写入数据时会用到
	buflength, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return
	}
	if buf[0] == 0xaa && buf[1] == 0x55 {
		var len int
		len = int(buf[3])*256 + int(buf[2])
		if len == buflength-10 {
			crc := gatewaycrc16hex(buf, buflength-3)
			if crc[0] == buf[buflength-2] && crc[1] == buf[buflength-3] {
				packnums := buf[4]
				packindex := buf[5]
				if packnums == packindex {
					//最后一个包了，可能传感器数量不为10

					bufdata := buf[7 : len+7]

					if buf[6] == 0x02 {
						nodenums := len / 44
						for i := 0; i < nodenums; i++ {
							if bufdata[i*44] != 0 {
								//数据包有有效数据
								s := GateWayMain.Node[(int(packindex)-1)*10+i]
								beego.Info(s)
								if s.Addr != 0 {
									//已存在节点对象
								} else {
									//不存在节点对象，说明此节点不需要被采集
									var nodebuf Getawaynode
									nodebuf.Nodetypeplus = make([]byte, 4)
									nodebuf.Node = make(map[string]Gatewaynodechild)
									nodebufdata := bufdata[i*44 : i*44+44]
									nodebuf.Addr = int(nodebufdata[0])
									nodebuf.Noderssi = int(nodebufdata[1])
									nodebuf.Datatype = int(nodebufdata[2])
									nodebuf.Nodetype = int(nodebufdata[3])
									nodebuf.Nodetypeplus = nodebufdata[4:8]
									nodebuf.Gatewayrssi = int(nodebufdata[36])
									nodebuf.Gatewatdb = int(nodebufdata[37])
									if nodebuf.Datatype == 0 {
										for n := 0; n < 8; n++ {
											if ((nodebuf.Nodetype)>>uint(n))&0x01 != 0 {
												sersornamebuf := SensorName[n]
												v := nodebuf.Node[sersornamebuf]
												if v.Name != sersornamebuf {
													var vbuf Gatewaynodechild
													vbuf.Name = sersornamebuf
													vbuf.Unit = SensorUnit[n]
													vbuf.Value = GetFloat32(nodebufdata[12+n*4 : 16+n*4])
													nodebuf.Node[sersornamebuf] = vbuf
												} else {
													v.Value = GetFloat32(nodebufdata[12+n*4 : 16+n*4])
												}
											}
										}
									}
									GateWayMain.Node[nodebuf.Addr] = nodebuf
								}
							}
						}
					}
				} else {

				}
			}
		}
	}
}
func GatewayLiteThread(port int) {
	GateWayMain.Node = make(map[int]Getawaynode)
	//监听端口
	udpconn, err2 := net.ListenUDP("udp4", &net.UDPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: port,
	})
	chkError(err2)
	//udp没有对客户端连接的Accept函数
	for {
		clientHandle(udpconn)
	}
}
