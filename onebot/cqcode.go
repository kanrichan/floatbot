package onebot

import (
	"fmt"
	"regexp"
	"strings"
	"yaya/core"
)

func cqCode2Array(text string) []map[string]interface{} {
	elems := SplitCQText(text)
	var (
		array = []map[string]interface{}{}

		isCQcode bool = false
		isFirst  bool = false
		isType   bool = false
		isKey    bool = false
		isValue  bool = false

		cqCodetype_ string
		cqCodeKey   []string
		cqCodeValue []string
	)
	for _, r := range elems {
		var temp = []rune{}
	elemLoop:
		for i := range r {
			switch {
			// TODO CQ码开始标记 []rune("[CQ:") = [91 67 81 58]
			case r[i] == 91 && r[i+1] == 67 && r[i+2] == 81 && r[i+3] == 58:
				isCQcode = true
				isFirst = true
			// TODO 不是CQ码
			case !isCQcode:
				array = append(
					array,
					map[string]interface{}{
						"type": "text",
						"data": map[string]string{
							"text": string(r),
						},
					},
				)
				break elemLoop
			// TODO type字段标记 []rune(":") = [58]
			case isCQcode && isFirst && r[i] == 58:
				isType = true
				isFirst = false
			// TODO 开始装载type字段 []rune(",") = [44]
			case isType && r[i] != 44:
				temp = append(temp, r[i])
			// TODO 结束装载type字段 key字段标记 []rune(",") = [44]
			case isType && r[i] == 44:
				cqCodetype_ = string(temp)
				temp = []rune{}
				isType = false
				isKey = true
			// TODO 开始装载key []rune("=") = [61]
			case isKey && r[i] != 61:
				temp = append(temp, r[i])
			// TODO 结束装载key字段 value字段标记 []rune("=") = [61]
			case isKey && r[i] == 61:
				cqCodeKey = append(cqCodeKey, string(temp))
				temp = []rune{}
				isKey = false
				isValue = true
			// TODO 开始装载value []rune(",") = [44] []rune("]") = [93]
			case isValue && r[i] != 44 && r[i] != 93:
				temp = append(temp, r[i])
			// TODO 结束装载value字段 key字段标记 []rune(",") = [44]
			case isValue && r[i] == 44:
				cqCodeValue = append(cqCodeValue, string(temp))
				temp = []rune{}
				isValue = false
				isKey = true
			// TODO 结束装载value字段 结束CQ码 []rune("]") = [93]
			case isValue && r[i] == 93:
				cqCodeValue = append(cqCodeValue, string(temp))
				temp = []rune{}
				cqCodeMap := map[string]interface{}{}
				cqCodeMap["type"] = cqCodetype_
				keyValue := map[string]string{}
				for i := range cqCodeKey {
					keyValue[cqCodeKey[i]] = cqCodeValue[i]
				}
				cqCodeMap["data"] = keyValue
				array = append(array, cqCodeMap)
				cqCodeKey = []string{}
				cqCodeValue = []string{}
				isValue = false
				isCQcode = false
			default:

				// TODO do nothing
			}
		}
	}
	return array
}

func SplitCQText(cqcode string) [][]rune {
	var (
		elems    [][]rune
		temp     []rune
		isCQcode bool = false
	)
	r := []rune(cqcode)
	for i := range r {
		switch {
		// TODO CQ码开始标记 []rune("[CQ:") = [91 67 81 58]
		case r[i] == 91 && r[i+1] == 67 && r[i+2] == 81 && r[i+3] == 58:
			isCQcode = true
			elems = append(elems, temp)
			// TODO 清空temp，开始装CQ码
			temp = []rune{}
			temp = append(temp, r[i])
		// TODO CQ码一直到出现"]"为止 []rune("]") = [93]
		case isCQcode && r[i] != 93:
			temp = append(temp, r[i])
		// TODO 出现"]"，开始下一个elem
		case isCQcode && r[i] == 93:
			isCQcode = false
			temp = append(temp, r[i])
			elems = append(elems, temp)
			temp = []rune{}
		default:
			temp = append(temp, r[i])
		case i == len(r)-1:
			temp = append(temp, r[i])
			elems = append(elems, temp)
		}
	}
	return elems
}

