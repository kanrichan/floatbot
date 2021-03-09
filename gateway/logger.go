package gateway

import (
	"fmt"
	core "onebot/core/xianqu"
	ser "onebot/server"
)

func init() {
	// 将其他package的LOG()连接到这里
	ser.CoreInfo = func(s string, v ...interface{}) { core.ApiOutPutLog(fmt.Sprintf(s, v...)) }
	ser.CoreDebug = func(s string, v ...interface{}) {
		if CONF.Debug {
			core.ApiOutPutLog(fmt.Sprintf(s, v...))
		}
	}
}

func INFO(s string, v ...interface{}) {
	core.ApiOutPutLog("[INFO] " + fmt.Sprintf(s, v...))
}

func WARN(s string, v ...interface{}) {
	core.ApiOutPutLog("[WARN] " + fmt.Sprintf(s, v...))
}

func DEBUG(s string, v ...interface{}) {
	core.ApiOutPutLog("[DEBUG] " + fmt.Sprintf(s, v...))
}

func ERROR(s string, v ...interface{}) {
	core.ApiOutPutLog("[ERROR] " + fmt.Sprintf(s, v...))
}
