package core

import "C"

var Create func(version string) string
var Event func(selfID int64, messageType int64, subType int64, groupID int64, userID int64, noticeID int64, message string, messageNum int64, messageID int64, rawMessage string, time int64, ret int64) int64
var DestroyPlugin func() int64
var SetUp func() int64

//export GO_Create
func GO_Create(version *C.char) *C.char {
	return CString(Create(GoString(version)))
}

// SelfID 机器人QQ, 多Q版用于判定哪个QQ接收到该消息
// MessageType 消息类型, 接收到消息类型，该类型可在常量表中查询具体定义，此处仅列举： -1 未定义事件 0,在线状态临时会话 1,好友信息 2,群信息 3,讨论组信息 4,群临时会话 5,讨论组临时会话 6,财付通转账 7,好友验证回复会话
// SubType 消息子类型, 此参数在不同消息类型下，有不同的定义，暂定：接收财付通转账时 1为好友 4为群临时会话 5为讨论组临时会话    有人请求入群时，不良成员这里为1
// GroupID 消息来源, 此消息的来源，如：群号、讨论组ID、临时会话QQ、好友QQ等
// UserID 触发对象_主动, 主动发送这条消息的QQ，踢人时为踢人管理员QQ
// NoticeID 触发对象_被动, 被动触发的QQ，如某人被踢出群，则此参数为被踢出人QQ
// Message 消息内容, 此参数有多重含义，常见为：对方发送的消息内容，但当消息类型为 某人申请入群，则为入群申请理由
// MessageNum 消息序号, 此参数暂定用于消息回复，消息撤回
// MessageID 消息ID, 此参数暂定用于消息回复，消息撤回
// RawMessage 原始信息, UDP收到的原始信息，特殊情况下会返回JSON结构（入群事件时，这里为该事件seq）
// Time 消息时间戳, 接受到消息的时间戳
// Ret 回传文本指针, 此参数用于插件加载拒绝理由
//export GO_Event
func GO_Event(selfID *C.char, messageType C.int, subType C.int, groupID *C.char, userID *C.char, noticeID *C.char, message *C.char, messageNum *C.char, messageID *C.char, rawMessage *C.char, time *C.char, ret *C.char) C.int {
	return C.int(Event(CStr2GoInt(selfID),
		int64(messageType),
		int64(subType),
		CStr2GoInt(groupID),
		CStr2GoInt(userID),
		CStr2GoInt(noticeID),
		// TODO 解决易语言的emoji到utf-8
		UnescapeEmoji(GoString(message)),
		CStr2GoInt(messageNum),
		CStr2GoInt(messageID),
		GoString(rawMessage),
		CStr2GoInt(time),
		CStr2GoInt(ret),
	))
}

//export GO_DestroyPlugin
func GO_DestroyPlugin() C.int {
	return C.int(DestroyPlugin())
}

//export GO_SetUp
func GO_SetUp() C.int {
	return C.int(SetUp())
}

func main() {
	//
}
