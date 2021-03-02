package server

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
)

var (
	HttpPostHandler = func(bot int64, send []byte, data []byte) { fmt.Println(string(data)) }
	HttpHandler     = func(bot int64, path string, data []byte) []byte { fmt.Println(string(data)); return []byte("ok") }
)

type HttpServer struct {
	// 参数
	id     int64
	addr   string
	token  string
	url    string
	secret string
	server *http.Server
}

func (s *HttpServer) Run(id int64, addr, token string) {
	defer func() {
		recover()
	}()
	s.id = id
	s.addr = addr
	s.token = token
	s.server = &http.Server{
		Addr:    addr,
		Handler: s,
	}
	if err := s.server.ListenAndServe(); err != nil {
		panic(err)
	}
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !(r.Header.Get("Authorization") == s.token || s.token == "") {
		// Token验证失败
		return
	}
	switch r.Header.Get("Content-Type") {
	case "application/json":
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		// 直接提交并返回数据
		ret := HttpHandler(s.id, r.URL.Path, buf.Bytes())
		w.Write(ret)
	default:
		r.ParseForm()
		dataMap := make(map[string]interface{})
		for k, v := range r.Form {
			dataMap[k] = v[0]
		}
		data, _ := json.Marshal(dataMap)
		// 直接提交并返回数据
		ret := HttpHandler(s.id, r.URL.Path, data)
		w.Write(ret)
	}
}

func (s *HttpServer) Send(data []byte) {
	client := &http.Client{}
	// TODO OneBot标准 HTTP POST 上报Header
	// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/http-post.md#%E4%B8%8A%E6%8A%A5
	req, _ := http.NewRequest("POST", s.url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Self-ID", strconv.FormatInt(s.id, 10))
	req.Header.Set("User-Agent", "CQHttp/4.15.0")
	if s.secret != "" {
		// TODO OneBot标准 HTTP POST 签名
		// https://github.com/howmanybots/onebot/blob/master/v11/specs/communication/http-post.md#%E7%AD%BE%E5%90%8D
		mac := hmac.New(sha1.New, []byte(s.secret))
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
		HttpPostHandler(s.id, data, body)
	}
	resp.Body.Close()
}
