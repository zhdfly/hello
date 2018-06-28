package tcpserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type AccessTokentmp struct {
	Key  string `json:"key"`
	Time string `json:"time"`
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

var VideoList Getlist

func Getvideolist() {
	resp, err := http.Post("https://open.ys7.com/api/lapp/live/video/list",
		"application/x-www-form-urlencoded",
		strings.NewReader("accessToken=at.do18wcz289ffyo5v5uiqu8fa9kyo03jq-39in4kqjap-1ca5m0g-rjqs1u5ub"))
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
	}
	fmt.Println(VideoList.Data[0].LiveAddress)
}
