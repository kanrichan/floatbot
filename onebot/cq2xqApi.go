package onebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"yaya/core"
)

type Result struct {
	Status  string      `json:"status"`
	Retcode int64       `json:"retcode"`
	Data    interface{} `json:"data"`
	Echo    interface{} `json:"echo"`
}

var apiList = map[string]func(int64, gjson.Result) Result{
	"send_msg": func(bot int64, p gjson.Result) Result {
		return cq2xqSendMsg(bot, p)
	},
	"send_private_msg": func(bot int64, p gjson.Result) Result {
		return cq2xqSendMsg(bot, p)
	},
	"send_group_msg": func(bot int64, p gjson.Result) Result {
		return cq2xqSendMsg(bot, p)
	},
	"delete_msg": func(bot int64, p gjson.Result) Result {
		return cq2xqDeleteMsg(bot, p)
	},
	"get_msg": func(bot int64, p gjson.Result) Result {
		return cq2xqGetMsg(bot, p)
	},
	"get_forward_msg": func(bot int64, p gjson.Result) Result {
		return cq2xqGetForwardMsg(bot, p)
	},
	"send_like": func(bot int64, p gjson.Result) Result {
		return cq2xqSendLike(bot, p)
	},
	"set_group_kick": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupKick(bot, p)
	},
	"set_group_ban": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupBan(bot, p)
	},
	"set_group_anonymous_ban": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupAnonymousBan(bot, p)
	},
	"set_group_whole_ban": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupWholeBan(bot, p)
	},
	"set_group_admin": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupAdmin(bot, p)
	},
	"set_group_anonymous": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupAnonymous(bot, p)
	},
	"set_group_card": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupCard(bot, p)
	},
	"set_group_name": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupName(bot, p)
	},
	"set_group_leave": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupLeave(bot, p)
	},
	"set_group_special_title": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupSpecialTitle(bot, p)
	},
	"set_friend_add_request": func(bot int64, p gjson.Result) Result {
		return cq2xqSetFriendAddRequest(bot, p)
	},
	"set_group_add_request": func(bot int64, p gjson.Result) Result {
		return cq2xqSetGroupAddRequest(bot, p)
	},
	"get_login_info": func(bot int64, p gjson.Result) Result {
		return cq2xqGetLoginInfo(bot, p)
	},
	"get_stranger_info": func(bot int64, p gjson.Result) Result {
		return cq2xqGetStrangerInfo(bot, p)
	},
	"get_friend_list": func(bot int64, p gjson.Result) Result {
		return cq2xqGetFriendList(bot, p)
	},
	"get_group_info": func(bot int64, p gjson.Result) Result {
		return cq2xqGetGroupInfo(bot, p)
	},
	"get_group_list": func(bot int64, p gjson.Result) Result {
		return cq2xqGetGroupList(bot, p)
	},
	"get_group_member_info": func(bot int64, p gjson.Result) Result {
		return cq2xqGetGroupMemberInfo(bot, p)
	},
	"get_group_member_list": func(bot int64, p gjson.Result) Result {
		return cq2xqGetGroupMemberList(bot, p)
	},
	"get_group_honor_info": func(bot int64, p gjson.Result) Result {
		return cq2xqGetGroupHonorInfo(bot, p)
	},
	"get_cookies": func(bot int64, p gjson.Result) Result {
		return cq2xqGetCookies(bot, p)
	},
	"get_csrf_token": func(bot int64, p gjson.Result) Result {
		return cq2xqGetCsrfToken(bot, p)
	},
	"get_credentials": func(bot int64, p gjson.Result) Result {
		return cq2xqGetCredentials(bot, p)
	},
	"get_record": func(bot int64, p gjson.Result) Result {
		return cq2xqGetRecord(bot, p)
	},
	"get_image": func(bot int64, p gjson.Result) Result {
		return cq2xqGetImage(bot, p)
	},
	"can_send_image": func(bot int64, p gjson.Result) Result {
		return cq2xqCanSendImage(bot, p)
	},
	"can_send_record": func(bot int64, p gjson.Result) Result {
		return cq2xqCanSendRecord(bot, p)
	},
	"get_status": func(bot int64, p gjson.Result) Result {
		return cq2xqGetStatus(bot, p)
	},
	"get_version_info": func(bot int64, p gjson.Result) Result {
		return cq2xqGetVersionInfo(bot, p)
	},
	"set_restart": func(bot int64, p gjson.Result) Result {
		return cq2xqSetRestart(bot, p)
	},
	"clean_cache": func(bot int64, p gjson.Result) Result {
		return cq2xqCleanCache(bot, p)
	},
	".handle_quick_operation": func(bot int64, p gjson.Result) Result {
		return cq2xqHandleQuickOperation(bot, p)
	},
	// 先驱新增
	"out_put_log": func(bot int64, p gjson.Result) Result {
		return cq2xqOutPutLog(bot, p)
	},
	"send_json": func(bot int64, p gjson.Result) Result {
		return cq2xqSendJson(bot, p)
	},
	"send_xml": func(bot int64, p gjson.Result) Result {
		return cq2xqSendXml(bot, p)
	},
}

