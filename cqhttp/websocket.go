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
var CQHttpOK bool

type WSC struct {
	Enable    bool
	Bot       int64
	Url       string
	Token     string
	Reconnect int64
	HeratBeat int64
	Conn      *websocket.Conn
	Send      chan []byte
	Api       chan []byte
	Quit      chan bool
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
			Bot:       c.Bot,
			Enable:    c.WSCConf.Enable,
			Url:       c.WSCConf.Url,
			Token:     c.WSCConf.AccessToken,
			Reconnect: c.WSCConf.ReconnectInterval,
			HeratBeat: heartbeat,
			Conn:      &websocket.Conn{},
			Send:      make(chan []byte, 100),
			Api:       make(chan []byte, 100),
			Quit:      make(chan bool),
		}
		WSCs = append(WSCs, newWSC)
	}
	CQHttpOK = true
}

func WSCStarts() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[反向WS] WSC Starts()发生错误 %v，无法启动，请到GitHub提交issue", err)
		}
	}()
	for i, _ := range WSCs {
		go WSCs[i].WSCStart()
		go WSCs[i].WSCApi()
		go WSCs[i].WSCHeartBeat()
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
	if c.Url != "" {
		c.WSCConnect()
		go c.WSCListen()
		go c.WSCSend()
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
	send, _ := json.Marshal(handshake)
	c.Send <- []byte(string(send))
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
			send, _ := json.Marshal(heartbeat)
			c.Send <- []byte(string(send))
		}
	}
}

func (c *WSC) WSCListen() {
	defer func() {
		go c.WSCStart()
		DEBUG("[监听服务] Bot %v 服务开始自闭并重新启动...... ", c.Bot)
	}()

	defer func() {
		c.Quit <- true
		_ = c.Conn.Close()
		if err := recover(); err != nil {
			ERROR("[监听服务] Bot %v 服务发生错误，正在自动恢复中...... %v，", c.Bot, err)
		}
	}()

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
		}
	}()

	DEBUG("[上报服务] Bot %v 服务开始启动...... ", c.Bot)
	for {
		select {
		case quit := <-c.Quit:
			if quit {
				DEBUG("[上报服务] Bot %v 服务被动自闭...... ", c.Bot)
				return
			}
		case send := <-c.Send:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			err := c.Conn.WriteMessage(websocket.TextMessage, send)
			if err != nil {
				panic(err)
			} else {
				DEBUG("[上报服务] Bot %v 上报至 %v ：%v", c.Bot, c.Url, string(send))
			}
		}
	}
}
