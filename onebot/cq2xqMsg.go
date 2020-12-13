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

type cq2xqMsgToWhere struct {
	BotID   int64
	Type_   int64
	GroupID int64
	UserID  int64
}

func cq2xqMessageType(p gjson.Result) int64 {
	switch {
	case p.Get("message_type").Str == "private":
		return 1
	case p.Get("message_type").Str == "group":
		return 2
	case p.Get("group_id").Int() != 0:
		return 2
	default:
		return 1
	}
}

func cq2xqSendMsg(bot int64, p gjson.Result) Result {
	target := cq2xqMsgToWhere{
		BotID:   bot,
		Type_:   cq2xqMessageType(p),
		GroupID: p.Get("group_id").Int(),
		UserID:  p.Get("user_id").Int(),
	}

	var out string = ""
	var bubble int64 = 0
	var msg gjson.Result

	if len(p.Get("message.#.type").Array()) == 0 {
		b, _ := json.Marshal(cqCode2Array(p.Get("message").Str))
		msg = gjson.ParseBytes(b)
		TEST("%v", msg)
	} else {
		msg = p.Get("message")
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
			TEST("id %v", bubble)
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
		data := core.SendMsgEX_V2(
			target.BotID,
			target.Type_,
			target.GroupID,
			target.UserID,
			out,
			bubble,
			false,
			" ",
		)
		if data != "" {
			p = gjson.Parse(data[:strings.LastIndex(data, "}")])
			if p.Get("sendok").Bool() {
				var xe XEvent
				xe.messageID = p.Get("msgid").Int()
				xe.messageNum = p.Get("msgno").Int()
				xe.cqID = 0
				for i, _ := range Conf.BotConfs {
					if bot == Conf.BotConfs[i].Bot && bot != 0 && Conf.BotConfs[i].DB != nil {
						time.Sleep(time.Millisecond * 100)
						xe.xq2cqid(Conf.BotConfs[i].DB)
						return resultOK(map[string]interface{}{"message_id": xe.cqID})
					}
				}
			}
		}
		send := ""
		for i, o := range strings.Split(out, "\n") {
			if (i%3) != 0 && i != 0 {
				send = send + o + "\n"
			} else {
				send = send + o + "[Next]"
			}
		}
		data = core.SendMsgEX_V2(
			target.BotID,
			target.Type_,
			target.GroupID,
			target.UserID,
			send,
			bubble,
			false,
			" ",
		)
		if data != "" {
			p = gjson.Parse(data[:strings.LastIndex(data, "}")])
			if p.Get("sendok").Bool() {
				var xe XEvent
				xe.messageID = p.Get("msgid").Int()
				xe.messageNum = p.Get("msgno").Int()
				xe.cqID = 0
				for i, _ := range Conf.BotConfs {
					if bot == Conf.BotConfs[i].Bot && bot != 0 && Conf.BotConfs[i].DB != nil {
						time.Sleep(time.Millisecond * 100)
						xe.xq2cqid(Conf.BotConfs[i].DB)
						return resultOK(map[string]interface{}{"message_id": xe.cqID})
					}
				}
			}
		}
		return resultOK(map[string]interface{}{"message_id": 0})
	}
	return resultOK(map[string]interface{}{"message_id": 0})
}

func (target cq2xqMsgToWhere) cq2xqText(message gjson.Result) string {
	return message.Get("data.*").Str
}

func (target cq2xqMsgToWhere) cq2xqFace(message gjson.Result) string {
	return fmt.Sprintf("[Face%s.gif]", message.Get("data.*").Str)
}

func (target cq2xqMsgToWhere) cq2xqAt(message gjson.Result) string {
	return fmt.Sprintf("[@%s] ", message.Get("data.*").Str)
}

func (target cq2xqMsgToWhere) cq2xqEmoji(message gjson.Result) string {
	return fmt.Sprintf("[emoji=%s]", message.Get("data.*").Str)
}

func (target cq2xqMsgToWhere) cq2xqRps(message gjson.Result) string {
	return []string{
		"[魔法猜拳] 石头",
		"[魔法猜拳] 剪刀",
		"[魔法猜拳] 布",
	}[rand.Intn(3)]
}