func resultOK(data interface{}) Result {
	return Result{
		Status:  "ok",
		Retcode: 0,
		Data:    data,
		Echo:    nil,
	}
}

func resultFail(data interface{}) Result {
	return Result{
		Status:  "failed",
		Retcode: 100,
		Data:    data,
		Echo:    nil,
	}
}

func cq2xqDeleteMsg(bot int64, p gjson.Result) Result {
	var xe XEvent
	for i, _ := range Conf.BotConfs {
		if bot == Conf.BotConfs[i].Bot && Conf.BotConfs[i].DB != nil {
			xe = db2Mseeage(Conf.BotConfs[i].DB, bot, p.Get("message_id").Int())
			break
		}
		return resultFail(map[string]interface{}{"error": "空指针"})
	}
	core.WithdrawMsgEX(
		xe.selfID,
		xe.mseeageType,
		xe.groupID,
		xe.userID,
		xe.messageNum,
		xe.messageID,
		xe.time,
	)
	return resultOK(map[string]interface{}{})
}

func cq2xqGetMsg(bot int64, p gjson.Result) Result {
	var xe XEvent
	for i, _ := range Conf.BotConfs {
		if bot == Conf.BotConfs[i].Bot {
			xe = db2Mseeage(Conf.BotConfs[i].DB, bot, p.Get("message_id").Int())
			break
		}
		return resultFail(map[string]interface{}{"error": "空指针"})
	}
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
	return resultOK(map[string]interface{}{
		"time":         xe.time,
		"message_type": Tsubtype,
		"message_id":   p.Get("message_id").Int(),
		"real_id":      xe.messageID,
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
		"message": xe.message,
	})
}

func cq2xqGetForwardMsg(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "先驱好像不支持"})
}

func cq2xqSendLike(bot int64, p gjson.Result) Result {
	core.UpVote(
		bot,
		p.Get("user_id").Int())
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupKick(bot int64, p gjson.Result) Result {
	core.KickGroupMBR(
		bot,
		p.Get("group_id").Int(),
		p.Get("user_id").Int(),
		p.Get("reject_add_request").Bool(),
	)
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupBan(bot int64, p gjson.Result) Result {
	core.ShutUP(
		bot,
		p.Get("group_id").Int(),
		p.Get("user_id").Int(),
		p.Get("duration").Int(),
	)
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupAnonymousBan(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "还没写，催更去GitHub提issue"})
}

func cq2xqSetGroupWholeBan(bot int64, p gjson.Result) Result {
	if p.Get("enable").Bool() {
		core.ShutUP(
			bot,
			p.Get("group_id").Int(),
			0,
			1,
		)
	} else {
		core.ShutUP(
			bot,
			p.Get("group_id").Int(),
			0,
			0,
		)
	}
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupAdmin(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "先驱好像不支持"})
}

func cq2xqSetGroupAnonymous(bot int64, p gjson.Result) Result {
	core.SetAnon(
		bot,
		p.Get("group_id").Int(),
		p.Get("enable").Bool(),
	)
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupCard(bot int64, p gjson.Result) Result {
	core.SetGroupCard(
		bot,
		p.Get("group_id").Int(),
		p.Get("user_id").Int(),
		p.Get("card").Str,
	)
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupName(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "先驱好像不支持"})
}

func cq2xqSetGroupLeave(bot int64, p gjson.Result) Result {
	core.QuitGroup(
		bot,
		p.Get("group_id").Int(),
	)
	return resultOK(map[string]interface{}{"message_id": 0})
}

func cq2xqSetGroupSpecialTitle(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "先驱好像不支持"})
}

