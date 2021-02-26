package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (s *WebSocketServer) Init() {

	if s.Conn == nil {
		s.Conn = []*WebSocketServerConn{}
	}
}

func (s *WebSocketServer) Run() {
	switch {
	case s.Enable && s.ConnectStatus == S_WAIT && s.Host != "":
		if s.Port == 0 {
			s.Port = 80
		}
		INFO(s, "Listen Start")
		go s.Connect()
		fallthrough
	case s.Enable && s.ConnectStatus == S_OK:
		for i := range s.Conn {
			switch {
			case s.Conn[i].ListenStatus == S_WAIT:
				go s.Listen(s.Conn[i])
				INFO(s, "Listen Start")
			case s.Conn[i].SendStatus == S_WAIT:
				go s.Send(s.Conn[i])
				INFO(s, "Send Start")
			}
		}
	}
}

func (s *WebSocketServer) Stop(force bool) {
	switch {
	case s.ConnectStatus == S_ERROR:
		if s.Server != nil {
			s.Server.Close()
		}
		time.Sleep(time.Second * 1)
		s.ConnectStatus = S_WAIT
		fallthrough
	case s.ConnectStatus == S_OK:
		for i := range s.Conn {
			switch {
			case s.Conn[i].ListenStatus == S_ERROR:
				if s.Conn[i].SendStatus == S_OK {
					s.Conn[i].StopSend <- true
				}
				if s.Conn[i].Conn != nil {
					s.Conn[i].Conn.Close()
					s.Conn[i].Conn = nil
				}
				time.Sleep(time.Second * 1)
				s.Mutex.Lock()
				s.Conn = append(s.Conn[:i], s.Conn[i+1:]...)
				s.Mutex.Unlock()
			case s.Conn[i].SendStatus == S_ERROR:
				if s.Conn[i].Conn != nil {
					s.Conn[i].Conn.Close()
					s.Conn[i].Conn = nil
				}
				time.Sleep(time.Second * 1)
				s.Mutex.Lock()
				s.Conn = append(s.Conn[:i], s.Conn[i+1:]...)
				s.Mutex.Unlock()
			}
		}
	case force:
		for i := range s.Conn {
			switch {
			case s.Conn[i].SendStatus == S_OK:
				s.Conn[i].StopSend <- true
			case s.Conn[i].ListenStatus == S_OK:
				if s.Conn[i].Conn != nil {
					s.Conn[i].Conn.Close()
					s.Conn[i].Conn = nil
				}
			}
		}
		time.Sleep(time.Second * 1)
		s.Mutex.Lock()
		s.Conn = nil
		s.Mutex.Unlock()
		if s.Server != nil {
			s.Server.Close()
			s.Server = nil
		}
	}
}

func (s *WebSocketServer) Connect() {
	// 开始监听
	s.ConnectStatus = S_OK
	defer func() {
		s.ConnectStatus = S_WAIT
		if err := recover(); err != nil {
			// 建立监听时或监听过程发生错误
			s.ConnectStatus = S_ERROR
			ERROR(s, err)
		}
	}()
	s.ConnectStatus = S_OK
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.Host, s.Port),
		Handler: s,
	}
	s.Server = server
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *WebSocketServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !(r.Header.Get("Authorization") == s.AccessToken || s.AccessToken == "") {
		// Token验证失败
		ERROR(s, "Token ERROR")
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		ERROR(s, err)
		return
	}
	// 元事件 表示OneBot启动成功
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
	handshake, _ := json.Marshal(map[string]string{
		"meta_event_type": "lifecycle",
		"post_type":       "meta_event",
		"self_id":         strconv.FormatInt(s.ID, 10),
		"sub_type":        "connect",
		"time":            strconv.FormatInt(time.Now().Unix(), 10),
	})
	if err := conn.WriteMessage(websocket.TextMessage, handshake); err != nil {
		ERROR(s, err)
		return
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.Conn = append(s.Conn, &WebSocketServerConn{
		Conn:         conn,
		StopSend:     make(chan bool, 1),
		ListenStatus: S_WAIT,
		SendStatus:   S_WAIT,
		SendChan:     make(chan []byte, 100),
		Handler:      WebSocketServerHandler,
	})
}

func (s *WebSocketServer) Listen(connect *WebSocketServerConn) {
	connect.ListenStatus = S_OK
	defer func() {
		connect.ListenStatus = S_WAIT
		if err := recover(); err != nil {
			// 读取信息或Handler传递时发生错误
			connect.ListenStatus = S_ERROR
			ERROR(s, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			ERROR(s, fmt.Sprintf("[TRACEBACK]:\n%v", string(buf)))
		}
	}()
	for {
		_, buf, err := connect.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		ret := connect.Handler(s.ID, buf)
		connect.SendChan <- ret
	}
}

func (s *WebSocketServer) Send(connect *WebSocketServerConn) {
	connect.SendStatus = S_OK
	defer func() {
		connect.SendStatus = S_WAIT
		if err := recover(); err != nil {
			connect.SendStatus = S_ERROR
			ERROR(s, err)
		}
	}()
	for {
		select {
		case <-connect.StopSend:
			ERROR(s, "Send 命令退出")
			return
		case <-time.After(time.Millisecond * time.Duration(s.HeartBeatInterval)):
			// OneBot标准 元事件 心跳
			// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E5%BF%83%E8%B7%B3
			heartbeat, _ := json.Marshal(
				map[string]interface{}{
					"interval":        strconv.FormatInt(s.HeartBeatInterval, 10),
					"meta_event_type": "heartbeat",
					"post_type":       "meta_event",
					"self_id":         strconv.FormatInt(s.ID, 10),
					"status": map[string]interface{}{
						"online": true,
						"good":   true,
					},
					"time": strconv.FormatInt(time.Now().Unix(), 10),
				},
			)
			connect.SendChan <- heartbeat
		case send := <-connect.SendChan:
			if err := connect.Conn.WriteMessage(websocket.TextMessage, send); err != nil {
				panic(err)
			}

		}
	}
}
