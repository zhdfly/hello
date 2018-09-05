package tcpserver

import (
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/robfig/cron"
)

var DrvDotHighValueMap MainMAPstringfloat //= make(map[string]float32) //记录每个数据点的单日最高值
var DrvDotLowValueMap MainMAPstringfloat  //= make(map[string]float32)  //记录每个数据点的单日最低值

var DrvDotHourAvgValueMap MainMAPstringfloat //= make(map[string]float32) //记录每个数据点的一小时平均值

func SetOneDotIntoMap(drvname string, dotname string) {
	DrvDotHighValueMap.M.Lock()
	DrvDotHighValueMap.U[drvname+"-"+dotname] = 0
	DrvDotHighValueMap.M.Unlock()
	DrvDotLowValueMap.M.Lock()
	DrvDotLowValueMap.U[drvname+"-"+dotname] = 0
	DrvDotLowValueMap.M.Unlock()
	DrvDotHourAvgValueMap.M.Lock()
	DrvDotHourAvgValueMap.U[drvname+"-"+dotname] = 0
	DrvDotHourAvgValueMap.M.Unlock()
}

func ClearOneDotIntoMap(drvname string, dotname string) {
	DrvDotHighValueMap.M.Lock()
	delete(DrvDotHighValueMap.U, drvname+"-"+dotname)
	DrvDotHighValueMap.M.Unlock()
	DrvDotLowValueMap.M.Lock()
	delete(DrvDotLowValueMap.U, drvname+"-"+dotname)
	DrvDotLowValueMap.M.Unlock()
	DrvDotHourAvgValueMap.M.Lock()
	delete(DrvDotHourAvgValueMap.U, drvname+"-"+dotname)
	DrvDotHourAvgValueMap.M.Unlock()
}

func SetOneDotValueToMap(drvname string, dotname string, value float32) {
	if value == 0 {
		return
	}
	DrvDotHighValueMap.M.Lock()
	HighOldValue := DrvDotHighValueMap.U[drvname+"-"+dotname]
	if value > HighOldValue {
		DrvDotHighValueMap.U[drvname+"-"+dotname] = value
	}
	DrvDotHighValueMap.M.Unlock()
	DrvDotLowValueMap.M.Lock()
	LowOldValue := DrvDotLowValueMap.U[drvname+"-"+dotname]
	if value < LowOldValue || LowOldValue == 0 {
		DrvDotLowValueMap.U[drvname+"-"+dotname] = value
	}
	DrvDotLowValueMap.M.Unlock()
	DrvDotHourAvgValueMap.M.Lock()
	AvgOldValue := DrvDotHourAvgValueMap.U[drvname+"-"+dotname]
	if AvgOldValue == 0 {
		DrvDotHourAvgValueMap.U[drvname+"-"+dotname] = value
	} else {
		DrvDotHourAvgValueMap.U[drvname+"-"+dotname] = (AvgOldValue + value) / 2
	}
	DrvDotHourAvgValueMap.M.Unlock()
}
func RunDailySaveThread() {
	DrvDotHighValueMap.M.Lock()
	DrvDotLowValueMap.M.Lock()
	for c, v := range DrvDotHighValueMap.U {
		drvname := strings.Split(c, "-")[0]
		dotname := strings.Split(c, "-")[1]
		x := DrvDotLowValueMap.U[c]
		SaveDailyValue(drvname, dotname, v, x)
		logs.Info(v, c, drvname)
	}
	DrvDotHighValueMap.M.Unlock()
	DrvDotLowValueMap.M.Unlock()
}
func RunHourlySaveThread() {
	DrvDotHourAvgValueMap.M.Lock()
	for c, v := range DrvDotHourAvgValueMap.U {
		drvname := strings.Split(c, "-")[0]
		dotname := strings.Split(c, "-")[1]
		SaveHourlyValue(drvname, dotname, v)
		logs.Info(v, c, drvname)
	}
	DrvDotHourAvgValueMap.M.Unlock()
}
func StatisDataConfig() {
	DrvDotHighValueMap.U = make(map[string]float32)
	DrvDotLowValueMap.U = make(map[string]float32)
	DrvDotHourAvgValueMap.U = make(map[string]float32)
	c := cron.New()
	c.AddFunc("@daily", func() { RunDailySaveThread() })
	c.AddFunc("@hourly", func() { RunHourlySaveThread() })
	c.Start()

}
