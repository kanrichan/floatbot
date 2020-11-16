package cqhttp

import (
	"encoding/json"

	"github.com/tidwall/gjson"

	"github.com/Yiwen-Chan/OneBot-YaYa/core"
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

func (c *WSC) WSCApi() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[响应服务] Bot %v 服务发生错误 %v，正在自动恢复中......", c.Bot, err)
			c.WSCApi()
		}
	}()

	DEBUG("[响应服务] Bot %v 服务开始启动...... ", c.Bot)
	for {
		select {
		case api := <-c.Api:
			req := gjson.ParseBytes(api)
			action := req.Get("action").Str
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
		return SendMessage(bot, message_type, group_id, user_id, messages)
	},
	"send_private_msg": func(bot int64, p gjson.Result) Result {
		user_id := p.Get("user_id").Int()
		messages := p.Get("message")
		return SendPrivateMessage(bot, user_id, messages)
	},
	"send_group_msg": func(bot int64, p gjson.Result) Result {
		group_id := p.Get("group_id").Int()
		messages := p.Get("message")
		return SendGroupMessage(bot, group_id, messages)
	},
	"out_put_log": func(bot int64, p gjson.Result) Result {
		text := p.Get("text").Str
		return OutPutLog(text)
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

func SendMessage(selfID int64, messageType string, groupID int64, userID int64, messages gjson.Result) Result {
	SendMessages(selfID, messageType, groupID, userID, messages)
	return ResultOK(Data{"message_id": 0})
}

func SendGroupMessage(selfID int64, groupID int64, messages gjson.Result) Result {
	SendMessages(selfID, "group", groupID, 0, messages)
	return ResultOK(Data{"message_id": 0})
}

func SendPrivateMessage(selfID int64, userID int64, messages gjson.Result) Result {
	SendMessages(selfID, "private", 0, userID, messages)
	return ResultOK(Data{"message_id": 0})

}

func OutPutLog(text string) Result {
	core.OutPutLog(text)
	return ResultOK(Data{})
}
