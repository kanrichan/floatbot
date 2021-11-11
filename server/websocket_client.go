package server

import (
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// 反向WS
type WSC struct {
	// Bot的qq号
	ID    int64
	Addr  string
	Token string

	mutex sync.Mutex
	conn  *websocket.Conn
	// 判断是否可以上报数据
	status bool
	// 关闭重连的goroutine
	stopconnect chan bool
}

// Run 建立反向WS
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
			// 立刻停止重连
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
			handshake := fmt.Sprintf(`{"meta_event_type":"lifecycle","post_type":"meta_event","self_id":%d,"sub_type":"connect","time":%d}`,
				s.ID, time.Now().Unix())
			if err := conn.WriteMessage(websocket.TextMessage, helper.StringToBytes(handshake)); err != nil {
				s.DEBUG(err)
				continue
			}
			// 连接成功
			s.conn = conn
			s.status = true
			s.INFO("反向WS连接服务器成功")
			s.listen() // 阻塞至监听错误，错误后继续重连
			s.INFO("反向WS掉线")
		}
	}
}

// listen 持续监听服务端数据
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
			go s.call(data)
		}
	}
	s.mutex.Lock()
	s.conn.Close()
	s.conn = nil
	s.status = false
	s.mutex.Unlock()
}

func (s *WSC) call(data []byte) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
	rep := WSCHandler(s.ID, data)
	s.Send(rep)
}

// Send 向服务端发送字节数组
func (s *WSC) Send(data []byte) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
	if !s.status {
		return
	}
	if s.conn != nil {
		s.mutex.Lock()
		err := s.conn.WriteMessage(websocket.TextMessage, data)
		s.mutex.Unlock()
		if err != nil {
			s.ERROR(err)
		}
	}
}

// 关闭反向WS的连接
func (s *WSC) Close() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
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