/*
// cqCode2Array 字符串CQ码转数组
func cqCode2Array(message string) []map[string]interface{} {
	array := []map[string]interface{}{}
	for _, elem := range cqCode2Elems(message) {
		array = append(array, map[string]interface{}{
			"type": whatCQcode(elem)["type"],
			"data": whatCQparams(elem)["data"],
		})
	}
	return array
}

// cqCode2Elems CQ 码转字符串数组
func cqCode2Elems(message string) []string {
	elems := []string{}
	start := 0
	for {
		if !strings.Contains(message[start:], "[CQ:") {
			// 如果没有 CQ 码了就直接把剩下的添加进去然后返回
			elems = append(elems, message[start:])
			break
		}
		index := strings.Index(message[start:], "[CQ:") + start
		// 保证不是 CQ 码开头导致无文本
		if start != index {
			elems = append(elems, message[start:index])
		}
		// 找 CQ 码结束位置并把 CQ 码添加进去
		end := strings.Index(message[start:], "]") + start
		elems = append(elems, message[index:end+1])
		start = end + 1
		if start == len(message) {
			break
		}
	}
	return elems
}

// whatCQcode 自动解析 CQ 码类型
func whatCQcode(code string) map[string]string {
	if !strings.Contains(code, "[CQ:") {
		return map[string]string{"type": "text"}
	} else {
		index := strings.Index(code, ",")
		if index == -1 {
			return map[string]string{"type": code[4:(len(code) - 1)]}
		} else {
			return map[string]string{"type": code[4:index]}
		}
	}
}

// whatCQcode 自动解析 CQ 码参数
func whatCQparams(code string) map[string]interface{} {
	elems := map[string]string{}
	start := 0
	for {
		if !strings.Contains(code, "[CQ:") {
			elems["text"] = code
			break
		}
		if !strings.Contains(code[start:], ",") {
			// 如果没有 , 那就是没有参数
			break
		}
		index := strings.Index(code[start:], ",") + start
		// 找 参数 结束位置并把 参数 添加进去
		equal := strings.Index(code[start:], "=") + start
		if !strings.Contains(code[equal:], ",") {
			// 如果没有 , 参数 遍历完毕
			elems[code[index+1:equal]] = escape(code[equal+1 : (len(code) - 1)])
			break
		}
		end := strings.Index(code[equal:], ",") + equal
		elems[code[index+1:equal]] = escape(code[equal+1 : end])
		start = end
	}
	return map[string]interface{}{"data": elems}
}
*/
// xq2cqCode 普通XQ码转CQ码
func xq2cqCode(message string) string {
	// 防止注入
	message = strings.ReplaceAll(message, "[CQ", "[YaYa")
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
	pic := regexp.MustCompile(`\[pic={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(\..*?)\]`)
	for _, p := range pic.FindAllStringSubmatch(message, -1) {
		oldpic := p[0]
		res := fmt.Sprintf("{%s-%s-%s-%s-%s}.jpg", p[1], p[2], p[3], p[4], p[5])
		md5 := fmt.Sprintf("%s%s%s%s%s", p[1], p[2], p[3], p[4], p[5])
		newpic := fmt.Sprintf("[CQ:image,file=%s,url=http://gchat.qpic.cn/gchatpic_new//--%s/0]", res, md5)
		message = strings.ReplaceAll(message, oldpic, newpic)
	}

	// 转语音
	voi := regexp.MustCompile(`\[Voi={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(\..*?),(.*?)\]`)
	for _, v := range voi.FindAllStringSubmatch(message, -1) {
		oldpic := v[0]
		res := fmt.Sprintf("{%s-%s-%s-%s-%s}.amr", v[1], v[2], v[3], v[4], v[5])
		url := core.GetVoiLink(DefaultQQ(), oldpic)
		newpic := fmt.Sprintf("[CQ:record,file=%s,url=%s]", res, url)
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

func escape(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&#44;", ",")
	text = strings.ReplaceAll(text, "&#91;", "[")
	text = strings.ReplaceAll(text, "&#93;", "]")
	return text
}
