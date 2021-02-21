package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketClient struct {
	ID                 int64  `yaml:"-"`
	Name               string `yaml:"name"`
	Enable             bool   `yaml:"enable"`
	Url                string `yaml:"url"`
	ApiUrl             string `yaml:"api_url"`
	EventUrl           string `yaml:"event_url"`
	UseUniversalClient bool   `yaml:"use_universal_client"`
	AccessToken        string `yaml:"access_token"`
	PostMessageFormat  string `yaml:"post_message_format"`
	HeartBeatInterval  int64  `yaml:"reconnect_interval"`
	ReconnectInterval  int64  `yaml:"reconnect_interval"`

	Conn *websocket.Conn `yaml:"-"`

	StopSend      chan bool `yaml:"-"`
	StopHeartBeat chan bool `yaml:"-"`

	ConnectStatus   string `yaml:"-"`
	ListenStatus    string `yaml:"-"`
	SendStatus      string `yaml:"-"`
	HeartBeatStatus string `yaml:"-"`

	SendChan chan []byte `yaml:"-"`
}

func (this *WebSocketClient) Connect() {
	this.ConnectStatus = "connect"
	defer func() {
		this.ConnectStatus = "wait"
		recover()
	}()
	for {
		// 连接
		header := http.Header{
			"X-Client-Role": []string{"Universal"},
			"X-Self-ID":     []string{strconv.FormatInt(this.ID, 10)},
			"User-Agent":    []string{"CQHttp/4.15.0"},
		}
		conn, _, err := websocket.DefaultDialer.Dial(this.Url, header)
		if err != nil {
			time.Sleep(time.Millisecond * time.Duration(this.ReconnectInterval))
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
		if err := this.Conn.WriteMessage(websocket.TextMessage, handshake); err != nil {
			time.Sleep(time.Millisecond * time.Duration(this.ReconnectInterval))
			continue
		}
		_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
		this.Conn = conn
		break
	}
	this.ConnectStatus = "ok"
}

func (this *WebSocketClient) Listen() {
	this.ListenStatus = "ok"
	defer func() {
		this.ListenStatus = "wait"
		recover()
	}()
	// TODO 监听wsc
	for {
		_, buf, err := this.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		ret := CallApi("WSC", buf)
		this.SendChan <- ret
	}
}

func (this *WebSocketClient) Send() {
	this.SendStatus = "ok"
	defer func() {
		this.SendStatus = "wait"
		recover()
	}()
	for {
		select {
		case send := <-this.SendChan:
			if err := this.Conn.WriteMessage(websocket.TextMessage, send); err != nil {
				panic(err)
			}
		case <-this.StopHeartBeat:
			this.HeartBeatStatus = "wait"
			return
		}
	}
}

func (this *WebSocketClient) HeartBeat() {
	this.HeartBeatStatus = "ok"
	for {
		select {
		case <-time.After(time.Second * time.Duration(this.HeartBeatInterval)):
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
		case <-this.StopHeartBeat:
			this.HeartBeatStatus = "wait"
			return
		}
	}
}
