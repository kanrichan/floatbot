package onebot

import (
	"encoding/json"
	"fmt"

	"yaya/core"
)

var AppInfoJson string

type Event map[string]interface{}

func init() {
	core.Create = XQCreate
	core.Event = XQEvent
	core.DestroyPlugin = XQDestroyPlugin
	core.SetUp = XQSetUp
}

type MessageData map[string]string
type Message struct {
	Type string      `json:"type"`
	Data MessageData `json:"data"`
}

type XEvent struct {
	selfID      int64
	mseeageType int64
	subType     int64
	groupID     int64
	userID      int64
	noticID     int64
	message     string
	messageNum  int64
	messageID   int64
	rawMessage  []byte
	time        int64
	ret         int64
	cqID        int64
}

func XQCreate(version string) string {
	return AppInfoJson
}

func XQEvent(selfID int64, mseeageType int64, subType int64, groupID int64, userID int64, noticID int64, message string, messageNum int64, messageID int64, rawMessage []byte, time int64, ret int64) int64 {
	xe := XEvent{
		selfID:      selfID,
		mseeageType: mseeageType,
		subType:     subType,
		groupID:     groupID,
		userID:      userID,
		noticID:     noticID,
		message:     message,
		messageNum:  messageNum,
		messageID:   messageID,
		rawMessage:  rawMessage,
		time:        time,
		ret:         ret,
		cqID:        0,
	}

	switch mseeageType {
	case 12001:
		go ProtectRun(func() { onStart() }, "onStart()")
	case 12002:
		go ProtectRun(func() { onDisable() }, "onDisable()")
	// 消息事件
	// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
	case 0, 1, 4, 5, 7:
		for i, _ := range Conf.BotConfs {
			if selfID == Conf.BotConfs[i].Bot && selfID != 0 {
				if Conf.BotConfs[i].DB != nil {
					xe.event2DB(Conf.BotConfs[i].DB)
				}
			}
		}
		go ProtectRun(func() { onPrivateMessage(xe) }, "onPrivateMessage()")
	// 2：群聊信息
	case 2, 3:
		for i, _ := range Conf.BotConfs {
			if selfID == Conf.BotConfs[i].Bot && selfID != 0 {
				if Conf.BotConfs[i].DB != nil {
					xe.event2DB(Conf.BotConfs[i].DB)
				}
			}
		}
		go ProtectRun(func() { onGroupMessage(xe) }, "onGroupMessage()")
	// 10：回音信息
	case 10:
		for i, _ := range Conf.BotConfs {
			if selfID == Conf.BotConfs[i].Bot && selfID != 0 {
				if Conf.BotConfs[i].DB != nil {
					xe.event2DB(Conf.BotConfs[i].DB)
				}
			}
		}
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
		go ProtectRun(func() { noticeGroupMenberDecrease(xe, "leave") }, "OnGroupMenberDecrease()")
	case 202:
		go ProtectRun(func() { noticeGroupMenberDecrease(xe, "kick") }, "noticeGroupMenberDecrease()")
	// 群成员增加
	case 212:
		go ProtectRun(func() { noticeGroupMenberIncrease(xe, "approve") }, "noticeGroupMenberIncrease()")
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
		for i, _ := range Conf.BotConfs {
			if selfID == Conf.BotConfs[i].Bot && selfID != 0 {
				if Conf.BotConfs[i].DB != nil {
					xe.xq2cqid(Conf.BotConfs[i].DB)
				}
			}
		}
		if xe.subType == 2 {
			go ProtectRun(func() { noticGroupMsgDelete(xe) }, "noticGroupMsgDelete()")
		} else {
			go ProtectRun(func() { noticFriendMsgDelete(xe) }, "noticFriendMsgDelete()")
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

func XQDestroyPlugin() int64 {
	return 0
}

func XQSetUp() int64 {
	return 0
}

func WSCPush(bot int64, e Event, c *Yaml) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[推送][%v] BOT =X=> =X=> OneBot Error: %v", bot, err)
		}
	}()

	for i, _ := range c.BotConfs {
		if bot == c.BotConfs[i].Bot {
			for j, _ := range c.BotConfs[i].WSSConf {
				if c.BotConfs[i].WSSConf[j].Status == 1 && c.BotConfs[i].WSSConf[j].Enable == true && c.BotConfs[i].WSSConf[j].Host != "" {
					ce := e
					if c.BotConfs[i].WSSConf[j].PostMessageFormat == "array" {
						ce["message"] = cqCode2Array(e["message"].(string))
					}
					send, _ := json.Marshal(ce)
					c.BotConfs[i].WSSConf[j].Event <- send
				}
			}
			for k, _ := range c.BotConfs[i].WSCConf {
				if c.BotConfs[i].WSCConf[k].Status == 1 && c.BotConfs[i].WSCConf[k].Enable == true && c.BotConfs[i].WSCConf[k].Url != "" {
					ce := e
					if c.BotConfs[i].WSCConf[k].PostMessageFormat == "array" {
						ce["message"] = cqCode2Array(e["message"].(string))
					}
					send, _ := json.Marshal(ce)
					c.BotConfs[i].WSCConf[k].Event <- send
				}
			}
			for l, _ := range c.BotConfs[i].HTTPConf {
				ce := e
				if c.BotConfs[i].HTTPConf[l].Status == 1 && c.BotConfs[i].HTTPConf[l].Enable == true && c.BotConfs[i].HTTPConf[l].Host != "" {
					if c.BotConfs[i].HTTPConf[l].PostMessageFormat == "array" {
						ce["message"] = cqCode2Array(e["message"].(string))
					}
					send, _ := json.Marshal(ce)
					c.BotConfs[i].HTTPConf[l].Event <- send
				}
			}
		}
	}

}

