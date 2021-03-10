package gateway

import (
	"encoding/json"

	core "onebot/core/xianqu"
	middle "onebot/middleware"
	ser "onebot/server"
)

func init() {
	// 将所有OneBot事件连接到这里
	core.OnEnable = OnEnable
	core.OnDisable = OnDisable
	core.OnSetting = OnSetting
	core.OnMessagePrivate = OnMessagePrivate
	core.OnMessageGroup = OnMessageGroup
	core.OnNoticeFileUpload = OnNoticeFileUpload
	core.OnNoticeAdminChange = OnNoticeAdminChange
	core.OnNoticeGroupDecrease = OnNoticeGroupDecrease
	core.OnNoticeGroupIncrease = OnNoticeGroupIncrease
	core.OnNoticeGroupBan = OnNoticeGroupBan
	core.OnNoticeFriendAdd = OnNoticeFriendAdd
	core.OnNoticeMessageRecall = OnNoticeMessageRecall
	core.OnRequestFriendAdd = OnRequestFriendAdd
	core.OnRequestGroupAdd = OnRequestGroupAdd

	// 将所有的连接的Handler连接到这里
	ser.HttpHandler = HttpHandler
	ser.HttpPostHandler = HttpPostHandler
	ser.WSCHandler = WebSocketClientHandler
	ser.WSSHandler = WebSocketServerHandler
}

// HttpHandler Http的Handler
func HttpHandler(bot int64, path string, data []byte) []byte {
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	ctx := &core.Context{
		Bot: bot,
		Request: map[string]interface{}{
			"action": path,
			"params": temp,
		},
	}
	middle.RequestToArray(ctx)
	callapi(ctx)
	rsp, _ := json.Marshal(ctx.Response)
	return rsp
}

// HttpPostHandler 快速回复的Handler
func HttpPostHandler(bot int64, send, data []byte) {
	rsp := map[string]interface{}{}
	json.Unmarshal(send, &rsp)
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	ctx := &core.Context{
		Bot:      bot,
		Response: rsp,
		Request:  temp,
	}
	// 将快速回复转化成正常的onebot标准报文
	middle.RequestFastReplyFormat(ctx)
	middle.RequestToArray(ctx)
	callapi(ctx)
}

// WebSocketClientHandler 反向ws的Handler
func WebSocketClientHandler(bot int64, data []byte) []byte {
	request := map[string]interface{}{}
	json.Unmarshal(data, &request)
	ctx := &core.Context{
		Bot:     bot,
		Request: request,
	}
	middle.RequestToArray(ctx)
	callapi(ctx)
	rsp, _ := json.Marshal(ctx.Response)
	return rsp
}

// WebSocketServerHandler 正向ws的Handler
func WebSocketServerHandler(bot int64, data []byte) []byte {
	request := map[string]interface{}{}
	json.Unmarshal(data, &request)
	ctx := &core.Context{
		Bot:     bot,
		Request: request,
	}
	middle.RequestToArray(ctx)
	callapi(ctx)
	rsp, _ := json.Marshal(ctx.Response)
	return rsp
}

// OnMessagePrivate 收到私聊信息事件被触发
func OnMessagePrivate(ctx *core.Context) {
	OnEvent(ctx)
}

// OnMessageGroup 收到群聊信息事件被触发
func OnMessageGroup(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeFileUpload 收到群文件上传事件被触发
func OnNoticeFileUpload(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeAdminChange 收到上下管理事件被触发
func OnNoticeAdminChange(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeGroupDecrease 收到群成员减少事件被触发
func OnNoticeGroupDecrease(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeGroupIncrease 收到群成员增加事件被触发
func OnNoticeGroupIncrease(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeGroupBan 收到群禁言事件被触发
func OnNoticeGroupBan(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeFriendAdd 收到好友增加事件被触发
func OnNoticeFriendAdd(ctx *core.Context) {
	OnEvent(ctx)
}

// OnNoticeMessageRecall 收到好友减少事件被触发
func OnNoticeMessageRecall(ctx *core.Context) {
	OnEvent(ctx)
}

// OnRequestFriendAdd 收到好友添加请求事件被触发
func OnRequestFriendAdd(ctx *core.Context) {
	OnEvent(ctx)
}

// OnRequestGroupAdd 收到群聊加入申请事件被触发
func OnRequestGroupAdd(ctx *core.Context) {
	OnEvent(ctx)
}
