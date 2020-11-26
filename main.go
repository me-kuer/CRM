package main

import (
	_ "crm-server/models"
	_ "crm-server/routers"
	"crm-server/utils"
	"github.com/astaxie/beego"
)

func main() {
	// 注册模板函数
	beego.AddFuncMap("str2date", utils.FormatDate)
	beego.AddFuncMap("checkPower", utils.CheckPower)
	beego.AddFuncMap("checkModule", utils.CheckModule)

	// session配置
	beego.BConfig.WebConfig.Session.SessionProvider = "file"
	beego.BConfig.WebConfig.Session.SessionProviderConfig = "./tmp"
	beego.Run()
}
