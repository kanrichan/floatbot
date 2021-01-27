package onebot

import (
	"encoding/json"
	"fmt"
	"time"

	"yaya/core"
)

// Event 封装好的onebot事件
type Event map[string]interface{}

// onEvent onebot标准的事件触发会调用此函数
func onEvent(xe *XEvent) int64 {
	switch xe.MessageType {
	// 消息事件
	// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
	case 0, 1, 4, 5, 7:
		go PicPool.addPicPool(xe.Message)
		go Conf.getBotConfig(xe.SelfID).dbInsert(xe)

		go ProtectRun(func() { onPrivateMessage(xe) }, "onPrivateMessage()")
	// 2：群聊信息
	case 2, 3:
		go PicPool.addPicPool(xe.Message)
		go Conf.getBotConfig(xe.SelfID).dbInsert(xe)
		go ProtectRun(func() { onGroupMessage(xe) }, "onGroupMessage()")
	// 10：回音信息
	case 10:
		go PicPool.addPicPool(xe.Message)
		go Conf.getBotConfig(xe.SelfID).dbInsert(xe)
	// 通知事件
	// 群文件接收
	case 218:
		go ProtectRun(func() { noticeFileUpload(xe) }, "noticeFileUpload()")
	// 管理员变动 210为有人升为管理 211为有人被取消管理
	case 210:
		go ProtectRun(func() { noticeAdminChange(xe, "set") }, "noticeAdminChange()")
	case 211:
		go ProtectRun(func() { noticeAdminChange(xe, "unset") }, "noticeAdminChange()")
	// 群成员减少 201为主动退群 202为被踢
	case 201:
		go ProtectRun(func() { noticeGroupMemberDecrease(xe, "leave") }, "OnGroupMemberDecrease()")
	case 202:
		go ProtectRun(func() { noticeGroupMemberDecrease(xe, "kick") }, "noticeGroupMemberDecrease()")
	// 群成员增加
	case 212:
		go ProtectRun(func() { noticeGroupMemberIncrease(xe, "approve") }, "noticeGroupMemberIncrease()")
	// 群禁言 203为禁言 204为解禁
	case 203:
		go ProtectRun(func() { noticeGroupBan(xe, "ban") }, "noticeGroupBan()")
	case 204:
		go ProtectRun(func() { noticeGroupBan(xe, "lift_ban") }, "noticeGroupBan()")
	// new
	// 好友添加 100 为单向 102 为标准
	case 100, 102:
		go ProtectRun(func() { noticeFriendAdd(xe) }, "noticeFriendAdd()")
	// 群消息撤回 subType 2
	// 好友消息撤回 subType 1
	case 9:
		for i := range Conf.BotConfs {
			if xe.SelfID == Conf.BotConfs[i].Bot && xe.SelfID != 0 {
				if Conf.BotConfs[i].DB != nil {
					Conf.BotConfs[i].dbSelect(&xe, "message_num="+core.Int2Str(xe.MessageNum))
				}
			}
		}
		if xe.SubType == 2 {
			go ProtectRun(func() { noticeGroupMsgDelete(xe) }, "noticeGroupMsgDelete()")
		} else {
			go ProtectRun(func() { noticeFriendMsgDelete(xe) }, "noticeFriendMsgDelete()")
		}
	// 群内戳一戳

	// 群红包运气王

	// 群成员荣誉变更

	// 请求事件
	// 加好友请求
	case 101:
		go ProtectRun(func() { requestFriendAdd(xe) }, "requestFriendAdd()")
	// 加群请求／邀请 213为请求 214为被邀
	case 213:
		go ProtectRun(func() { requestGroupAdd(xe, "add") }, "requestGroupAdd()")
	case 214:
		go ProtectRun(func() { requestGroupAdd(xe, "invite") }, "requestGroupAdd()")
	default:
		//
	}
	return 0
}

