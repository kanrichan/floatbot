package xianqu

import "C"
import (
	"encoding/json"
	"fmt"
	"runtime"
)

const (
	EnableE         = 10000
	PrivateMessageE = 1
)

var (
	AppInfo = newAppInfo()

	OneBotPath = PathExecute() + "OneBot\\"

	OnMessagePrivate      = func(ctx *Context) {}
	OnMessageGroup        = func(ctx *Context) {}
	OnNoticeFileUpload    = func(ctx *Context) {}
	OnNoticeAdminChange   = func(ctx *Context) {}
	OnNoticeGroupDecrease = func(ctx *Context) {}
	OnNoticeGroupIncrease = func(ctx *Context) {}
	OnNoticeGroupBan      = func(ctx *Context) {}
	OnNoticeFriendAdd     = func(ctx *Context) {}
	OnNoticeMessageRecall = func(ctx *Context) {}
	OnRequestFriendAdd    = func(ctx *Context) {}
	OnRequestGroupAdd     = func(ctx *Context) {}
	OnEnable              = func(ctx *Context) {}
	OnDisable             = func(ctx *Context) {}
	OnSetting             = func(ctx *Context) {}

	MessageIDCache        = &CacheData{Max: 1000, Key: []interface{}{}, Value: []interface{}{}}
	MessageCache          = &CacheData{Max: 1000, Key: []interface{}{}, Value: []interface{}{}}
	GroupDataCache        = &CacheGroupsData{Group: []*GroupData{}}
	PicPoolCache          = &CacheData{Max: 1000, Key: []interface{}{}, Value: []interface{}{}}
	TemporarySessionCache = &CacheData{Max: 50, Key: []interface{}{}, Value: []interface{}{}}
)

func init() {
	CreatePath(OneBotPath + "\\image\\")
	CreatePath(OneBotPath + "\\record\\")
}

// App XQ要求的插件信息
type App struct {
	Name   string `json:"name"`   // 插件名字
	Pver   string `json:"pver"`   // 插件版本
	Sver   int    `json:"sver"`   // 框架版本
	Author string `json:"author"` // 作者名字
	Desc   string `json:"desc"`   // 插件说明
}

// newAppInfo 返回插件信息
func newAppInfo() *App {
	return &App{
		Name:   "OneBot-YaYa",
		Pver:   "1.2.0",
		Sver:   3,
		Author: "kanri",
		Desc:   "OneBot标准的先驱实现 项目地址: http://github.com/Yiwen-Chan/OneBot-YaYa",
	}
}

type Context struct {
	Bot      int64
	Request  map[string]interface{}
	Response map[string]interface{}
}

//export GoCreate
func GoCreate(version *C.char) *C.char {
	data, _ := json.Marshal(AppInfo)
	return CString(string(data))
}

//export GoSetUp
func GoSetUp() C.int {
	OnSetting(nil)
	return C.int(0)
}

//export GoDestroyPlugin
func GoDestroyPlugin() C.int {
	return C.int(0)
}

