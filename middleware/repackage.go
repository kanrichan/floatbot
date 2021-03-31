package middleware

import (
	"encoding/json"
	"fmt"
	"strings"

	core "onebot/core/xianqu"
)

// PackWSRequest 封装 Websocket 数据 返回 ctx
func PackWSRequest(bot int64, data []byte) (ctx *core.Context) {
	request := map[string]interface{}{}
	json.Unmarshal(data, &request)
	if !strings.Contains(request["action"].(string), ".handle_quick_operation") {
		return &core.Context{
			Bot:     bot,
			Request: request,
		}
	}
	send, _ := json.Marshal(request["params"].(map[string]interface{})["context"])
	recv, _ := json.Marshal(request["params"].(map[string]interface{})["operation"])
	ctx = PackPOSTRequest(bot, send, recv)
	ctx.Request["echo"] = request["echo"]
	return ctx
}

// PackHTTPRequest 封装 HTTP 数据 返回 ctx
func PackHTTPRequest(bot int64, action string, data []byte) (ctx *core.Context) {
	params := map[string]interface{}{}
	json.Unmarshal(data, &params)
	return &core.Context{
		Bot: bot,
		Request: map[string]interface{}{
			"action": action,
			"params": params,
		},
	}
}

// PackWSRequest 封装 POST 数据 返回 ctx
func PackPOSTRequest(bot int64, send, data []byte) (ctx *core.Context) {
	var (
		rsp = map[string]interface{}{}
		req = map[string]interface{}{}
	)
	json.Unmarshal(send, &rsp)
	json.Unmarshal(data, &req)
	ctx = &core.Context{}
	ctx.Bot = bot
	ctx.Response = rsp
	switch rsp["post_type"].(string) {
	case "message":
		switch {
		case core.Parse(req).Exist("reply"):
			var text string
			if core.Parse(req).Bool("at_sender") {
				text = fmt.Sprintf("[CQ:at,qq=%d]", core.Parse(rsp).Int("user_id"))
			}
			text += core.Parse(req).Str("reply")
			ctx.Request = map[string]interface{}{
				"action": "send_msg",
				"params": map[string]interface{}{
					"message_type": core.Parse(rsp).Str("message_type"),
					"group_id":     core.Parse(rsp).Int("group_id"),
					"user_id":      core.Parse(rsp).Int("user_id"),
					"message":      text,
				},
				"echo": 0,
			}
		case core.Parse(req).Bool("delete"):
			ctx.Request = map[string]interface{}{
				"action": "delete_msg",
				"params": map[string]interface{}{
					"message_id": rsp["message_id"].(int64),
				},
				"echo": 0,
			}
		case core.Parse(req).Bool("kick"):
			ctx.Request = map[string]interface{}{
				"action": "set_group_kick",
				"params": map[string]interface{}{
					"group_id":           core.Parse(rsp).Int("group_id"),
					"user_id":            core.Parse(rsp).Int("user_id"),
					"reject_add_request": false,
				},
				"echo": 0,
			}
		case core.Parse(req).Bool("ban"):
			ctx.Request = map[string]interface{}{
				"action": "set_group_ban",
				"params": map[string]interface{}{
					"group_id": core.Parse(rsp).Int("group_id"),
					"user_id":  core.Parse(rsp).Int("user_id"),
					"duration": core.Parse(req).Int("duration"),
				},
				"echo": 0,
			}
		}
	case "request":
		if core.Parse(rsp).Exist("approve") {
			switch {
			case rsp["request_type"].(string) == "friend":
				ctx.Request = map[string]interface{}{
					"action": "set_friend_add_request",
					"params": map[string]interface{}{
						"flag":    core.Parse(rsp).Bool("flag"),
						"approve": core.Parse(rsp).Bool("approve"),
						"remark":  core.Parse(req).Str("remark"),
					},
					"echo": 0,
				}
			case rsp["request_type"].(string) == "group":
				ctx.Request = map[string]interface{}{
					"action": "set_group_add_request",
					"params": map[string]interface{}{
						"flag":    rsp["flag"].(string),
						"approve": core.Parse(rsp).Bool("approve"),
						"reason":  core.Parse(req).Str("reason"),
					},
					"echo": 0,
				}
			}
		}
	}
	return ctx
}

// UnPackResponse 返回 Response 的字节数组
func UnPackResponse(ctx *core.Context) []byte {
	rsp, _ := json.Marshal(ctx.Response)
	return rsp
}
