package cqhttp

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/Yiwen-Chan/OneBot-YaYa/core"
)

func SendMessages(selfID int64, messageType string, groupID int64, userID int64, messages gjson.Result) {
	out := ""
	var target int64 = 2
	switch messageType {
	case "group":
		target = 2
	case "private":
		target = 1
	default:
		if groupID != 0 {
			target = 2
		} else if userID != 0 {
			target = 0
		}
	}
	for _, message := range messages.Array() {
		switch message.Get("type").Str {
		case "text":
			out += message.Get("data.*").Str
		case "at":
			out += fmt.Sprintf("[@%s]", message.Get("data.*").Str)
		case "image":
			pic := message.Get("data.*").Str
			if strings.Contains(pic, "base64://") {
				path := Base64SavePic(strings.ReplaceAll(pic, "base64://", ""))
				out += fmt.Sprintf("[pic=%s]", path)
			} else if strings.Contains(pic, "file:///") {
				out += fmt.Sprintf("[pic=%s]", strings.ReplaceAll(pic, "file:///", ""))
			} else if strings.Contains(pic, "http://") {
				out += fmt.Sprintf("[pic=%s]", pic)
			} else {
				out += fmt.Sprintf("[pic=%s]", "error")
			}
		case "record":
			rec := message.Get("data.*").Str
			if strings.Contains(rec, "base64://") {
				path := Base64SaveRec(strings.ReplaceAll(rec, "base64://", ""))
				out += fmt.Sprintf("[Voi=%s]", path)
			} else if strings.Contains(rec, "file:///") {
				out += fmt.Sprintf("[Voi=%s]", strings.ReplaceAll(rec, "file:///", ""))
			} else if strings.Contains(rec, "http://") {
				out += fmt.Sprintf("[Voi=%s]", rec)
			} else {
				out += fmt.Sprintf("[Voi=%s]", "error")
			}
		case "emoji":
			out += fmt.Sprintf("[emoji=%s]", message.Get("data.*").Str)
		case "face":
			out += fmt.Sprintf("[Face%s.gif]", message.Get("data.*").Str)
		default:
			out += "[Element Error]"
		}

	}
	core.SendMsg(selfID, target, groupID, userID, out, 0)
}
