package tcpserver

import (
	"encoding/json"

	"github.com/astaxie/beego/logs"
)

var UserAreaList []Userarea

func AreaConfig() {
	UserAreaList = GetArea()
	//logs.Info(UserAreaList)
}
func GetUserArea(user interface{}) string {
	var userlist []Userarea
	for i := 0; i < len(UserAreaList); i++ {
		if UserAreaList[i].User == user {
			userlist = append(userlist, UserAreaList[i])
		}
	}
	str, _ := json.Marshal(map[string]interface{}{"User": user, "Area": userlist})
	logs.Info(string(str))
	return string(str)
}
