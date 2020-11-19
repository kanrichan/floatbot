package cqhttp

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/tidwall/gjson"

	"yaya/core"
)

type Data map[string]interface{}
type Echo struct {
	Seq int64 `json:"seq"`
}

type Result struct {
	Status  string `json:"status"`
	Retcode int64  `json:"retcode"`
	Data    Data   `json:"data"`
	Echo    Echo   `json:"echo"`
}

type Reply []map[string]interface{}

func (c *WSC) WSCApi() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[响应服务] Bot %v 服务发生错误 %v，正在自动恢复中......", c.Bot, err)
			c.WSCApi()
		}
	}()
	for {
		if c.Status == 1 {
			break
		}
	}
	DEBUG("[响应服务] Bot %v 服务开始启动...... ", c.Bot)
	for {
		select {
		case api := <-c.Api:
			req := gjson.ParseBytes(api)
			action := strings.ReplaceAll(req.Get("action").Str, "_async", "")
			params := req.Get("params")
			echo := req.Get("echo.seq").Int()

			DEBUG("[响应服务] Bot %v 接收到API调用: %v 参数: %v", c.Bot, req.Get("action").Str, string(api))
			if f, ok := wsApi[action]; ok {
				ret := f(c.Bot, params)
				ret.Echo.Seq = echo
				send, _ := json.Marshal(ret)

				c.Send <- []byte(string(send))
			} else {
				ret := ResultFail()
				ret.Echo.Seq = echo
				send, _ := json.Marshal(ret)

				c.Send <- []byte(string(send))
			}
		}
	}
}

var wsApi = map[string]func(int64, gjson.Result) Result{
	"send_msg": func(bot int64, p gjson.Result) Result {
		message_type := p.Get("message_type").Str
		group_id := p.Get("group_id").Int()
		user_id := p.Get("user_id").Int()
		messages := p.Get("message")
		switch message_type {
		case "group":
			return SendMessage(bot, 2, group_id, 0, messages)
		case "private":
			return SendMessage(bot, 1, 0, user_id, messages)
		default:
			if group_id != 0 {
				return SendMessage(bot, 2, group_id, 0, messages)
			} else {
				return SendMessage(bot, 1, 0, user_id, messages)
			}
		}
	},
	"send_private_msg": func(bot int64, p gjson.Result) Result {
		user_id := p.Get("user_id").Int()
		messages := p.Get("message")
		return SendMessage(bot, 1, 0, user_id, messages)
	},
	"send_group_msg": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		messages := p.Get("message")
		return SendMessage(bot, 2, group_id, 0, messages)
	},
	"send_like": func(bot int64, p gjson.Result) Result {
		user_id := p.Get("user_id").Int()
		core.UpVote(bot, user_id)
		return ResultOK(Data{"message_id": 0})
	},
	"set_group_kick": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		user_id := p.Get("user_id").Int()
		reject_add_request := p.Get("reject_add_request").Bool()
		core.KickGroupMBR(bot, group_id, user_id, reject_add_request)
		return ResultOK(Data{"message_id": 0})
	},
	"set_group_ban": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		user_id := p.Get("user_id").Int()
		duration := p.Get("duration").Int()
		core.ShutUP(bot, group_id, user_id, duration)
		return ResultOK(Data{"message_id": 0})
	},
	"set_group_whole_ban": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		enable := p.Get("enable").Bool()
		if enable {
			core.ShutUP(bot, group_id, 0, 1)
		} else {
			core.ShutUP(bot, group_id, 0, 0)
		}
		return ResultOK(Data{"message_id": 0})
	},
	"set_group_anonymous": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		enable := p.Get("enable").Bool()
		core.SetAnon(bot, group_id, enable)
		return ResultOK(Data{"message_id": 0})
	},
	"set_group_card": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		user_id := p.Get("user_id").Int()
		card := p.Get("card").Str
		core.SetGroupCard(bot, group_id, user_id, card)
		return ResultOK(Data{"message_id": 0})
	},
	"set_group_leave": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		core.QuitGroup(bot, group_id)
		return ResultOK(Data{"message_id": 0})
	},
	"out_put_log": func(bot int64, p gjson.Result) Result {
		text := p.Get("text").Str
		return OutPutLog(text)
	},
	"get_login_info": func(bot int64, p gjson.Result) Result {
		return GetLoginInfo(bot)
	},
}

func ResultOK(data map[string]interface{}) Result {
	return Result{
		Status:  "ok",
		Retcode: 200,
		Data:    data,
		Echo: Echo{
			Seq: 0,
		},
	}
}

func ResultFail() Result {
	return Result{
		Status:  "failed",
		Retcode: 100,
		Data:    Data{},
		Echo: Echo{
			Seq: 0,
		},
	}
}

