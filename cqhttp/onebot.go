package cqhttp

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

func (conf *Yaml) runOnebot() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[OneBot] runOnebot()发生不可自动修复错误，无法启动OneBot，请到GitHub提交issue %v", err)
		}
	}()
	for i, _ := range conf.BotConfs {
		for j, _ := range conf.BotConfs[i].WSSConf {
			if conf.BotConfs[i].WSSConf[j].Status == 0 && conf.BotConfs[i].WSSConf[j].Enable == true && conf.BotConfs[i].WSSConf[j].Host != "" {
				go conf.BotConfs[i].WSSConf[j].start()
				go conf.BotConfs[i].WSSConf[j].listen()
				go conf.BotConfs[i].WSSConf[j].send()
			}
		}
		for k, _ := range conf.BotConfs[i].WSCConf {
			if conf.BotConfs[i].WSCConf[k].Status == 0 && conf.BotConfs[i].WSCConf[k].Enable == true && conf.BotConfs[i].WSCConf[k].Url != "" {
				go conf.BotConfs[i].WSCConf[k].WSCListen()
				go conf.BotConfs[i].WSCConf[k].WSCSend()
			}
		}
		for l, _ := range conf.BotConfs[i].HTTPConf {
			if conf.BotConfs[i].HTTPConf[l].Status == 0 && conf.BotConfs[i].HTTPConf[l].Enable == true && conf.BotConfs[i].HTTPConf[l].Host != "" {
				//go conf.BotConfs[i].HTTPConf[l].listen()
				go conf.BotConfs[i].HTTPConf[l].send()
			}
		}
	}
}

func (conf *Yaml) heartBeat() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[OneBot] heartBeat()发生不可自动修复错误，心跳停止，请到GitHub提交issue %v", err)
		}
	}()
	DEBUG("[心跳服务] 开始对已连接端发送心跳...... ")
	for {
		time.Sleep(time.Millisecond * time.Duration(conf.HeratBeatConf.Interval))
		if conf.HeratBeatConf.Enable && conf.HeratBeatConf.Interval != 0 {
			for i, _ := range conf.BotConfs {
				for j, _ := range conf.BotConfs[i].WSSConf {
					if conf.BotConfs[i].WSSConf[j].Status == 1 {
						conf.BotConfs[i].WSSConf[j].Event <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for k, _ := range conf.BotConfs[i].WSCConf {
					if conf.BotConfs[i].WSCConf[k].Status == 1 {
						conf.BotConfs[i].WSCConf[k].Event <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for l, _ := range conf.BotConfs[i].HTTPConf {
					if conf.BotConfs[i].HTTPConf[l].Status == 1 {
						conf.BotConfs[i].HTTPConf[l].Event <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
			}
		}
	}
}

func heartEvent(interval int64, bot int64) []byte {
	heartbeat := map[string]string{
		"interval":        fmt.Sprint(interval),
		"meta_event_type": "heartbeat",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(bot),
		"status":          "null",
		"time":            fmt.Sprint(time.Now().Unix()),
	}
	event, _ := json.Marshal(heartbeat)
	return event
}

func (c *WSCYaml) WSCConnect() {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(c.Url, c.WSCHeader())
		if err != nil {
			DEBUG("[反向WS] Bot %v 与 %v 服务器连接出现错误: %v ", c.BotID, c.Url, err)
			time.Sleep(time.Millisecond * time.Duration(c.ReconnectInterval))
			continue
		} else {
			c.Conn = conn
			c.Status = 1
			c.WSCHandShake()
			INFO("[反向WS] Bot %v 与 %v 服务器连接成功", c.BotID, c.Url)
			break
		}
	}

}

func (c *WSCYaml) WSCHeader() http.Header {
	header := http.Header{
		"X-Client-Role": []string{"Universal"},
		"X-Self-ID":     []string{strconv.FormatInt(c.BotID, 10)},
		"User-Agent":    []string{"CQHttp/4.15.0"},
	}
	if c.AccessToken != "" {
		header["Authorization"] = []string{"Token " + c.AccessToken}
	}
	return header
}

func (c *WSCYaml) WSCHandShake() {
	handshake := map[string]string{
		"meta_event_type": "lifecycle",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(c.BotID),
		"sub_type":        "connect",
		"time":            fmt.Sprint(time.Now().Unix()),
	}
	event, _ := json.Marshal(handshake)
	if c.Status == 1 {
		c.Event <- event
	}
}

func (c *WSCYaml) WSCListen() {
	defer func() {
		if err := recover(); err != nil {
			c.Status = 0
			ERROR("[监听服务] Bot %v 监听 %v 服务发生错误，正在自动恢复中...... %v，", c.BotID, c.Url, err)
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
	DEBUG("[监听服务] Bot %v 开始监听 %v ...... ", c.BotID, c.Url)
	for {
		_, buf, err := c.Conn.ReadMessage()
		if err != nil {
			panic(err)
		}
		go c.wscApi(buf)
	}
}

func (c *WSCYaml) WSCSend() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[上报服务] Bot %v 向 %v 上报发生错误，正在自动恢复中...... %v，", c.BotID, c.Url, err)
			c.WSCSend()
		}
	}()
	// 等待wsc连接成功
	for {
		if c.Status == 1 {
			break
		}
	}
	DEBUG("[上报服务] Bot %v 向 %v 开始推送...... ", c.BotID, c.Url)
	for {
		select {
		case send := <-c.Event:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(time.Second * 15))
			err := c.Conn.WriteMessage(websocket.TextMessage, send)
			if err != nil {
				panic(err)
			} else {
				DEBUG("[上报服务] Bot %v 向 %v 上报：%v", c.BotID, c.Url, string(send))
			}
		}
	}
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

func (s *WSSYaml) start() {
	http.ListenAndServe(fmt.Sprintf("%v:%v", s.Host, s.Port), s)
}

func (s *WSSYaml) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[正向WS] Bot %v 在 %v:%v 升级发生错误 %v，", s.BotID, s.Host, s.Port, err)
		}
	}()
	if r.Header.Get("Authorization") == "Token "+s.AccessToken || s.AccessToken == "" {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			panic(err)
		}
		s.Conn = append(s.Conn, conn)
		s.Status = 1
		s.handShake()
	} else {
		WARN("[正向WS] Bot %v 拒绝与 %v:%v 连接 token error", s.BotID, s.Host, s.Port)
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
			ERROR("[监听服务] Bot %v 在 %v:%v 监听发生错误 %v，", s.BotID, s.Host, s.Port, err)
			s.listen()
		}
	}()
	for {
		if s.Status == 1 {
			break
		}
	}
	DEBUG("[监听服务] Bot %v 在 %v:%v 开始监听...... ", s.BotID, s.Host, s.Port)
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
			go s.wscApi(buf)
		}
	}
}

