package gateway

import (
	"fmt"
	core "onebot/core/xianqu"
	"onebot/server"
)

func init() {
	// 将其他package的LOG()连接到这里
	server.LOG = INFO
}

func INFO(s string, v ...interface{}) {
	core.XQApiOutPutLog("[INFO] " + fmt.Sprintf(s, v...))
}

func WARN(s string, v ...interface{}) {
	core.XQApiOutPutLog("[WARN] " + fmt.Sprintf(s, v...))
}

func DEBUG(s string, v ...interface{}) {
	core.XQApiOutPutLog("[DEBUG] " + fmt.Sprintf(s, v...))
}

func ERROR(s string, v ...interface{}) {
	core.XQApiOutPutLog("[ERROR] " + fmt.Sprintf(s, v...))
}