func SendMessage(selfID int64, messageType int64, groupID int64, userID int64, messages gjson.Result) Result {
	out := ""
	messages = ApiToArray(messages)
	for _, message := range messages.Array() {
		switch message.Get("type").Str {
		case "text":
			out += message.Get("data.*").Str
		case "face":
			out += fmt.Sprintf("[Face%s.gif]", message.Get("data.*").Str)
		case "image":
			image := message.Get("data.*").Str
			if strings.Contains(image, "base64://") {
				path := Base64SaveImage(strings.ReplaceAll(image, "base64://", ""))
				out += fmt.Sprintf("[pic=%s]", path)
			} else if strings.Contains(image, "file:///") {
				out += fmt.Sprintf("[pic=%s]", strings.ReplaceAll(image, "file:///", ""))
			} else if strings.Contains(image, "http://") {
				out += fmt.Sprintf("[pic=%s]", image)
			} else if strings.Contains(image, "https://") {
				out += fmt.Sprintf("[pic=%s]", image)
			} else {
				out += fmt.Sprintf("[pic=%s]", "error")
			}
		case "record":
			record := message.Get("data.*").Str
			if strings.Contains(record, "base64://") {
				path := Base64SaveRecord(strings.ReplaceAll(record, "base64://", ""))
				out += fmt.Sprintf("[Voi=%s]", path)
			} else if strings.Contains(record, "file:///") {
				out += fmt.Sprintf("[Voi=%s]", strings.ReplaceAll(record, "file:///", ""))
			} else if strings.Contains(record, "http://") {
				out += fmt.Sprintf("[Voi=%s]", record)
			} else {
				out += fmt.Sprintf("[Voi=%s]", "error")
			}
		case "video":
			video := message.Get("data.*").Str
			if strings.Contains(video, "base64://") {
				path := Base64SaveVideo(strings.ReplaceAll(video, "base64://", ""))
				out += fmt.Sprintf("[Voi=%s]", path)
			} else if strings.Contains(video, "file:///") {
				out += fmt.Sprintf("[Voi=%s]", strings.ReplaceAll(video, "file:///", ""))
			} else if strings.Contains(video, "http://") {
				out += fmt.Sprintf("[Voi=%s]", video)
			} else {
				out += fmt.Sprintf("[Voi=%s]", "error")
			}
		case "at":
			out += fmt.Sprintf("[@%s]", message.Get("data.*").Str)
		case "rps":
			out += "[no such element]"
		case "dice":
			out += "[no such element]"
		case "shake":
			out += "[no such element]"
		case "poke":
			out += "[no such element]"
		case "anonymous":
			out += "[no such element]"
		case "share":
			out += "[no such element]"
		case "contact":
			out += "[no such element]"
		case "location":
			out += "[no such element]"
		case "music":
			typ := message.Get("data.type").Str
			if typ == "custom" {
				url := message.Get("data.url").Str
				audio := message.Get("data.audio").Str
				title := message.Get("data.title").Str
				content := message.Get("data.content").Str
				image := message.Get("data.image").Str
				json := SendCustomMusic(url, audio, title, content, image)
				TEST("json格式为%v", json)
				core.SendJSON(selfID, 1, 2, groupID, userID, json)
			} else {
				out += "[no such element]"
			}
		case "reply":
			out += "[no such element]"
		case "forward":
			out += "[no such element]"
		case "node":
			out += "[no such element]"
		case "xml":
			xml := message.Get("data.*").Str
			core.SendJSON(selfID, 1, 2, groupID, userID, xml)
		case "json":
			json := message.Get("data.*").Str
			core.SendJSON(selfID, 1, 2, groupID, userID, json)
		case "emoji":
			out += fmt.Sprintf("[emoji=%s]", message.Get("data.*").Str)
		default:
			WARN("CQ码解析失败，将以原格式返回：%v", message.Str)
			out += message.Str
		}
	}
	messageID := "0"
	if out != "" {
		messageID = core.SendMsgEX_V2(selfID, messageType, groupID, userID, out, 0, false, " ")
	}
	return ResultOK(Data{"message_id": messageID})
}

func Base64SaveImage(res string) string {
	data, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		ERROR("base64编码解码失败")
	}
	name := fmt.Sprintf("%x", md5.Sum(data))
	path := ImagePath + name + ".jpg"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("base64编码保存图片失败")
	} else {
		_, err = f.Write(data)
		if err != nil {
			ERROR("base64编码写入图片失败")
		}
	}
	return path
}

func Base64SaveRecord(res string) string {
	data, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		ERROR("base64编码解码失败")
	}
	name := fmt.Sprintf("%x", md5.Sum(data))
	path := RecordPath + name + ".mp3"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("base64编码保存语音失败")
	} else {
		_, err = f.Write(data)
		if err != nil {
			ERROR("base64编码写入语音失败")
		}
	}
	return path
}

func Base64SaveVideo(res string) string {
	data, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		ERROR("base64编码解码失败")
	}
	name := fmt.Sprintf("%x", md5.Sum(data))
	path := VideoPath + name + ".mp4"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("base64编码保存视频失败")
	} else {
		_, err = f.Write(data)
		if err != nil {
			ERROR("base64编码写入视频失败")
		}
	}
	return path
}

func SendCustomMusic(url string, audio string, title string, content string, image string) string {
	music := map[string]interface{}{
		"app":    "com.tencent.structmsg",
		"desc":   "音乐",
		"view":   "music",
		"ver":    "0.0.0.1",
		"prompt": "[分享]" + title,
		"meta": map[string]interface{}{
			"music": map[string]interface{}{
				"action":           "",
				"android_pkg_name": "",
				"app_type":         1,
				"appid":            100495085,
				"desc":             content,
				"jumpUrl":          url,
				"musicUrl":         audio,
				"preview":          image,
				"sourceMsgId":      "0",
				"source_icon":      "",
				"source_url":       "",
				"tag":              "OneBot",
				"title":            title,
			},
		},
	}
	data, _ := json.Marshal(music)
	return string(data)
}

func GetLoginInfo(selfID int64) Result {
	nickname := core.GetNick(selfID, selfID)
	return ResultOK(Data{
		"user_id":  selfID,
		"nickname": nickname,
	})
}

func OutPutLog(text string) Result {
	core.OutPutLog(text)
	return ResultOK(Data{})
}
