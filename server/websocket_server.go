package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	ID                int64  `yaml:"-"`
	Name              string `yaml:"name"`
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	AccessToken       string `yaml:"access_token"`
	PostMessageFormat string `yaml:"post_message_format"`
	HeartBeatInterval int64  `yaml:"reconnect_interval"`

	ConnectStatus string `yaml:"-"`

	Connects []Connect `yaml:"-"`
}

type Connect struct {
	Conn *websocket.Conn `yaml:"-"`

	StopSend      chan bool `yaml:"-"`
	StopHeartBeat chan bool `yaml:"-"`

	ListenStatus    string `yaml:"-"`
	SendStatus      string `yaml:"-"`
	HeartBeatStatus string `yaml:"-"`

	SendChan chan []byte `yaml:"-"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (this *WebSocketServer) Connect() {
	this.ConnectStatus = "connect"
	http.ListenAndServe(fmt.Sprintf("%v:%v", this.Host, this.Port), this)
	this.ConnectStatus = "ok"
}

func (this *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		recover()
	}()
	if !(r.Header.Get("Authorization") == "Token "+this.AccessToken || this.AccessToken == "") {
		panic(errors.New("token error"))
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
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
	_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
	if err := conn.WriteMessage(websocket.TextMessage, handshake); err != nil {
		panic(err)
	}
	this.Connects = append(this.Connects, Connect{
		Conn:            conn,
		StopSend:        make(chan bool, 1),
		StopHeartBeat:   make(chan bool, 1),
		ListenStatus:    "wait",
		SendStatus:      "wait",
		HeartBeatStatus: "wait",
		SendChan:        make(chan []byte, 50),
	})
}

func (this *WebSocketServer) Listen(connect Connect) {
	connect.ListenStatus = "ok"
	defer func() {
		connect.ListenStatus = "wait"
		recover()
	}()
	for {
		_, buf, err := connect.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		ret := CallApi("WSS", buf)
		connect.SendChan <- ret
	}
}

func (this *WebSocketServer) Send(connect Connect) {
	connect.SendStatus = "ok"
	defer func() {
		connect.SendStatus = "wait"
		recover()
	}()
	for {
		select {
		case send := <-connect.SendChan:
			if err := connect.Conn.WriteMessage(websocket.TextMessage, send); err != nil {
				panic(err)
			}
		case <-connect.StopSend:
			return
		}
	}
}

func (this *WebSocketServer) HeartBeat(connect Connect) {
	connect.HeartBeatStatus = "ok"
	defer func() {
		connect.HeartBeatStatus = "wait"
		recover()
	}()
	for {
		select {
		case <-time.After(time.Second * time.Duration(this.HeartBeatInterval)):
			// OneBot标准 元事件 心跳
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
			connect.SendChan <- heartbeat
		case <-connect.StopHeartBeat:
			return
		}
	}
}
