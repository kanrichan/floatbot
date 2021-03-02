package middleware

import (
	"fmt"
	core "onebot/core/xianqu"
)

// 将快速回复的报文的Request转换为标准OneBot报文
func RequestFastReplyFormat(ctx *core.Context) {
	switch ctx.Response["post_type"].(string) {
	case "message":
		switch {
		case core.Parse(ctx.Request).Exist("reply"):
			var text string
			if core.Parse(ctx.Request).Bool("at_sender") {
				text = fmt.Sprintf("[CQ:at,qq=%d]", core.Parse(ctx.Response).Int("user_id"))
			}
			text += core.Parse(ctx.Request).Str("reply")
			ctx.Request = map[string]interface{}{
				"action": "send_msg",
				"params": map[string]interface{}{
					"message_type": core.Parse(ctx.Response).Str("message_type"),
					"group_id":     core.Parse(ctx.Response).Int("group_id"),
					"user_id":      core.Parse(ctx.Response).Int("user_id"),
					"message":      text,
				},
				"echo": 0,
			}
			return
		case core.Parse(ctx.Request).Bool("delete"):
			ctx.Request = map[string]interface{}{
				"action": "delete_msg",
				"params": map[string]interface{}{
					"message_id": ctx.Response["message_id"].(int64),
				},
				"echo": 0,
			}
			return
		case core.Parse(ctx.Request).Bool("kick"):
			ctx.Request = map[string]interface{}{
				"action": "set_group_kick",
				"params": map[string]interface{}{
					"group_id":           core.Parse(ctx.Response).Int("group_id"),
					"user_id":            core.Parse(ctx.Response).Int("user_id"),
					"reject_add_request": false,
				},
				"echo": 0,
			}
			return
		case core.Parse(ctx.Request).Bool("ban"):
			ctx.Request = map[string]interface{}{
				"action": "set_group_ban",
				"params": map[string]interface{}{
					"group_id": core.Parse(ctx.Response).Int("group_id"),
					"user_id":  core.Parse(ctx.Response).Int("user_id"),
					"duration": core.Parse(ctx.Request).Int("duration"),
				},
				"echo": 0,
			}
			return
		}
	case "request":
		if core.Parse(ctx.Response).Exist("approve") {
			switch {
			case ctx.Response["request_type"].(string) == "friend":
				ctx.Request = map[string]interface{}{
					"action": "set_friend_add_request",
					"params": map[string]interface{}{
						"flag":    core.Parse(ctx.Response).Bool("flag"),
						"approve": core.Parse(ctx.Response).Bool("approve"),
						"remark":  core.Parse(ctx.Request).Str("remark"),
					},
					"echo": 0,
				}
				return
			case ctx.Response["request_type"].(string) == "group":
				ctx.Request = map[string]interface{}{
					"action": "set_group_add_request",
					"params": map[string]interface{}{
						"flag":    ctx.Response["flag"].(string),
						"approve": core.Parse(ctx.Response).Bool("approve"),
						"reason":  core.Parse(ctx.Request).Str("reason"),
					},
					"echo": 0,
				}
				return
			}
		}
	}
}
