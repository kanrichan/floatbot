package cqhttp

import (
	"encoding/json"
)

type Event map[string]interface{}

func WSCPush(bot int64, e Event) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[推送服务] Bot %v 服务发生错误 %v，将忽略本次推送......", bot, err)
		}
	}()

	if CQHttpOK {
		send, _ := json.Marshal(e)
		for _, c := range WSCs {
			if bot == c.Bot {
				c.Send <- []byte(string(send))
			}
		}
	}
}

func OnPrivateMessage(xe XEvent) {
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
		"message_id":   xe.messageID,
		"user_id":      xe.userID,
		"message":      CQ(xe.message),
		"raw_message":  CQ(xe.message),
		"font":         0,
		"sender": Event{
			"user_id":  xe.userID,
			"nickname": "unkown",
			"sex":      "unkown",
			"age":      "unkown",
		},
	}
	WSCPush(xe.selfID, e)
}

func OnGroupMessage(xe XEvent) {
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
		"message_id":   xe.messageID,
		"group_id":     xe.groupID,
		"user_id":      xe.userID,
		"anonymous": Event{
			"id":   0,
			"name": "none",
			"flag": "none"},
		"message":     CQ(xe.message),
		"raw_message": CQ(xe.message),
		"font":        0,
		"sender": Event{
			"user_id":  xe.userID,
			"nickname": "unkown",
			"sex":      "unkown",
			"age":      "unkown",
		},
	}
	WSCPush(xe.selfID, e)
}

func OnEnable(xe XEvent) {
	e := Event{
		"time":            xe.time,
		"self_id":         xe.selfID,
		"post_type":       "meta_event",
		"meta_event_type": "lifecycle",
		"sub_type":        "connect",
	}
	WSCPush(xe.selfID, e)
}
