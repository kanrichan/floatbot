package server

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// 正向WS
type WSS struct {
	// Bot的qq号
	ID    int64
	Token string
	Addr  string

	server *http.Server
	mutex  sync.Mutex
	conn   []*WSSConn
}

// 正向WS的连接
type WSSConn struct {
	mutex sync.Mutex
	conn  *websocket.Conn
}

// Run 建立正向WS
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
	handshake := fmt.Sprintf(`{"meta_event_type":"lifecycle","post_type":"meta_event","self_id":%d,"sub_type":"connect","time":%d}`,
		s.ID, time.Now().Unix())
	if err := conn.WriteMessage(websocket.TextMessage, helper.StringToBytes(handshake)); err != nil {
		return
	}
	s.INFO("正向WS客户端连接成功")
	c := &WSSConn{conn: conn}
	s.mutex.Lock()
	s.conn = append(s.conn, c)
	s.mutex.Unlock()
	go s.listen(c)
}

// listen 持续监听客户端数据
func (s *WSS) listen(conn *WSSConn) {
	for {
		if conn.conn == nil || s.server == nil {
			break
		}
		type_, data, err := conn.conn.ReadMessage()
		if err != nil {
			break // 发生错误跳出监听
		}
		if type_ == websocket.TextMessage {
			go s.call(conn, data) // 收到数据调用，不阻塞
		}
	}
}

func (s *WSS) call(conn *WSSConn, data []byte) {
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
	// 将返回的数据向conn发送
	conn.mutex.Lock()
	conn.conn.WriteMessage(websocket.TextMessage, rep)
	conn.mutex.Unlock()
}

// Send 向所有客户端发送字节数组
func (s *WSS) Send(data []byte) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
	var close []int
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
		// 记录该删除的conn
		close = append(close, i)
	}
	// 锁上conn列表并从列表上删除该conn
	for i := range close {
		s.mutex.Lock()
		s.conn = append(s.conn[:close[i]-i], s.conn[close[i]-i+1:]...)
		s.mutex.Unlock()
	}
}

// Close 关闭正向WS的连接
func (s *WSS) Close() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
	// 关闭所有的客户端连接
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
	s.mutex.Lock()
	s.conn = nil
	// 关闭WS服务器
	if s.server != nil {
		s.server.Close()
	}
	s.mutex.Unlock()
}
