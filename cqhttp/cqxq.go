package cqhttp

import (
	"fmt"
	"regexp"
	"strings"
)

func CQ(message string) string {
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
		newpic := fmt.Sprintf("[CQ:image,file=%s%s%s%s%s%s]", p[1], p[2], p[3], p[4], p[5], p[6])
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

func XQ(message string) string {
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