func cq2xqSetFriendAddRequest(bot int64, p gjson.Result) Result {
	if p.Get("approve").Bool() {
		core.HandleFriendEvent(
			bot,
			core.Str2Int(p.Get("flag").Str),
			10,
			p.Get("remark").Str,
		)
	} else {
		core.HandleFriendEvent(
			bot,
			core.Str2Int(p.Get("flag").Str),
			20,
			p.Get("remark").Str,
		)
	}
	return resultOK(map[string]interface{}{})
}

func cq2xqSetGroupAddRequest(bot int64, p gjson.Result) Result {
	if p.Get("approve").Bool() {
		core.HandleGroupEvent(bot,
			core.Str2Int(strings.Split(p.Get("flag").Str, "|")[0]),
			0,
			core.Str2Int(strings.Split(p.Get("flag").Str, "|")[1]),
			core.Str2Int(strings.Split(p.Get("flag").Str, "|")[2]),
			10,
			p.Get("reason").Str,
		)
	} else {
		core.HandleGroupEvent(bot,
			core.Str2Int(strings.Split(p.Get("flag").Str, "|")[0]),
			0,
			core.Str2Int(strings.Split(p.Get("flag").Str, "|")[1]),
			core.Str2Int(strings.Split(p.Get("flag").Str, "|")[2]),
			20,
			p.Get("reason").Str,
		)
	}
	return resultOK(map[string]interface{}{})
}

func cq2xqGetLoginInfo(bot int64, p gjson.Result) Result {
	nickname := strings.Split(core.GetNick(bot, bot), "\n")[0]
	return resultOK(map[string]interface{}{"user_id": bot, "nickname": nickname})
}

func cq2xqGetStrangerInfo(bot int64, p gjson.Result) Result {
	return resultOK(map[string]interface{}{
		"user_id":  p.Get("user_id").Int(),
		"nickname": core.GetNick(bot, p.Get("user_id").Int()),
		"sex":      []string{"unknown", "male", "female"}[core.GetGender(bot, p.Get("user_id").Int())],
		"age":      core.GetAge(bot, p.Get("user_id").Int()),
	})
}

func cq2xqGetFriendList(bot int64, p gjson.Result) Result {
	list := core.GetFriendList(bot)
	if list == "" {
		return resultFail(map[string]interface{}{"error": "无法获取到好友列表"})
	}
	g := gjson.Parse(list)
	friendList := []map[string]interface{}{}
	for _, o := range g.Get("result.0.mems").Array() {
		info := map[string]interface{}{
			"user_id":  o.Get("uin").Int(),
			"nickname": unicode2chinese(o.Get("name").Str),
			"remark":   "unknown",
		}
		friendList = append(friendList, info)
	}
	return resultOK(friendList)
}

func cq2xqGetGroupInfo(bot int64, p gjson.Result) Result {
	name := core.GetGroupName(
		bot,
		p.Get("group_id").Int(),
	)
	members := strings.Split(core.GetGroupMemberNum(
		bot,
		p.Get("group_id").Int(),
	), "\n")
	var (
		count int64
		max   int64
	)
	if len(members) != 2 {
		count = -1
		max = -1
	} else {
		count = core.Str2Int(members[0])
		max = core.Str2Int(members[1])
	}
	return resultOK(map[string]interface{}{
		"group_id":         p.Get("group_id").Int(),
		"group_name":       name,
		"member_count":     count,
		"max_member_count": max,
	})
}

func cq2xqGetGroupList(bot int64, p gjson.Result) Result {
	list := core.GetGroupList(bot)
	if list == "" {
		return resultFail(map[string]interface{}{"error": "无法获取到群列表"})
	}
	g := gjson.Parse(list)
	groupList := []map[string]interface{}{}
	for _, o := range g.Get("create").Array() {
		info := map[string]interface{}{
			"group_id":         o.Get("gc").Int(),
			"group_name":       unicode2chinese(o.Get("gn").Str),
			"member_count":     0,
			"max_member_count": 0,
		}
		groupList = append(groupList, info)
	}
	for _, o := range g.Get("manage").Array() {
		info := map[string]interface{}{
			"group_id":         o.Get("gc").Int(),
			"group_name":       unicode2chinese(o.Get("gn").Str),
			"member_count":     0,
			"max_member_count": 0,
		}
		groupList = append(groupList, info)
	}
	for _, o := range g.Get("join").Array() {
		info := map[string]interface{}{
			"group_id":         o.Get("gc").Int(),
			"group_name":       unicode2chinese(o.Get("gn").Str),
			"member_count":     0,
			"max_member_count": 0,
		}
		groupList = append(groupList, info)
	}
	return resultOK(groupList)
}

