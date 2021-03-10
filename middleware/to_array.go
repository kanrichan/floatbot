package middleware

import (
	"fmt"
	"strings"

	core "onebot/core/xianqu"
)

const (
	LEFT_COLON  = 91 // "["
	RIGHT_COLON = 93 // "]"
	SEMI_COLON  = 58 // ":"
	COMMA       = 44 // ","
	EQUAL       = 61 // "="
)

// ResponseToArray 将报文中Response的message转换为array格式
func ResponseToArray(ctx *core.Context) {
	message := ctx.Response["message"]
	switch message.(type) {
	case string:
		//
	default:
		return
	}
	ctx.Response["message"] = toArray(message.(string))
	return
}

// RequestToArray 将报文中Request的message转换为array格式
func RequestToArray(ctx *core.Context) {
	request := core.Parse(ctx.Request)
	if !request.Exist("params") {
		return
	}
	params := request.Get("params")
	if !params.Exist("message") {
		return
	}
	// 如果本来就是数组格式则不转化
	if len(params.Array("message")) != 0 {
		return
	}
	message := params.Str("message")
	fmt.Println(message)
	if message == "" {
		return
	}
	// 保证不是拷贝的
	ctx.Request["params"].(map[string]interface{})["message"] = toArray(message)
	return
}

// 缓存字节数组
type temp struct {
	data []byte
}

// newTemp 返回一个空 temp
func newTemp() *temp {
	return &temp{}
}

// push 将元素放到数组的最后
func (b *temp) push(v byte) {
	b.data = append(b.data, v)
}

// pop 取出所有的元素
func (b *temp) pop() (v []byte) {
	v = b.data
	b.data = nil
	return v
}

func (b *temp) size() int {
	return len(b.data)
}

// 先进先出
type heap struct {
	data []interface{}
}

// newHeap 返回一个空 heap
func newHeap() *heap {
	return &heap{}
}

// push 将元素放到数组的最后
func (t *heap) push(v interface{}) {
	t.data = append(t.data, v)
}

// pop 取出首元素，后面元素往前移
func (t *heap) pop() (v interface{}) {
	switch len(t.data) {
	case 0:
		return nil
	case 1:
		v = t.data[0]
		t.data = nil
	default:
		v = t.data[0]
		t.data = t.data[1:]
	}
	return v
}

// size 返回 heap 的大小
func (t *heap) size() int {
	return len(t.data)
}

// 存 message 的数组 map
type maps struct {
	data []map[string]interface{}
}

// newMaps 返回一个 maps
func newMaps() *maps {
	return &maps{}
}

// buildMaps 返回 message 的数组map
func (b *maps) buildMaps(type_ *heap, key *heap, val *heap) {
	kv := map[string]interface{}{}
	if key.data != nil {
		size := key.size()
		for i := 0; i < size; i++ {
			kv[key.pop().(string)] = escape(val.pop().(string))
		}
	}
	b.data = append(
		b.data,
		map[string]interface{}{
			"type": type_.pop(),
			"data": kv,
		},
	)
}

// toArray 快速解析 message 字符串 --> 数组
func toArray(message string) []map[string]interface{} {
	data := []byte(message)
	var (
		top   = len(data) - 1
		build = newMaps() // 输出的数组格式的message
		text  = newTemp() // 字符串message的缓存
		type_ = newHeap() // cq码中的type
		key   = newHeap() // cq码中的key
		val   = newHeap() // cq码中的val
	)
	for i := range data {
		switch data[i] {
		case LEFT_COLON:
			if text.size() == 0 {
				break // "[" 前面没有文本
			}
			// "[" 前面有文本
			type_.push("text")
			key.push("text")
			val.push(string(text.pop()))
			build.buildMaps(type_, key, val)
		case SEMI_COLON:
			text.pop() // 删除 ":" 前的 temp
		case COMMA:
			switch len(type_.data) {
			case 0: // 没有 type ，所以 "," 前的是 type
				type_.push(string(text.pop()))
			default: // 有 type ，所以 "," 前的是 val
				val.push(string(text.pop()))
			}
		case EQUAL:
			// "=" 前面的是 key
			key.push(string(text.pop()))
		case RIGHT_COLON:
			switch len(type_.data) {
			case 0: // 没有 type ，所以 "]" 前的是 type
				type_.push(string(text.pop()))
			default: // 有 type ，所以 "]" 前的是 val
				val.push(string(text.pop()))
			}
			build.buildMaps(type_, key, val)
		default:
			if i == top {
				// 结束前有文本
				text.push(data[i])
				type_.push("text")
				key.push("text")
				val.push(string(text.pop()))
				build.buildMaps(type_, key, val)
			}
			text.push(data[i])
		}
	}
	return build.data
}

// escape CQ码转义
func escape(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&#44;", ",")
	text = strings.ReplaceAll(text, "&#91;", "[")
	text = strings.ReplaceAll(text, "&#93;", "]")
	return text
}
