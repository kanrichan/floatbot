package gateway

import (
	"encoding/json"
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

// 向框架发送 INFO 日志
func INFO(s string, v ...interface{}) {
	core.ApiOutPutLog("[I]" + fmt.Sprintf(s, v...))
}

// 向框架发送 WARN 日志
func WARN(s string, v ...interface{}) {
	core.ApiOutPutLog("[W]" + fmt.Sprintf(s, v...))
}

// 向框架发送 DEBUG 日志
func DEBUG(s string, v ...interface{}) {
	if CONF.Debug {
		core.ApiOutPutLog("[D]" + fmt.Sprintf(s, v...))
	}
}

// 向框架发送 ERROR 日志
func ERROR(s string, v ...interface{}) {
	core.ApiOutPutLog("[E]" + fmt.Sprintf(s, v...))
}

// printCall 打印调用信息
func printCall(ctx *core.Context) {
	params, _ := json.Marshal(ctx.Request["params"])
	INFO("[收到调用][%d] 标记: %v API: %s 参数: %s",
		ctx.Bot,
		ctx.Request["echo"],
		ctx.Request["action"],
		params,
	)
}

// printBack 打印返回信息
func printBack(ctx *core.Context) {
	back, _ := json.Marshal(ctx.Response)
	DEBUG("[返回响应][%d] 标记: %v 返回: %s",
		ctx.Bot,
		ctx.Request["echo"],
		back,
	)
}

// printSend 打印上报信息
func printSend(ctx *core.Context) {
	send, _ := json.Marshal(ctx.Response)
	INFO("[上报信息][%d] 事件: %s",
		ctx.Bot,
		send,
	)
}
