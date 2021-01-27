package onebot

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/tidwall/gjson"

	"yaya/core"
)

var Split bool

type msgTarget struct {
	BotID   int64
	Type_   int64
	GroupID int64
	UserID  int64
}

func xq2cqMsgType(type_ int64) string {
	switch type_ {
	default:
		return "friend"
	case 1:
		return "friend"
	case 2:
		return "group"
	}
}

func cq2xqMsgType(type_ string) int64 {
	switch type_ {
	default:
		return 1
	case "friend":
		return 1
	case "group":
		return 2
	}
}

func xq2cqSex(sex int64) string {
	switch sex {
	default:
		return "unknown"
	case 1:
		return "male"
	case 2:
		return "female"
	}
}

func (this *Routers) SendMsg(bot *BotYaml, params gjson.Result) Result {
	var type_ string = params.Get("message_type").Str
	var groupID int64 = params.Get("group_id").Int()
	var userID int64 = params.Get("user_id").Int()
	var message = params.Get("message")
	if type_ == "group" && groupID == 0 {
		return makeError("无效'group_id'")
	}
	if type_ == "private" && userID == 0 {
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
	target := msgTarget{
		BotID:   bot.Bot,
		Type_:   cq2xqMsgType(type_),
		GroupID: groupID,
		UserID:  userID,
	}

	var out string = ""
	var bubble int64 = 0
	var msg gjson.Result

	if len(message.Get("#.type").Array()) == 0 {
		b, _ := json.Marshal(cqCode2Array(message.Str))
		msg = gjson.ParseBytes(b)
	} else {
		msg = message
	}

	for _, message := range msg.Array() {
		switch message.Get("type").Str {
		// 文字
		case "text":
			out += target.cq2xqText(message)
		case "at":
			out += target.cq2xqAt(message)
		case "face":
			out += target.cq2xqFace(message)
		case "emoji":
			out += target.cq2xqEmoji(message)
		case "rps":
			out += target.cq2xqRps(message)
		case "dice":
			out += target.cq2xqDice(message)
		case "bubble":
			bubble = message.Get("data.id").Int()
		// 媒体
		case "image":
			out += target.cq2xqImage(message)
		case "record":
			out += target.cq2xqRecord(message)
		case "video":
			out += target.cq2xqVideo(message)
		// 富文本
		case "xml":
			out += target.cq2xqXml(message)
		case "json":
			out += target.cq2xqJson(message)
		case "share":
			out += target.cq2xqShare(message)
		case "music":
			out += target.cq2xqMusic(message)
		case "weather":
			out += target.cq2xqWeather(message)
		case "contact":
			out += target.cq2xqContact(message)
		case "location":
			out += target.cq2xqLocation(message)
		// 其他
		case "shake":
			out += target.cq2xqShake(message)
		case "poke":
			out += target.cq2xqPoke(message)
		case "anonymous":
			out += target.cq2xqAnonymous(message)
		case "reply":
			out += target.cq2xqReply(message)
		case "forward":
			out += target.cq2xqForward(message)
		case "node":
			out += target.cq2xqNode(message)
		default:
			out += target.cq2xqDefault(message)
		}
	}
	if out != "" {
		// 如果开了分片就切割信息
		if Conf.Cache.Video {
			out = messageSplit(out)
			// 调用core发送信息
			for _, o := range strings.Split(out, "[Next]") {
				core.SendMsgEX_V2(
					target.BotID,
					target.Type_,
					target.GroupID,
					target.UserID,
					o,
					bubble,
					false,
					"",
				)
				if strings.Contains(o, "[pic") {
					time.Sleep(time.Millisecond * 1000)
				} else {
					time.Sleep(time.Millisecond * 200)
				}
			}
			return makeOk(map[string]interface{}{"message_id": 0})
		}
		data := core.SendMsgEX_V2(
			target.BotID,
			target.Type_,
			target.GroupID,
			target.UserID,
			out,
			bubble,
			false,
			"",
		)
		// 获取ID返回
		if data != "" {
			ret := gjson.Parse(data[:strings.LastIndex(data, "}")])
			if ret.Get("sendok").Bool() {
				var xe XEvent
				xe.MessageID = ret.Get("msgid").Int()
				xe.MessageNum = ret.Get("msgno").Int()
				xe.ID = 0
				for i := range Conf.BotConfs {
					if bot.Bot == Conf.BotConfs[i].Bot && bot.Bot != 0 && Conf.BotConfs[i].DB != nil {
						time.Sleep(time.Millisecond * 100)
						Conf.BotConfs[i].dbSelect(&xe, "message_num="+core.Int2Str(xe.MessageNum))
						return makeOk(map[string]interface{}{"message_id": xe.ID})
					}
				}
			}
		}
		return makeError("可能受到风控")
	}
	return makeOk(map[string]interface{}{"message_id": 0})
}

func messageSplit(texts string) string {
	var (
		send  string   = ""
		split []string = strings.Split(texts, "\n")
	)
	for i, text := range split {
		if i+1 == len(split) {
			send = send + text
			break
		}
		now := strings.Split(send, "[Next]")
		if strings.Contains(text, "[pic") && strings.Contains(text, "]") {
			send = send + text + "[Next]"
		} else if len(now[len(now)-1])+len(text) > 120 && len(split[i+1]) > 60 {
			send = send + text + "[Next]"
		} else if len(now[len(now)-1])+len(text) > 180 {
			send = send + text + "\n"
		} else {
			send = send + text + "\n"
		}
	}
	return send
}

func (target msgTarget) cq2xqText(message gjson.Result) string {
	return message.Get("data.*").Str
}

func (target msgTarget) cq2xqFace(message gjson.Result) string {
	return fmt.Sprintf("[Face%s.gif]", message.Get("data.*").Str)
}

func (target msgTarget) cq2xqAt(message gjson.Result) string {
	return fmt.Sprintf("[@%s] ", message.Get("data.*").Str)
}

func (target msgTarget) cq2xqEmoji(message gjson.Result) string {
	return fmt.Sprintf("[emoji=%s]", message.Get("data.*").Str)
}

func (target msgTarget) cq2xqRps(message gjson.Result) string {
	return []string{
		"[魔法猜拳] 石头",
		"[魔法猜拳] 剪刀",
		"[魔法猜拳] 布",
	}[rand.Intn(3)]
}

func (target msgTarget) cq2xqDice(message gjson.Result) string {
	return []string{
		"[魔法骰子] 1",
		"[魔法骰子] 2",
		"[魔法骰子] 3",
		"[魔法骰子] 4",
		"[魔法骰子] 5",
		"[魔法骰子] 6",
	}[rand.Intn(6)]
}

func (target msgTarget) cq2xqImage(message gjson.Result) string {
	url := strings.ReplaceAll(message.Get("data.url").Str, `\/`, `/`)
	file := strings.ReplaceAll(message.Get("data.file").Str, `\/`, `/`)
	cache := true
	if message.Get("data.cache").Exists() {
		cache = message.Get("data.cache").Bool()
	}
	showID := message.Get("data.id").Int() - 40000
	pic := picDownloader{
		file:     file,
		url:      url,
		suffix:   ".jpg",
		savePath: ImagePath,
		iscache:  cache,
	}
	if message.Get("data.type").Str == "show" {
		return fmt.Sprintf("[ShowPic=%s,type=%d]", pic.path(), showID)
	} else {
		return fmt.Sprintf("[pic=%s]", pic.path())
	}
}

func (target msgTarget) cq2xqRecord(message gjson.Result) string {
	url := strings.ReplaceAll(message.Get("data.url").Str, `\/`, `/`)
	file := strings.ReplaceAll(message.Get("data.file").Str, `\/`, `/`)
	cache := true
	if message.Get("data.cache").Exists() {
		cache = message.Get("data.cache").Bool()
	}
	rec := picDownloader{
		file:     file,
		url:      url,
		suffix:   ".mp3",
		savePath: RecordPath,
		iscache:  cache,
	}
	return fmt.Sprintf("[pic=%s]", rec.path())
}

func (target msgTarget) cq2xqVideo(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		message.Get("data.*").Str,
		0,
		false,
		" ",
	)
	return ""
}