// Push 将封装好的Event上报到 http|wss|wsc
func Push(bot int64, e Event, c *Yaml) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[推送][%v] BOT =X=> =X=> OneBot Error: %v", bot, err)
		}
	}()
	// TODO 获取当前Bot的Conf
	thisBotConf := Conf.getBotConfig(bot)
	// TODO 推送到WSS
	for j := range thisBotConf.WSSConf {
		if thisBotConf.WSSConf[j].Status == 1 && thisBotConf.WSSConf[j].Enable == true && thisBotConf.WSSConf[j].Host != "" {
			ce := e
			if thisBotConf.WSSConf[j].PostMessageFormat == "array" {
				ce["message"] = cqCode2Array(e["message"].(string))
			}
			send, _ := json.Marshal(ce)
			thisBotConf.WSSConf[j].Event <- send
		}
	}
	// TODO 推送到WSC
	for k := range thisBotConf.WSCConf {
		if thisBotConf.WSCConf[k].Status == 1 && thisBotConf.WSCConf[k].Enable == true && thisBotConf.WSCConf[k].Url != "" {
			ce := e
			if thisBotConf.WSCConf[k].PostMessageFormat == "array" {
				ce["message"] = cqCode2Array(e["message"].(string))
			}
			send, _ := json.Marshal(ce)
			thisBotConf.WSCConf[k].Event <- send
		}
	}
	// TODO 推送到HTTP
	for l := range thisBotConf.HTTPConf {
		ce := e
		if thisBotConf.HTTPConf[l].Status == 1 && thisBotConf.HTTPConf[l].Enable == true && thisBotConf.HTTPConf[l].Host != "" {
			if thisBotConf.HTTPConf[l].PostMessageFormat == "array" {
				ce["message"] = cqCode2Array(e["message"].(string))
			}
			send, _ := json.Marshal(ce)
			thisBotConf.HTTPConf[l].Event <- send
		}
	}
}

// 私聊信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#私聊消息
func onPrivateMessage(xe *XEvent) {
	var type_ string
	switch xe.MessageType {
	default:
		return
	case 1:
		type_ = "friend"
	case 4:
		type_ = "group"
	}
	e := Event{
		"time":         xe.Time,
		"self_id":      xe.SelfID,
		"post_type":    "message",
		"message_type": "private",
		"sub_type":     type_,
		"message_id":   xe.ID,
		"user_id":      xe.UserID,
		"message":      xq2cqCode(xe.Message),
		"raw_message":  xq2cqCode(xe.Message),
		"font":         0,
		"sender": Event{
			"user_id":  xe.UserID,
			"nickname": "unknown",
			"sex":      "unknown",
			"age":      "unknown",
		},
	}
	Push(xe.SelfID, e, Conf)
}

// 群信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群消息
func onGroupMessage(xe *XEvent) {
	var type_ string
	switch xe.MessageType {
	default:
		return
	case 2:
		type_ = "group"
	case 3:
		type_ = "discuss"
	}
	e := Event{
		"time":         xe.Time,
		"self_id":      xe.SelfID,
		"post_type":    "message",
		"message_type": type_,
		"sub_type":     "normal",
		"message_id":   xe.ID,
		"group_id":     xe.GroupID,
		"user_id":      xe.UserID,
		"anonymous":    nil,
		"message":      xq2cqCode(xe.Message),
		"raw_message":  xq2cqCode(xe.Message),
		"font":         0,
		"sender": Event{
			"user_id":  xe.UserID,
			"nickname": "unknown",
			"sex":      "unknown",
			"age":      0,
			"area":     "",
			"card":     "",
			"level":    "",
			"role":     "admin",
			"title":    "unknown",
		},
	}
	Push(xe.SelfID, e, Conf)
}

// 群文件上传
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群文件上传
func noticeFileUpload(xe *XEvent) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "group_upload",
		"group_id":    xe.GroupID,
		"user_id":     xe.UserID,
		"file": Event{
			"id":    "unknown",
			"name":  xe.Message,
			"size":  "unknown",
			"busid": "unknown",
		},
	}
	Push(xe.SelfID, e, Conf)
}

// 群管理员变动
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群管理员变动
func noticeAdminChange(xe *XEvent, typ string) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "group_admin",
		"sub_type":    typ,
		"group_id":    xe.GroupID,
		"user_id":     xe.UserID,
	}
	Push(xe.SelfID, e, Conf)
}

// 群成员减少
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群成员减少
func noticeGroupMemberDecrease(xe *XEvent, typ string) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "group_decrease",
		"sub_type":    typ,
		"group_id":    xe.GroupID,
		"operator_id": xe.UserID,
		"user_id":     xe.NoticeID,
	}
	Push(xe.SelfID, e, Conf)
}

// 群成员增加
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群成员增加
func noticeGroupMemberIncrease(xe *XEvent, typ string) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "group_increase",
		"sub_type":    "unknown",
		"group_id":    xe.GroupID,
		"operator_id": xe.UserID,
		"user_id":     xe.NoticeID,
	}
	Push(xe.SelfID, e, Conf)
}

// 群禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群禁言
func noticeGroupBan(xe *XEvent, typ string) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "group_ban",
		"sub_type":    typ,
		"group_id":    xe.GroupID,
		"operator_id": xe.UserID,
		"user_id":     xe.NoticeID,
		"duration":    "unknown",
	}
	Push(xe.SelfID, e, Conf)
}

