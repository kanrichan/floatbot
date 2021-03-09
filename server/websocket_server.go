package server

import (
	"fmt"
	"net/http"
	"runtime"
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
	ID    int64
	Token string
	Addr  string

	server *http.Server
	mutex  sync.Mutex
	conn   []*WSSConn
}

type WSSConn struct {
	mutex sync.Mutex
	conn  *websocket.Conn
}

func (s *WSS) Run() {
	defer func() {
		s.server = nil
		recover()
	}()
	s.server = &http.Server{Addr: s.Addr, Handler: s}
	s.INFO("正向WS服务建立，等待客户端连接")
	if err := s.server.ListenAndServe(); err != nil {
		s.ERROR(err)
	}
}

func (s *WSS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !(r.Header.Get("Authorization") == s.Token || s.Token == "") {
		// Token验证失败
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	// 元事件 OneBot连接
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
	handshake := fmt.Sprintf(`{"meta_event_type":"lifecycle","post_type":"meta_event","self_ID":%d,"sub_type":"connect","time":%d}`,
		s.ID, time.Now().Unix())
	if err := conn.WriteMessage(websocket.TextMessage, []byte(handshake)); err != nil {
		return
	}
	c := &WSSConn{conn: conn}
	s.mutex.Lock()
	s.conn = append(s.conn, c)
	s.mutex.Unlock()
	go s.listen(c)
	s.INFO("正向WS客户端连接成功")
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
				defer func() {
					if err := recover(); err != nil {
						buf := make([]byte, 1<<16)
						runtime.Stack(buf, true)
						s.PANIC(err, buf)
					}
				}()
				rep := WSSHandler(s.ID, data)
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
	var temp = s.conn
	for i, conn := range s.conn {
		if conn.conn != nil && s.server != nil {
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
		temp = append(temp[:i], temp[i+1:]...)
	}
	s.mutex.Lock()
	s.conn = temp
	s.mutex.Unlock()
}

func (s *WSS) Close() {
	for _, conn := range s.conn {
		if conn.conn == nil {
			continue
		}
		// 锁上当前conn并关闭
		conn.mutex.Lock()
		conn.conn.Close()
		conn.conn = nil
		conn.mutex.Unlock()
	}
	s.conn = nil
	if s.server != nil {
		s.server.Close()
	}
}
