package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var (
	WSSHandler = func(bot int64, data []byte) []byte { fmt.Println(string(data)); return []byte("ok") }
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// 正向WS
type WSS struct {
	id     int64
	token  string
	addr   string
	server *http.Server
	mutex  sync.Mutex
	conn   []*WSSConn
}

type WSSConn struct {
	mutex sync.Mutex
	conn  *websocket.Conn
}

func (s *WSS) Run(id int64, addr, token string) {
	defer func() {
		s.server = nil
		recover()
	}()
	s.id = id
	s.addr = addr
	s.token = token
	s.server = &http.Server{Addr: addr, Handler: s}
	go s.heartbeat()
	s.server.ListenAndServe()
	s.server = nil
}

func (s *WSS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !(r.Header.Get("Authorization") == s.token || s.token == "") {
		// Token验证失败
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// 元事件 OneBot连接
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
	handshake := fmt.Sprintf(`{"meta_event_type":"lifecycle","post_type":"meta_event","self_id":%d,"sub_type":"connect","time":%d}`,
		s.id, time.Now().Unix())
	if err := conn.WriteMessage(websocket.TextMessage, []byte(handshake)); err != nil {
		return
	}
	c := &WSSConn{conn: conn}
	s.mutex.Lock()
	s.conn = append(s.conn, c)
	s.mutex.Unlock()
	go s.listen(c)
}

func (s *WSS) listen(conn *WSSConn) {
	for {
		if conn.conn == nil || s.server == nil {
			break
		}
		type_, data, err := conn.conn.ReadMessage()
		if err != nil {
			break
		}
		if type_ == websocket.TextMessage {
			go func() {
				rep := WSSHandler(s.id, data)
				if conn.conn == nil {
					return
				}
				conn.mutex.Lock()
				conn.conn.WriteMessage(websocket.TextMessage, rep)
				conn.mutex.Unlock()
			}()
		}
	}
}

func (s *WSS) Send(data []byte) {
	fmt.Println(s.conn)
	for i, conn := range s.conn {
		if conn.conn != nil && s.server != nil {
			fmt.Println("loop")
			conn.mutex.Lock()
			err := conn.conn.WriteMessage(websocket.TextMessage, data)
			conn.mutex.Unlock()
			if err == nil {
				continue // 没有任何错误即进入下一个循环
			}
		}
		// 锁上当前conn并关闭
		conn.mutex.Lock()
		conn.conn.Close()
		conn.conn = nil
		conn.mutex.Unlock()
		// 锁上conn列表并从列表上删除该conn
		s.mutex.Lock()
		s.conn = append(s.conn[:i], s.conn[i+1:]...)
		s.mutex.Unlock()
	}
}

func (s *WSS) heartbeat() {
	for {
		time.Sleep(time.Second * 3)
		heartbeat := fmt.Sprintf(`{"interval":%d,"meta_event_type":"heartbeat","post_type":"meta_event","self_id":%d,"status":{"good":true,"online":true},"time":%d}`,
			3000, s.id, time.Now().Unix())
		s.Send([]byte(heartbeat))
	}
}
