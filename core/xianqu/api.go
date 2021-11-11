package xianqu

//#include "api.h"
import "C"

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Yiwen-Chan/go-silk/silk"
	"github.com/tidwall/gjson"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

type Message struct {
	Bot       int64
	Send      string
	Type_     int64
	GroupID   int64
	UserID    int64
	Bubble    int64
	Anonymous bool
}

func (ctx *Context) MakeOkResponse(data interface{}) {
	ctx.Response = map[string]interface{}{
		"status":  "ok",
		"retcode": 0,
		"data":    data,
		"echo":    ctx.Request["echo"],
	}
}

func (ctx *Context) MakeFailResponse(err string) {
	ctx.Response = map[string]interface{}{
		"status":  "fail",
		"retcode": 100,
		"data":    err,
		"echo":    ctx.Request["echo"],
	}
}

// ApiSendMsg 发送私聊消息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func ApiSendMsg(ctx *Context) {
	var (
		params = Parse(ctx.Request).Get("params")
		sender = newMessage(ctx)
	)
	if !params.Exist("message") {
		return
	}
	messages := params.Array("message")
	for i := range messages {
		data := messages[i].Get("data")
		switch messages[i].Str("type") {
		default:
			//
		case "text":
			sender.text(data)
		case "at":
			sender.at(data)
		case "face":
			sender.face(data)
		case "emoji":
			sender.emoji(data)
		case "rps":
			sender.rps(data)
		case "dice":
			sender.dice(data)
		case "bubble":
			sender.bubble(data)
		// 媒体
		case "image":
			sender.image(data)
		case "record":
			sender.record(data)
		case "video":
			// TODO
		// 富文本
		case "xml":
			sender.xml(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "json":
			sender.json(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "share":
			sender.share(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "music":
			sender.music(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "weather":
			sender.weather(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "contact":
			sender.contact(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "location":
			sender.location(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "shake":
			sender.shake(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "poke":
			sender.poke(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "anonymous":
			sender.anonymous(data)
		case "reply":
			sender.reply(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "forward":
			sender.forward(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		case "node":
			sender.node(data)
			ctx.MakeOkResponse(map[string]interface{}{"message_id": 1})
			return
		}
	}
	if sender.Send == "" {
		ctx.MakeFailResponse("invalid message")
		return
	}
	// 发送信息
	id := sender.send()
	if id == 0 {
		ctx.MakeFailResponse("send failed")
	}
	ctx.MakeOkResponse(map[string]interface{}{"message_id": id})
}

// SendPrivateMsg 发送私聊消息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func ApiSendPrivateMsg(ctx *Context) {
	ApiSendMsg(ctx)
}

// SendGroupMsg 发送群消息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func ApiSendGroupMsg(ctx *Context) {
	ApiSendMsg(ctx)
}

// DeleteMsg 撤回消息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
func ApiDeleteMsg(ctx *Context) {
	id := Parse(ctx.Request).Get("params").Int("message_id")
	// 获取 先驱的 message num
	num := MessageIDCache.Search(id)
	if num == nil {
		ctx.MakeFailResponse("invalid message id")
		return
	}
	// 获取 OneBot 的 message 报文
	res := MessageCache.Search(id)
	if res == nil {
		ctx.MakeFailResponse("invalid message id")
		return
	}
	ctx.Response = res.(map[string]interface{})
	C.S3_Api_WithdrawMsgEX(
		goInt2CStr(ctx.Bot),
		C.int(ctx.GetResponseType()),
		goInt2CStr(Parse(ctx.Response).Int("group_id")),
		goInt2CStr(ctx.GetUserID()),
		goInt2CStr(num.(int64)),
		goInt2CStr(Parse(ctx.Response).Int("message_id")),
		goInt2CStr(Parse(ctx.Response).Int("time")),
	)
	ctx.MakeOkResponse(nil)
}

// GetMsg 获取消息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
func ApiGetMsg(ctx *Context) {
	id := Parse(ctx.Request).Get("params").Int("message_id")
	// 获取 OneBot 的 message 报文
	res := MessageCache.Search(id)
	if res == nil {
		ctx.MakeFailResponse("invalid message id")
		return
	}
	ctx.MakeOkResponse(res.(map[string]interface{}))
}

// GetForwardMsg 获取合并转发消息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func ApiGetForwardMsg(ctx *Context) {
	ctx.MakeFailResponse("xq not support")
}

// SendLike 发送好友赞
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#send_like-%E5%8F%91%E9%80%81%E5%A5%BD%E5%8F%8B%E8%B5%9E
func ApiSendLike(ctx *Context) {
	C.S3_Api_UpVote(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupKick 群组踢人
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func ApiSetGroupKick(ctx *Context) {
	C.S3_Api_KickGroupMBR(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		cBool(Parse(ctx.Request).Get("params").Bool("reject_add_request")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupBan 群组单人禁言
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func ApiSetGroupBan(ctx *Context) {
	C.S3_Api_ShutUp(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		C.int(Parse(ctx.Request).Get("params").Int("duration")/60),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupAnonymousBan 群组匿名用户禁言
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_anonymous_ban-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D%E7%94%A8%E6%88%B7%E7%A6%81%E8%A8%80
func ApiSetGroupAnonymousBan(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// SetGroupWholeBan 群组全员禁言
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func ApiSetGroupWholeBan(ctx *Context) {
	C.S3_Api_ShutUpAll(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		cBool(Parse(ctx.Request).Get("params").Bool("enable")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupAdmin 群组设置管理员
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_admin-%E7%BE%A4%E7%BB%84%E8%AE%BE%E7%BD%AE%E7%AE%A1%E7%90%86%E5%91%98
func ApiSetGroupAdmin(ctx *Context) {
	// 获取网页 cookie 和计算 bnk
	cookie1 := goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot)))
	cookie2 := goString(C.S3_Api_GetGroupPsKey(goInt2CStr(ctx.Bot)))
	cookie := cookie1 + cookie2
	bnk := getBnk(cookie1)
	// 提交 post 请求
	client := &http.Client{}
	dataUrl := url.Values{}
	dataUrl.Add("gc", Parse(ctx.Request).Get("params").Str("group_id"))
	dataUrl.Add("op", int2Str(Parse(ctx.Request).Get("params").Int("enable")))
	dataUrl.Add("ul", Parse(ctx.Request).Get("params").Str("user_id"))
	dataUrl.Add("bkn", strconv.Itoa(bnk))
	reqest, _ := http.NewRequest("POST", "https://qun.qq.com/cgi-bin/qun_mgr/set_group_admin", strings.NewReader(dataUrl.Encode()))
	reqest.Header.Set("Cookie", cookie)
	resp, err := client.Do(reqest)
	if err != nil {
		ApiOutPutLog(err)
		return
	}
	data, _ := ioutil.ReadAll(resp.Body)
	// 判断是否返回成功
	ul := gjson.ParseBytes(data).Get("ul").Str
	em := gjson.ParseBytes(data).Get("em").Str
	if ul != Parse(ctx.Request).Get("params").Str("user_id") {
		ctx.MakeFailResponse(em)
		return
	}
	ctx.MakeOkResponse(nil)
}

// SetGroupAnonymous 群组匿名
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func ApiSetGroupAnonymous(ctx *Context) {
	C.S3_Api_SetAnon(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		cBool(Parse(ctx.Request).Get("params").Bool("enable")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupCard 设置群名片（群备注）
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func ApiSetGroupCard(ctx *Context) {
	C.S3_Api_SetGroupCard(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		cString(Parse(ctx.Request).Get("params").Str("card")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupName 设置群名
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func ApiSetGroupName(ctx *Context) {
	ctx.MakeFailResponse("先驱不支持")
}

// SetGroupLeave 退出群组
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func ApiSetGroupLeave(ctx *Context) {
	C.S3_Api_QuitGroup(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func ApiSetGroupSpecialTitle(ctx *Context) {
	ctx.MakeFailResponse("先驱不支持")
}

// SetFriendAddRequest 处理加好友请求
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func ApiSetFriendAddRequest(ctx *Context) {
	C.S3_Api_HandleFriendEvent(
		goInt2CStr(ctx.Bot),
		goInt2CStr(Parse(ctx.Request).Get("params").Int("flag")),
		cBool(Parse(ctx.Request).Get("params").Bool("approve")),
		cString(Parse(ctx.Request).Get("params").Str("remark")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupAddRequest 处理加群请求／邀请
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func ApiSetGroupAddRequest(ctx *Context) {
	flag := strings.Split(Parse(ctx.Request).Get("params").Str("flag"), "|")
	if len(flag) != 4 {
		ctx.MakeFailResponse("invalid flag")
		return
	}
	C.S3_Api_HandleGroupEvent(
		goInt2CStr(ctx.Bot),
		C.int(str2Int(flag[0])),
		cString(flag[2]),
		cString(flag[1]),
		cString(flag[3]),
		cBool(Parse(ctx.Request).Get("params").Bool("approve")),
		cString(Parse(ctx.Request).Get("params").Str("reason")),
	)
}

// GetLoginInfo 获取登录号信息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func ApiGetLoginInfo(ctx *Context) {
	nickname := goString(
		C.S3_Api_GetNick(
			goInt2CStr(ctx.Bot),
			goInt2CStr(ctx.Bot),
		),
	)
	ctx.MakeOkResponse(
		map[string]interface{}{
			"user_id":  ctx.Bot,
			"nickname": nickname,
		},
	)
}

// GetStrangerInfo 获取陌生人信息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func ApiGetStrangerInfo(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetFriendList 获取好友列表
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func ApiGetFriendList(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetGroupInfo 获取群信息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func ApiGetGroupInfo(ctx *Context) {
	ctx.MakeOkResponse(
		GroupDataCache.GetCacheGroup(
			ctx.Bot,
			Parse(ctx.Request).Get("params").Int("group_id"),
			!Parse(ctx.Request).Get("params").Bool("cache"),
		).GroupInfo,
	)
}

// GetGroupList 获取群列表
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func ApiGetGroupList(ctx *Context) {
	var temp []GroupInfo
	list := strings.Split(
		goString(
			C.S3_Api_GetGroupList_B(
				goInt2CStr(ctx.Bot),
			),
		),
		"\r\n",
	)
	for _, groupID := range list {
		temp = append(temp, GroupInfo{
			GroupID: str2Int(groupID),
		})
	}
	ctx.MakeOkResponse(temp)
}

// GetGroupMemberInfo 获取群成员信息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func ApiGetGroupMemberInfo(ctx *Context) {
	ctx.MakeOkResponse(
		GroupDataCache.GetCacheGroupMember(
			ctx.Bot,
			Parse(ctx.Request).Get("params").Int("group_id"),
			Parse(ctx.Request).Get("params").Int("user_id"),
			!Parse(ctx.Request).Get("params").Bool("cache"),
		),
	)
}

// GetGroupMemberList 获取群成员列表
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func ApiGetGroupMemberList(ctx *Context) {
	ctx.MakeOkResponse(
		GroupDataCache.GetCacheGroup(
			ctx.Bot,
			Parse(ctx.Request).Get("params").Int("group_id"),
			!Parse(ctx.Request).Get("params").Bool("cache"),
		).GroupMembers,
	)
}

// GetGroupHonorInfo 获取群荣誉信息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func ApiGetGroupHonorInfo(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetCookies 获取 Cookies
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_cookies-%E8%8E%B7%E5%8F%96-cookies
func ApiGetCookies(ctx *Context) {
	switch Parse(ctx.Request).Get("params").Str("domain") {
	case "qun.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot))) + goString(C.S3_Api_GetGroupPsKey(goInt2CStr(ctx.Bot)))})
		return
	case "qzone.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot))) + goString(C.S3_Api_GetZonePsKey(goInt2CStr(ctx.Bot)))})
		return
	default:
		ctx.MakeOkResponse(map[string]interface{}{"cookies": goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot)))})
		return
	}
}

// GetCsrfToken 获取 CSRF Token
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_csrf_token-%E8%8E%B7%E5%8F%96-csrf-token
func ApiGetCsrfToken(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetCredentials 获取 QQ 相关接口凭证
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_credentials-%E8%8E%B7%E5%8F%96-qq-%E7%9B%B8%E5%85%B3%E6%8E%A5%E5%8F%A3%E5%87%AD%E8%AF%81
func ApiGetCredentials(ctx *Context) {
	switch Parse(ctx.Request).Get("params").Str("domain") {
	case "qun.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot))) + goString(C.S3_Api_GetGroupPsKey(goInt2CStr(ctx.Bot)))})
		return
	case "qzone.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot))) + goString(C.S3_Api_GetZonePsKey(goInt2CStr(ctx.Bot)))})
		return
	default:
		ctx.MakeOkResponse(map[string]interface{}{"cookies": goString(C.S3_Api_GetCookies(goInt2CStr(ctx.Bot)))})
		return
	}
}

// GetRecord 获取语音
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func ApiGetRecord(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetImage 获取图片
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func ApiGetImage(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// CanSendImage 检查是否可以发送图片
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#can_send_image-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E5%9B%BE%E7%89%87
func ApiCanSendImage(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"yes": true,
		},
	)
}

// CanSendRecord 检查是否可以发送语音
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#can_send_record-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E8%AF%AD%E9%9F%B3
func ApiCanSendRecord(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"yes": true,
		},
	)
}

// GetStatus 获取运行状态
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func ApiGetStatus(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"online": true,
			"good":   true,
		},
	)
}

// GetVersionInfo 获取版本信息
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#get_version_info-%E8%8E%B7%E5%8F%96%E7%89%88%E6%9C%AC%E4%BF%A1%E6%81%AF
func ApiGetVersionInfo(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"app_name":         "OneBot-YaYa",
			"app_version":      AppInfo.Pver,
			"protocol_version": "v11",
		},
	)
}

// SetRestart 重启 OneBot 实现
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#set_restart-%E9%87%8D%E5%90%AF-onebot-%E5%AE%9E%E7%8E%B0
func ApiSetRestart(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// CleanCache 清理缓存
// https://github.com/botuniverse/onebot-11/tree/master/api/public.md#clean_cache-%E6%B8%85%E7%90%86%E7%BC%93%E5%AD%98
func ApiCleanCache(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// SendXml 发送xml信息，相比cqcode方式不用处理转义
// ex: send_xml(group_id=xxx,data=xxx)
func ApiSendXml(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// SendJson 发送json信息，相比cqcode方式不用处理转义
// ex: send_json(group_id=xxx,data=xxx)
func ApiSendJson(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// ApiNotFound 没有这样的API
func ApiNotFound(ctx *Context) {
	ctx.MakeFailResponse("API NOT FOUND")
}

func newMessage(ctx *Context) *Message {
	return &Message{
		Bot:       ctx.Bot,
		Send:      "",
		Type_:     ctx.XQMessageType(),
		GroupID:   Parse(ctx.Request).Get("params").Int("group_id"),
		UserID:    Parse(ctx.Request).Get("params").Int("user_id"),
		Bubble:    0,
		Anonymous: false,
	}
}

func (m *Message) send() int64 {
	ret := goString(
		C.S3_Api_SendMsgEX_V2(
			goInt2CStr(m.Bot),
			C.int(m.Type_),
			goInt2CStr(m.GroupID),
			goInt2CStr(m.UserID),
			cString(escapeEmoji(m.Send)),
			C.int(m.Bubble),
			cBool(m.Anonymous),
			cString(""),
		),
	)
	// 处理返回的 message_id
	var temp map[string]interface{}
	json.Unmarshal(helper.StringToBytes(ret), &temp)
	num := Parse(temp).Int("msgno")
	if num == 0 {
		return 0
	}
	id := MessageIDCache.Hcraes(num)
	if id == nil {
		return 0
	}
	return id.(int64)
}

func (m *Message) json(data value) {
	C.S3_Api_SendJSON(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(escape(data.Str("data"))), // nonebot 的数组没有转义，这里暂时这么解决
	)
}

func (m *Message) xml(data value) {
	C.S3_Api_SendXML(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(escape(data.Str("data"))), // nonebot 的数组没有转义，这里暂时这么解决
		0,
	)
}

func (m *Message) text(data value) {
	m.Send += data.Str("text")
}

func (m *Message) at(data value) {
	m.Send += fmt.Sprintf("[@%s] ", data.Str("qq"))
}

func (m *Message) face(data value) {
	m.Send += fmt.Sprintf("[Face%s.gif]", data.Str("id"))
}

func (m *Message) emoji(data value) {
	m.Send += fmt.Sprintf("[emoji=%s]", data.Str("id"))
}

func (m *Message) rps(data value) {
	m.Send += []string{"[魔法猜拳] 石头", "[魔法猜拳] 剪刀", "[魔法猜拳] 布"}[rand.Intn(3)]
}

func (m *Message) dice(data value) {
	m.Send += []string{"[魔法骰子] 1", "[魔法骰子] 2", "[魔法骰子] 3", "[魔法骰子] 4", "[魔法骰子] 5", "[魔法骰子] 6"}[rand.Intn(6)]
}

func (m *Message) bubble(data value) {
	m.Bubble += data.Int("id")
}

type image struct {
	type_  string
	res    string
	showID int64
}

func newImage(data value) *image {
	return &image{
		type_:  data.Str("type"),
		res:    "",
		showID: 0,
	}
}

func (i *image) insert(m *Message) {
	switch i.type_ {
	default:
		m.Send += fmt.Sprintf("[pic=%s]", i.res)
	case "show":
		m.Send += fmt.Sprintf("[ShowPic=%s,type=%d]", i.res, i.showID+40000)
	}
}

func (m *Message) image(data value) {
	var (
		file  = data.Str("file")
		cache = true
		image = newImage(data)
	)
	if data.Exist("cache") {
		cache = data.Bool("cache")
	}
	if data.Str("url") != "" {
		// 解决tx图片链接不落地
		temp := PicPoolCache.Search(textMD5(data.Str("url")))
		if temp != nil && cache {
			image.res = temp.(string)
			image.insert(m)
			return
		}
	}
	// 判断file字段为哪种类型
	switch {
	// 链接方式
	case strings.Contains(file, "http://") || strings.Contains(file, "https://"):
		path := OneBotPath + "image\\" + textMD5(file) + ".jpg"
		if !pathExists(path) || !cache {
			// 下载图片
			if err := Download(file, path); err != nil {
				panic(err)
			}
		}
		image.res = path
	// base64方式
	case strings.Contains(file, "base64://"):
		path := OneBotPath + "image\\" + textMD5(file[9:]) + ".jpg"
		if err := DecodeBase64(file[9:], path); err != nil {
			panic(err)
		}
		image.res = path
	// 本地文件方式
	case strings.Contains(file, "file:///"):
		image.res = strings.ReplaceAll(file[8:], "/", "\\")
	// 默认方式
	default:
		image.res = file
	}
	// 用文件md5判断缓冲池是否存在该图片
	temp := PicPoolCache.Search(fileMD5(image.res))
	if temp != nil && cache {
		image.res = temp.(string)
	}
	image.insert(m)
}

type record struct {
	res string
}

func newRecord(data value) *record {
	return &record{
		res: "",
	}
}

func (i *record) insert(m *Message) {
	m.Send += fmt.Sprintf("[Voi=%s]", i.res)
}

func (m *Message) record(data value) {
	var (
		file   = data.Str("file")
		cache  = true
		record = newRecord(data)
	)
	if data.Exist("cache") {
		cache = data.Bool("cache")
	}
	// 判断file字段为哪种类型
	switch {
	// 链接方式
	case strings.Contains(file, "http://") || strings.Contains(file, "https://"):
		path := OneBotPath + "record\\" + textMD5(file)
		if !pathExists(path) || !cache {
			// 下载音频
			if err := Download(file, path); err != nil {
				panic(err)
			}
		}
		record.res = path
	// base64方式
	case strings.Contains(file, "base64://"):
		path := OneBotPath + "record\\" + textMD5(file[9:])
		if err := DecodeBase64(file[9:], path); err != nil {
			panic(err)
		}
		record.res = path
	// 本地文件方式
	case strings.Contains(file, "file:///"):
		record.res = strings.ReplaceAll(file[8:], "/", "\\")
	// 默认方式
	default:
		record.res = file
	}
	name := strings.ReplaceAll(filepath.Base(record.res), path.Ext(record.res), "")
	silkEncoder := &silk.Encoder{}
	if err := silkEncoder.Init("OneBot/record", "OneBot/codec"); err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(record.res)
	if err != nil {
		panic(err)
	}
	_, err = silkEncoder.EncodeToSilk(b, name, true)
	if err != nil {
		panic(err)
	}
	record.res = OneBotPath + "record/" + name + ".silk"
	record.insert(m)
}

func (m *Message) share(data value) {
	temp := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
<msg serviceID="33" templateID="123" action="web" brief="%s" 
sourceMsgId="0" url="%s" 
flag="8" adverSign="0" multiMsgFlag="0"><item layout="2" 
advertiser_id="0" aid="0"><picture cover="%s" w="0" h="0" />
<title>%s</title><summary>%s</summary>
</item><source name="" icon="" action="" appid="-1" /></msg>`,
		data.Str("brief"),
		data.Str("url"),
		data.Str("image"),
		data.Str("title"),
		data.Str("content"),
	)
	C.S3_Api_SendXML(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(strings.ReplaceAll(temp, "\n", "")),
		0,
	)
}

func (m *Message) music(data value) {
	temp := fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
<msg serviceID="2" templateID="1" action="web" brief="[分享] %s" 
sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0">
<item layout="2"><audio cover="%s" src="%s"/><title>%s</title>
<summary>%s</summary></item><source name="音乐" 
icon="https://i.gtimg.cn/open/app_icon/01/07/98/56/1101079856_100_m.png" 
url="http://web.p.qq.com/qqmpmobile/aio/app.html?id=1101079856" 
action="app" a_actionData="com.tencent.qqmusic" 
i_actionData="tencent1101079856://" appid="1101079856" /></msg>`,
		xmlEscape(data.Str("title")),
		data.Str("url"),
		data.Str("image"),
		data.Str("audio"),
		xmlEscape(data.Str("title")),
		xmlEscape(data.Str("content")),
	)
	C.S3_Api_SendXML(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(strings.ReplaceAll(temp, "\n", "")),
		0,
	)
}

func (m *Message) weather(data value) {
	temp := fmt.Sprintf(`{"app":"com.tencent.weather","desc":"天气",
"view":"RichInfoView","ver":"0.0.0.1","prompt":"[应用]天气",
"appID":"","sourceName":"","actionData":"","actionData_A":"",
"sourceUrl":"","meta":{"richinfo":{"adcode":"","air":"%s",
"city":"%s","date":"%s","max":"%s","min":"%s",
"ts":"15158613","type":"%s","wind":""}},"text":"","sourceAd":"","extra":""}`,
		data.Str("air"),
		data.Str("city"),
		data.Str("date"),
		data.Str("max"),
		data.Str("min"),
		data.Str("type"),
	)
	C.S3_Api_SendJSON(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(strings.ReplaceAll(temp, "\n", "")),
	)
}

func (m *Message) location(data value) {
	temp := fmt.Sprintf(`{"app":"com.tencent.map","desc":"","view":"Share",
"ver":"0.0.0.1","prompt":"[应用]地图","appID":"","sourceName":"",
"actionData":"","actionData_A":"","sourceUrl":"","meta":{"Share":{"locSub":"%s",
"lng":%s,"lat":%s,"zoom":15,"locName":"%s"}},
"config":{"forward":true,"autosize":1},"text":"","extraApps":[],
"sourceAd":"","extra":""}`,
		data.Str("content"),
		data.Str("lon"),
		data.Str("lat"),
		data.Str("title"),
	)
	C.S3_Api_SendJSON(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(strings.ReplaceAll(temp, "\n", "")),
	)
}
func (m *Message) shake(data value) {
	C.S3_Api_ShakeWindow(
		goInt2CStr(m.Bot),
		goInt2CStr(m.UserID),
	)
}
func (m *Message) poke(data value) {
	C.S3_Api_SendMsgEX_V2(
		goInt2CStr(m.Bot),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(fmt.Sprintf(
			"[系统提示] %v 戳了一下 %v",
			m.Bot,
			m.UserID,
		)),
		C.int(0),
		cBool(false),
		cString(""),
	)
}

func (m *Message) anonymous(data value) {
	m.Anonymous = true
}

func (m *Message) reply(data value) {
	C.S3_Api_SendMsgEX_V2(
		goInt2CStr(m.Bot),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(fmt.Sprintf(
			"[系统消息] %v 尝试回复一条消息并失败了",
			m.Bot,
		)),
		C.int(0),
		cBool(false),
		cString(""),
	)
}

func (m *Message) forward(data value) {
	C.S3_Api_SendMsgEX_V2(
		goInt2CStr(m.Bot),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(fmt.Sprintf(
			"[系统消息] %v 尝试合并转发一条消息并失败了",
			m.Bot,
		)),
		C.int(0),
		cBool(false),
		cString(""),
	)
}
func (m *Message) node(data value) {
	C.S3_Api_SendMsgEX_V2(
		goInt2CStr(m.Bot),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(fmt.Sprintf(
			"[系统消息] %v 尝试合并转发节点并失败了",
			m.Bot,
		)),
		C.int(0),
		cBool(false),
		cString(""),
	)
}
func (m *Message) contact(data value) {
	temp := ""
	switch data.Str("type") {
	case "qq":
		temp = fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
<msg serviceID="14" templateID="1" action="plugin" 
actionData="AppCmd://OpenContactInfo/?uin=%s" 
a_actionData="mqqapi://card/show_pslcard?src_type=internal&amp;
source=sharecard&amp;version=1&amp;uin=%s" 
i_actionData="mqqapi://card/show_pslcard?src_type=internal&amp;
source=sharecard&amp;version=1&amp;uin=%s" 
brief="推荐了%s" sourceMsgId="0" url="" flag="1" 
adverSign="0" multiMsgFlag="0"><item layout="0" 
mode="1" advertiser_id="0" aid="0">
<summary>推荐联系人</summary><hr hidden="false" style="0" />
</item><item layout="2" mode="1" advertiser_id="0" aid="0">
<picture cover="mqqapi://card/show_pslcard?src_type=internal&amp;
source=sharecard&amp;version=1&amp;uin=%s" w="0" h="0" />
<title>%s</title><summary>帐号:%s</summary>
</item><source name="" icon="" action="" appid="-1" /></msg>`,
			data.Str("id"),
			data.Str("id"),
			data.Str("id"),
			data.Str("name"),
			data.Str("id"),
			data.Str("name"),
			data.Str("id"),
		)
	case "group":
		temp = fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
<msg serviceID="15" templateID="1" action="web" 
actionData="group:%s" a_actionData="group:%s" 
i_actionData="group:%s" brief="推荐群聊：%s" 
sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0">
<item layout="0" mode="1" advertiser_id="0" aid="0">
<summary>推荐群聊</summary><hr hidden="false" style="0" />
</item><item layout="2" mode="1" advertiser_id="0" aid="0">
<picture cover="https://p.qlogo.cn/gh/%s/%s/100" w="0" h="0" needRoundView="0" />
<title>%s</title><summary>创建人：%s</summary></item>
<source name="" icon="" action="" appid="-1" /></msg>`,
			data.Str("id"),
			data.Str("id"),
			data.Str("id"),
			data.Str("name"),
			data.Str("url"),
			data.Str("id"),
			data.Str("id"),
			data.Str("name"),
			data.Str("owner"),
		)
	}
	C.S3_Api_SendXML(
		goInt2CStr(m.Bot),
		C.int(1),
		C.int(m.Type_),
		goInt2CStr(m.GroupID),
		goInt2CStr(m.UserID),
		cString(strings.ReplaceAll(temp, "\n", "")),
		0,
	)
}

func ApiOutPutLog(text interface{}) {
	C.S3_Api_OutPutLog(
		cString(fmt.Sprintln(text)),
	)
}

func ApiOutPutLog1(text interface{}) {
	fmt.Println(text)
}

func XQApiGroupName(bot, groupID int64) string {
	return goString(
		C.S3_Api_GetGroupName(
			goInt2CStr(bot),
			goInt2CStr(groupID),
		),
	)
}

func XQApiGroupMemberListB(bot, groupID int64) string {
	return goString(
		C.S3_Api_GetGroupMemberList_B(
			goInt2CStr(bot),
			goInt2CStr(groupID),
		),
	)
}

func XQApiGroupMemberListC(bot, groupID int64) string {
	return goString(
		C.S3_Api_GetGroupMemberList_C(
			goInt2CStr(bot),
			goInt2CStr(groupID),
		),
	)
}

func XQApiGetNick(bot, userID int64) string {
	return strings.Split(
		goString(
			C.S3_Api_GetNick(
				goInt2CStr(bot),
				goInt2CStr(userID),
			),
		),
		"\n",
	)[0]
}

func XQApiGetAge(bot, userID int64) int64 {
	return int64(
		C.S3_Api_GetAge(
			goInt2CStr(bot),
			goInt2CStr(userID),
		),
	)
}

func XQApiGetGender(bot, userID int64) string {
	return []string{"unknown", "male", "female"}[int64(
		C.S3_Api_GetGender(
			goInt2CStr(bot),
			goInt2CStr(userID),
		),
	)]
}

func XQApiIsFriend(bot, userID int64) bool {
	return goBool(
		C.S3_Api_IfFriend(
			goInt2CStr(bot),
			goInt2CStr(userID),
		),
	)
}

func ApiCallMessageBox(text string) {
	C.S3_Api_CallMessageBox(
		cString(text),
	)
}

func ApiMessageBoxButton(text string) int64 {
	// 6 为是
	// 7 为否
	return int64(
		C.S3_Api_MessageBoxButton(
			cString(text),
		),
	)
}

func ApiDefaultQQ() int64 {
	botList := strings.Split(
		goString(C.S3_Api_GetQQList()),
		"/n",
	)
	if len(botList) < 0 {
		return 0
	}
	return str2Int(botList[0])
}
