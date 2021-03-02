package gateway

import (
	"encoding/json"
	core "onebot/core/xianqu"
	middle "onebot/middleware"
	"onebot/server"
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
	server.HttpHandler = HttpHandler
	server.HttpPostHandler = HttpPostHandler
	server.WSCHandler = WebSocketClientHandler
	server.WSSHandler = WebSocketServerHandler
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
func HttpPostHandler(bot int64, rsp, data []byte) {
	temp := map[string]interface{}{}
	json.Unmarshal(data, &temp)
	ctx := &core.Context{
		Bot:     bot,
		Request: temp,
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

func OnMessagePrivate(ctx *core.Context) {
	broadcast(ctx)
}
func OnMessageGroup(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeFileUpload(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeAdminChange(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeGroupDecrease(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeGroupIncrease(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeGroupBan(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeFriendAdd(ctx *core.Context) {
	broadcast(ctx)
}
func OnNoticeMessageRecall(ctx *core.Context) {
	broadcast(ctx)
}
func OnRequestFriendAdd(ctx *core.Context) {
	broadcast(ctx)
}
func OnRequestGroupAdd(ctx *core.Context) {
	broadcast(ctx)
}
