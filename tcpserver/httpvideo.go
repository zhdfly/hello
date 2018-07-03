package tcpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/astaxie/beego/orm"
)

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

var VideoList Getlist
var Videodrvlist []Videodrv
var Tokenlist []Videodrv

func Getdrvvedio(drv string) (string, error) {
	var tl []Drvvideo
	for i := 0; i < len(Videodrvlist); i++ {
		if Videodrvlist[i].Drv == drv {
			var tmp Drvvideo
			tmp.Name = Videodrvlist[i].Name
			tmp.Liveurl = Videodrvlist[i].Liveurl
			tl = append(tl, tmp)
		}
	}
	data, err := json.Marshal(tl)
	return string(data), err
}

func Getdrvvedionum(drv string) int {
	num := 0
	for i := 0; i < len(Videodrvlist); i++ {
		if Videodrvlist[i].Drv == drv {
			num++
		}
	}
	return num
}

func Getappaccesstekon(appkey string, appaccess string) (string, string, error) {
	var para = url.Values{}
	para.Add("appKey", appkey)
	para.Add("appSecret", appaccess)
	data := para.Encode()
	resp, err := http.Post("https://open.ys7.com/api/lapp/token/get",
		"application/x-www-form-urlencoded",
		strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// handle error
		return "请求错误", "ERR", err
	}
	var token AccessToken
	fmt.Println(string(body))
	err = json.Unmarshal(body, &token)
	if err != nil {
		fmt.Println(err)
	}
	if token.Code == "200" {
		return token.Data.Key, "OK", err
	}
	if token.Code == "10005" {
		return "APPKEY错误", "ERR", err
	}
	if token.Code == "10017" {
		return "APPKEY不存在", "ERR", err
	}
	if token.Code == "10030" {
		return "APPKEY与APPACCESS不匹配", "ERR", err
	}
	return "请求错误", "ERR", err
}
func Getvideolist() {
	o := orm.NewOrm()
	_, err := o.Raw("select *, count(distinct accesstoken) from videodrv group by accesstoken").QueryRows(&Tokenlist)
	if err == nil {
		fmt.Println(Tokenlist)
	}
	_, err = o.Raw("select * from videodrv").QueryRows(&Videodrvlist)
	if err == nil {
		fmt.Println(Videodrvlist)
	}
	for i := 0; i < len(Tokenlist); i++ {
		resp, err := http.Post("https://open.ys7.com/api/lapp/live/video/list",
			"application/x-www-form-urlencoded",
			strings.NewReader("accessToken="+Tokenlist[i].Accesstoken))
		if err != nil {
			fmt.Println(err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			// handle error
		}
		err = json.Unmarshal(body, &VideoList)
		if err != nil {
			fmt.Println("Get Video LIst Error")
			var com Comrlt
			err = json.Unmarshal(body, &com)
		} else {
			switch VideoList.Code {
			case "200":
				for j := 0; j < len(VideoList.Data); j++ {
					for k := 0; k < len(Videodrvlist); k++ {
						if VideoList.Data[j].DeviceSerial == Videodrvlist[k].Sn {
							Videodrvlist[k].Liveurl = VideoList.Data[j].LiveAddress
							break
						}
					}
				}
				break
			case "10002":
				token, sta, err := Getappaccesstekon(Tokenlist[i].Appkey, Tokenlist[i].Appsecret)
				if sta == "OK" {
					Tokenlist[i].Accesstoken = token
					_, err = o.Raw("update videodrv set accesstoken = ? where appkey = ? and appaccess = ?", token, Tokenlist[i].Appkey, Tokenlist[i].Appsecret).Exec()
					if err == nil {

					}
					i--
				}
				break //TOKEN过期，需要重新获取TOKEN
			case "10001":
				token, sta, err := Getappaccesstekon(Tokenlist[i].Appkey, Tokenlist[i].Appsecret)
				if sta == "OK" {
					Tokenlist[i].Accesstoken = token
					_, err = o.Raw("update videodrv set accesstoken = ? where appkey = ? and appsecret = ?", token, Tokenlist[i].Appkey, Tokenlist[i].Appsecret).Exec()
					if err == nil {

					}
					i--
				}
				break //TOKEN过期，需要重新获取TOKEN
			case "10005":
				break //APPKEY被冻结，需要和萤石云联系
			}
		}
		fmt.Println(VideoList)
	}
}
