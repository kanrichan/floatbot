package onebot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/tidwall/gjson"
)

// cqCode2Array 字符串CQ码转数组
func cqCode2Array(text gjson.Result) gjson.Result {
	if len(text.Get("#.type").Array()) == 0 {
		message := text.Str

		cqcode := regexp.MustCompile(`\[CQ:(.*?)\]`)
		codeList := cqcode.FindAllStringSubmatch(message, -1)
		codeLen := len(codeList)
		messageElem := []string{}
		if codeLen == 0 {
			messageElem = append(messageElem, message)
		} else {
			sMSGe := "start<-" + message + "<-end"
			codeElem := ""
			preElem := ""
			endElem := ""
			for i, c := range codeList {
				codeElem = c[0]
				split := strings.Split(sMSGe, codeElem)
				preElem = split[0]
				endElem = "start<-" + split[1]
				if preElem != "start<-" {
					messageElem = append(messageElem, strings.ReplaceAll(preElem, "start<-", ""))
				}
				messageElem = append(messageElem, codeElem)
				if i+1 == codeLen {
					if endElem != "start<-<-end" {
						messageElem = append(messageElem, strings.ReplaceAll(strings.ReplaceAll(endElem, "start<-", ""), "<-end", ""))
					}
				}
				sMSGe = endElem
			}
		}

		paramsArray := []map[string]interface{}{}
		for _, e := range messageElem {
			if len(cqcode.FindAllStringSubmatch(e, -1)) == 0 {
				paramsArray = append(paramsArray, map[string]interface{}{"type": "text", "data": map[string]interface{}{"text": e}})
			} else {
				codeR1 := regexp.MustCompile(`\[CQ:(.*?),(.*)\]`)
				codeR2 := regexp.MustCompile(`\[CQ:(.*)\]`)
				code := codeR1.FindAllStringSubmatch(e, -1)
				codeType := ""
				codeParm := ""
				if len(code) != 0 {
					codeType = code[0][1]
					codeParm = code[0][2]
				} else {
					code = codeR2.FindAllStringSubmatch(e, -1)
					codeType = code[0][1]
				}

				switch codeType {
				case "face":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "face", "data": map[string]interface{}{
						"id": cqCodeParm(codeParm, "id"),
					}})
				case "image":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "image", "data": map[string]interface{}{
						"file":    cqCodeParm(codeParm, "file"),
						"type":    cqCodeParm(codeParm, "type"),
						"url":     cqCodeParm(codeParm, "url"),
						"cache":   cqCodeParm(codeParm, "cache"),
						"proxy":   cqCodeParm(codeParm, "proxy"),
						"timeout": cqCodeParm(codeParm, "timeout"),
					}})
				case "record":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "record", "data": map[string]interface{}{
						"file":    cqCodeParm(codeParm, "file"),
						"magic":   cqCodeParm(codeParm, "magic"),
						"url":     cqCodeParm(codeParm, "url"),
						"cache":   cqCodeParm(codeParm, "cache"),
						"proxy":   cqCodeParm(codeParm, "proxy"),
						"timeout": cqCodeParm(codeParm, "timeout"),
					}})
				case "video":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "video", "data": map[string]interface{}{
						"file":    cqCodeParm(codeParm, "file"),
						"url":     cqCodeParm(codeParm, "url"),
						"cache":   cqCodeParm(codeParm, "cache"),
						"proxy":   cqCodeParm(codeParm, "proxy"),
						"timeout": cqCodeParm(codeParm, "timeout"),
					}})
				case "at":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "at", "data": map[string]interface{}{
						"qq": cqCodeParm(codeParm, "qq"),
					}})
				case "rps":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "rps", "data": map[string]interface{}{}})
				case "dice":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "dice", "data": map[string]interface{}{}})
				case "shake":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "shake", "data": map[string]interface{}{}})
				case "poke":
					DEBUG("[CQ码解析] %v 不支持，将以原样发送", code[0][0])
					paramsArray = append(paramsArray, map[string]interface{}{"type": "text", "data": map[string]interface{}{
						"text": code[0][0],
					}})
				case "anonymous":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "anonymous", "data": map[string]interface{}{}})
				case "share":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "share", "data": map[string]interface{}{
						"url":     cqCodeParm(codeParm, "url"),
						"title":   cqCodeParm(codeParm, "title"),
						"content": cqCodeParm(codeParm, "content"),
						"image":   cqCodeParm(codeParm, "image"),
					}})
				case "contact":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "contact", "data": map[string]interface{}{
						"type": cqCodeParm(codeParm, "type"),
						"id":   cqCodeParm(codeParm, "id"),
					}})
				case "location":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "location", "data": map[string]interface{}{
						"lat":     cqCodeParm(codeParm, "lat"),
						"lon":     cqCodeParm(codeParm, "lon"),
						"title":   cqCodeParm(codeParm, "title"),
						"content": cqCodeParm(codeParm, "content"),
					}})
				case "music":
					if cqCodeParm(codeParm, "id") != "" {
						paramsArray = append(paramsArray, map[string]interface{}{"type": "music", "data": map[string]interface{}{
							"type": cqCodeParm(codeParm, "type"),
							"id":   cqCodeParm(codeParm, "id"),
						}})
					} else {
						paramsArray = append(paramsArray, map[string]interface{}{"type": "music", "data": map[string]interface{}{
							"type":    cqCodeParm(codeParm, "type"),
							"url":     cqCodeParm(codeParm, "url"),
							"audio":   cqCodeParm(codeParm, "audio"),
							"title":   cqCodeParm(codeParm, "title"),
							"content": cqCodeParm(codeParm, "content"),
							"image":   cqCodeParm(codeParm, "image"),
						}})
					}
				case "reply":
					DEBUG("[CQ码解析] %v 不支持，将以原样发送", code[0][0])
					paramsArray = append(paramsArray, map[string]interface{}{"type": "text", "data": map[string]interface{}{
						"text": code[0][0],
					}})
				case "forward":
					DEBUG("[CQ码解析] %v 不支持，将以原样发送", code[0][0])
					paramsArray = append(paramsArray, map[string]interface{}{"type": "text", "data": map[string]interface{}{
						"text": code[0][0],
					}})
				case "node":
					DEBUG("[CQ码解析] %v 不支持，将以原样发送", code[0][0])
					paramsArray = append(paramsArray, map[string]interface{}{"type": "text", "data": map[string]interface{}{
						"text": code[0][0],
					}})
				case "xml":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "xml", "data": map[string]interface{}{
						"data": cqCodeParm(codeParm, "data"),
					}})
				case "json":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "json", "data": map[string]interface{}{
						"data": cqCodeParm(codeParm, "data"),
					}})
				case "emoji":
					paramsArray = append(paramsArray, map[string]interface{}{"type": "emoji", "data": map[string]interface{}{
						"id": cqCodeParm(codeParm, "id"),
					}})
				default:
					DEBUG("[CQ码解析] %v 解析失败，将以原样发送", code[0][0])
					paramsArray = append(paramsArray, map[string]interface{}{"type": "text", "data": map[string]interface{}{
						"text": code[0][0],
					}})
				}
			}
		}
		data, _ := json.Marshal(paramsArray)
		return gjson.Parse(string(data))
	}
	return text
}