// 好友添加
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#好友添加
func noticeFriendAdd(xe *XEvent) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "friend_add",
		"user_id":     xe.UserID,
	}
	Push(xe.SelfID, e, Conf)
}

// 群消息撤回
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群消息撤回
func noticeGroupMsgDelete(xe *XEvent) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "group_recall",
		"group_id":    xe.GroupID,
		"user_id":     xe.NoticeID,
		"operator_id": xe.UserID,
		"message_id":  xe.ID,
	}
	Push(xe.SelfID, e, Conf)
}

// 好友消息撤回
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#好友消息撤回
func noticeFriendMsgDelete(xe *XEvent) {
	e := Event{
		"time":        xe.Time,
		"self_id":     xe.SelfID,
		"post_type":   "notice",
		"notice_type": "friend_recall",
		"user_id":     xe.NoticeID,
		"message_id":  xe.ID,
	}
	Push(xe.SelfID, e, Conf)
}

// 群内戳一戳
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群内戳一戳

// 群红包运气王
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群红包运气王

// 群成员荣誉变更
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#群成员荣誉变更

// 加好友请求
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#加好友请求
func requestFriendAdd(xe *XEvent) {
	e := Event{
		"time":         xe.Time,
		"self_id":      xe.SelfID,
		"post_type":    "request",
		"request_type": "friend",
		"user_id":      xe.NoticeID,
		"comment":      xe.Message,
		"flag":         xe.UserID,
	}
	Push(xe.SelfID, e, Conf)
}

// 加群请求/邀请
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#加群请求邀请
func requestGroupAdd(xe *XEvent, typ string) {
	e := Event{
		"time":      xe.Time,
		"self_id":   xe.SelfID,
		"post_type": "request",
		"sub_type":  typ,
		"group_id":  xe.GroupID,
		"user_id":   xe.NoticeID,
		"comment":   xe.Message,
		"flag":      fmt.Sprintf("%v|%v|%v", xe.SubType, xe.GroupID, xe.RawMessage),
	}
	Push(xe.SelfID, e, Conf)
}

// 生命周期
// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/message.md#生命周期
func onEnable(xe *XEvent) {
	e := Event{
		"time":            xe.Time,
		"self_id":         xe.SelfID,
		"post_type":       "meta_event",
		"meta_event_type": "lifecycle",
		"sub_type":        "connect",
	}
	Push(xe.SelfID, e, Conf)
}

// heartBeat HeartBeat --> ALL PLUGINS
func (conf *Yaml) heartBeat() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[心跳] XQ =X=> =X=> Plugins Error: %v", err)
		}
	}()
	if conf.HeratBeatConf.Interval == 0 || !conf.HeratBeatConf.Enable {
		return
	}
	if conf.HeratBeatConf.Interval < 1000 {
		INFO("[心跳] Interval %v -> 1000", conf.HeratBeatConf.Interval)
		conf.HeratBeatConf.Interval = 1000
	}
	INFO("[心跳] XQ ==> ==> Plugins")
	for {
		time.Sleep(time.Millisecond * time.Duration(conf.HeratBeatConf.Interval))
		if conf.HeratBeatConf.Enable && conf.HeratBeatConf.Interval != 0 {
			for i := range conf.BotConfs {
				for j := range conf.BotConfs[i].WSSConf {
					if conf.BotConfs[i].WSSConf[j].Status == 1 && conf.BotConfs[i].WSSConf[j].Enable {
						conf.BotConfs[i].WSSConf[j].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for k := range conf.BotConfs[i].WSCConf {
					if conf.BotConfs[i].WSCConf[k].Status == 1 && conf.BotConfs[i].WSCConf[k].Enable {
						conf.BotConfs[i].WSCConf[k].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for l := range conf.BotConfs[i].HTTPConf {
					if conf.BotConfs[i].HTTPConf[l].Status == 1 && conf.BotConfs[i].HTTPConf[l].Enable {
						conf.BotConfs[i].HTTPConf[l].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
			}
		}
	}
}

func heartEvent(interval int64, bot int64) []byte {
	heartbeat := map[string]interface{}{
		"interval":        fmt.Sprint(interval),
		"meta_event_type": "heartbeat",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(bot),
		"status": map[string]interface{}{
			"online": true,
			"good":   true,
		},
		"time": fmt.Sprint(time.Now().Unix()),
	}
	event, _ := json.Marshal(heartbeat)
	return event
}