//export GoEvent
func GoEvent(cBot *C.char, cMessageType, cSubType C.int, cGroupID, cUserID, cNoticeID, cMessage, cMessageNum, cMessageID, cRawMessage, cTime *C.char, cRet C.int) C.int {
	var (
		bot         = CStr2GoInt(cBot)
		messageType = int64(cMessageType)
		subType     = int64(cSubType)
		groupID     = CStr2GoInt(cGroupID)
		userID      = CStr2GoInt(cUserID)
		noticeID    = CStr2GoInt(cNoticeID)
		message     = UnescapeEmoji(GoString(cMessage)) // 解决易语言的emoji到utf-8
		messageNum  = CStr2GoInt(cMessageNum)
		messageID   = CStr2GoInt(cMessageID)
		rawMessage  = GoString(cRawMessage)
		time        = CStr2GoInt(cTime)
		// ret         = CStr2GoInt(cRet)
	)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				ApiOutPutLog(fmt.Sprintf("[PANIC] 发生了不可预知的错误，请在GitHub提交issue：%v", err))
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				ApiOutPutLog(fmt.Sprintf("[TRACEBACK]:\n%v", string(buf)))
			}
		}()
		switch messageType {
		// 消息事件
		// 0：临时会话 1：好友会话 4：群临时会话 7：好友验证会话
		case 1:
			cqMessage := xq2cqCode(message)
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "message",
					"message_type": "private",
					"sub_type":     "friend",
					"message_id":   messageID,
					"user_id":      userID,
					"message":      cqMessage,
					"raw_message":  cqMessage,
					"font":         0,
					"sender": map[string]interface{}{
						"user_id":  userID,
						"nickname": XQApiGetNick(bot, userID),
						"sex":      XQApiGetGender(bot, userID),
						"age":      XQApiGetAge(bot, userID),
					},
				},
			}
			MessageIDCache.Insert(messageID, messageNum)
			MessageCache.Insert(messageID, ctx.Response)
			OnMessagePrivate(ctx)
		// 2：群聊信息
		case 2:
			cqMessage := xq2cqCode(message)
			info := GroupDataCache.GetCacheGroupMember(
				bot,
				groupID,
				userID,
				true,
			)
			sender := map[string]interface{}{
				"user_id":  userID,
				"nickname": "unknown",
				"sex":      "unknown",
				"age":      0,
				"area":     "",
				"card":     "",
				"level":    "",
				"role":     "unknown",
				"title":    "unknown",
			}
			if info != nil {
				sender = map[string]interface{}{
					"user_id":  userID,
					"nickname": info.Nickname,
					"sex":      info.Sex,
					"age":      info.Age,
					"area":     "",
					"card":     "",
					"level":    "",
					"role":     info.Role,
					"title":    "unknown",
				}
			}
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "message",
					"message_type": "group",
					"sub_type":     "normal",
					"message_id":   messageID,
					"group_id":     groupID,
					"user_id":      userID,
					"anonymous":    nil,
					"message":      cqMessage,
					"raw_message":  cqMessage,
					"font":         0,
					"sender":       sender,
				},
			}
			MessageIDCache.Insert(messageID, messageNum)
			MessageCache.Insert(messageID, ctx.Response)
			OnMessageGroup(ctx)
		// 4：临时群聊信息
		case 4:
			cqMessage := xq2cqCode(message)
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "message",
					"message_type": "private",
					"sub_type":     "group",
					"message_id":   messageID,
					"user_id":      userID,
					"message":      cqMessage,
					"raw_message":  cqMessage,
					"font":         0,
					"sender": map[string]interface{}{
						"user_id":  userID,
						"nickname": "unknown",
						"sex":      "unknown",
						"age":      0,
					},
				},
			}
			TemporarySessionCache.Insert(userID, groupID)
			MessageIDCache.Insert(messageID, messageNum)
			MessageCache.Insert(messageID, ctx.Response)
			OnMessageGroup(ctx)
		// 10：回音信息
		case 10:
			cqMessage := xq2cqCode(message)
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "message",
					"message_type": "group",
					"sub_type":     "normal",
					"message_id":   messageID,
					"group_id":     groupID,
					"user_id":      userID,
					"anonymous":    nil,
					"message":      cqMessage,
					"raw_message":  cqMessage,
					"font":         0,
					"sender": map[string]interface{}{
						"user_id":  userID,
						"nickname": "unknown",
						"sex":      "unknown",
						"age":      0,
						"area":     "",
						"card":     "",
						"level":    "",
						"role":     "admin",
						"title":    "unknown",
					},
				},
			}
			MessageIDCache.Insert(messageID, messageNum)
			MessageCache.Insert(messageID, ctx.Response)
		// 通知事件
		// 群文件接收
		case 218:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_upload",
					"group_id":    groupID,
					"user_id":     userID,
					"file": map[string]interface{}{
						"id":    "unknown",
						"name":  message,
						"size":  "unknown",
						"busid": "unknown",
					},
				},
			}
			OnNoticeFileUpload(ctx)
		// 管理员变动 210为有人升为管理 211为有人被取消管理
		case 210:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_admin",
					"sub_type":    "set",
					"group_id":    groupID,
					"user_id":     userID,
				},
			}
			OnNoticeAdminChange(ctx)
		case 211:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_admin",
					"sub_type":    "unset",
					"group_id":    groupID,
					"user_id":     userID,
				},
			}
			OnNoticeAdminChange(ctx)
		// 群成员减少 201为某人退出群 202为某人被管理移除群
		case 201:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_decrease",
					"sub_type":    "leave",
					"group_id":    groupID,
					"operator_id": userID,
					"user_id":     noticeID,
				},
			}
			OnNoticeGroupDecrease(ctx)
		case 202:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_decrease",
					"sub_type":    "kick",
					"group_id":    groupID,
					"operator_id": userID,
					"user_id":     noticeID,
				},
			}
			OnNoticeGroupDecrease(ctx)
		// 群成员增加
		case 212:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_increase",
					"sub_type":    "unknown",
					"group_id":    groupID,
					"operator_id": userID,
					"user_id":     noticeID,
				},
			}
			OnNoticeGroupIncrease(ctx)
		// 群禁言 203为禁言 204为解禁
		case 203:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_ban",
					"sub_type":    "ban",
					"group_id":    groupID,
					"operator_id": userID,
					"user_id":     noticeID,
					"duration":    1,
				},
			}
			OnNoticeGroupBan(ctx)
		case 204:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "group_ban",
					"sub_type":    "lift_ban",
					"group_id":    groupID,
					"operator_id": userID,
					"user_id":     noticeID,
					"duration":    0,
				},
			}
			OnNoticeGroupBan(ctx)
		// new
		// 好友添加 100 为单向 102 为标准
		case 100, 102:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":        time,
					"self_id":     bot,
					"post_type":   "notice",
					"notice_type": "friend_add",
					"user_id":     userID,
				},
			}
			OnNoticeFriendAdd(ctx)
		// 群消息撤回 subType 2
		// 好友消息撤回 subType 1
		case 9:
			switch subType {
			case 1:
				ctx := &Context{
					Bot: bot,
					Response: map[string]interface{}{
						"time":        time,
						"self_id":     bot,
						"post_type":   "notice",
						"notice_type": "group_recall",
						"group_id":    groupID,
						"user_id":     noticeID,
						"operator_id": userID,
						"message_id":  messageID,
					},
				}
				OnNoticeMessageRecall(ctx)
			case 2:
				ctx := &Context{
					Bot: bot,
					Response: map[string]interface{}{
						"time":        time,
						"self_id":     bot,
						"post_type":   "notice",
						"notice_type": "friend_recall",
						"group_id":    groupID,
						"user_id":     noticeID,
						"operator_id": userID,
						"message_id":  messageID,
					},
				}
				OnNoticeMessageRecall(ctx)
			}

		// 群内戳一戳

		// 群红包运气王

		// 群成员荣誉变更

		// 请求事件
		// 加好友请求
		case 101:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "request",
					"request_type": "friend",
					"user_id":      noticeID,
					"comment":      message,
					"flag":         userID,
				},
			}
			OnRequestFriendAdd(ctx)
		// 加群请求／邀请 213为请求 214为被邀
		case 213:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":      time,
					"self_id":   bot,
					"post_type": "request",
					"sub_type":  "add",
					"group_id":  groupID,
					"user_id":   noticeID,
					"comment":   message,
					"flag":      fmt.Sprintf("%v|%v|%v", subType, groupID, rawMessage),
				},
			}
			OnRequestGroupAdd(ctx)
		case 214:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":      time,
					"self_id":   bot,
					"post_type": "request",
					"sub_type":  "invite",
					"group_id":  groupID,
					"user_id":   noticeID,
					"comment":   message,
					"flag":      fmt.Sprintf("%v|%v|%v", subType, groupID, rawMessage),
				},
			}
			OnRequestGroupAdd(ctx)
		case 12001:
			OnEnable(nil)
		case 12002:
			OnDisable(nil)
		default:
			//
		}
	}()
	return 0
}
