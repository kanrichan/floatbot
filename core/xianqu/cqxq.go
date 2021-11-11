package xianqu

import (
	"fmt"
	"regexp"
	"strings"
)

// xq2cqCode 普通XQ码转CQ码
func xq2cqCode(message string) string {
	// TODO 有cq码注入风险
	// message = strings.ReplaceAll(message, "[CQ", "[YaYa")
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
	pic := regexp.MustCompile(`\[pic={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(.*?)\]`)
	for _, p := range pic.FindAllStringSubmatch(message, -1) {
		oldpic := p[0]
		res := fmt.Sprintf("{%s-%s-%s-%s-%s}.jpg", p[1], p[2], p[3], p[4], p[5])
		md5 := fmt.Sprintf("%s%s%s%s%s", p[1], p[2], p[3], p[4], p[5])
		newpic := fmt.Sprintf("[CQ:image,file=%s.image,url=http://gchat.qpic.cn/gchatpic_new//--%s/0]", md5, md5)
		message = strings.ReplaceAll(message, oldpic, newpic)
		// 记录收到过的图片
		hash := textMD5(fmt.Sprintf("http://gchat.qpic.cn/gchatpic_new//--%s/0", md5))
		PicPoolCache.Insert(strings.ToLower(hash), res)
		PicPoolCache.Insert(md5, res)
	}
	pic2 := regexp.MustCompile(`\[pic={(.*?)-(.*?)-(.*?)}(.*?)\]`)
	for _, p := range pic2.FindAllStringSubmatch(message, -1) {
		oldpic := p[0]
		res := fmt.Sprintf("{%s}.jpg", p[3])
		md5 := p[3]
		newpic := fmt.Sprintf("[CQ:image,file=%s.image,url=http://gchat.qpic.cn/gchatpic_new//--%s/0]", md5, md5)
		message = strings.ReplaceAll(message, oldpic, newpic)
		// 记录收到过的图片
		hash := textMD5(fmt.Sprintf("http://gchat.qpic.cn/gchatpic_new//--%s/0", md5))
		PicPoolCache.Insert(strings.ToLower(hash), res)
		PicPoolCache.Insert(md5, res)
	}

	// 转语音
	voi := regexp.MustCompile(`\[Voi={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(.*?)\]`)
	for _, v := range voi.FindAllStringSubmatch(message, -1) {
		oldpic := v[0]
		res := fmt.Sprintf("{%s-%s-%s-%s-%s}.amr", v[1], v[2], v[3], v[4], v[5])
		url := "" // TODO
		newpic := fmt.Sprintf("[CQ:record,file=%s,url=%s]", res, url)
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	return message
}

func (ctx *Context) XQMessageType() int64 {
	// 1 为先驱私聊代码
	// 2 为先驱群聊代码
	// 4 为先驱群临时代码
	var (
		params = Parse(ctx.Request).Get("params")
		action = Parse(ctx.Request).Str("action")
	)
	switch {
	case params.Str("message_type") == "group":
		return 2
	case params.Str("message_type") == "private":
		//
	case action == "send_group_msg":
		return 2
	case action == "send_private_msg":
		//
	case params.Exist("group_id"):
		return 2
	default:
		//
	}
	tempGroup := TemporarySessionCache.Search(params.Int("user_id"))
	if tempGroup == nil {
		return 1
	}
	ctx.Request["params"].(map[string]interface{})["group_id"] = tempGroup.(int64)
	return 4
}

func (ctx *Context) GetResponseType() int64 {
	// 1 为先驱私聊代码
	// 2 为先驱群聊代码
	// 4 为先驱群临时代码
	switch {
	case Parse(ctx.Response).Str("message_type") == "group":
		return 2
	case Parse(ctx.Response).Str("message_type") == "private":
		return 1
	case Parse(ctx.Response).Exist("group_id"):
		return 2
	default:
		return 1
	}
}

func (ctx *Context) GetUserID() int64 {
	// 1 为先驱私聊代码
	// 2 为先驱群聊代码
	// 4 为先驱群临时代码
	switch ctx.GetResponseType() {
	case 1:
		return 0
	default:
		return Parse(ctx.Response).Int("user_id")
	}
}
