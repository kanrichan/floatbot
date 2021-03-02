package gateway

import (
	"fmt"
	core "onebot/core/xianqu"
	ser "onebot/server"
	"time"
)

var (
	// 原始配置
	CONF = &Yaml{}
	// 配置上已注册的bot
	BotsMap = []int64{}
	// 配置上所有的连接
	Connects = make(map[int64]BotsConnect)
	// 关闭所有连接的chan
	StopServer = make(chan bool, 1)
	IsRunning  = false
)

// 所有连接的表
type BotsConnect struct {
	HttpServers      []*ser.HttpServer
	WebSocketClients []*ser.WebSocketClient
	WebSocketServers *ser.WSS
}

// core触发启动事件
func OnEnable(_ *core.Context) {
	if IsRunning {
		core.XQApiCallMessageBox("禁止多次启动！")
		return // 防止多次运行
	}
	CONF = configLoad(core.OneBotPath + "config.yml")
	if CONF != nil {
		IsRunning = true
		if CONF.BotConfs[0].Bot == 0 {
			ERROR("配置文件中未设置姬气人账号！")
		}
		connInit()
		INFO("[初始化] 所有初始化准备就绪，开始执行连接")
		controller()
	}
}

// core触发停止事件
func OnDisable(_ *core.Context) {
	if IsRunning {
		StopServer <- true
	}
	time.Sleep(time.Second * 2)
	BotsMap = []int64{}
	Connects = make(map[int64]BotsConnect)
	StopServer = make(chan bool, 1)
	IsRunning = false
	INFO("已停止运行！")
}

// 将配置上所有连接都注册到表 Connects 上
func connInit() {
	for _, b := range CONF.BotConfs {
		if b.WSCConf
		if b.WSSConf != nil {
			c := b.WSSConf
			wss := &ser.WSS{}
			wss.Run(b.Bot, c.HeartBeatInterval, fmt.Sprintf("%s:%d", c.Host, c.Port), c.AccessToken)
		}

		if len(httpServers)+len(webSocketClients) == 0 {
			// 这个bot一个连接都没有，不注册
			continue
		}
		BotsMap = append(BotsMap, b.Bot)
		Connects[b.Bot] = BotsConnect{
			httpServers,
			webSocketClients,
			wss,
		}
	}
}
func (s *WSS) heartbeat() {
	defer func() {
		recover()
	}()
	for {
		time.Sleep(time.Millisecond * time.Duration(s.heart))
		heartbeat := fmt.Sprintf(`{"interval":%d,"meta_event_type":"heartbeat","post_type":"meta_event","self_id":%d,"status":{"good":true,"online":true},"time":%d}`,
			s.heart, s.id, time.Now().Unix())
		for _, conn := range s.conn {
			conn.mutex.Lock()
			if err := conn.conn.WriteMessage(websocket.TextMessage, []byte(heartbeat)); err != nil {
				conn.conn.Close()
			}
			conn.mutex.Unlock()
		}
	}
}

// 管理所有连接的管理器
func controller() {
	for {
		select {
		// 关闭所有的连接并退出此goroutine
		case <-StopServer:
			for _, bot := range BotsMap {
				INFO("%v", Connects[bot].WebSocketServers[0].Name)
				INFO("停止！")
				for _, s := range Connects[bot].HttpServers {
					s.Stop(true)
				}
				for _, s := range Connects[bot].WebSocketServers {
					INFO("停止！")
					INFO("???")
					s.Stop(true)
				}
				for _, s := range Connects[bot].WebSocketClients {
					s.Stop(true)
				}

			}
			return
		// 开启所有的连接并清除错误连接
		case <-time.After(time.Second * 1):
			for _, bot := range BotsMap {
				for _, s := range Connects[bot].HttpServers {
					s.Stop(false) // 检查发生错误的连接并进行关闭
					s.Run()       // 检查等待状态的连接并进行连接
				}
				for _, s := range Connects[bot].WebSocketClients {
					s.Stop(false)
					s.Run()
				}
				for _, s := range Connects[bot].WebSocketServers {
					s.Stop(false)
					s.Run()
				}
			}
		}
	}
}
