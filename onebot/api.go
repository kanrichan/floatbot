package onebot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/tidwall/gjson"

	"yaya/core"
)

/*
api请求的json对象（发送私聊信息）：
{
    "action": "send_private_msg",
    "params": {
        "user_id": 10001000,
        "message": "你好"
    },
    "echo": "123"
}

下面api函数中的 (params gjson.value) 为 params 的 gjson.value 对象
/* 例：
{
    "user_id": 10001000,
    "message": "你好"
},
*/

// apiMap XQApiMap，通过函数名获取函数
var apiMap ApiMap

// Routers XQApi路由
type Routers struct {
}

// ApiMap XQApi的name与method对应表
type ApiMap struct {
	this     Routers
	name     []string
	function []func(this *Routers, bot *BotYaml, params gjson.Result) Result
}

// Get 二分法高效获得字符串对应的method
func (apiMap *ApiMap) Get(name string) func(this *Routers, bot *BotYaml, params gjson.Result) Result {
	length := len(apiMap.name)
	low := 0
	high := length - 1
	var fun func(this *Routers, bot *BotYaml, params gjson.Result) Result
	for low <= high {
		mid := (low + high) / 2
		switch {
		default:
			fun = apiMap.function[mid]
			return fun
		case apiMap.name[mid] > name:
			high = mid - 1
		case apiMap.name[mid] < name:
			low = mid + 1
		}
	}
	return nil
}

// Register 反射所有Routers的method并注册
func (apiMap *ApiMap) Register(this *Routers) {
	obj := reflect.TypeOf(this)
	var i int = 0
	for i < obj.NumMethod() {
		apiMap.name = append(apiMap.name, obj.Method(i).Name)
		apiMap.function = append(apiMap.function, obj.Method(i).Func.Interface().(func(this *Routers, bot *BotYaml, params gjson.Result) Result))
		i++
	}
	sort.Sort(apiMap)
}

func (apiMap *ApiMap) Len() int { return len(apiMap.name) }

func (apiMap *ApiMap) Less(i, j int) bool { return apiMap.name[i] < apiMap.name[j] }

func (apiMap *ApiMap) Swap(i, j int) {
	apiMap.name[i], apiMap.name[j] = apiMap.name[j], apiMap.name[i]
	apiMap.function[i], apiMap.function[j] = apiMap.function[j], apiMap.function[i]
}

// CallApi 调用XQApi
func (apiMap *ApiMap) CallApi(action string, bot int64, params gjson.Result) Result {
	name := action2funName(action)
	if apiMap.Get(name) == nil {
		return makeError("no such api")
	}
	botConfig := Conf.getBotConfig(bot)
	return apiMap.Get(name)(&apiMap.this, botConfig, params)
}

// action2funName OneBot的action转驼峰命名
func action2funName(action string) string {
	up := true
	name := ""
	for _, r := range action {
		if up {
			name += strings.ToUpper(string(r))
			up = false
			continue
		}
		if string(r) == "_" {
			up = true
		} else {
			name += string(r)
		}
	}
	return name
}

// Result api请求返回的数据
type Result struct {
	Status  string      `json:"status"`
	Retcode int64       `json:"retcode"`
	Data    interface{} `json:"data"`
	Echo    interface{} `json:"echo"`
}

// makeError 返回api请求错误
func makeError(err string) Result {
	return Result{
		Status:  "failed",
		Retcode: 100,
		Data:    map[string]interface{}{"error": err},
		Echo:    nil,
	}
}

// makeOk 返回api请求正常
func makeOk(data interface{}) Result {
	return Result{
		Status:  "ok",
		Retcode: 0,
		Data:    data,
		Echo:    nil,
	}
}

// SendPrivateMsg 发送私聊消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func (this *Routers) SendPrivateMsg(bot *BotYaml, params gjson.Result) Result {
	return this.SendMsg(bot, params)
}

