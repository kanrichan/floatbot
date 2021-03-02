package server

import (
	"fmt"
	"net/http"
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
	id     int64
	status int32
	addr   string
	token  string

	mutex       sync.Mutex
	conn        *websocket.Conn
	stopconnect chan bool
}

func (s *WSC) Run(id int64, addr, token string) {
	defer func() {
		recover()
	}()
	s.id = id
	s.addr = addr
	s.token = token
	s.stopconnect = make(chan bool, 1)
	// OneBot协议
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(s.id, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}
	if s.token != "" {
		header["Authorization"] = []string{"Token " + s.token}
	}
	for {
		select {
		case <-s.stopconnect:
			fmt.Println("stop6")
			return
		// 重连定时发生
		case <-time.After(time.Second * 1):
			// 连接
			conn, _, err := websocket.DefaultDialer.Dial(s.addr, header)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// 元事件 OneBot连接
			// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
			handshake := fmt.Sprintf(`{"meta_event_type":"lifecycle","post_type":"meta_event","self_id":%d,"sub_type":"connect","time":%d}`,
				s.id, time.Now().Unix())
			if err := conn.WriteMessage(websocket.TextMessage, []byte(handshake)); err != nil {
				fmt.Println(err)
				continue
			}
			// 连接成功
			s.conn = conn
			go s.heartbeat()
			s.listen()
		}
	}
}

func (s *WSC) listen() {
	for {
		if s.conn == nil {
			break
		}
		type_, data, err := s.conn.ReadMessage()
		if err != nil {
			break
		}
		if type_ == websocket.TextMessage {
			go func() {
				rep := WSCHandler(s.id, data)
				s.Send(rep)
			}()
		}
	}
}

func (s *WSC) Send(data []byte) {
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
	s.mutex.Unlock()
}

func (s *WSC) heartbeat() {
	for {
		time.Sleep(time.Second * 3)
		heartbeat := fmt.Sprintf(`{"interval":%d,"meta_event_type":"heartbeat","post_type":"meta_event","self_id":%d,"status":{"good":true,"online":true},"time":%d}`,
			3000, s.id, time.Now().Unix())
		s.Send([]byte(heartbeat))
	}
}

func (s *WSC) Stop() {
	time.Sleep(time.Second * 1)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.conn != nil {
		s.conn.Close()
		s.conn = nil
	}
	s.stopconnect <- true
}
