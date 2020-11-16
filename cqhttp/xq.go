package cqhttp

import (
	"github.com/Yiwen-Chan/OneBot-YaYa/core"
)

var AppInfoJson string

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
	rawMessage  string
	time        int64
	ret         int64
}

func XQCreate(version string) string {
	return AppInfoJson
}

func XQEvent(selfID int64, mseeageType int64, subType int64, groupID int64, userID int64, noticID int64, message string, messageNum int64, messageID int64, rawMessage string, time int64, ret int64) int64 {
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
	}
	switch mseeageType {
	case 12001:
		go ProtectRun(func() { OnStart() }, "OnStart()")
	// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
	case 0, 1, 4, 5, 7:
		go ProtectRun(func() { OnPrivateMessage(xe) }, "OnPrivateMessage()")
	// 2：群聊信息
	case 2, 3:
		go ProtectRun(func() { OnGroupMessage(xe) }, "OnGroupMessage()")
	// notice 信息撤回
	case 9:
		//
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
