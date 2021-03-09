package xianqu

//#include <api.h>
import "C"

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Yiwen-Chan/go-silk/silk"
	"github.com/tidwall/gjson"
)

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
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func ApiSendMsg(ctx *Context) {
	var (
		params = Parse(ctx.Request).Get("params")
		// 先计算好发送信息类型以便更新group_id
		type_ = ctx.XQMessageType()

		out       string = ""
		bubble    int64  = 0
		anonymous bool   = false
	)
	if !params.Exist("message") {
		return
	}
	messages := params.Array("message")
	for i := range messages {
		message := messages[i]
		data := message.Get("data")
		switch message.Str("type") {
		default:
			//
		case "text":
			out += data.Str("text")
		case "at":
			out += fmt.Sprintf("[@%s]", data.Str("qq"))
		case "face":
			out += fmt.Sprintf("[Face%s.gif]", data.Str("id"))
		case "emoji":
			out += fmt.Sprintf("[emoji=%s]", data.Str("id"))
		case "rps":
			out += []string{"[魔法猜拳] 石头", "[魔法猜拳] 剪刀", "[魔法猜拳] 布"}[rand.Intn(3)]
		case "dice":
			out += []string{"[魔法骰子] 1", "[魔法骰子] 2", "[魔法骰子] 3", "[魔法骰子] 4", "[魔法骰子] 5", "[魔法骰子] 6"}[rand.Intn(6)]
		case "bubble":
			bubble = data.Int("id")
		// 媒体
		case "image":
			// 是否使用缓存
			cache := true
			if data.Exist("cache") {
				cache = data.Bool("cache")
			}
			// OneBot标准的字段
			file := data.Str("file")
			url := data.Str("url")
			if url != "" {
				file = url
			}
			var (
				path string
				res  string
			)
			// 判断file字段为哪种类型
			switch {
			// 链接方式
			case strings.Contains(file, "http://") || strings.Contains(file, "https://"):
				// 解决tx图片链接不落地
				temp := PicPoolCache.Search(hashText(file))
				if temp != nil && cache {
					res = temp.(string)
					break // 存在tx图片链接并且使用缓存
				}
				path = OneBotPath + "image\\" + hashText(file) + ".jpg"
				if PathExists(path) && cache {
					break // 存在链接图片
				}
				// 下载图片
				client := &http.Client{}
				reqest, _ := http.NewRequest("GET", file, nil)
				reqest.Header.Set("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
				reqest.Header.Set("Net-Type", "Wifi")
				resp, err := client.Do(reqest)
				if err != nil {
					panic(err)
				}
				data, _ := ioutil.ReadAll(resp.Body)
				f, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				f.Write(data)
				f.Close()
				resp.Body.Close()
			// base64方式
			case strings.Contains(file, "base64://"):
				path = OneBotPath + "image\\" + hashText(file[9:]) + ".jpg"
				data, err := base64.StdEncoding.DecodeString(file[9:])
				if err != nil {
					continue
				}
				f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				if err != nil {
					continue
				}
				f.Write(data)
				f.Close()
			// 本地文件方式
			case strings.Contains(file, "file:///"):
				path = file[8:]
			// 默认方式
			default:
				path = file
			}
			if res == "" {
				// 用文件md5判断缓冲池是否存在该图片
				temp := PicPoolCache.Search(FileMD5(path))
				if temp != nil && cache {
					res = temp.(string)
				} else {
					res = path
				}
			}
			switch data.Str("type") {
			default:
				out += fmt.Sprintf("[pic=%s]", res)
			case "show":
				out += fmt.Sprintf("[ShowPic=%s,type=%d]", res, data.Int("id")+40000)
			}
		case "record":
			// 是否使用缓存
			cache := true
			if data.Exist("cache") {
				cache = data.Bool("cache")
			}
			// OneBot标准的字段
			file := data.Str("file")
			// 判断file字段为哪种类型
			switch {
			// 链接方式
			case strings.Contains(file, "http://") || strings.Contains(file, "https://"):
				file = OneBotPath + "image\\" + hashText(file) + ".mp3"
				if PathExists(file) && cache {
					break // 存在链接音频
				}
				// 下载音频
				client := &http.Client{}
				reqest, _ := http.NewRequest("GET", file, nil)
				reqest.Header.Set("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
				reqest.Header.Set("Net-Type", "Wifi")
				link, _ := url.Parse(file)
				reqest.Header.Set("Host", link.Hostname())
				resp, err := client.Do(reqest)
				if err != nil {
					panic(err)
				}
				data, _ := ioutil.ReadAll(resp.Body)
				f, _ := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				f.Write(data)
				f.Close()
				resp.Body.Close()
			// base64方式
			case strings.Contains(file, "base64://"):
				file = OneBotPath + "image\\" + hashText(file[9:]) + ".jpg"
				data, err := base64.StdEncoding.DecodeString(file[9:])
				if err != nil {
					continue
				}
				f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
				if err != nil {
					continue
				}
				f.Write(data)
				f.Close()
			// 本地文件方式
			case strings.Contains(file, "file:///"):
				file = file[8:]
			// 默认方式
			default:
			}
			// 获取本地文件的文件名
			name := strings.ReplaceAll(filepath.Base(file), path.Ext(file), "")
			silkEncoder := &silk.Encoder{}
			err := silkEncoder.Init("OneBot/record", "OneBot/codec")
			if err != nil {
				continue
			}
			data, err := ioutil.ReadFile(file + ".mp3")
			if err != nil {
				continue
			}
			_, err = silkEncoder.EncodeToSilk(data, name, true)
			if err != nil {
				continue
			}
			res := OneBotPath + "record/" + name + ".silk"
			out += fmt.Sprintf("[Voi=%s]", res)
		case "video":
			// TODO
		// 富文本
		case "xml":
			C.S3_Api_SendXML(
				GoInt2CStr(ctx.Bot),
				C.int(1),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(data.Str("data")),
				0,
			)
		case "json":
			C.S3_Api_SendJSON(
				GoInt2CStr(ctx.Bot),
				C.int(1),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(data.Str("data")),
			)
		case "share":
			C.S3_Api_SendXML(
				GoInt2CStr(ctx.Bot),
				C.int(1),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
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
				)),
				0,
			)
		case "music":
			C.S3_Api_SendXML(
				GoInt2CStr(ctx.Bot),
				C.int(1),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
					<msg serviceID="2" templateID="1" action="web" brief="[分享] %s" 
					sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0">
					<item layout="2"><audio cover="%s" src="%s"/><title>%s</title>
					<summary>%s</summary></item><source name="音乐" 
					icon="https://i.gtimg.cn/open/app_icon/01/07/98/56/1101079856_100_m.png" 
					url="http://web.p.qq.com/qqmpmobile/aio/app.html?id=1101079856" 
					action="app" a_actionData="com.tencent.qqmusic" 
					i_actionData="tencent1101079856://" appid="1101079856" /></msg>`,
					XmlEscape(data.Str("title")),
					data.Str("url"),
					data.Str("image"),
					data.Str("audio"),
					XmlEscape(data.Str("title")),
					XmlEscape(data.Str("content")),
				)),
				0,
			)
		case "weather":
			C.S3_Api_SendJSON(
				GoInt2CStr(ctx.Bot),
				C.int(1),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(`{"app":"com.tencent.weather","desc":"天气",
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
				)),
			)
		case "contact":
			switch data.Str("type") {
			case "qq":
				C.S3_Api_SendXML(
					GoInt2CStr(ctx.Bot),
					C.int(1),
					C.int(type_),
					GoInt2CStr(params.Int("group_id")),
					GoInt2CStr(params.Int("user_id")),
					CString(fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
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
					)),
					0,
				)
			case "group":
				C.S3_Api_SendXML(
					GoInt2CStr(ctx.Bot),
					C.int(1),
					C.int(type_),
					GoInt2CStr(params.Int("group_id")),
					GoInt2CStr(params.Int("user_id")),
					CString(fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
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
					)),
					0,
				)
			}
		case "location":
			C.S3_Api_SendJSON(
				GoInt2CStr(ctx.Bot),
				C.int(1),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(`{"app":"com.tencent.map","desc":"","view":"Share",
					"ver":"0.0.0.1","prompt":"[应用]地图","appID":"","sourceName":"",
					"actionData":"","actionData_A":"","sourceUrl":"","meta":{"Share":{"locSub":"%s",
					"lng":%s,"lat":%s,"zoom":15,"locName":"%s"}},
					"config":{"forward":true,"autosize":1},"text":"","extraApps":[],
					"sourceAd":"","extra":""}`,
					data.Str("content"),
					data.Str("lon"),
					data.Str("lat"),
					data.Str("title"),
				)),
			)
		// 其他
		case "shake":
			C.S3_Api_ShakeWindow(
				GoInt2CStr(ctx.Bot),
				GoInt2CStr(params.Int("user_id")),
			)
		case "poke":
			C.S3_Api_SendMsgEX_V2(
				GoInt2CStr(ctx.Bot),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(
					"[系统提示] %v 戳了一下 %v",
					ctx.Bot,
					params.Int("user_id"),
				)),
				C.int(0),
				CBool(false),
				CString(""),
			)
		case "anonymous":
			anonymous = true
		case "reply":
			C.S3_Api_SendMsgEX_V2(
				GoInt2CStr(ctx.Bot),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(
					"[系统消息] %v 尝试回复一条消息并失败了",
					ctx.Bot,
				)),
				C.int(0),
				CBool(false),
				CString(""),
			)
		case "forward":
			C.S3_Api_SendMsgEX_V2(
				GoInt2CStr(ctx.Bot),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(
					"[系统消息] %v 尝试合并转发一条消息并失败了",
					ctx.Bot,
				)),
				C.int(0),
				CBool(false),
				CString(""),
			)
		case "node":
			C.S3_Api_SendMsgEX_V2(
				GoInt2CStr(ctx.Bot),
				C.int(type_),
				GoInt2CStr(params.Int("group_id")),
				GoInt2CStr(params.Int("user_id")),
				CString(fmt.Sprintf(
					"[系统消息] %v 尝试合并转发节点并失败了",
					ctx.Bot,
				)),
				C.int(0),
				CBool(false),
				CString(""),
			)
		}
	}
	if out == "" {
		// Xml json 信息无法返回有效的id
		ctx.MakeOkResponse(map[string]interface{}{"message_id": 0})
		return
	}
	ret := CPtr2GoStr(
		C.S3_Api_SendMsgEX_V2(
			GoInt2CStr(ctx.Bot),
			C.int(type_),
			GoInt2CStr(params.Int("group_id")),
			GoInt2CStr(params.Int("user_id")),
			CString(EscapeEmoji(out)),
			C.int(bubble),
			CBool(anonymous),
			CString(""),
		),
	)
	// 处理返回的 message_id
	var temp map[string]interface{}
	json.Unmarshal([]byte(ret), &temp)
	if !Parse(temp).Bool("sendok") {
		ctx.MakeFailResponse("可能受到风控")
		return
	}
	num := Parse(temp).Int("msgno")
	if num == 0 {
		ctx.MakeFailResponse("可能受到风控")
		return
	}
	id := MessageIDCache.Hcraes(num)
	if id == nil {
		ctx.MakeFailResponse("可能受到风控")
		return
	}
	ctx.MakeOkResponse(map[string]interface{}{"message_id": id.(int64)})
}