// xq2cqCode 普通XQ码转CQ码
func xq2cqCode(message string) string {
	// 转艾特
	message = strings.ReplaceAll(message, "[@", "[CQ:at,qq=")
	// 转emoji
	message = strings.ReplaceAll(message, "[emoji", "[CQ:emoji,id=")

	// 转face
	face := regexp.MustCompile(`\[Face(.*?)\.gif]`)
	for _, f := range face.FindAllStringSubmatch(message, -1) {
		oldpic := f[0]
		newpic := fmt.Sprintf("[CQ:face,id=%s]", f[1])
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	// 转图片
	pic := regexp.MustCompile(`\[pic={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(\..*?),(.*?)\]`)
	for _, p := range pic.FindAllStringSubmatch(message, -1) {
		oldpic := p[0]
		md5 := strings.ToUpper(fmt.Sprintf("%s%s%s%s%s", p[1], p[2], p[3], p[4], p[5]))
		newpic := fmt.Sprintf("[CQ:image,file=%s.image,url=http://gchat.qpic.cn/gchatpic_new//--%s/0]", md5, md5)
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	pic2 := regexp.MustCompile(`\[pic={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(\..*?)]`)
	for _, p := range pic2.FindAllStringSubmatch(message, -1) {
		oldpic := p[0]
		md5 := strings.ToUpper(fmt.Sprintf("%s%s%s%s%s", p[1], p[2], p[3], p[4], p[5]))
		newpic := fmt.Sprintf("[CQ:image,file=%s.image,url=http://gchat.qpic.cn/gchatpic_new//--%s/0]", md5, md5)
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	// 转语音
	voi := regexp.MustCompile(`\[Voi={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(\..*?),(.*?)\]`)
	for _, v := range voi.FindAllStringSubmatch(message, -1) {
		oldpic := v[0]
		newpic := fmt.Sprintf("[CQ:record,file=%s%s%s%s%s]", v[1], v[2], v[3], v[4], v[5])
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	return message
}

// cq2xqCode 普通CQ码转XQ码
func cq2xqCode(message string) string {
	// 转艾特
	message = strings.ReplaceAll(message, "[CQ:at,qq=", "[@")
	// 转emoji
	message = strings.ReplaceAll(message, "[CQ:emoji,id=", "[emoji")

	// 转face
	face := regexp.MustCompile(`\[CQ:face,id=(.*?)\]`)
	for _, f := range face.FindAllStringSubmatch(message, -1) {
		oldpic := f[0]
		newpic := fmt.Sprintf("[Face%s.gif]", f[1])
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	// 转图片
	pic := regexp.MustCompile(`\[CQ:image,file=(.*?)\]`)
	for _, p := range pic.FindAllStringSubmatch(message, -1) {
		oldpic := p[0]
		newpic := fmt.Sprintf("[pic=%s]", p[1])
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	// 转语音
	voi := regexp.MustCompile(`\[CQ:record,file=(.*?)\]`)
	for _, v := range voi.FindAllStringSubmatch(message, -1) {
		oldpic := v[0]
		newpic := fmt.Sprintf("[Voi=%s]", v[1])
		message = strings.ReplaceAll(message, oldpic, newpic)
	}
	return message
}

func cqCodeParm(codeParm string, field string) string {
	if !strings.Contains(codeParm, field+"=") {
		return ""
	}
	p := strings.Index(codeParm, field+"=") + len(field) + 1
	temp := codeParm[p:]
	s := strings.Index(temp, ",")
	if s == -1 {
		return escape(temp)
	}
	return escape(temp[:s])
}

func escape(text string) string {
	text = strings.ReplaceAll(text, "&#44;", ",")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&#91;", "[")
	text = strings.ReplaceAll(text, "&#93;", "]")
	return text
}
