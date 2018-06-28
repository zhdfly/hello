package tcpserver

func Addnewuser(name string) {
	var tmp ModBusUserDrv
	tmp.User = name
	Mbdrv = append(Mbdrv, tmp)
}
func Addnewdrv(user string, drv string) {
	var tmp ModBusDrv
	tmp.Drvname = drv
	for i := 0; i < len(Mbdrv); i++ {
		if user == Mbdrv[i].User {
			Mbdrv[i].Drv = append(Mbdrv[i].Drv, tmp)
		}
	}
}
func Addnewuserdrv(user string, drv string) {
	var tmp ModBusDrv
	tmp.Drvname = drv
	for i := 0; i < len(dotob); i++ {
		if dotob[i].Drv == drv {
			var dottmp ModBusDot
			dottmp.Dotname = dotob[i].Name
			dottmp.Dottype = dotob[i].Datatype
			dottmp.Value = dotob[i].Val
			tmp.Dot = append(tmp.Dot, dottmp)
		}
	}
	for i := 0; i < len(Mbdrv); i++ {
		if user == Mbdrv[i].User {
			Mbdrv[i].Drv = append(Mbdrv[i].Drv, tmp)
		}
	}
}
func Addnewdot(name string, datatype string, dottype string, drv string) {
	var tmp ModBusDot
	tmp.Dotname = name
	tmp.Dottype = datatype
	var tmps Dot
	tmps.Name = name
	tmps.Datatype = datatype
	tmps.Dottype = dottype
	tmps.Drv = drv
	dotob = append(dotob, tmps)
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drv); j++ {
			if drv == Mbdrv[i].Drv[j].Drvname {
				Mbdrv[i].Drv[j].Dot = append(Mbdrv[i].Drv[j].Dot, tmp)
			}
		}
	}
	Creaturl()
}
func Dltdot(name string, drv string) {
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drv); j++ {
			if drv == Mbdrv[i].Drv[j].Drvname {
				for k := 0; k < len(Mbdrv[i].Drv[j].Dot); k++ {
					if name == Mbdrv[i].Drv[j].Dot[k].Dotname {
						Mbdrv[i].Drv[j].Dot = append(Mbdrv[i].Drv[j].Dot[:k], Mbdrv[i].Drv[j].Dot[k+1:]...)
					}
				}
			}
		}
	}
	for i := 0; i < len(dotob); i++ {
		if dotob[i].Name == name && dotob[i].Drv == drv {
			dotob = append(dotob[:i], dotob[i+1:]...)
		}
	}
	Creaturl()
}
func Dltdrv(name string) {
	for i := 0; i < len(Mbdrv); i++ {
		for j := 0; j < len(Mbdrv[i].Drv); j++ {
			if name == Mbdrv[i].Drv[j].Drvname {
				Mbdrv[i].Drv = append(Mbdrv[i].Drv[:j], Mbdrv[i].Drv[j+1:]...)
			}
		}
	}
	for i := 0; i < len(dotob); i++ {
		if dotob[i].Drv == name {
			dotob = append(dotob[:i], dotob[i+1:]...)
		}
	}
	Creaturl()
}
func Dltuser(name string) {
	for i := 0; i < len(Mbdrv); i++ {
		if name == Mbdrv[i].User {
			Mbdrv = append(Mbdrv[:i], Mbdrv[i+1:]...)
		}
	}
}