// SendPrivateMsg 发送私聊消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_private_msg-%E5%8F%91%E9%80%81%E7%A7%81%E8%81%8A%E6%B6%88%E6%81%AF
func ApiSendPrivateMsg(ctx *Context) {
	ApiSendMsg(ctx)
}

// SendGroupMsg 发送群消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_group_msg-%E5%8F%91%E9%80%81%E7%BE%A4%E6%B6%88%E6%81%AF
func ApiSendGroupMsg(ctx *Context) {
	ApiSendMsg(ctx)
}

// DeleteMsg 撤回消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#delete_msg-%E6%92%A4%E5%9B%9E%E6%B6%88%E6%81%AF
func ApiDeleteMsg(ctx *Context) {
	var (
		id    = Parse(ctx.Request).Get("params").Int("message_id")
		num   = MessageIDCache.Search(id).(int64)
		type_ = ctx.XQMessageType()
	)
	C.S3_Api_WithdrawMsgEX(
		GoInt2CStr(ctx.Bot),
		C.int(type_),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		GoInt2CStr(num),
		GoInt2CStr(id),
		GoInt2CStr(0),
	)
	ctx.MakeOkResponse(nil)
}

// GetMsg 获取消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_msg-%E8%8E%B7%E5%8F%96%E6%B6%88%E6%81%AF
func ApiGetMsg(ctx *Context) {
	ctx.MakeOkResponse(MessageCache.Search(Parse(ctx.Request).Get("params").Int("message_id")).(map[string]interface{}))
}