// SendGroupMsg 发送群消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func (this *Routers) SendGroupMsg(bot *BotYaml, params gjson.Result) Result {
	return this.SendMsg(bot, params)
}

// DeleteMsg 撤回消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
func (this *Routers) DeleteMsg(bot *BotYaml, params gjson.Result) Result {
	var id int64 = params.Get("message_id").Int()
	if id == 0 {
		return makeError("无效'message_id'")
	}
	var xe XEvent
	if bot.DB != nil {
		bot.dbSelect(&xe, "id="+core.Int2Str(id))
	}
	if xe.ID == 0 {
		return makeError("查询无此消息")
	}
	core.WithdrawMsgEX(
		xe.SelfID,
		xe.MessageType,
		xe.GroupID,
		xe.UserID,
		xe.MessageNum,
		xe.MessageID,
		xe.Time,
	)
	return makeOk(nil)
}

// GetMsg 获取消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
func (this *Routers) GetMsg(bot *BotYaml, params gjson.Result) Result {
	var id int64 = params.Get("message_id").Int()
	if id == 0 {
		return makeError("无效'message_id'")
	}
	var xe XEvent
	if bot.DB != nil {
		bot.dbSelect(&xe, "id="+core.Int2Str(id))
	}
	if xe.ID == 0 {
		return makeError("查询无此消息")
	}
	return makeOk(map[string]interface{}{
		"time":         xe.Time,
		"message_type": xq2cqMsgType(xe.MessageType),
		"message_id":   xe.ID,
		"real_id":      xe.MessageID,
		"sender": Event{
			"user_id":  xe.UserID,
			"nickname": "unknown",
			"sex":      "unknown",
			"age":      0,
			"area":     "",
			"card":     "",
			"level":    "",
			"role":     "unknown",
			"title":    "unknown",
		},
		"message": xe.Message,
	})
}

// GetForwardMsg 获取合并转发消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func (this *Routers) GetForwardMsg(bot *BotYaml, params gjson.Result) Result {
	return makeError("先驱不支持")
}

// SendLike 发送好友赞
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_like-%E5%8F%91%E9%80%81%E5%A5%BD%E5%8F%8B%E8%B5%9E
func (this *Routers) SendLike(bot *BotYaml, params gjson.Result) Result {
	var userID int64 = params.Get("user_id").Int()
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	core.UpVote(
		bot.Bot,
		userID,
	)
	return makeOk(nil)
}

// SetGroupKick 群组踢人
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func (this *Routers) SetGroupKick(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	var rejectAddRequest bool = params.Get("reject_add_request").Bool()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	core.KickGroupMBR(
		bot.Bot,
		groupID,
		userID,
		rejectAddRequest,
	)
	return makeOk(nil)
}

// SetGroupBan 群组单人禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func (this *Routers) SetGroupBan(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	var duration int64 = params.Get("duration").Int()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	core.ShutUP(
		bot.Bot,
		groupID,
		userID,
		duration,
	)
	return makeOk(nil)
}

// SetGroupAnonymousBan 群组匿名用户禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous_ban-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D%E7%94%A8%E6%88%B7%E7%A6%81%E8%A8%80
func (this *Routers) SetGroupAnonymousBan(bot *BotYaml, params gjson.Result) Result {
	return makeError("先驱不支持")
}

// SetGroupWholeBan 群组全员禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func (this *Routers) SetGroupWholeBan(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var enable bool = params.Get("enable").Bool()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if enable {
		core.ShutUP(
			bot.Bot,
			groupID,
			0,
			1,
		)
	} else {
		core.ShutUP(
			bot.Bot,
			groupID,
			0,
			0,
		)
	}
	return makeOk(nil)
}

// SetGroupAdmin 群组设置管理员
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_admin-%E7%BE%A4%E7%BB%84%E8%AE%BE%E7%BD%AE%E7%AE%A1%E7%90%86%E5%91%98
func (this *Routers) SetGroupAdmin(bot *BotYaml, params gjson.Result) Result {
	return makeError("先驱不支持")
}