func cq2xqGetGroupMemberInfo(bot int64, p gjson.Result) Result {
	return resultOK(map[string]interface{}{
		"group_id":          p.Get("group_id").Int(),
		"user_id":           p.Get("user_id").Int(),
		"nickname":          core.GetNick(bot, p.Get("user_id").Int()),
		"card":              core.GetNick(bot, p.Get("user_id").Int()),
		"sex":               []string{"unknown", "male", "female"}[core.GetGender(bot, p.Get("user_id").Int())],
		"age":               core.GetAge(bot, p.Get("user_id").Int()),
		"area":              "unknown",
		"join_time":         0,
		"last_sent_time":    0,
		"level":             "unknown",
		"role":              "unknown",
		"unfriendly":        false,
		"title":             "unknown",
		"title_expire_time": 0,
		"card_changeable":   true,
	})
}

func cq2xqGetGroupMemberList(bot int64, p gjson.Result) Result {
	list := core.GetGroupMemberList_C(
		bot,
		p.Get("group_id").Int(),
	)
	if list == "" {
		return resultFail(map[string]interface{}{"error": "无法获取到群成员列表"})
	}
	g := gjson.Parse(list)
	memberList := []map[string]interface{}{}
	for _, o := range g.Get("list").Array() {
		member := map[string]interface{}{
			"group_id":          p.Get("group_id").Int(),
			"user_id":           o.Get("QQ").Int(),
			"nickname":          "unknown",
			"card":              "unknown",
			"sex":               "unknown",
			"age":               0,
			"area":              "unknown",
			"join_time":         0,
			"last_sent_time":    0,
			"level":             o.Get("lv").Int(),
			"role":              "unknown",
			"unfriendly":        false,
			"title":             "unknown",
			"title_expire_time": 0,
			"card_changeable":   true,
		}
		memberList = append(memberList, member)
	}
	return resultOK(memberList)
}

func cq2xqGetGroupHonorInfo(bot int64, p gjson.Result) Result {
	groupID := p.Get("group_id").Int()
	cookie := fmt.Sprintf("%s%s", core.GetCookies(bot), core.GetGroupPsKey(bot))
	var honorType int64 = 1
	switch p.Get("type").Str {
	case "talkative":
		honorType = 1
	case "performer":
		honorType = 2
	case "legend":
		honorType = 3
	case "strong_newbie":
		honorType = 5
	case "emotion":
		honorType = 6
	}
	data := groupHonor(groupID, honorType, cookie)
	if data != nil {
		data = data[bytes.Index(data, []byte(`window.__INITIAL_STATE__=`))+25:]
		data = data[:bytes.Index(data, []byte("</script>"))]
		ret := GroupHonorInfo{}
		json.Unmarshal(data, &ret)
		return resultOK(ret)
	} else {
		return resultFail(map[string]interface{}{"data": "error"})
	}
}

func cq2xqGetCookies(bot int64, p gjson.Result) Result {
	switch p.Get("domain").Str {
	case "qun.qq.com":
		return resultOK(map[string]interface{}{"cookies": core.GetCookies(bot) + core.GetGroupPsKey(bot)})
	case "qzone.qq.com":
		return resultOK(map[string]interface{}{"cookies": core.GetCookies(bot) + core.GetZonePsKey(bot)})
	default:
		return resultOK(map[string]interface{}{"cookies": core.GetCookies(bot)})
	}
}

func cq2xqGetCsrfToken(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "还没写，催更去GitHub提issue"})
}

func cq2xqGetCredentials(bot int64, p gjson.Result) Result {
	switch p.Get("domain").Str {
	case "qun.qq.com":
		return resultOK(map[string]interface{}{"cookies": core.GetCookies(bot) + core.GetGroupPsKey(bot)})
	case "qzone.qq.com":
		return resultOK(map[string]interface{}{"cookies": core.GetCookies(bot) + core.GetZonePsKey(bot)})
	default:
		return resultOK(map[string]interface{}{"cookies": core.GetCookies(bot)})
	}
}

func cq2xqGetRecord(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "还没写，催更去GitHub提issue"})
}

func cq2xqGetImage(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "还没写，催更去GitHub提issue"})
}

func cq2xqCanSendImage(bot int64, p gjson.Result) Result {
	return resultOK(map[string]interface{}{"yes": true})
}