func (target cq2xqMsgToWhere) cq2xqDice(message gjson.Result) string {
	return []string{
		"[魔法骰子] 1",
		"[魔法骰子] 2",
		"[魔法骰子] 3",
		"[魔法骰子] 4",
		"[魔法骰子] 5",
		"[魔法骰子] 6",
	}[rand.Intn(6)]
}

func (target cq2xqMsgToWhere) cq2xqImage(message gjson.Result) string {
	url := strings.ReplaceAll(message.Get("data.url").Str, `\/`, `/`)
	image := strings.ReplaceAll(message.Get("data.file").Str, `\/`, `/`)
	showID := message.Get("data.id").Int() - 40000
	switch message.Get("data.type").Str {
	case "show":
		switch {
		case url != "":
			return fmt.Sprintf("[ShowPic=%s,type=%d]", Url2Image(url), showID)
		case strings.Contains(image, "base64://"):
			return fmt.Sprintf("[ShowPic=%s,type=%d]", Base642Image(image[9:]), showID)
		case strings.Contains(image, "file:///"):
			return fmt.Sprintf("[ShowPic=%s,type=%d]", image[8:], showID)
		case strings.Contains(image, "http://"):
			return fmt.Sprintf("[ShowPic=%s,type=%d]", Url2Image(image), showID)
		case strings.Contains(image, "https://"):
			return fmt.Sprintf("[ShowPic=%s,type=%d]", Url2Image(image), showID)
		default:
			return fmt.Sprintf("[ShowPic=%s,type=%d]", "error", showID)
		}
	default:
		switch {
		case url != "":
			return fmt.Sprintf("[pic=%s]", Url2Image(url))
		case strings.Contains(image, "base64://"):
			return fmt.Sprintf("[pic=%s]", Base642Image(image[9:]))
		case strings.Contains(image, "file:///"):
			return fmt.Sprintf("[pic=%s]", image[8:])
		case strings.Contains(image, "http://"):
			return fmt.Sprintf("[pic=%s]", image)
		case strings.Contains(image, "https://"):
			return fmt.Sprintf("[pic=%s]", image)
		default:
			return fmt.Sprintf("[pic=%s]", "error")
		}
	}
}

func (target cq2xqMsgToWhere) cq2xqRecord(message gjson.Result) string {
	record := strings.ReplaceAll(message.Get("data.file").Str, `\/`, `/`)
	switch {
	case strings.Contains(record, "file:///"):
		return fmt.Sprintf("[Voi=%s]", record[8:])
	default:
		return fmt.Sprintf("[Voi=%s]", "error")
	}
}

func (target cq2xqMsgToWhere) cq2xqVideo(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqMusic(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqXml(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqJson(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqShare(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 暂未实现", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf("分享：%s\n%s",
			message.Get("data.title").Str,
			message.Get("data.url").Str,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target cq2xqMsgToWhere) cq2xqContact(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 暂未实现", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf("分享：%s: %s",
			message.Get("data.type").Str,
			message.Get("data.id").Str,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target cq2xqMsgToWhere) cq2xqLocation(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 暂未实现", message.Str)
	core.SendMsgEX_V2(
		target.BotID,
		target.Type_,
		target.GroupID,
		target.UserID,
		fmt.Sprintf("位置分享：%s/n%s/n经度：%s 纬度：%s",
			message.Get("data.title").Str,
			message.Get("data.content").Str,
			message.Get("data.lon").Str,
			message.Get("data.lat").Str,
		),
		0,
		false,
		" ",
	)
	return ""
}

func (target cq2xqMsgToWhere) cq2xqShake(message gjson.Result) string {
	core.ShakeWindow(
		target.BotID,
		target.UserID,
	)
	return ""
}

func (target cq2xqMsgToWhere) cq2xqPoke(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqAnonymous(message gjson.Result) string {
	DEBUG("[CQ码解析] %v 不支持", message.Str)
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

func (target cq2xqMsgToWhere) cq2xqReply(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqForward(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqNode(message gjson.Result) string {
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

func (target cq2xqMsgToWhere) cq2xqDefault(message gjson.Result) string {
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
