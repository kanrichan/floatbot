package xianqu

import "C"

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	// 先驱插件信息
	AppInfo = &App{
		Name:   "OneBot-YaYa",
		Pver:   "1.2.10",
		Sver:   3,
		Author: "kanri",
		Desc:   "OneBot标准的先驱实现 项目地址: http://github.com/Yiwen-Chan/OneBot-YaYa",
	}
	// 当前OneBot目录
	OneBotPath = pathExecute() + "OneBot\\"

	OnMessagePrivate      func(ctx *Context)
	OnMessageGroup        func(ctx *Context)
	OnNoticeFileUpload    func(ctx *Context)
	OnNoticeAdminChange   func(ctx *Context)
	OnNoticeGroupDecrease func(ctx *Context)
	OnNoticeGroupIncrease func(ctx *Context)
	OnNoticeGroupBan      func(ctx *Context)
	OnNoticeFriendAdd     func(ctx *Context)
	OnNoticeMessageRecall func(ctx *Context)
	OnRequestFriendAdd    func(ctx *Context)
	OnRequestGroupAdd     func(ctx *Context)
	OnEnable              func(ctx *Context)
	OnDisable             func(ctx *Context)
	OnSetting             func(ctx *Context)

	// 信息 id 与 num 对应缓冲池
	MessageIDCache = &CacheData{Max: 1000, Key: []interface{}{}, Value: []interface{}{}}
	// 信息 id 与 数据 对应缓冲池
	MessageCache = &CacheData{Max: 1000, Key: []interface{}{}, Value: []interface{}{}}
	// 群数据缓冲池
	GroupDataCache = &CacheGroupsData{Group: []*GroupData{}}
	// 图片链接 或 md5 对应 xq 资源缓冲池
	PicPoolCache = &CacheData{Max: 1000, Key: []interface{}{}, Value: []interface{}{}}
	// 群临时会话缓冲此
	TemporarySessionCache = &CacheData{Max: 50, Key: []interface{}{}, Value: []interface{}{}}
)

func init() {
	// 创建数据目录
	createPath(OneBotPath + "image\\")
	createPath(OneBotPath + "record\\")
}

// App XQ要求的插件信息
type App struct {
	Name   string `json:"name"`   // 插件名字
	Pver   string `json:"pver"`   // 插件版本
	Sver   int    `json:"sver"`   // 框架版本
	Author string `json:"author"` // 作者名字
	Desc   string `json:"desc"`   // 插件说明
}

// 上下报文，包括了上报数据以及api调用数据
type Context struct {
	Bot      int64
	Request  map[string]interface{}
	Response map[string]interface{}
}

//export GoCreate
func GoCreate(version *C.char) *C.char {
	data, _ := json.Marshal(AppInfo)
	return cString(helper.BytesToString(data))
}

//export GoSetUp
func GoSetUp() C.int {
	OnSetting(nil)
	return 0
}

//export GoDestroyPlugin
func GoDestroyPlugin() C.int {
	runtime.GC()
	return 0
}

//export GoEvent
func GoEvent(cBot *C.char, cMessageType, cSubType C.int, cGroupID, cUserID, cNoticeID, cMessage, cMessageNum, cMessageID, cRawMessage, cTime *C.char, cRet C.int) C.int {
	var (
		bot         = cStr2GoInt(cBot)
		messageType = int64(cMessageType)
		subType     = int64(cSubType)
		groupID     = cStr2GoInt(cGroupID)
		userID      = cStr2GoInt(cUserID)
		noticeID    = cStr2GoInt(cNoticeID)
		message     = unescapeEmoji(goString(cMessage)) // 解决易语言的emoji到utf-8
		messageNum  = cStr2GoInt(cMessageNum)
		messageID   = cStr2GoInt(cMessageID)
		rawMessage  = goString(cRawMessage)
		time        = cStr2GoInt(cTime)
		// ret         = CStr2GoInt(cRet)
	)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				ApiOutPutLog("发生不可预知错误，请[右键↓错误信息↓]并[点击查看完整消息]，截图提交到 GitHub issue 或者到 QQ群 1048452984")
				ApiOutPutLog(fmt.Sprintf("[PANIC] [错误]：%v \n[TRACEBACK]:\n%v", err, helper.BytesToString(buf)))
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
						"nickname": info.GetNick(),
						"sex":      info.GetSex(),
						"age":      info.GetAge(),
						"area":     "",
						"card":     "",
						"level":    "",
						"role":     info.GetRole(),
						"title":    "unknown",
					},
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
						"nickname": XQApiGetNick(bot, userID),
						"sex":      XQApiGetGender(bot, userID),
						"age":      XQApiGetAge(bot, userID),
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
						"role":     "unknown",
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
					"user_id":     noticeID,
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
					"user_id":     noticeID,
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
					"time":         time,
					"self_id":      bot,
					"post_type":    "request",
					"request_type": "group",
					"sub_type":     "add",
					"group_id":     groupID,
					"user_id":      noticeID,
					"comment":      message,
					"flag":         fmt.Sprintf("%d|%d|%d|%s", messageType, groupID, 0, rawMessage),
				},
			}
			OnRequestGroupAdd(ctx)
		case 214:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "request",
					"request_type": "group",
					"sub_type":     "invite",
					"group_id":     groupID,
					"user_id":      noticeID,
					"comment":      message,
					"flag":         fmt.Sprintf("%d|%d|%d|%s", messageType, groupID, userID, rawMessage),
				},
			}
			OnRequestGroupAdd(ctx)
		case 215:
			ctx := &Context{
				Bot: bot,
				Response: map[string]interface{}{
					"time":         time,
					"self_id":      bot,
					"post_type":    "request",
					"request_type": "group",
					"sub_type":     "add",
					"group_id":     groupID,
					"user_id":      noticeID,
					"comment":      message,
					"flag":         fmt.Sprintf("%d|%d|%d|%s", messageType, groupID, 0, rawMessage),
				},
			}
			OnRequestGroupAdd(ctx)
		case 10000:
			go update(AppInfo.Pver, pathExecute())
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
