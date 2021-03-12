package middleware

import (
	core "onebot/core/xianqu"
	"strings"
)

// RemoveAsync 去除 action 中的 async
func RemoveAsync(ctx *core.Context) {
	request := core.Parse(ctx.Request)
	if !request.Exist("action") {
		return
	}
	action := request.Str("action")
	// 保证不是拷贝的
	ctx.Request["action"] = strings.ReplaceAll(action, "_async", "")
	return
}

// IsAsync 返回 action 是否为异步调用
func IsAsync(ctx *core.Context) bool {
	request := core.Parse(ctx.Request)
	if !request.Exist("action") {
		return false
	}
	action := request.Str("action")
	if strings.Contains(action, "_async") {
		return true
	}
	return false
}
