package onebot

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tidwall/gjson"
)

// runOnebot run all server in config
func (conf *Yaml) runOnebot() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[OneBot] OneBot =X=> =X=> Start Error: %v", err)
			WARN("[OneBot] OneBot ==> ==> Sleep")
		}
	}()
	go Conf.heartBeat()
	for i, _ := range conf.BotConfs {
		for j, _ := range conf.BotConfs[i].WSSConf {
			if conf.BotConfs[i].WSSConf[j].Status == 0 && conf.BotConfs[i].WSSConf[j].Enable == true && conf.BotConfs[i].WSSConf[j].Host != "" {
				go conf.BotConfs[i].WSSConf[j].start()
				go conf.BotConfs[i].WSSConf[j].listen()
				go conf.BotConfs[i].WSSConf[j].send()
			}
		}
		for k, _ := range conf.BotConfs[i].WSCConf {
			if conf.BotConfs[i].WSCConf[k].Status == 0 && conf.BotConfs[i].WSCConf[k].Enable == true && conf.BotConfs[i].WSCConf[k].Url != "" {
				go conf.BotConfs[i].WSCConf[k].listen()
				go conf.BotConfs[i].WSCConf[k].send()
			}
		}
		for l, _ := range conf.BotConfs[i].HTTPConf {
			if conf.BotConfs[i].HTTPConf[l].Status == 0 && conf.BotConfs[i].HTTPConf[l].Enable == true {
				if conf.BotConfs[i].HTTPConf[l].Host != "" {
					go conf.BotConfs[i].HTTPConf[l].listen()
				}
				go conf.BotConfs[i].HTTPConf[l].send()
			}
		}
	}
}

// heartBeat HeartBeat --> ALL PLUGINS
func (conf *Yaml) heartBeat() {
	defer func() {
		if err := recover(); err != nil {
			ERROR("[心跳] XQ =X=> =X=> Plugins Error: %v", err)
		}
	}()
	if conf.HeratBeatConf.Interval == 0 || !conf.HeratBeatConf.Enable {
		return
	}
	if conf.HeratBeatConf.Interval < 1000 {
		INFO("[心跳] Interval %v -> 1000", conf.HeratBeatConf.Interval)
		conf.HeratBeatConf.Interval = 1000
	}
	INFO("[心跳] XQ ==> ==> Plugins")
	for {
		time.Sleep(time.Millisecond * time.Duration(conf.HeratBeatConf.Interval))
		if conf.HeratBeatConf.Enable && conf.HeratBeatConf.Interval != 0 {
			for i, _ := range conf.BotConfs {
				for j, _ := range conf.BotConfs[i].WSSConf {
					if conf.BotConfs[i].WSSConf[j].Status == 1 && conf.BotConfs[i].WSSConf[j].Enable {
						conf.BotConfs[i].WSSConf[j].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for k, _ := range conf.BotConfs[i].WSCConf {
					if conf.BotConfs[i].WSCConf[k].Status == 1 && conf.BotConfs[i].WSCConf[k].Enable {
						conf.BotConfs[i].WSCConf[k].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
				for l, _ := range conf.BotConfs[i].HTTPConf {
					if conf.BotConfs[i].HTTPConf[l].Status == 1 && conf.BotConfs[i].HTTPConf[l].Enable {
						conf.BotConfs[i].HTTPConf[l].Heart <- heartEvent(conf.HeratBeatConf.Interval, conf.BotConfs[i].Bot)
					}
				}
			}
		}
	}
}

func heartEvent(interval int64, bot int64) []byte {
	heartbeat := map[string]string{
		"interval":        fmt.Sprint(interval),
		"meta_event_type": "heartbeat",
		"post_type":       "meta_event",
		"self_id":         fmt.Sprint(bot),
		"status":          "null",
		"time":            fmt.Sprint(time.Now().Unix()),
	}
	event, _ := json.Marshal(heartbeat)
	return event
}

func tieEcho(ret Result, req gjson.Result) Result {
	if req.Get("echo").Int() != 0 {
		ret.Echo = req.Get("echo").Int()
	} else if req.Get("echo").Str != "" {
		ret.Echo = req.Get("echo").Str
	} else {
		ret.Echo, _ = req.Get("echo").Value().(map[string]interface{})
	}
	return ret
}