// GetForwardMsg 获取合并转发消息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_forward_msg-%E8%8E%B7%E5%8F%96%E5%90%88%E5%B9%B6%E8%BD%AC%E5%8F%91%E6%B6%88%E6%81%AF
func ApiGetForwardMsg(ctx *Context) {
	ctx.MakeFailResponse("先驱不支持")
}

// SendLike 发送好友赞
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#send_like-%E5%8F%91%E9%80%81%E5%A5%BD%E5%8F%8B%E8%B5%9E
func ApiSendLike(ctx *Context) {
	C.S3_Api_UpVote(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupKick 群组踢人
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_kick-%E7%BE%A4%E7%BB%84%E8%B8%A2%E4%BA%BA
func ApiSetGroupKick(ctx *Context) {
	C.S3_Api_KickGroupMBR(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		CBool(Parse(ctx.Request).Get("params").Bool("reject_add_request")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupBan 群组单人禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_ban-%E7%BE%A4%E7%BB%84%E5%8D%95%E4%BA%BA%E7%A6%81%E8%A8%80
func ApiSetGroupBan(ctx *Context) {
	C.S3_Api_ShutUP(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		C.int(Parse(ctx.Request).Get("params").Int("duration")/60),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupAnonymousBan 群组匿名用户禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous_ban-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D%E7%94%A8%E6%88%B7%E7%A6%81%E8%A8%80
func ApiSetGroupAnonymousBan(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// SetGroupWholeBan 群组全员禁言
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_whole_ban-%E7%BE%A4%E7%BB%84%E5%85%A8%E5%91%98%E7%A6%81%E8%A8%80
func ApiSetGroupWholeBan(ctx *Context) {
	var (
		temp         = Parse(ctx.Request).Get("params").Bool("enable")
		enable int64 = 0 // 解除禁言
	)
	if temp {
		enable = 1 // 设置禁言
	}
	C.S3_Api_ShutUP(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		GoInt2CStr(0),
		C.int(enable),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupAdmin 群组设置管理员
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_admin-%E7%BE%A4%E7%BB%84%E8%AE%BE%E7%BD%AE%E7%AE%A1%E7%90%86%E5%91%98
func ApiSetGroupAdmin(ctx *Context) {
	var (
		groupID = Parse(ctx.Request).Get("params").Str("group_id")
		userID  = Parse(ctx.Request).Get("params").Str("user_id")
		enable  = Parse(ctx.Request).Get("params").Int("enable")
	)
	temp1 := CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot)))
	temp2 := CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot)))
	temp3 := []byte{}
	for i := range temp1 {
		if temp1[i] != temp2[i] || temp1[i] == 92 {
			break
		}
		temp3 = append(temp3, temp1[i])
	}
	cookie := string(temp3) + CPtr2GoStr(C.S3_Api_GetGroupPsKey(GoInt2CStr(ctx.Bot)))
	skey := string(temp3)[strings.Index(string(temp3), "skey=")+5:]
	bnk := 5381
	for i := range skey {
		bnk += (bnk << 5) + int(skey[i])
	}
	bnk = bnk & 2147483647
	client := &http.Client{}
	dataUrl := url.Values{}
	dataUrl.Add("gc", groupID)
	dataUrl.Add("op", Int2Str(enable))
	dataUrl.Add("ul", userID)
	dataUrl.Add("bkn", strconv.Itoa(bnk))
	reqest, _ := http.NewRequest("POST", "https://qun.qq.com/cgi-bin/qun_mgr/set_group_admin", strings.NewReader(dataUrl.Encode()))
	reqest.Header.Set("Cookie", cookie)
	resp, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	ul := gjson.ParseBytes(data).Get("ul").Str
	em := gjson.ParseBytes(data).Get("em").Str
	if ul == userID {
		ctx.MakeOkResponse(nil)
	}
	ctx.MakeFailResponse(em)
}

// SetGroupAnonymous 群组匿名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_anonymous-%E7%BE%A4%E7%BB%84%E5%8C%BF%E5%90%8D
func ApiSetGroupAnonymous(ctx *Context) {
	C.S3_Api_SetAnon(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		CBool(Parse(ctx.Request).Get("params").Bool("enable")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupCard 设置群名片（群备注）
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_card-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D%E7%89%87%E7%BE%A4%E5%A4%87%E6%B3%A8
func ApiSetGroupCard(ctx *Context) {
	C.S3_Api_SetGroupCard(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("user_id")),
		CString(Parse(ctx.Request).Get("params").Str("card")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupName 设置群名
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_name-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E5%90%8D
func ApiSetGroupName(ctx *Context) {
	ctx.MakeFailResponse("先驱不支持")
}

// SetGroupLeave 退出群组
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_leave-%E9%80%80%E5%87%BA%E7%BE%A4%E7%BB%84
func ApiSetGroupLeave(ctx *Context) {
	C.S3_Api_QuitGroup(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("group_id")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupSpecialTitle 设置群组专属头衔
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_special_title-%E8%AE%BE%E7%BD%AE%E7%BE%A4%E7%BB%84%E4%B8%93%E5%B1%9E%E5%A4%B4%E8%A1%94
func ApiSetGroupSpecialTitle(ctx *Context) {
	ctx.MakeFailResponse("先驱不支持")
}

// SetFriendAddRequest 处理加好友请求
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_friend_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E5%A5%BD%E5%8F%8B%E8%AF%B7%E6%B1%82
func ApiSetFriendAddRequest(ctx *Context) {
	var (
		temp          = Parse(ctx.Request).Get("params").Bool("approve")
		approve int64 = 20
	)
	if temp {
		approve = 10
	}
	C.S3_Api_HandleFriendEvent(
		GoInt2CStr(ctx.Bot),
		GoInt2CStr(Parse(ctx.Request).Get("params").Int("flag")),
		C.int(approve),
		CString(Parse(ctx.Request).Get("params").Str("remark")),
	)
	ctx.MakeOkResponse(nil)
}

// SetGroupAddRequest 处理加群请求／邀请
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_group_add_request-%E5%A4%84%E7%90%86%E5%8A%A0%E7%BE%A4%E8%AF%B7%E6%B1%82%E9%82%80%E8%AF%B7
func ApiSetGroupAddRequest(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetLoginInfo 获取登录号信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_login_info-%E8%8E%B7%E5%8F%96%E7%99%BB%E5%BD%95%E5%8F%B7%E4%BF%A1%E6%81%AF
func ApiGetLoginInfo(ctx *Context) {
	nickname := strings.Split(
		CPtr2GoStr(
			C.S3_Api_GetNick(
				GoInt2CStr(ctx.Bot),
				GoInt2CStr(ctx.Bot),
			),
		),
		"\n",
	)[0]
	ctx.MakeOkResponse(
		map[string]interface{}{
			"user_id":  ctx.Bot,
			"nickname": nickname,
		},
	)
}

// GetStrangerInfo 获取陌生人信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_stranger_info-%E8%8E%B7%E5%8F%96%E9%99%8C%E7%94%9F%E4%BA%BA%E4%BF%A1%E6%81%AF
func ApiGetStrangerInfo(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetFriendList 获取好友列表
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_friend_list-%E8%8E%B7%E5%8F%96%E5%A5%BD%E5%8F%8B%E5%88%97%E8%A1%A8
func ApiGetFriendList(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetGroupInfo 获取群信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E4%BF%A1%E6%81%AF
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
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E5%88%97%E8%A1%A8
func ApiGetGroupList(ctx *Context) {
	var temp []GroupInfo
	list := strings.Split(
		CPtr2GoStr(
			C.S3_Api_GetGroupList_B(
				GoInt2CStr(ctx.Bot),
			),
		),
		"\r\n",
	)
	for _, groupID := range list {
		temp = append(temp, GroupInfo{
			GroupID: Str2Int(groupID),
		})
	}
	ctx.MakeOkResponse(temp)
}

// GetGroupMemberInfo 获取群成员信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E4%BF%A1%E6%81%AF
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
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_member_list-%E8%8E%B7%E5%8F%96%E7%BE%A4%E6%88%90%E5%91%98%E5%88%97%E8%A1%A8
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
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_group_honor_info-%E8%8E%B7%E5%8F%96%E7%BE%A4%E8%8D%A3%E8%AA%89%E4%BF%A1%E6%81%AF
func ApiGetGroupHonorInfo(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetCookies 获取 Cookies
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_cookies-%E8%8E%B7%E5%8F%96-cookies
func ApiGetCookies(ctx *Context) {
	switch Parse(ctx.Request).Get("params").Str("domain") {
	case "qun.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot))) + CPtr2GoStr(C.S3_Api_GetGroupPsKey(GoInt2CStr(ctx.Bot)))})
		return
	case "qzone.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot))) + CPtr2GoStr(C.S3_Api_GetZonePsKey(GoInt2CStr(ctx.Bot)))})
		return
	default:
		ctx.MakeOkResponse(map[string]interface{}{"cookies": CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot)))})
		return
	}
}

// GetCsrfToken 获取 CSRF Token
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_csrf_token-%E8%8E%B7%E5%8F%96-csrf-token
func ApiGetCsrfToken(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetCredentials 获取 QQ 相关接口凭证
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_credentials-%E8%8E%B7%E5%8F%96-qq-%E7%9B%B8%E5%85%B3%E6%8E%A5%E5%8F%A3%E5%87%AD%E8%AF%81
func ApiGetCredentials(ctx *Context) {
	switch Parse(ctx.Request).Get("params").Str("domain") {
	case "qun.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot))) + CPtr2GoStr(C.S3_Api_GetGroupPsKey(GoInt2CStr(ctx.Bot)))})
		return
	case "qzone.qq.com":
		ctx.MakeOkResponse(map[string]interface{}{"cookies": CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot))) + CPtr2GoStr(C.S3_Api_GetZonePsKey(GoInt2CStr(ctx.Bot)))})
		return
	default:
		ctx.MakeOkResponse(map[string]interface{}{"cookies": CPtr2GoStr(C.S3_Api_GetCookies(GoInt2CStr(ctx.Bot)))})
		return
	}
}

// GetRecord 获取语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_record-%E8%8E%B7%E5%8F%96%E8%AF%AD%E9%9F%B3
func ApiGetRecord(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// GetImage 获取图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_image-%E8%8E%B7%E5%8F%96%E5%9B%BE%E7%89%87
func ApiGetImage(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// CanSendImage 检查是否可以发送图片
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#can_send_image-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E5%9B%BE%E7%89%87
func ApiCanSendImage(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"yes": true,
		},
	)
}

// CanSendRecord 检查是否可以发送语音
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#can_send_record-%E6%A3%80%E6%9F%A5%E6%98%AF%E5%90%A6%E5%8F%AF%E4%BB%A5%E5%8F%91%E9%80%81%E8%AF%AD%E9%9F%B3
func ApiCanSendRecord(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"yes": true,
		},
	)
}

// GetStatus 获取运行状态
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_status-%E8%8E%B7%E5%8F%96%E8%BF%90%E8%A1%8C%E7%8A%B6%E6%80%81
func ApiGetStatus(ctx *Context) {
	ctx.MakeOkResponse(
		map[string]interface{}{
			"online": true,
			"good":   true,
		},
	)
}

// GetVersionInfo 获取版本信息
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#get_version_info-%E8%8E%B7%E5%8F%96%E7%89%88%E6%9C%AC%E4%BF%A1%E6%81%AF
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
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#set_restart-%E9%87%8D%E5%90%AF-onebot-%E5%AE%9E%E7%8E%B0
func ApiSetRestart(ctx *Context) {
	ctx.MakeFailResponse("暂时不支持")
}

// CleanCache 清理缓存
// https://github.com/howmanybots/onebot/blob/master/v11/specs/api/public.md#clean_cache-%E6%B8%85%E7%90%86%E7%BC%93%E5%AD%98
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

func ApiOutPutLog(text interface{}) {
	C.S3_Api_OutPutLog(
		CString(fmt.Sprintln(text)),
	)
}

func ApiOutPutLog1(text interface{}) {
	fmt.Println(text)
}

func XQApiGroupName(bot, groupID int64) string {
	return CPtr2GoStr(
		C.S3_Api_GetGroupName(
			GoInt2CStr(bot),
			GoInt2CStr(groupID),
		),
	)
}

func XQApiGroupMemberListB(bot, groupID int64) string {
	return CPtr2GoStr(
		C.S3_Api_GetGroupMemberList_B(
			GoInt2CStr(bot),
			GoInt2CStr(groupID),
		),
	)
}

func XQApiGroupMemberListC(bot, groupID int64) string {
	return CPtr2GoStr(
		C.S3_Api_GetGroupMemberList_C(
			GoInt2CStr(bot),
			GoInt2CStr(groupID),
		),
	)
}

func XQApiGetNick(bot, userID int64) string {
	return strings.Split(
		CPtr2GoStr(
			C.S3_Api_GetNick(
				GoInt2CStr(bot),
				GoInt2CStr(userID),
			),
		),
		"\n",
	)[0]
}

func XQApiGetAge(bot, userID int64) int64 {
	return int64(
		C.S3_Api_GetAge(
			GoInt2CStr(bot),
			GoInt2CStr(userID),
		),
	)
}

func XQApiGetGender(bot, userID int64) string {
	return []string{"unknown", "male", "female"}[int64(
		C.S3_Api_GetGender(
			GoInt2CStr(bot),
			GoInt2CStr(userID),
		),
	)]
}

func XQApiIsFriend(bot, userID int64) bool {
	return GoBool(
		C.S3_Api_IfFriend(
			GoInt2CStr(bot),
			GoInt2CStr(userID),
		),
	)
}

func ApiCallMessageBox(text string) {
	C.S3_Api_CallMessageBox(
		CString(text),
	)
}

func ApiMessageBoxButton(text string) int64 {
	// 6 为是
	// 7 为否
	return int64(
		C.S3_Api_MessageBoxButton(
			CString(text),
		),
	)
}

func ApiDefaultQQ() int64 {
	botList := strings.Split(
		GoString(C.S3_Api_GetQQList()),
		"/n",
	)
	if len(botList) < 0 {
		return 0
	}
	return Str2Int(botList[0])
}
