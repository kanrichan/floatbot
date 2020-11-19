package cqhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var WSCs []*WSC

type WSC struct {
	Enable    bool
	Bot       int64
	Status    int64
	Url       string
	Token     string
	Reconnect int64
	HeratBeat int64
	Conn      *websocket.Conn
	Send      chan []byte
	Heart     chan []byte
	Api       chan []byte
}

func WSCInit(conf *Yaml) {
	for _, c := range conf.BotConfs {
		var heartbeat int64
		if !c.HeratBeatConf.Enable {
			heartbeat = 0
		} else {
			heartbeat = c.HeratBeatConf.Interval
		}
		newWSC := &WSC{
			Enable:    c.WSCConf.Enable,
			Bot:       c.Bot,
			Status:    0,
			Url:       c.WSCConf.Url,
			Token:     c.WSCConf.AccessToken,
			Reconnect: c.WSCConf.ReconnectInterval,
			HeratBeat: heartbeat,
			Conn:      &websocket.Conn{},
			Send:      make(chan []byte, 100),
			Heart:     make(chan []byte, 1),
			Api:       make(chan []byte, 100),
		}
		WSCs = append(WSCs, newWSC)
	}
}

func WSCStarts() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[反向WS] WSC Starts()发生错误 %v，无法启动，请到GitHub提交issue", err)
		}
	}()
	for i, _ := range WSCs {
		if WSCs[i].Status == 0 {
			go WSCs[i].WSCStart()
		}
	}
}

func (c *WSC) WSCStart() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[反向WS] Bot %v Start()发生错误 %v，正在自动恢复中......", c.Bot, err)
		}
	}()

	if !c.Enable {
		return
	}
	if c.Url != "" && c.Status == 0 {
		go c.WSCListen()
		go c.WSCSend()
		go c.WSCHeartBeat()
		go c.WSCApi()
	}
}

func (c *WSC) WSCConnect() {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.Url, c.WSCHeader())
		if err != nil {
			DEBUG("[反向WS] Bot %v 与 %v 服务器连接出现错误: %v ", c.Bot, c.Url, err)
			time.Sleep(time.Millisecond * time.Duration(c.Reconnect))
			continue
		} else {
			c.Conn = conn
			c.Status = 1
			c.WSCHandShake()
			INFO("[反向WS] Bot %v 与 %v 服务器连接成功", c.Bot, c.Url)
			break
		}
	}

}

func (c *WSC) WSCHeader() http.Header {
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(c.Bot, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}
	if c.Token != "" {
		header["Authorization"] = []string{"Token " + c.Token}
	}
	return header
}

func (c *WSC) WSCHandShake() {
	handshake := map[string]string{
		"meta_event_type": "lifecycle",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(c.Bot),
		"sub_type":        "connect",
		"time":            fmt.Sprint(time.Now().Unix()),
	}
	heart, _ := json.Marshal(handshake)
	c.Heart <- []byte(string(heart))
}

func (c *WSC) WSCHeartBeat() {
	if c.HeratBeat != 0 {
		for {
			time.Sleep(time.Millisecond * time.Duration(c.HeratBeat))
			heartbeat := map[string]string{
				"interval":        fmt.Sprint(c.HeratBeat),
				"meta_event_type": "heartbeat",
				"post_type":       "meta_event",
				"self_id":         fmt.Sprint(c.Bot),
				"status":          "null",
				"time":            fmt.Sprint(time.Now().Unix()),
			}
			heart, _ := json.Marshal(heartbeat)
			if c.Status == 1 {
				c.Heart <- []byte(string(heart))
			}
		}
	}
}

func (c *WSC) WSCListen() {
	defer func() {
		if err := recover(); err != nil {
			c.Status = 0
			ERROR("[监听服务] Bot %v 服务发生错误，正在自动恢复中...... %v，", c.Bot, err)
			c.WSCListen()
		}
	}()
	c.WSCConnect()
	// 等待wsc连接成功
	for {
		if c.Status == 1 {
			break
		}
	}
	// 监听wsc
	DEBUG("[监听服务] Bot %v 开始监听...... ", c.Bot)
	for {
		_, buf, err := c.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		c.Api <- buf
	}
}

func (c *WSC) WSCSend() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[上报服务] Bot %v 服务发生错误，正在自动恢复中...... %v，", c.Bot, err)
			c.WSCSend()
		}
	}()
	// 等待wsc连接成功
	for {
		if c.Status == 1 {
			break
		}
	}
	DEBUG("[上报服务] Bot %v 服务开始启动...... ", c.Bot)
	for {
		select {
		case send := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			err := c.Conn.WriteMessage(websocket.TextMessage, send)
			if err != nil {
				panic(err)
			} else {
				DEBUG("[上报服务] Bot %v 上报至 %v ：%v", c.Bot, c.Url, string(send))
			}
		case heart := <-c.Heart:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			err := c.Conn.WriteMessage(websocket.TextMessage, heart)
			if err != nil {
				panic(err)
			} else {
				DEBUG("[心跳服务] Bot %v 连接至 %v ：%v", c.Bot, c.Url, string(heart))
			}
		}
	}
}