// SetGroupAnonymous 群组匿名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func (this *Routers) SetGroupAnonymous(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var enable bool = params.Get("enable").Bool()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	core.SetAnon(
		bot.Bot,
		groupID,
		enable,
	)
	return makeOk(nil)
}

// SetGroupCard 设置群名片（群备注）
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func (this *Routers) SetGroupCard(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	var card string = params.Get("enable").Str
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	core.SetGroupCard(
		bot.Bot,
		groupID,
		userID,
		card,
	)
	return makeOk(nil)
}

// SetGroupName 设置群名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func (this *Routers) SetGroupName(bot *BotYaml, params gjson.Result) Result {
	return makeError("先驱不支持")
}

// SetGroupLeave 退出群组
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func (this *Routers) SetGroupLeave(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	core.QuitGroup(
		bot.Bot,
		groupID,
	)
	return makeOk(nil)
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func (this *Routers) SetGroupSpecialTitle(bot *BotYaml, params gjson.Result) Result {
	return makeError("先驱不支持")
}

// SetFriendAddRequest 处理加好友请求
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func (this *Routers) SetFriendAddRequest(bot *BotYaml, params gjson.Result) Result {
	var flag int64 = params.Get("flag").Int()
	var approve bool = params.Get("approve").Bool()
	var remark string = params.Get("remark").Str
	if flag == 0 {
		return makeError("无效'flag'")
	}
	if approve {
		core.HandleFriendEvent(
			bot.Bot,
			flag,
			10,
			remark,
		)
	} else {
		core.HandleFriendEvent(
			bot.Bot,
			flag,
			20,
			remark,
		)
	}
	return makeOk(nil)
}

// SetGroupAddRequest 处理加群请求／邀请
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func (this *Routers) SetGroupAddRequest(bot *BotYaml, params gjson.Result) Result {
	flag := params.Get("flag").Str
	if flag == "" {
		return makeError("无效'flag'")
	}
	var approve int64
	if params.Get("approve").Bool() {
		approve = 10
	} else {
		approve = 20
	}
	reason := params.Get("reason").Str

	split := strings.Split(flag, "|")
	core.HandleGroupEvent(bot.Bot,
		213,
		params.Get("user_id").Int(),
		core.Str2Int(split[1]),
		core.Str2Int(split[2]),
		approve,
		reason,
	)
	return makeOk(nil)
}

// GetLoginInfo 获取登录号信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func (this *Routers) GetLoginInfo(bot *BotYaml, params gjson.Result) Result {
	nickname := strings.Split(core.GetNick(
		bot.Bot,
		bot.Bot,
	), "\n")[0]
	return makeOk(map[string]interface{}{
		"user_id":  bot.Bot,
		"nickname": nickname,
	})
}

// GetStrangerInfo 获取陌生人信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func (this *Routers) GetStrangerInfo(bot *BotYaml, params gjson.Result) Result {
	var userID int64 = params.Get("user_id").Int()
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	var nickname string = core.GetNick(
		bot.Bot,
		userID,
	)
	var sex string = xq2cqSex(
		core.GetGender(
			bot.Bot,
			userID,
		),
	)
	var age int64 = core.GetAge(
		bot.Bot,
		userID,
	)
	return makeOk(map[string]interface{}{
		"user_id":  userID,
		"nickname": nickname,
		"sex":      sex,
		"age":      age,
	})
}

// GetFriendList 获取好友列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func (this *Routers) GetFriendList(bot *BotYaml, params gjson.Result) Result {
	var list string = core.GetFriendList(bot.Bot)
	if list == "" {
		return makeError("获取好友列表失败")
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
	return makeOk(friendList)
}

