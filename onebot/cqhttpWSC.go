package onebot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
)

func (c *WSCYaml) connect() {
	if c.ReconnectInterval < 1000 {
		INFO("[连接][反向WS][%v] ReconnectInterval %v -> 1000", c.BotID, c.ReconnectInterval)
		c.ReconnectInterval = 1000
	}
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(c.BotID, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}
	if c.AccessToken != "" {
		header["Authorization"] = []string{"Token " + c.AccessToken}
	}
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.Url, header)
		if err != nil {
			DEBUG("[连接][反向WS][%v] BOT =X=> ==> %v ", c.BotID, c.Url)
			time.Sleep(time.Millisecond * time.Duration(c.ReconnectInterval))
			continue
		} else {
			c.Conn = conn
			c.Status = 1
			c.handShake()
			INFO("[连接][反向WS][%v] BOT ==> ==> %v ", c.BotID, c.Url)
			break
		}
	}

}

func (c *WSCYaml) handShake() {
	handshake := map[string]string{
		"meta_event_type": "lifecycle",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(c.BotID),
		"sub_type":        "connect",
		"time":            fmt.Sprint(time.Now().Unix()),
	}
	event, _ := json.Marshal(handshake)
	if c.Status == 1 {
		c.Heart <- event
	}
}

func (c *WSCYaml) listen() {
	defer func() {
		if err := recover(); err != nil {
			c.Status = 0
			WARN("[监听][反向WS][%v] BOT =X=> ==> %v ERROR: %v", c.BotID, c.Url, err)
			c.listen()
		}
	}()
	c.connect()
	// 等待wsc连接成功
	for {
		if c.Status == 1 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	// 监听wsc
	INFO("[监听][反向WS][%v] BOT ==> ==> %v ", c.BotID, c.Url)
	for {
		_, buf, err := c.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		go c.apiReply(buf)
	}
}

func (c *WSCYaml) send() {
	defer func() {
		if err := recover(); err != nil {
			WARN("[上报][反向WS][%v] BOT =X=> ==> %v ERROR: %v", c.BotID, c.Url, err)
			c.send()
		}
	}()
	// 等待wsc连接成功
	for {
		if c.Status == 1 {
			break
		}
		time.Sleep(time.Second * 1)
	}
	INFO("[上报][反向WS][%v] BOT ==> ==> %v ", c.BotID, c.Url)
	for {
		select {
		case send := <-c.Event:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			err := c.Conn.WriteMessage(websocket.TextMessage, send)
			if err != nil {
				panic(err)
			} else {
				DEBUG("[上报][反向WS][%v] %v <- %v", c.BotID, c.Url, string(send))
			}
		case send := <-c.Heart:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			err := c.Conn.WriteMessage(websocket.TextMessage, send)
			if err != nil {
				panic(err)
			} else {
				META("[心跳][反向WS][%v] %v <- %v", c.BotID, c.Url, string(send))
			}
		}
	}
}

func (c *WSCYaml) apiReply(api []byte) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[响应][HTTP][%v] BOT X %v Error: %v", c.BotID, c.Url, err)
		}
	}()

	req := gjson.ParseBytes(api)
	action := strings.ReplaceAll(req.Get("action").Str, "_async", "")
	params := req.Get("params")
	DEBUG("[响应][HTTP][%v] BOT <- %v API: %v Params: %v", c.BotID, c.Url, action, string(api))

	if f, ok := apiList[action]; ok {
		ret := tieEcho(f(c.BotID, params), req)
		send, _ := json.Marshal(ret)
		c.Event <- send
	} else {
		ret := tieEcho(resultFail("no such api"), req)
		send, _ := json.Marshal(ret)
		c.Event <- send
	}
}