func cq2xqCanSendRecord(bot int64, p gjson.Result) Result {
	return resultOK(map[string]interface{}{"yes": true})
}

func cq2xqGetStatus(bot int64, p gjson.Result) Result {
	return resultOK(map[string]interface{}{
		"online": core.IsOnline(bot, bot),
		"good":   true,
	})
}

func cq2xqGetVersionInfo(bot int64, p gjson.Result) Result {
	return resultOK(map[string]interface{}{
		"app_name":         "OneBot-YaYa",
		"app_version":      gjson.Parse(AppInfoJson).Get("pver"),
		"protocol_version": "v11",
	})
}

func cq2xqSetRestart(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "还没写，催更去GitHub提issue"})
}

func cq2xqCleanCache(bot int64, p gjson.Result) Result {
	return resultFail(map[string]interface{}{"data": "还没写，催更去GitHub提issue"})
}

func cq2xqOutPutLog(bot int64, p gjson.Result) Result {
	core.OutPutLog(p.Get("text").Str)
	return resultOK(map[string]interface{}{})
}

func cq2xqHandleQuickOperation(bot int64, p gjson.Result) Result {
	var reply []byte
	switch p.Get("context.post_type").Str {
	case "message":
		if p.Get("operation.reply").Exists() {
			var text string
			if p.Get("operation.at_sender").Bool() {
				text += fmt.Sprintf("[CQ:at,qq=%d]", p.Get("context.user_id").Int())
			}
			if p.Get("operation.reply").Str != "" {
				text += unicode2chinese(p.Get("operation.reply").Str)
			}
			reply, _ = json.Marshal(
				map[string]interface{}{
					"message_type": p.Get("context.message_type").Str,
					"group_id":     p.Get("context.group_id").Int(),
					"user_id":      p.Get("context.user_id").Int(),
					"message":      text,
				})
		}
		if p.Get("operation.delete").Bool() {
			parmas, _ := json.Marshal(
				map[string]interface{}{
					"message_id": p.Get("context.message_id").Int(),
				})
			cq2xqDeleteMsg(bot, gjson.ParseBytes(parmas))
		}
		if p.Get("operation.kick").Bool() {
			parmas, _ := json.Marshal(
				map[string]interface{}{
					"group_id":           p.Get("context.group_id").Int(),
					"user_id":            p.Get("context.user_id").Int(),
					"reject_add_request": false,
				})
			cq2xqSetGroupKick(bot, gjson.ParseBytes(parmas))
		}
		if p.Get("operation.ban").Bool() {
			parmas, _ := json.Marshal(
				map[string]interface{}{
					"group_id": p.Get("context.group_id").Int(),
					"user_id":  p.Get("context.user_id").Int(),
					"duration": p.Get("operation.duration").Int(),
				})
			cq2xqDeleteMsg(bot, gjson.ParseBytes(parmas))
		}
	case "request":
		if p.Get("operation.approve").Exists() {
			if p.Get("operation.request_type").Str == "friend" {
				parmas, _ := json.Marshal(
					map[string]interface{}{
						"flag":    p.Get("context.flag").Str,
						"approve": p.Get("operation.approve").Bool(),
						"remark":  p.Get("operation.remark").Str,
					})
				cq2xqSetFriendAddRequest(bot, gjson.ParseBytes(parmas))
			}
			if p.Get("operation.request_type").Str == "group" {
				parmas, _ := json.Marshal(
					map[string]interface{}{
						"flag":    p.Get("context.flag").Str,
						"approve": p.Get("operation.approve").Str,
						"reason":  p.Get("operation.reason").Bool(),
					})
				cq2xqSetGroupAddRequest(bot, gjson.ParseBytes(parmas))
			}
		}
	default:
		//
	}
	return cq2xqSendMsg(bot, gjson.ParseBytes(reply))
}

func cq2xqSendJson(bot int64, p gjson.Result) Result {
	core.SendJSON(
		bot,
		1,
		cq2xqMessageType(p),
		p.Get("group_id").Int(),
		p.Get("user_id").Int(),
		p.Get("data").Str,
	)
	return resultOK(map[string]interface{}{})
}

func cq2xqSendXml(bot int64, p gjson.Result) Result {
	core.SendXML(
		bot,
		1,
		cq2xqMessageType(p),
		p.Get("group_id").Int(),
		p.Get("user_id").Int(),
		p.Get("data").Str,
		0,
	)
	return resultOK(map[string]interface{}{})
}
