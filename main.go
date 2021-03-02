package main

import (
	"errors"
	"fmt"
	core "onebot/core/xianqu"
	"onebot/gateway"
	_ "onebot/gateway"
	_ "onebot/middleware"
	"runtime"
	"time"
)

func main() {
	/*
		ctx := &core.Context{
			Bot: int64(123),
			Response: map[string]interface{}{
				"time":         int64(123),
				"self_id":      int64(123),
				"post_type":    "message",
				"message_type": "group",
				"sub_type":     "normal",
				"message_id":   int64(123),
				"group_id":     int64(123),
				"user_id":      int64(123),
				"anonymous":    nil,
				"message":      "2333",
				"raw_message":  "2333",
				"font":         int64(123),
				"sender": map[string]interface{}{
					"user_id":  int64(123),
					"nickname": "unknown",
					"sex":      "unknown",
					"age":      "unknown",
					"area":     "",
					"card":     "",
					"level":    "",
					"role":     "unknown",
					"title":    "unknown",
				},
			},
			Request: map[string]interface{}{
				"reply":     "hello",
				"at_sender": true,
			},
		}
		middle.ResponeFastReplyFormat(ctx)
		fmt.Println(ctx.Request)
	*/
	//gateway.Controller()
	/*
		gateway.WebSocketServerHandler(123, []byte(`{
				"action": "send_msg",
				"params": {
					"user_id": 10001000,
					"message": "你好"
				},
				"echo": "123"
			}`))
	*/
	core.OneBotPath = "C:\\Users\\kanri\\Desktop\\XQ\\OneBot\\"
	go gateway.OnEnable(nil)
	time.Sleep(time.Second * 10)
	gateway.OnDisable(nil)
	time.Sleep(time.Second * 10)
	go gateway.OnEnable(nil)
	time.Sleep(time.Second * 10)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[PANIC] 发生了不可预知的错误，请在GitHub提交issue：%v", err)
				buf := make([]byte, 1<<16)
				runtime.Stack(buf, true)
				fmt.Printf("[TRACEBACK]:\n%v", string(buf))
			}
		}()
		fmt.Println("111")
		panic(errors.New("?"))
	}()
	time.Sleep(time.Second * 10)
}
