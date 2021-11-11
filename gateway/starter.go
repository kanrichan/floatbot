package gateway

import (
	"fmt"
	"time"

	core "onebot/core/xianqu"

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
)

var (
	// 原始配置
	CONF = &Yaml{}
	// 关闭所有连接的chan
	StopServer = make(chan bool, 1)
	IsRunning  = false
)

// core触发启动事件
func OnEnable(_ *core.Context) {
	if IsRunning {
		core.ApiCallMessageBox("禁止多次启动！")
		return // 防止多次运行
	}
	IsRunning = true
	time.Sleep(time.Second * 1)
	// 启动 显示 ONEBOT
	core.ApiOutPutLog(`   ____    _      __ ________ _______     _____  __________`)
	core.ApiOutPutLog(` /  __  \ |  \   |   |   ______|   ___   ) /  __    |___    ___|`)
	core.ApiOutPutLog(`|   |   |   |    \ |   |    __|   |    __   \|   |   |   |    |  |`)
	core.ApiOutPutLog(`|   |__|   |   | \    |   |_____|    |_)    ||   |__|   |    |  |`)
	core.ApiOutPutLog(` \_____/ |__|   \__|________|________/ \_______/    |_ |`)
	CONF = configLoad(core.OneBotPath + "config.yml")
	if len(CONF.BotConfs) == 1 && CONF.BotConfs[0].Bot == 0 {
		go core.ApiMessageBoxButton("配置文件中未设置姬气人账号!请自行修改配置文件并热重载")
	}
	INFO("[初始化] 所有初始化准备就绪，开始执行连接")
	table := NewServersTable()
	for i := range CONF.BotConfs {
		for j := range CONF.BotConfs[i].WSCConf {
			if !CONF.BotConfs[i].WSCConf[j].Enable {
				continue
			}
			s, format := CONF.BotConfs[i].WSCConf[j].GetWSCServer(CONF.BotConfs[i].Bot)
			table.Add(CONF.BotConfs[i].Bot, s, format)
		}
		for j := range CONF.BotConfs[i].WSSConf {
			if !CONF.BotConfs[i].WSSConf[j].Enable {
				continue
			}
			s, format := CONF.BotConfs[i].WSSConf[j].GetWSSServer(CONF.BotConfs[i].Bot)
			table.Add(CONF.BotConfs[i].Bot, s, format)
		}
		for j := range CONF.BotConfs[i].HTTPConf {
			if !CONF.BotConfs[i].HTTPConf[j].Enable {
				continue
			}
			s, format := CONF.BotConfs[i].HTTPConf[j].GetHTTPServer(CONF.BotConfs[i].Bot)
			table.Add(CONF.BotConfs[i].Bot, s, format)
		}
	}
	for _, bot := range table.Bots {
		INFO("[初始化] [%d] 连接列表%s", bot, table.Servers[bot])
	}
	table.Run()
	// 检查配置文件并热重启
	checker := NewConfChecker(core.OneBotPath + "config.yml")
	// 阻塞至关闭
BLOCK:
	for {
		select {
		case <-StopServer:
			break BLOCK
		case <-time.After(time.Millisecond * 761):
			if checker.Check() {
				if core.ApiMessageBoxButton("检测到配置文件变化，是否热重载？") == 6 {
					go OnRestart(nil)
				}
			}
		case <-time.After(time.Second * 3):
			for _, bot := range table.Bots {
				heartbeat := fmt.Sprintf(`{"interval":%d,"meta_event_type":"heartbeat","post_type":"meta_event","self_id":%d,"status":{"good":true,"online":true},"time":%d}`,
					3000, bot, time.Now().Unix())
				table.SendByte(bot, helper.StringToBytes(heartbeat))
			}
		}
	}
	table.Close()
}

// core触发停止事件
func OnDisable(_ *core.Context) {
	if IsRunning {
		StopServer <- true
	}
	IsRunning = false
}

// core触发热重载事件
func OnRestart(_ *core.Context) {
	OnDisable(nil)
	time.Sleep(time.Second * 1)
	OnEnable(nil)
}

// core触发设置
func OnSetting(_ *core.Context) {
	core.ApiCallMessageBox(fmt.Sprintf("等个好心人写UI，修改配置请到 %sconfig.yml\n", core.OneBotPath))
}
