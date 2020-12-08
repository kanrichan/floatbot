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

// runOnebot run all server in config
func (conf *Yaml) runOnebot() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[OneBot] OneBot =X=> =X=> Start Error: %v", err)
			WARN("[OneBot] OneBot ==> ==> Sleep")
		}
	}()
	go Conf.heartBeat()
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
				go conf.BotConfs[i].WSCConf[k].listen()
				go conf.BotConfs[i].WSCConf[k].send()
			}
		}
		for l, _ := range conf.BotConfs[i].HTTPConf {
			if conf.BotConfs[i].HTTPConf[l].Status == 0 && conf.BotConfs[i].HTTPConf[l].Enable == true && conf.BotConfs[i].HTTPConf[l].Host != "" {
				go conf.BotConfs[i].HTTPConf[l].listen()
				go conf.BotConfs[i].HTTPConf[l].send()
			}
		}
	}
}

// heartBeat HeartBeat --> ALL PLUGINS
func (conf *Yaml) heartBeat() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[心跳] XQ =X=> =X=> Plugins Error: %v", err)
		}
	}()
	if conf.HeratBeatConf.Interval == 0 || !conf.HeratBeatConf.Enable {
		return
	}
	if conf.HeratBeatConf.Interval < 1000 {
		INFO("[心跳] Interval %v -> 1000", conf.HeratBeatConf.Interval)
		conf.HeratBeatConf.Interval = 1000
	}
	INFO("[心跳] XQ ==> ==> Plugins")
	for {
		time.Sleep(time.Millisecond * time.Duration(conf.HeratBeatConf.Interval))
		if conf.HeratBeatConf.Enable && conf.HeratBeatConf.Interval != 0 {
			for i, _ := range conf.BotConfs {
				for j, _ := range conf.BotConfs[i].WSSConf {
					if conf.BotConfs[i].WSSConf[j].Status == 1 && conf.BotConfs[i].WSSConf[j].Enable {
						conf.BotConfs[i].WSSConf[j].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for k, _ := range conf.BotConfs[i].WSCConf {
					if conf.BotConfs[i].WSCConf[k].Status == 1 && conf.BotConfs[i].WSCConf[k].Enable {
						conf.BotConfs[i].WSCConf[k].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for l, _ := range conf.BotConfs[i].HTTPConf {
					if conf.BotConfs[i].HTTPConf[l].Status == 1 && conf.BotConfs[i].HTTPConf[l].Enable {
						conf.BotConfs[i].HTTPConf[l].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
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

// 反向WS

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
		go c.wscApi(buf)
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

// 正向WS

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
		s.Conn = append(s.Conn, conn)
		s.Status = 1
		s.handShake()
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
			go s.wscApi(buf)
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

func (h *HTTPYaml) listen() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[连接][HTTP][%v] BOT =X=> ==> %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
			h.listen()
		}
	}()
	INFO("[连接][HTTP][%v] BOT ==> ==> %v:%v", h.BotID, h.Host, h.Port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", h.Host, h.Port), h)
}

func (h *HTTPYaml) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[监听][HTTP][%v] BOT =X=> ==> %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
		}
	}()

	if r.Header.Get("Authorization") == h.AccessToken || h.AccessToken == "" {
		h.Status = 1
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		w.Write(h.apiReply(r.URL.Path, buf.Bytes()))
	} else {
		WARN("[监听][HTTP][%v] BOT X Secret X %v:%v", h.BotID, h.Host, h.Port)
	}
}

func (h *HTTPYaml) send() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[上报][HTTP][%v] BOT =X=> ==> %v Error: %v", h.BotID, h.PostUrl, err)
			h.Status = 0
			h.send()
		}
	}()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	INFO("[上报][HTTP][%v] BOT ==> ==> %v", h.BotID, h.PostUrl)
	for {
		select {
		case send := <-h.Event:
			if h.PostUrl != "" {
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
					panic(err)
				} else {
					DEBUG("[上报][HTTP][%v] %v <- %v", h.BotID, h.PostUrl, string(send))
					body, _ := ioutil.ReadAll(resp.Body)
					if string(body) != "" {
						h.reply(send, body)
					}
				}
				resp.Body.Close()
			}
		case send := <-h.Heart:
			if h.PostUrl != "" {
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
					panic(err)
				}
				META("[心跳][HTTP][%v] %v <- %v", h.BotID, h.PostUrl, string(send))
				resp.Body.Close()
			}
		}
	}
}
