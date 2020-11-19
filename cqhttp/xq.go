package cqhttp

import (
	"yaya/core"
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
	case 12002:
		go ProtectRun(func() { OnDisable() }, "OnDisable()")
	// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
	case 0, 1, 4, 5, 7:
		go ProtectRun(func() { OnPrivateMessage(xe) }, "OnPrivateMessage()")
		go ProtectRun(func() { CommandHandle(xe) }, "CommandHandle()")
	// 2：群聊信息
	case 2, 3:
		go ProtectRun(func() { OnGroupMessage(xe) }, "OnGroupMessage()")
		go ProtectRun(func() { CommandHandle(xe) }, "CommandHandle()")
	// notice 信息撤回
	case 9:
		//
	// 群员退群
	case 201:
		go ProtectRun(func() { OnGroupMenberDecrease(xe, "leave") }, "OnGroupMenberDecrease()")
	// 群员被踢
	case 202:
		go ProtectRun(func() { OnGroupMenberDecrease(xe, "kick") }, "OnGroupMenberDecrease()")
	// 群员被禁言
	case 203:
		go ProtectRun(func() { OnGroupBan(xe, "ban") }, "OnGroupBan()")
	// 群员被解除禁言
	case 204:
		go ProtectRun(func() { OnGroupBan(xe, "lift_ban") }, "OnGroupBan()")
	// 群员被升为管理
	case 210:
		go ProtectRun(func() { OnAdminChange(xe, "set") }, "OnAdminChange()")
	// 群员被取消管理
	case 211:
		go ProtectRun(func() { OnAdminChange(xe, "unset") }, "OnAdminChange()")
	// 某人被批准加群
	case 212:
		go ProtectRun(func() { OnGroupMenberIncrease(xe, "approve") }, "OnGroupMenberIncrease()")
	// 某人被邀请加入群
	case 219:
		go ProtectRun(func() { OnGroupMenberIncrease(xe, "invite") }, "OnGroupMenberIncrease()")
	// 群文件接收
	case 218:
		go ProtectRun(func() { OnFileUpload(xe) }, "OnFileUpload()")
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