func (s *WSSYaml) send() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[上报服务] Bot %v 在 %v:%v 上报发生错误 %v，", s.BotID, s.Host, s.Port, err)
			s.send()
		}
	}()
	for {
		if s.Status == 1 {
			break
		}
	}
	DEBUG("[上报服务] Bot %v 在 %v:%v 开始推送...... ", s.BotID, s.Host, s.Port)
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
					DEBUG("[上报服务] Bot %v 在 %v:%v 上报：%v", s.BotID, s.Host, s.Port, string(send))
				}
			}
		}
	}
}

func (h *HTTPYaml) listen() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[监听服务] Bot %v 在 %v:%v 监听发生错误 %v，", h.BotID, h.Host, h.Port, err)
			h.listen()
		}
	}()
	DEBUG("[监听服务] Bot %v 在 %v:%v 开始监听...... ", h.BotID, h.Host, h.Port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", h.Host, h.Port), h)
}

func (h *HTTPYaml) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[HTTP] Bot %v 在 %v:%v 监听发生错误 %v，", h.BotID, h.Host, h.Port, err)
		}
	}()

	if r.Header.Get("Authorization") == h.AccessToken || h.AccessToken == "" {
		h.Status = 1
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		w.Write(h.wscApi(r.URL.Path, buf.Bytes()))
	} else {
		WARN("[正向WS] Bot %v 拒绝与 %v:%v 连接 secret error", h.BotID, h.Host, h.Port)
	}
}

func (h *HTTPYaml) send() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[上报服务] Bot %v 在 %v 上报发生错误 %v，", h.BotID, h.PostUrl, err)
			h.Status = 0
			h.send()
		}
	}()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	h.Status = 1
	DEBUG("[上报服务] Bot %v 在 %v 开始推送...... ", h.BotID, h.PostUrl)
	for {
		select {
		case send := <-h.Event:
			req, _ := http.NewRequest("POST", h.PostUrl, bytes.NewBuffer(send))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("X-Self-ID", strconv.FormatInt(h.BotID, 10))
			req.Header.Set("User-Agent", "CQHttp/4.15.0")
			if h.Secret != "" {
				mac := hmac.New(sha1.New, []byte(h.Secret))
				mac.Write(send)
				req.Header.Set("X-Signature", "sha1="+hex.EncodeToString(mac.Sum(nil)))
			}
			resp, err := client.Do(req)
			if err != nil {
				ERROR("[HTTP POST] Bot %v 在 %v 上报出错", h.BotID, h.PostUrl)
			} else {
				DEBUG("[HTTP POST] Bot %v 在 %v 上报：%v", h.BotID, h.PostUrl, string(send))
				body, _ := ioutil.ReadAll(resp.Body)
				if string(body) != "" {
					h.reply(send, body)
				}
			}
			resp.Body.Close()
		}
	}
}
