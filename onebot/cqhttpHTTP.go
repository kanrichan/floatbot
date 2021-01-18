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
	"runtime"
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

func (h *HTTPYaml) apiReply(path string, data []byte) []byte {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[响应][HTTP][%v] BOT X %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			ERROR("traceback:\n%v", string(buf))
		}
	}()

	action := strings.ReplaceAll(path, "/", "")
	DEBUG("[响应][HTTP][%v] BOT <- %v:%v API: %v Params: %v", h.BotID, h.Host, h.Port, action, string(data))

	params := gjson.ParseBytes(data)
	ret := apiMap.CallApi(action, h.BotID, params)
	send, _ := json.Marshal(ret)
	return send
}

func (h *HTTPYaml) fastReply(send []byte, reply []byte) {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[快速回复][HTTP][%v] BOT X %v:%v Error: %v", h.BotID, h.Host, h.Port, err)
		}
	}()
	DEBUG("[快速回复][HTTP][%v] BOT <- %v:%v API: %v", h.BotID, h.Host, h.Port, string(reply))

	context := gjson.ParseBytes(send)
	operation := gjson.ParseBytes(reply)

	switch context.Get("path post_type").Str {
	case "message":
		switch {
		case operation.Get("reply").Exists():
			var text string
			if operation.Get("at_sender").Bool() {
				text += fmt.Sprintf("[CQ:at,qq=%s]", context.Get("user_id").Str)
			}
			text += unicode2chinese(operation.Get("reply").Str)
			data, _ := json.Marshal(
				map[string]interface{}{
					"message_type": context.Get("message_type").Str,
					"group_id":     context.Get("group_id").Int(),
					"user_id":      context.Get("user_id").Int(),
					"message":      text,
				},
			)
			apiMap.CallApi(
				"send_msg",
				h.BotID,
				gjson.ParseBytes(data),
			)
			return
		case operation.Get("delete").Bool():
			data, _ := json.Marshal(
				map[string]interface{}{
					"message_id": context.Get("message_id").Int(),
				},
			)
			apiMap.CallApi(
				"delete_msg",
				h.BotID,
				gjson.ParseBytes(data),
			)
			return
		case operation.Get("kick").Bool():
			data, _ := json.Marshal(
				map[string]interface{}{
					"group_id":           context.Get("group_id").Int(),
					"user_id":            context.Get("user_id").Int(),
					"reject_add_request": false,
				},
			)
			apiMap.CallApi(
				"set_group_kick",
				h.BotID,
				gjson.ParseBytes(data),
			)
			return
		case operation.Get("ban").Bool():
			data, _ := json.Marshal(
				map[string]interface{}{
					"group_id": context.Get("group_id").Int(),
					"user_id":  context.Get("user_id").Int(),
					"duration": context.Get("duration").Int(),
				},
			)
			apiMap.CallApi(
				"set_group_ban",
				h.BotID,
				gjson.ParseBytes(data),
			)
			return
		}
	case "request":
		if operation.Get("approve").Exists() {
			switch {
			case operation.Get("request_type").Str == "friend":
				data, _ := json.Marshal(
					map[string]interface{}{
						"flag":    context.Get("flag").Str,
						"approve": context.Get("approve").Bool(),
						"remark":  context.Get("remark").Str,
					},
				)
				apiMap.CallApi(
					"set_friend_add_request",
					h.BotID,
					gjson.ParseBytes(data),
				)
				return
			case operation.Get("request_type").Str == "group":
				data, _ := json.Marshal(
					map[string]interface{}{
						"flag":    context.Get("flag").Str,
						"approve": context.Get("approve").Bool(),
						"reason":  context.Get("reason").Str,
					},
				)
				apiMap.CallApi(
					"set_group_add_request",
					h.BotID,
					gjson.ParseBytes(data),
				)
				return
			}
		}
	}
}