func (target msgTarget) cq2xqMusic(message gjson.Result) string {
	switch {
	case message.Get("data.type").Str == "custom":
		core.SendXML(
			target.BotID,
			1,
			target.Type_,
			target.GroupID,
			target.UserID,
			fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
				<msg serviceID="2" templateID="1" action="web" brief="[分享] %s" 
				sourceMsgId="0" url="%s" flag="0" adverSign="0" multiMsgFlag="0">
				<item layout="2"><audio cover="%s" src="%s"/><title>%s</title>
				<summary>%s</summary></item><source name="音乐" 
				icon="https://i.gtimg.cn/open/app_icon/01/07/98/56/1101079856_100_m.png" 
				url="http://web.p.qq.com/qqmpmobile/aio/app.html?id=1101079856" 
				action="app" a_actionData="com.tencent.qqmusic" 
				i_actionData="tencent1101079856://" appid="1101079856" /></msg>`,
				XmlEscape(message.Get("data.title").Str),
				message.Get("data.url").Str,
				message.Get("data.image").Str,
				message.Get("data.audio").Str,
				XmlEscape(message.Get("data.title").Str),
				XmlEscape(message.Get("data.content").Str),
			),
			0,
		)
	default:
		DEBUG("[CQ码解析] %v 暂未实现", message.Str)
		core.SendMsgEX_V2(
			target.BotID,
			target.Type_,
			target.GroupID,
			target.UserID,
			fmt.Sprintf("音乐分享：%s %s",
				message.Get("data.type").Str,
				message.Get("data.id").Str,
			),
			0,
			false,
			" ",
		)
	}
	return ""
}

func (target msgTarget) cq2xqWeather(message gjson.Result) string {
	core.SendJSON(
		target.BotID,
		1,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(`{"app":"com.tencent.weather","desc":"天气",
			"view":"RichInfoView","ver":"0.0.0.1","prompt":"[应用]天气",
			"appID":"","sourceName":"","actionData":"","actionData_A":"",
			"sourceUrl":"","meta":{"richinfo":{"adcode":"","air":"%s",
			"city":"%s","date":"%s","max":"%s","min":"%s",
			"ts":"15158613","type":"%s","wind":""}},"text":"","sourceAd":"","extra":""}`,
			message.Get("data.air").Str,
			message.Get("data.city").Str,
			message.Get("data.date").Str,
			message.Get("data.max").Str,
			message.Get("data.min").Str,
			message.Get("data.type").Str,
		),
	)
	return ""
}

func (target msgTarget) cq2xqXml(message gjson.Result) string {
	core.SendXML(
		target.BotID,
		1,
		target.Type_,
		target.GroupID,
		target.UserID,
		message.Get("data.*").Str,
		0,
	)
	return ""
}

func (target msgTarget) cq2xqJson(message gjson.Result) string {
	core.SendJSON(
		target.BotID,
		1,
		target.Type_,
		target.GroupID,
		target.UserID,
		message.Get("data.*").Str,
	)
	return ""
}

func (target msgTarget) cq2xqShare(message gjson.Result) string {
	core.SendXML(
		target.BotID,
		1,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
				<msg serviceID="33" templateID="123" action="web" brief="%s" 
				sourceMsgId="0" url="%s" 
				flag="8" adverSign="0" multiMsgFlag="0"><item layout="2" 
				advertiser_id="0" aid="0"><picture cover="%s" w="0" h="0" />
				<title>%s</title><summary>%s</summary>
				</item><source name="" icon="" action="" appid="-1" /></msg>`,
			message.Get("data.brief").Str,
			message.Get("data.url").Str,
			message.Get("data.image").Str,
			message.Get("data.title").Str,
			message.Get("data.content").Str,
		),
		0,
	)
	return ""
}

