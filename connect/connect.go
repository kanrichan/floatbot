package connect

import (
	"fmt"
	"time"
	"yaya/server"
)

var QQMap []int64

type BotConnect struct {
	Bot int64 `yaml:"bot"`

	HttpServers      []*server.HttpServer      `yaml:"http"`
	WebSocketClients []*server.WebSocketClient `yaml:"websocket"`
	WebSocketServers []*server.WebSocketServer `yaml:"websocket_reverse"`
}

var StopServer chan bool = make(chan bool, 1)

var Connects []BotConnect

func Run() {
	for {
		select {
		case <-time.After(time.Second * 1):
			for i, _ := range Connects {
				for _, s := range Connects[i].WebSocketClients {
					// 错误处理
					if s.ListenStatus == "error" && s.SendStatus == "ok" {
						s.StopSend <- true
						s.Conn.Close()
						s.ConnectStatus = "wait"
						s.ListenStatus = "wait"
					}
					if s.SendStatus == "error" && s.ListenStatus == "ok" {
						s.Conn.Close()
						s.ConnectStatus = "wait"
						s.ListenStatus = "wait"
						s.SendStatus = "wait"
					}
					if s.ListenStatus == "error" && s.SendStatus == "error" {
						s.Conn.Close()
						s.ConnectStatus = "wait"
						s.ListenStatus = "wait"
						s.SendStatus = "wait"
					}
					if s.ConnectStatus == "wait" {
						fmt.Println("Connect")
						go s.Connect()
					}
					if s.ConnectStatus == "ok" && s.ListenStatus == "wait" {
						fmt.Println("Listen")
						go s.Listen()
					}
					if s.ConnectStatus == "ok" && s.SendStatus == "wait" {
						fmt.Println("Send")
						go s.Send()
					}
				}
			}

		case <-StopServer:
			for i, _ := range Connects {
				for _, s := range Connects[i].WebSocketClients {
					s.StopSend <- true
					s.Conn.Close()
					s.ConnectStatus = "wait"
					s.ListenStatus = "wait"
				}
			}
			return
		}
	}

}
