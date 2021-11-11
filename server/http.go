package server

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

// HTTP & HTTP POST
type HTTP struct {
	// Bot的qq号
	ID int64
	// 监听的地址
	Addr  string
	Token string
	// 上报数据的地址
	URL    string
	Secret string

	server *http.Server
}

func (s *HTTP) Run() {
	defer func() {
		recover()
	}()
	if s.Addr == "" {
		return
	}
	s.server = &http.Server{Addr: s.Addr, Handler: s}
	s.INFO("HTTP服务建立，等待API调用")
	if err := s.server.ListenAndServe(); err != nil {
		s.ERROR(err)
	}
}

func (s *HTTP) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !(r.Header.Get("Authorization") == s.Token || s.Token == "") {
		// Token验证失败
		return
	}
	if r.URL.Path == "/favicon.ico" {
		return
	}
	switch r.Header.Get("Content-Type") {
	case "application/json":
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		// 直接提交并返回数据
		if len(buf.Bytes()) == 0 {
			return
		}
		ret := HttpHandler(s.ID, r.URL.Path[1:], buf.Bytes())
		w.Write(ret)
	default:
		r.ParseForm()
		dataMap := make(map[string]interface{})
		for k, v := range r.Form {
			dataMap[k] = v[0]
		}
		data, _ := json.Marshal(dataMap)
		// 直接提交并返回数据
		if len(data) == 0 {
			return
		}
		ret := HttpHandler(s.ID, r.URL.Path[1:], data)
		w.Write(ret)
	}
}

// Send 向 POST_URL 发送上报数据
func (s *HTTP) Send(data []byte) {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
	if s.URL == "" {
		return
	}
	client := &http.Client{}
	// TODO OneBot标准 HTTP POST 上报Header
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/http-post.md#%E4%B8%8A%E6%8A%A5
	req, _ := http.NewRequest("POST", s.URL, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Self-ID", strconv.FormatInt(s.ID, 10))
	req.Header.Set("User-Agent", "CQHttp/4.15.0")
	if s.Secret != "" {
		// TODO OneBot标准 HTTP POST 签名
		// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/http-post.md#%E7%AD%BE%E5%90%8D
		mac := hmac.New(sha1.New, helper.StringToBytes(s.Secret))
		mac.Write(data)
		req.Header.Set("X-Signature", "sha1="+hex.EncodeToString(mac.Sum(nil)))
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	if len(body) != 0 {
		// 快速回复
		HttpPostHandler(s.ID, data, body)
	}
	resp.Body.Close()
}

// Close 关闭HTTP监听
func (s *HTTP) Close() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			s.PANIC(err, buf)
		}
	}()
	if s.server != nil {
		s.server.Close()
	}
}
