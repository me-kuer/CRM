package utils

import (
	"github.com/astaxie/beego/logs"
)

var Logger *logs.BeeLogger

func init() {
	// 初始化 log
	Logger = logs.NewLogger(10000)

	config := `{
		"filename": "logs/customer.log",
        "maxlines": 1000,
        "maxsize": 10240
	}`

	// 配置log
	logs.Async()	// 设置为异步
	Logger.SetLogger("file", config)
	Logger.SetLevel(logs.LevelDebug)
	Logger.EnableFuncCallDepth(true)
}
