package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func (this *WebSocketClient) Connect() {
	this.ConnectStatus = "connect"
	defer func() {
		recover()
	}()
	for {
		select {
		case <-time.After(time.Millisecond * time.Duration(this.ReconnectInterval)):
			// 连接
			header := http.Header{
				"X-Client-Role": []string{"Universal"},
				"X-Self-ID":     []string{strconv.FormatInt(this.ID, 10)},
				"User-Agent":    []string{"CQHttp/4.15.0"},
			}
			conn, _, err := websocket.DefaultDialer.Dial(this.Url, header)
			if err != nil {
				continue
			}
			// 元事件 表示OneBot启动成功
			// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
			handshake, _ := json.Marshal(map[string]string{
				"meta_event_type": "lifecycle",
				"post_type":       "meta_event",
				"self_id":         strconv.FormatInt(this.ID, 10),
				"sub_type":        "connect",
				"time":            strconv.FormatInt(time.Now().Unix(), 10),
			})
			if err := conn.WriteMessage(websocket.TextMessage, handshake); err != nil {
				continue
			}
			_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			this.Conn = conn
			this.ConnectStatus = "ok"
			return
		case <-this.StopConnect:
			this.ConnectStatus = "wait"
			return
		}
	}
}

func (this *WebSocketClient) Listen() {
	this.ListenStatus = "ok"
	defer func() {
		this.ListenStatus = "error"
		recover()
	}()
	// TODO 监听wsc
	for {
		_, buf, err := this.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		ret := WebsocketClientCall(buf)
		this.SendChan <- ret
	}
}

func (this *WebSocketClient) Send() {
	this.SendStatus = "ok"
	defer func() {
		this.SendStatus = "wait"
		if err := recover(); err != nil {
			this.SendStatus = "error"
		}
	}()
	for {
		select {
		case <-this.StopSend:
			return
		case send := <-this.SendChan:
			if err := this.Conn.WriteMessage(websocket.TextMessage, send); err != nil {
				panic(err)
			}
		case <-time.After(time.Millisecond * time.Duration(this.HeartBeatInterval)):
			// TODO OneBot标准 元事件 心跳
			// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E5%BF%83%E8%B7%B3
			heartbeat, _ := json.Marshal(
				map[string]interface{}{
					"interval":        strconv.FormatInt(this.HeartBeatInterval, 10),
					"meta_event_type": "heartbeat",
					"post_type":       "meta_event",
					"self_id":         strconv.FormatInt(this.ID, 10),
					"status": map[string]interface{}{
						"online": true,
						"good":   true,
					},
					"time": strconv.FormatInt(time.Now().Unix(), 10),
				},
			)
			this.SendChan <- heartbeat
		}
	}
}