func xq2cqMsgID(xqid int64, xqnum int64) int64 {
	return core.Str2Int(fmt.Sprintf("%01d%02d%06d%010d", len(core.Int2Str(xqid)), len(core.Int2Str(xqnum)), xqid, xqnum))
}

func cq2xqMsgID(cqid int64) (int64, int64) {
	idLen := core.Str2Int(core.Int2Str(cqid)[0:1])
	numLen := core.Str2Int(core.Int2Str(cqid)[1:3])
	return core.Str2Int(core.Int2Str(cqid)[(9 - idLen):9]),
		core.Str2Int(core.Int2Str(cqid)[(19 - numLen):19])
}

func onPrivateMessage(xe XEvent) {
	Tsubtype := "error"
	switch xe.mseeageType {
	case 0:
		Tsubtype = "other"
	case 1:
		Tsubtype = "friend"
	case 4:
		Tsubtype = "group"
	case 5:
		Tsubtype = "discuss"
	case 7:
		Tsubtype = "other"
	default:
		Tsubtype = "error"
	}
	e := Event{
		"time":         xe.time,
		"self_id":      xe.selfID,
		"post_type":    "message",
		"message_type": "private",
		"sub_type":     Tsubtype,
		"message_id":   xe.cqID,
		"user_id":      xe.userID,
		"message":      xq2cqCode(xe.message),
		"raw_message":  xq2cqCode(xe.message),
		"font":         0,
		"sender": Event{
			"user_id":  xe.userID,
			"nickname": "unknown",
			"sex":      "unknown",
			"age":      "unknown",
		},
	}
	WSCPush(xe.selfID, e, Conf)
}

func onGroupMessage(xe XEvent) {
	Tmessagetype := "error"
	switch xe.mseeageType {
	case 2:
		Tmessagetype = "group"
	case 3:
		Tmessagetype = "discuss"
	default:
		Tmessagetype = "error"
	}
	e := Event{
		"time":         xe.time,
		"self_id":      xe.selfID,
		"post_type":    "message",
		"message_type": Tmessagetype,
		"sub_type":     "normal",
		"message_id":   xe.cqID,
		"group_id":     xe.groupID,
		"user_id":      xe.userID,
		"anonymous":    nil,
		"message":      xq2cqCode(xe.message),
		"raw_message":  xq2cqCode(xe.message),
		"font":         0,
		"sender": Event{
			"user_id":  xe.userID,
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
	WSCPush(xe.selfID, e, Conf)
}

func onEnable(xe XEvent) {
	e := Event{
		"time":            xe.time,
		"self_id":         xe.selfID,
		"post_type":       "meta_event",
		"meta_event_type": "lifecycle",
		"sub_type":        "connect",
	}
	WSCPush(xe.selfID, e, Conf)
}

// 群文件上传
func noticeFileUpload(xe XEvent) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "group_upload",
		"group_id":    xe.groupID,
		"user_id":     xe.userID,
		"file": Event{
			"id":    "unknow",
			"name":  xe.message,
			"size":  "unknow",
			"busid": "unknow",
		},
	}
	WSCPush(xe.selfID, e, Conf)
}

// 管理员变动
func noticeAdminChange(xe XEvent, typ string) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "group_admin",
		"sub_type":    typ,
		"group_id":    xe.groupID,
		"user_id":     xe.userID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 群成员减少
func noticeGroupMenberDecrease(xe XEvent, typ string) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "group_decrease",
		"sub_type":    typ,
		"group_id":    xe.groupID,
		"operator_id": xe.userID,
		"user_id":     xe.noticID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 群成员增加
func noticeGroupMenberIncrease(xe XEvent, typ string) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "group_increase",
		"sub_type":    "unknow",
		"group_id":    xe.groupID,
		"operator_id": xe.userID,
		"user_id":     xe.noticID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 群禁言
func noticeGroupBan(xe XEvent, typ string) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "group_ban",
		"sub_type":    typ,
		"group_id":    xe.groupID,
		"operator_id": xe.userID,
		"user_id":     xe.noticID,
		"duration":    "unknow",
	}
	WSCPush(xe.selfID, e, Conf)
}

// 好友添加
func noticeFriendAdd(xe XEvent) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "friend_add",
		"user_id":     xe.userID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 群消息撤回
func noticGroupMsgDelete(xe XEvent) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "group_recall",
		"group_id":    xe.groupID,
		"user_id":     xe.noticID,
		"operator_id": xe.userID,
		"message_id":  xe.cqID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 好友消息撤回
func noticFriendMsgDelete(xe XEvent) {
	e := Event{
		"time":        xe.time,
		"self_id":     xe.selfID,
		"post_type":   "notice",
		"notice_type": "friend_recall",
		"user_id":     xe.noticID,
		"message_id":  xe.cqID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 群内戳一戳

// 群红包运气王

// 群成员荣誉变更

// 加好友请求
func requestFriendAdd(xe XEvent) {
	e := Event{
		"time":         xe.time,
		"self_id":      xe.selfID,
		"post_type":    "request",
		"request_type": "friend",
		"user_id":      xe.noticID,
		"comment":      xe.message,
		"flag":         xe.userID,
	}
	WSCPush(xe.selfID, e, Conf)
}

// 加群请求
func requestGroupAdd(xe XEvent, typ string) {
	e := Event{
		"time":      xe.time,
		"self_id":   xe.selfID,
		"post_type": "request",
		"sub_type":  typ,
		"group_id":  xe.groupID,
		"user_id":   xe.noticID,
		"comment":   xe.message,
		"flag":      fmt.Sprintf("%v|%v|%v", xe.subType, xe.groupID, xe.rawMessage),
	}
	WSCPush(xe.selfID, e, Conf)
}
