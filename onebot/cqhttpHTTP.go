package onebot

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
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

func (h *HTTPYaml) listen() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[连接][HTTP][%v] BOT =X=> ==> %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
			time.Sleep(time.Second * 1)
			h.listen()
		}
	}()
	INFO("[连接][HTTP][%v] BOT ==> ==> %v:%v", h.BotID, h.Host, h.Port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", h.Host, h.Port), h)
	time.Sleep(time.Second * 1)
	h.listen()
}

func (h *HTTPYaml) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[监听][HTTP][%v] BOT =X=> ==> %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
		}
	}()

	if r.Header.Get("Authorization") == h.AccessToken || h.AccessToken == "" {
		if r.Header.Get("Content-Type") == "application/json" {
			buf := new(bytes.Buffer)
			buf.ReadFrom(r.Body)
			w.Write(h.apiReply(r.URL.Path, buf.Bytes()))
		} else {
			r.ParseForm()
			dataMap := make(map[string]interface{})
			for k, v := range r.Form {
				dataMap[k] = v[0]
			}
			data, _ := json.Marshal(dataMap)
			w.Write(h.apiReply(r.URL.Path, data))
		}
	} else {
		WARN("[监听][HTTP][%v] BOT X Secret X %v:%v", h.BotID, h.Host, h.Port)
	}
}

func (h *HTTPYaml) send() {
	defer func() {
		if err := recover(); err != nil {
			h.Status = 0
			ERROR("[上报][HTTP][%v] BOT =X=> ==> %v Error: %v", h.BotID, h.PostUrl, err)
			time.Sleep(time.Second * 1)
			h.send()
		}
	}()
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}
	INFO("[上报][HTTP][%v] BOT ==> ==> %v", h.BotID, h.PostUrl)
	h.Status = 1
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
						h.fastReply(send, body)
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

func (h *HTTPYaml) apiReply(path string, api []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[响应][HTTP][%v] BOT X %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
		}
	}()

	action := strings.ReplaceAll(path, "/", "")
	params := gjson.ParseBytes(api)
	DEBUG("[响应][HTTP][%v] BOT <- %v:%v API: %v Params: %v", h.BotID, h.Host, h.Port, action, string(api))

	if f, ok := apiList[action]; ok {
		ret := f(h.BotID, params)
		send, _ := json.Marshal(ret)
		return send
	} else {
		ret := resultFail("no such api")
		send, _ := json.Marshal(ret)
		return send
	}
}

func (h *HTTPYaml) fastReply(send []byte, reply []byte) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[快速回复][HTTP][%v] BOT X %v:%v Error: %v", h.BotID, h.Host, err)
		}
	}()
	DEBUG("[快速回复][HTTP][%v] BOT <- %v:%v API: %v", h.BotID, h.Host, h.Port, string(reply))

	req := gjson.ParseBytes(send)
	res := gjson.ParseBytes(reply)

	params := []map[string]interface{}{}
	if res.Get("at_sender").Bool() {
		elem := map[string]interface{}{
			"type": "at",
			"data": map[string]interface{}{
				"qq": req.Get("user_id").Int(),
			},
		}
		params = append(params, elem)
	}
	if res.Get("reply").Str != "" {
		elem := map[string]interface{}{
			"type": "text",
			"data": map[string]interface{}{
				"text": res.Get("reply"),
			},
		}
		params = append(params, elem)
	}
	p := map[string]interface{}{
		"message_type": req.Get("message_type").Str,
		"group_id":     req.Get("group_id").Int(),
		"user_id":      req.Get("user_id").Int(),
		"message":      params,
	}
	data, _ := json.Marshal(p)
	cq2xqSendMsg(h.BotID, gjson.Parse(string(data)))
}