// GetGroupInfo 获取群信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
func (this *Routers) GetGroupInfo(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	var name string = core.GetGroupName(
		bot.Bot,
		groupID,
	)
	members := strings.Split(core.GetGroupMemberNum(
		bot.Bot,
		groupID,
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
	return makeOk(map[string]interface{}{
		"group_id":         groupID,
		"group_name":       name,
		"member_count":     count,
		"max_member_count": max,
	})
}

// GetGroupList 获取群列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func (this *Routers) GetGroupList(bot *BotYaml, params gjson.Result) Result {
	list := core.GetGroupList(bot.Bot)
	if list == "" {
		return makeError("获取群列表失败")
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
	return makeOk(groupList)
}

// GetGroupMemberInfo 获取群成员信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
func (this *Routers) GetGroupMemberInfo(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	return makeOk(map[string]interface{}{
		"group_id":          groupID,
		"user_id":           userID,
		"nickname":          core.GetNick(bot.Bot, userID),
		"card":              core.GetNick(bot.Bot, userID),
		"sex":               []string{"unknown", "male", "female"}[core.GetGender(bot.Bot, userID)],
		"age":               core.GetAge(bot.Bot, userID),
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

// GetGroupMemberList 获取群成员列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
func (this *Routers) GetGroupMemberList(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	list := core.GetGroupMemberList_C(
		bot.Bot,
		groupID,
	)
	if list == "" {
		return makeError("获取群员列表失败")
	}
	g := gjson.Parse(list)
	memberList := []map[string]interface{}{}
	for _, o := range g.Get("list").Array() {
		member := map[string]interface{}{
			"group_id":          groupID,
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
	return makeOk(memberList)
}

// GetGroupHonorInfo 获取群荣誉信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func (this *Routers) GetGroupHonorInfo(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var type_ string = params.Get("message_type").Str
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	cookie := fmt.Sprintf("%s%s", core.GetCookies(bot.Bot), core.GetGroupPsKey(bot.Bot))
	var honorType int64 = 1
	switch type_ {
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
		return makeOk(ret)
	} else {
		return makeError("error")
	}
}

// GetCookies 获取 Cookies
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_cookies-%E8%8E%B7%E5%8F%96-cookies
func (this *Routers) GetCookies(bot *BotYaml, params gjson.Result) Result {
	var domain string = params.Get("domain").Str
	switch domain {
	case "qun.qq.com":
		return makeOk(map[string]interface{}{"cookies": core.GetCookies(bot.Bot) + core.GetGroupPsKey(bot.Bot)})
	case "qzone.qq.com":
		return makeOk(map[string]interface{}{"cookies": core.GetCookies(bot.Bot) + core.GetZonePsKey(bot.Bot)})
	default:
		return makeOk(map[string]interface{}{"cookies": core.GetCookies(bot.Bot)})
	}
}

// GetCsrfToken 获取 CSRF Token
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_csrf_token-%E8%8E%B7%E5%8F%96-csrf-token
func (this *Routers) GetCsrfToken(bot *BotYaml, params gjson.Result) Result {
	return makeError("暂未实现")
}

// GetCredentials 获取 QQ 相关接口凭证
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_credentials-%E8%8E%B7%E5%8F%96-qq-%E7%9B%B8%E5%85%B3%E6%8E%A5%E5%8F%A3%E5%87%AD%E8%AF%81
func (this *Routers) GetCredentials(bot *BotYaml, params gjson.Result) Result {
	var domain string = params.Get("domain").Str
	switch domain {
	case "qun.qq.com":
		return makeOk(map[string]interface{}{"cookies": core.GetCookies(bot.Bot) + core.GetGroupPsKey(bot.Bot)})
	case "qzone.qq.com":
		return makeOk(map[string]interface{}{"cookies": core.GetCookies(bot.Bot) + core.GetZonePsKey(bot.Bot)})
	default:
		return makeOk(map[string]interface{}{"cookies": core.GetCookies(bot.Bot)})
	}
}

// GetRecord 获取语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func (this *Routers) GetRecord(bot *BotYaml, params gjson.Result) Result {
	return makeError("暂未实现")
}

// GetImage 获取图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func (this *Routers) GetImage(bot *BotYaml, params gjson.Result) Result {
	return makeError("暂未实现")
}

// CanSendImage 检查是否可以发送图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#can_send_image-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E5%9B%BE%E7%89%87
func (this *Routers) CanSendImage(bot *BotYaml, params gjson.Result) Result {
	return makeOk(map[string]interface{}{"yes": true})
}

// CanSendRecord 检查是否可以发送语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#can_send_record-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E8%AF%AD%E9%9F%B3
func (this *Routers) CanSendRecord(bot *BotYaml, params gjson.Result) Result {
	return makeOk(map[string]interface{}{"yes": true})
}

// GetStatus 获取运行状态
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func (this *Routers) GetStatus(bot *BotYaml, params gjson.Result) Result {
	return makeOk(map[string]interface{}{
		"online": core.IsOnline(bot.Bot, bot.Bot),
		"good":   true,
	})
}

// GetVersionInfo 获取版本信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_version_info-%E8%8E%B7%E5%8F%96%E7%89%88%E6%9C%AC%E4%BF%A1%E6%81%AF
func (this *Routers) GetVersionInfo(bot *BotYaml, params gjson.Result) Result {
	return makeOk(map[string]interface{}{
		"app_name":         "OneBot-YaYa",
		"app_version":      gjson.Parse(AppInfoJson).Get("pver"),
		"protocol_version": "v11",
	})
}

// SetRestart 重启 OneBot 实现
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_restart-%E9%87%8D%E5%90%AF-onebot-%E5%AE%9E%E7%8E%B0
func (this *Routers) SetRestart(bot *BotYaml, params gjson.Result) Result {
	return makeError("暂未实现")
}

// CleanCache 清理缓存
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#clean_cache-%E6%B8%85%E7%90%86%E7%BC%93%E5%AD%98
func (this *Routers) CleanCache(bot *BotYaml, params gjson.Result) Result {
	return makeError("暂未实现")
}

// OutPutLog 向XQ框架发送日志记录
// ex out_put_log(text=xxx)
func (this *Routers) OutPutLog(bot *BotYaml, params gjson.Result) Result {
	var text string = params.Get("text").Str
	core.OutPutLog(text)
	return makeOk(nil)
}

// SendXml 发送xml信息，相比cqcode方式不用处理转义
// ex: send_xml(group_id=xxx,data=xxx)
func (this *Routers) SendXml(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	var type_ string = params.Get("message_type").Str
	var data string = params.Get("data").Str
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	if groupID == 0 && userID == 0 {
		return makeError("无效'group_id'或'user_id'")
	}
	if type_ == "" {
		if groupID != 0 {
			type_ = "group"
		} else {
			type_ = "private"
		}
	}
	core.SendXML(
		bot.Bot,
		1,
		cq2xqMsgType(type_),
		groupID,
		userID,
		data,
		0,
	)
	return makeOk(map[string]interface{}{})
}

// SendJson 发送json信息，相比cqcode方式不用处理转义
// ex: send_json(group_id=xxx,data=xxx)
func (this *Routers) SendJson(bot *BotYaml, params gjson.Result) Result {
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	var type_ string = params.Get("message_type").Str
	var data string = params.Get("data").Str
	if groupID == 0 {
		return makeError("无效'group_id'")
	}
	if userID == 0 {
		return makeError("无效'user_id'")
	}
	if groupID == 0 && userID == 0 {
		return makeError("无效'group_id'或'user_id'")
	}
	if type_ == "" {
		if groupID != 0 {
			type_ = "group"
		} else {
			type_ = "private"
		}
	}
	core.SendJSON(
		bot.Bot,
		1,
		cq2xqMsgType(type_),
		groupID,
		userID,
		data,
	)
	return makeOk(map[string]interface{}{})
}
