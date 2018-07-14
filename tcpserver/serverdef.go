package tcpserver

const (
	outio   = 1
	inio    = 2
	holdreg = 3
	inreg   = 4
)

//用户的基本结构
type Usr struct {
	Id   int
	Name string
	Pass string
}

//态神设备的基本结构
// type Drv struct {
// 	Id   int
// 	Name string
// 	Port string
// 	Type string
// 	Info string
// 	Time string
// }

//态神数据点的基本结构
// type Dot struct {
// 	Id         int
// 	Name       string
// 	Dottype    string
// 	Datatype   string
// 	Info       string
// 	Val        float32
// 	Drv        string
// 	Warningtop float32
// 	Warningbot float32
// }

//用户设备的结构
type Usrdrv struct {
	Id      int
	Usrname string
	Drvname string
}

//数据点值的储存结构
type Dotvalue struct {
	Id      int
	Drvname string
	Dotname string
	Value   float32
	Status  string
	Time    string
}

//视频设备的基本参数
type Videodrv struct {
	Id          int
	Name        string
	Appkey      string
	Appsecret   string
	Accesstoken string
	Sn          string
	Vercode     string
	Drv         string
	Liveurl     string
}

//设备点的基本参数
type Maindot struct {
	ID          int
	Name        string
	Type        int     //寄存器类型
	Addr        int     //寄存器地址
	Rw          int     //读写方式，默认只读
	Alarmtop    float32 //报警值参数，默认为0，当四个参数都不相同是才会启动报警功能
	Alarmtoptop float32
	Alarmbot    float32
	Alarmbotbot float32 //报警值定义完毕
	Savetime    int     //数据保存时间间隔，默认10S
	Unit        string  //数据单位
	Drvname     string  //隶属于设备
	Value       float32
	Status      string
}

//设备的基本参数
type Maindrv struct {
	Id           int
	Name         string //设备名称
	Addr         int    //地址，通用地址，无论是什么数据包类型
	Port         int    //TCP监听端口
	Packtype     string //MODBUS,MQTT,CONFIG,WHATEVER
	Samplingtime int    //采样间隔时间,默认1000ms
	Retrytime    int    //通讯超时时间，单位秒，超过此时间无信息流，则认为通讯超时
	Retrycount   int    //通讯失败重试次数，0为无限次
	Status       int    //设备状态
	Time         string //泛指添加时间
}
type MainDrvType struct {
	Drv       Maindrv
	Dot       []Maindot
	Sensornum int
	IOnum     int
	Logicnum  int
	Videonum  int
	Flashtime string
	TSUrls    string //态神设备专用读取参数的URL
}

//态神的设备和数据点
// type ModBusDot struct {
// 	Dotname       string
// 	Dottype       string
// 	Dotwarningtop float32
// 	Dotwarningbot float32
// 	Dotstatus     string
// 	Value         float32
// }
// type ModBusDrv struct {
// 	Drvname   string
// 	Sensornum int
// 	IOnum     int
// 	Logicnum  int
// 	Videonum  int
// 	Flashtime string
// 	Dot       []ModBusDot
// }
type MainUserDrv struct {
	User string
	Drv  []string
}

//萤石云JSON数据格式
type AccessTokentmp struct {
	Key  string `json:"accessToken"`
	Time int    `json:"expireTime"`
}
type AccessToken struct {
	Data AccessTokentmp `json:"data"`
	Code string         `json:"code"`
	Msg  string         `json:"msg"`
}
type Comrlt struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}
type Listpage struct {
	Total int `json:"total"`
	Page  int `json:"page"`
	Size  int `json:"size"`
}
type Lsitvideo struct {
	DeviceSerial string `json:"deviceSerial"`
	ChannelNo    int    `json:"channelNo"`
	LiveAddress  string `json:"liveAddress"`
	HdAddress    string `json:"hdAddress"`
	Rtmp         string `json:"rtmp"`
	RtmpHd       string `json:"rtmpHd"`
	Status       int    `json:"status"`
	Exception    int    `json:"exception"`
	BeginTime    int    `json:"beginTime"`
	EndTime      int    `json:"endTime"`
}
type Getlist struct {
	Page Listpage    `json:"page"`
	Data []Lsitvideo `json:"data"`
	Code string      `json:"code"`
	Msg  string      `json:"mag"`
}
type Drvvideo struct {
	Name    string
	Liveurl string
}
