package onebot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (s *WSSYaml) start() {
	http.ListenAndServe(fmt.Sprintf("%v:%v", s.Host, s.Port), s)
}

func (s *WSSYaml) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[连接][正向WS][%v] BOT =X=> ==> %v:%v Error: %v", s.BotID, s.Host, s.Port, err)
		}
	}()
	if r.Header.Get("Authorization") == "Token "+s.AccessToken || s.AccessToken == "" {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		if r.URL.Path != "/api" {
			s.Conn = append(s.Conn, conn)
			s.Status = 1
			s.handShake()
		}
	} else {
		ERROR("[连接][正向WS][%v] BOT X Token X %v:%v", s.BotID, s.Host, s.Port)
	}
}

func (s *WSSYaml) handShake() {
	handshake := map[string]string{
		"meta_event_type": "lifecycle",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(s.BotID),
		"sub_type":        "connect",
		"time":            fmt.Sprint(time.Now().Unix()),
	}
	event, _ := json.Marshal(handshake)
	if s.Status == 1 {
		s.Event <- event
	}
}

func (s *WSSYaml) listen() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[监听][正向WS][%v] BOT =X=> ==> %v:%v Error: %v", s.BotID, s.Host, s.Port, err)
			s.Status = 0
			s.listen()
		}
	}()
	for {
		if s.Status == 1 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	INFO("[监听][正向WS][%v] BOT ==> ==> %v:%v", s.BotID, s.Host, s.Port)
	for {
		for i, conn := range s.Conn {
			_, buf, err := conn.ReadMessage()
			if err != nil {
				s.Conn = append(s.Conn[:i], s.Conn[i+1:]...)
				if len(s.Conn) < 1 {
					s.Status = 0
				}
				panic(err)
			}
			go s.apiReply(buf)
		}
	}
}

func (s *WSSYaml) send() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[上报][正向WS][%v] BOT =X=> ==> %v:%v Error: %v", s.BotID, s.Host, s.Port, err)
			s.Status = 0
			s.send()
		}
	}()
	for {
		if s.Status == 1 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	INFO("[上报][正向WS][%v] BOT ==> ==> %v:%v", s.BotID, s.Host, s.Port)
	for {
		select {
		case send := <-s.Event:
			for i, conn := range s.Conn {
				_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
				err := conn.WriteMessage(websocket.TextMessage, send)
				if err != nil {
					s.Conn = append(s.Conn[:i], s.Conn[i+1:]...)
					if len(s.Conn) < 1 {
						s.Status = 0
					}
					panic(err)
				} else {
					DEBUG("[上报][正向WS][%v] %v:%v <- %v", s.BotID, s.Host, s.Port, string(send))
				}
			}
		case send := <-s.Heart:
			for i, conn := range s.Conn {
				_ = conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
				err := conn.WriteMessage(websocket.TextMessage, send)
				if err != nil {
					s.Conn = append(s.Conn[:i], s.Conn[i+1:]...)
					if len(s.Conn) < 1 {
						s.Status = 0
					}
					panic(err)
				} else {
					META("[心跳][正向WS][%v] %v:%v <- %v", s.BotID, s.Host, s.Port, string(send))
				}
			}
		}
	}
}

func (s *WSSYaml) apiReply(api []byte) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[响应][HTTP][%v] BOT X %v:%v Error: %v", s.BotID, s.Host, s.Port, err)
		}
	}()

	req := gjson.ParseBytes(api)
	action := strings.ReplaceAll(req.Get("action").Str, "_async", "")
	params := req.Get("params")
	DEBUG("[响应][HTTP][%v] BOT <- %v:%v API: %v Params: %v", s.BotID, s.Host, s.Port, action, string(api))

	if f, ok := apiList[action]; ok {
		ret := tieEcho(f(s.BotID, params), req)
		send, _ := json.Marshal(ret)
		s.Event <- send
	} else {
		ret := tieEcho(resultFail("no such api"), req)
		send, _ := json.Marshal(ret)
		s.Event <- send
	}
}
