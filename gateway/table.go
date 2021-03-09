package gateway

import (
	"encoding/json"
	"fmt"
	core "onebot/core/xianqu"
	middle "onebot/middleware"
	ser "onebot/server"
)

var (
	Servers = &ServersTable{}
)

// 所有bot的所有连接的表
type ServersTable struct {
	// 所有注册的bot
	Bots []int64
	// 上报的连接列表
	Servers map[int64][]ser.Server
	// 上报数据格式
	Format map[int64][]string
}

func NewServersTable() *ServersTable {
	Servers = &ServersTable{
		Servers: make(map[int64][]ser.Server),
		Format:  make(map[int64][]string),
	}
	return Servers
}

func GetServersTable() *ServersTable {
	return Servers
}

func (t *ServersTable) Add(id int64, server ser.Server, format string) {
	func() {
		for i := range t.Bots {
			if t.Bots[i] == id {
				return
			}
		}
		// 遍历没有就添加
		t.Bots = append(t.Bots, id)
	}()
	t.Servers[id] = append(t.Servers[id], server)
	t.Format[id] = append(t.Format[id], format)
}

func (t *ServersTable) Run() {
	for i := range t.Bots {
		for j := range t.Servers[t.Bots[i]] {
			go t.Servers[t.Bots[i]][j].Run()
		}
	}
}

func (t *ServersTable) Close() {
	for i := range t.Bots {
		for j := range t.Servers[t.Bots[i]] {
			go t.Servers[t.Bots[i]][j].Close()
		}
	}
}

func (t *ServersTable) Send(id int64, ctx *core.Context) {
	dataString, _ := json.Marshal(ctx.Response)
	middle.ResponseToArray(ctx)
	dataArray, _ := json.Marshal(ctx.Response)
	for i := range t.Servers[id] {
		if t.Format[id][i] == "array" {
			t.Servers[id][i].Send(dataArray)
			continue
		}
		t.Servers[id][i].Send(dataString)
	}
}

func (t *ServersTable) SendByte(id int64, data []byte) {
	for i := range t.Servers[id] {
		t.Servers[id][i].Send(data)
	}
}

func (y *WSCYaml) GetWSCServer(id int64) (s *ser.WSC, format string) {
	s = &ser.WSC{}
	s.ID = id
	s.Addr = y.Url
	s.Token = y.AccessToken
	return s, y.PostMessageFormat
}

func (y *WSSYaml) GetWSSServer(id int64) (s *ser.WSS, format string) {
	s = &ser.WSS{}
	s.ID = id
	s.Addr = fmt.Sprintf("%s:%d", y.Host, y.Port)
	s.Token = y.AccessToken
	return s, y.PostMessageFormat
}

func (y *HTTPYaml) GetHTTPServer(id int64) (s *ser.HTTP, format string) {
	s = &ser.HTTP{}
	s.ID = id
	if y.Port == 0 {
		y.Port = 80
	}
	if y.Host != "" {
		s.Addr = fmt.Sprintf("%s:%d", y.Host, y.Port)
	}
	s.Token = y.AccessToken
	s.URL = y.PostUrl
	s.Secret = y.Secret
	return s, y.PostMessageFormat
}
