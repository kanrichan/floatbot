package server

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	WSCHandler = func(bot int64, data []byte) []byte { fmt.Println(string(data)); return []byte("ok") }
)

// 反向WS
type WSC struct {
	ID    int64
	Addr  string
	Token string

	mutex       sync.Mutex
	conn        *websocket.Conn
	status      bool
	stopconnect chan bool
}

func (s *WSC) Run() {
	defer func() {
		recover()
	}()
	s.stopconnect = make(chan bool, 1)
	// OneBot协议
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(s.ID, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}
	if s.Token != "" {
		header["Authorization"] = []string{"Token " + s.Token}
	}
	s.INFO("反向WS正在尝试连接服务器")
	for {
		select {
		case <-s.stopconnect:
			return
		// 重连定时发生
		case <-time.After(time.Second * 1):
			// 连接
			conn, _, err := websocket.DefaultDialer.Dial(s.Addr, header)
			if err != nil {
				s.DEBUG(err)
				continue
			}
			// 元事件 OneBot连接
			// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
			handshake := fmt.Sprintf(`{"meta_event_type":"lifecycle","post_type":"meta_event","self_ID":%d,"sub_type":"connect","time":%d}`,
				s.ID, time.Now().Unix())
			if err := conn.WriteMessage(websocket.TextMessage, []byte(handshake)); err != nil {
				s.DEBUG(err)
				continue
			}
			// 连接成功
			s.conn = conn
			s.status = true
			s.INFO("反向WS连接服务器成功")
			s.listen()
		}
	}
}

func (s *WSC) listen() {
	for {
		if !s.status || s.conn == nil {
			break
		}
		type_, data, err := s.conn.ReadMessage()
		if err != nil {
			break
		}
		if type_ == websocket.TextMessage {
			go func() {
				defer func() {
					if err := recover(); err != nil {
						buf := make([]byte, 1<<16)
						runtime.Stack(buf, true)
						s.PANIC(err, buf)
					}
				}()
				rep := WSCHandler(s.ID, data)
				s.Send(rep)
			}()
		}
	}
	if s.conn == nil {
		return
	}
	s.mutex.Lock()
	s.conn.Close()
	s.conn = nil
	s.status = false
	s.mutex.Unlock()
}

func (s *WSC) Send(data []byte) {
	if !s.status {
		return
	}
	if s.conn != nil {
		s.mutex.Lock()
		err := s.conn.WriteMessage(websocket.TextMessage, data)
		s.mutex.Unlock()
		if err == nil {
			return
		}
	}
	s.mutex.Lock()
	s.conn.Close()
	s.conn = nil
	s.status = false
	s.mutex.Unlock()
}

func (s *WSC) Close() {
	if !s.status || s.conn == nil {
		s.stopconnect <- true
		return
	}
	s.mutex.Lock()
	s.conn.Close()
	s.conn = nil
	s.status = false
	s.mutex.Unlock()
	s.stopconnect <- true
}
