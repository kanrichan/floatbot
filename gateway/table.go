package gateway

import (
	"encoding/json"
	"fmt"

	core "onebot/core/xianqu"
	middle "onebot/middleware"
	ser "onebot/server"
)

var (
	// 当前的table
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

// NewServersTable 新建一个table
func NewServersTable() *ServersTable {
	Servers = &ServersTable{
		Servers: make(map[int64][]ser.Server),
		Format:  make(map[int64][]string),
	}
	return Servers
}

// GetServersTable 返回当前的table
func GetServersTable() *ServersTable {
	return Servers
}

// Add 向table里面增加一个连接
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

// Run 建立table里面所有qq的所有连接
func (t *ServersTable) Run() {
	for i := range t.Bots {
		for j := range t.Servers[t.Bots[i]] {
			go t.Servers[t.Bots[i]][j].Run()
		}
	}
}

// Close 关闭table里面所有qq的所有连接
func (t *ServersTable) Close() {
	for i := range t.Bots {
		for j := range t.Servers[t.Bots[i]] {
			go t.Servers[t.Bots[i]][j].Close()
		}
	}
}

// Send 向table里面所有qq的所有连接发送上报
func (t *ServersTable) Send(id int64, ctx *core.Context) {
	dataString, _ := json.Marshal(ctx.Response)
	middle.ResponseToArray(ctx)
	dataArray, _ := json.Marshal(ctx.Response)
	for i := range t.Servers[id] {
		if t.Format[id][i] == "array" {
			go t.Servers[id][i].Send(dataArray)
			return
		}
		go t.Servers[id][i].Send(dataString)
	}
}

// SendByte 向table里面所有qq的所有连接发送数据
func (t *ServersTable) SendByte(id int64, data []byte) {
	for i := range t.Servers[id] {
		t.Servers[id][i].Send(data)
	}
}

// GetWSCServer 返回 WSCServer
func (y *WSCYaml) GetWSCServer(id int64) (s *ser.WSC, format string) {
	s = &ser.WSC{}
	s.ID = id
	s.Addr = y.Url
	s.Token = y.AccessToken
	return s, y.PostMessageFormat
}

// GetWSSServer 返回 WSSServer
func (y *WSSYaml) GetWSSServer(id int64) (s *ser.WSS, format string) {
	s = &ser.WSS{}
	s.ID = id
	s.Addr = fmt.Sprintf("%s:%d", y.Host, y.Port)
	s.Token = y.AccessToken
	return s, y.PostMessageFormat
}

// GetHTTPServer 返回HTTPServer
func (y *HTTPYaml) GetHTTPServer(id int64) (s *ser.HTTP, format string) {
	s = &ser.HTTP{}
	s.ID = id
	if y.Host != "" {
		s.Addr = fmt.Sprintf("%s:%d", y.Host, y.Port)
	}
	s.Token = y.AccessToken
	s.URL = y.PostUrl
	s.Secret = y.Secret
	return s, y.PostMessageFormat
}