func (target msgTarget) cq2xqContact(message gjson.Result) string {
	switch message.Get("data.type").Str {
	case "qq":
		core.SendXML(
			target.BotID,
			1,
			target.Type_,
			target.GroupID,
			target.UserID,
			fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
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
				message.Get("data.id").Str,
				message.Get("data.id").Str,
				message.Get("data.id").Str,
				message.Get("data.name").Str,
				message.Get("data.id").Str,
				message.Get("data.name").Str,
				message.Get("data.id").Str,
			),
			0,
		)
	case "group":
		core.SendXML(
			target.BotID,
			1,
			target.Type_,
			target.GroupID,
			target.UserID,
			fmt.Sprintf(`<?xml version='1.0' encoding='UTF-8' standalone='yes' ?>
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
				message.Get("data.id").Str,
				message.Get("data.id").Str,
				message.Get("data.id").Str,
				message.Get("data.name").Str,
				message.Get("data.url").Str,
				message.Get("data.id").Str,
				message.Get("data.id").Str,
				message.Get("data.name").Str,
				message.Get("data.owner").Str,
			),
			0,
		)
	}
	return ""
}

func (target msgTarget) cq2xqLocation(message gjson.Result) string {
	core.SendJSON(
		target.BotID,
		1,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(`{"app":"com.tencent.map","desc":"","view":"Share",
			"ver":"0.0.0.1","prompt":"[应用]地图","appID":"","sourceName":"",
			"actionData":"","actionData_A":"","sourceUrl":"","meta":{"Share":{"locSub":"%s",
			"lng":%s,"lat":%s,"zoom":15,"locName":"%s"}},
			"config":{"forward":true,"autosize":1},"text":"","extraApps":[],
			"sourceAd":"","extra":""}`,
			message.Get("data.content").Str,
			message.Get("data.lon").Str,
			message.Get("data.lat").Str,
			message.Get("data.title").Str,
		),
	)
	return ""
}

func (target msgTarget) cq2xqShake(message gjson.Result) string {
	core.ShakeWindow(
		target.BotID,
		target.UserID,
	)
	return ""
}

func (target msgTarget) cq2xqPoke(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(
			"[系统提示] %v 戳了一下 %v",
			target.BotID,
			target.UserID,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target msgTarget) cq2xqAnonymous(message gjson.Result) string {
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		message.Get("data.*").Str,
		0,
		true,
		" ",
	)
	return ""
}

func (target msgTarget) cq2xqReply(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(
			"[系统消息] %v 尝试回复一条消息并失败了",
			target.BotID,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target msgTarget) cq2xqForward(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(
			"[系统消息] %v 尝试合并转发一条消息并失败了",
			target.BotID,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target msgTarget) cq2xqNode(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf(
			"[系统消息] %v 尝试合并转发节点并失败了",
			target.BotID,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target msgTarget) cq2xqDefault(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		message.Get("data.*").Str,
		0,
		false,
		" ",
	)
	return ""
}
