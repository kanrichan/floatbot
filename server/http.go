package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HttpServer struct {
	ID                int64  `yaml:"-"`
	Name              string `yaml:"name"`
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	AccessToken       string `yaml:"token"`
	PostUrl           string `yaml:"post_url"`
	Secret            string `yaml:"secret"`
	TimeOut           int64  `yaml:"time_out"`
	PostMessageFormat string `yaml:"post_message_format"`
	HeartBeatInterval int64  `yaml:"reconnect_interval"`
	ReconnectInterval int64  `yaml:"reconnect_interval"`

	StopSend      chan bool `yaml:"-"`
	StopHeartBeat chan bool `yaml:"-"`

	ConnectStatus   string `yaml:"-"`
	SendStatus      string `yaml:"-"`
	HeartBeatStatus string `yaml:"-"`

	SendChan chan []byte `yaml:"-"`
}

func HttpCall(path string, req []byte) (ret []byte) {
	return req
}

func FastReply(req, ret []byte) {
	//
}

func (this *HttpServer) Connect() {
	this.ConnectStatus = "ok"
	defer func() {
		this.ConnectStatus = "wait"
		recover()
	}()
	http.ListenAndServe(fmt.Sprintf("%v:%v", this.Host, this.Port), this)
}

func (this *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			//
		}
	}()
	if r.Header.Get("Authorization") == this.AccessToken || this.AccessToken == "" {
		if r.Header.Get("Content-Type") == "application/json" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			ret := HttpCall(r.URL.Path, buf.Bytes())
			w.Write(ret)
		} else {
			r.ParseForm()
			dataMap := make(map[string]interface{})
			for k, v := range r.Form {
				dataMap[k] = v[0]
			}
			data, _ := json.Marshal(dataMap)
			ret := HttpCall(r.URL.Path, data)
			w.Write(ret)
		}
	} else {
		panic(errors.New("token error"))
	}
}

func (this *HttpServer) Send() {
	this.SendStatus = "ok"
	defer func() {
		this.SendStatus = "wait"
		recover()
	}()
	// TODO 元事件 表示OneBot启动成功
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E7%94%9F%E5%91%BD%E5%91%A8%E6%9C%9F
	handshake, _ := json.Marshal(map[string]string{
		"meta_event_type": "lifecycle",
		"post_type":       "meta_event",
		"self_id":         strconv.FormatInt(this.ID, 10),
		"sub_type":        "connect",
		"time":            strconv.FormatInt(time.Now().Unix(), 10),
	})
	this.SendChan <- handshake

	client := &http.Client{}
	for {
		select {
		case send := <-this.SendChan:
			if this.PostUrl != "" {
				// TODO OneBot标准 HTTP POST 上报Header
				// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/http-post.md#%E4%B8%8A%E6%8A%A5
				req, _ := http.NewRequest("POST", this.PostUrl, bytes.NewBuffer(send))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Self-ID", strconv.FormatInt(this.ID, 10))
				req.Header.Set("User-Agent", "CQHttp/4.15.0")
				if this.Secret != "" {
					// TODO OneBot标准 HTTP POST 签名
					// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/http-post.md#%E7%AD%BE%E5%90%8D
					mac := hmac.New(sha1.New, []byte(this.Secret))
					mac.Write(send)
					req.Header.Set("X-Signature", "sha1="+hex.EncodeToString(mac.Sum(nil)))
				}

				resp, err := client.Do(req)
				if err != nil {
					panic(err)
				}
				body, _ := ioutil.ReadAll(resp.Body)
				if len(body) != 0 {
					FastReply(send, body)
				}
				resp.Body.Close()
			}
		case <-this.StopSend:
			return
		}
	}
}

func (this *HttpServer) HeartBeat() {
	this.HeartBeatStatus = "ok"
	for {
		select {
		case <-time.After(time.Second * time.Duration(this.HeartBeatInterval)):
			// TODO OneBot标准 元事件 心跳
			// https://github.com/howmanybots/onebot/blob/master/v11/specs/event/meta.md#%E5%BF%83%E8%B7%B3
			heartbeat, _ := json.Marshal(
				map[string]interface{}{
					"interval":        strconv.FormatInt(this.HeartBeatInterval, 10),
					"meta_event_type": "heartbeat",
					"post_type":       "meta_event",
					"self_id":         strconv.FormatInt(this.ID, 10),
					"status": map[string]interface{}{
						"online": true,
						"good":   true,
					},
					"time": strconv.FormatInt(time.Now().Unix(), 10),
				},
			)
			this.SendChan <- heartbeat
		case <-this.StopHeartBeat:
			this.HeartBeatStatus = "wait"
			return
		}
	}
}
