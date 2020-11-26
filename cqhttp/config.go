package cqhttp

import (
	"gopkg.in/yaml.v2"
	"os"
	//"regexp"
	"strconv"
	//"strings"
	"time"
	//"yaya/core"

	"github.com/gorilla/websocket"
)

var Conf *Yaml

type Yaml struct {
	Version       string         `yaml:"version"`
	Master        int64          `yaml:"master"`
	Debug         bool           `yaml:"debug"`
	HeratBeatConf *HeratBeatYaml `yaml:"heratbeat"`
	Cache         *CacheYaml     `yaml:"cache"`
	BotConfs      []*BotYaml     `yaml:"bots"`
}

type CacheYaml struct {
	DataBase bool `yaml:"database"`
	Image    bool `yaml:"image"`
	Record   bool `yaml:"record"`
	Video    bool `yaml:"video"`
}

type HeratBeatYaml struct {
	Enable   bool  `yaml:"enable"`
	Interval int64 `yaml:"interval"`
}

type BotYaml struct {
	Bot      int64       `yaml:"bot"`
	WSSConf  []*WSSYaml  `yaml:"websocket"`
	WSCConf  []*WSCYaml  `yaml:"websocket_reverse"`
	HTTPConf []*HTTPYaml `yaml:"http"`
}

type HTTPYaml struct {
	Enable            bool        `yaml:"enable"`
	Host              string      `yaml:"host"`
	Port              int64       `yaml:"port"`
	AccessToken       string      `yaml:"token"`
	PostUrl           string      `yaml:"post_url"`
	Secret            string      `yaml:"secret"`
	TimeOut           int64       `yaml:"time_out"`
	PostMessageFormat string      `yaml:"post_message_format"`
	BotID             int64       `yaml:"-"`
	Status            int64       `yaml:"-"`
	Event             chan []byte `yaml:"-"`
}

type WSCYaml struct {
	Enable             bool            `yaml:"enable"`
	Url                string          `yaml:"url"`
	ApiUrl             string          `yaml:"api_url"`
	EventUrl           string          `yaml:"event_url"`
	UseUniversalClient bool            `yaml:"use_universal_client"`
	AccessToken        string          `yaml:"access_token"`
	PostMessageFormat  string          `yaml:"post_message_format"`
	ReconnectInterval  int64           `yaml:"reconnect_interval"`
	BotID              int64           `yaml:"-"`
	Status             int64           `yaml:"-"`
	Conn               *websocket.Conn `yaml:"-"`
	Event              chan []byte     `yaml:"-"`
}

type WSSYaml struct {
	Enable            bool              `yaml:"enable"`
	Host              string            `yaml:"host"`
	Port              int64             `yaml:"port"`
	AccessToken       string            `yaml:"access_token"`
	PostMessageFormat string            `yaml:"post_message_format"`
	BotID             int64             `yaml:"-"`
	Status            int64             `yaml:"-"`
	Conn              []*websocket.Conn `yaml:"-"`
	Event             chan []byte       `yaml:"-"`
}

func DefaultConfig() *Yaml {
	return &Yaml{
		Version: "1.0.5",
		Master:  12345678,
		Debug:   true,
		Cache: &CacheYaml{
			DataBase: false,
			Image:    false,
			Record:   false,
			Video:    false,
		},
		HeratBeatConf: &HeratBeatYaml{
			Enable:   true,
			Interval: 10000,
		},
		BotConfs: []*BotYaml{DefaultBotConfig()},
	}
}

func DefaultBotConfig() *BotYaml {
	return &BotYaml{
		Bot: 0,
		WSSConf: []*WSSYaml{
			&WSSYaml{
				Enable:            false,
				Host:              "127.0.0.1",
				Port:              6700,
				AccessToken:       "",
				PostMessageFormat: "string",
			},
		},
		WSCConf: []*WSCYaml{
			&WSCYaml{
				Enable:             false,
				Url:                "ws://127.0.0.1:8080/ws",
				ApiUrl:             "ws://127.0.0.1:8080/api",
				EventUrl:           "ws://127.0.0.1:8080/event",
				UseUniversalClient: true,
				AccessToken:        "",
				PostMessageFormat:  "string",
				ReconnectInterval:  3000,
			},
		},
		HTTPConf: []*HTTPYaml{
			&HTTPYaml{
				Enable:            false,
				Host:              "127.0.0.1",
				Port:              5700,
				AccessToken:       "",
				PostUrl:           "http://127.0.0.1:5705/",
				Secret:            "",
				TimeOut:           0,
				PostMessageFormat: "string",
			},
		},
	}
}

func Load(p string) *Yaml {
	if !PathExists(p) {
		c := DefaultConfig()
		c.Save(p)
	}
	c := Yaml{}
	err := yaml.Unmarshal([]byte(ReadAllText(p)), &c)
	if err != nil {
		ERROR("Emmm，夜夜觉得配置文件有问题")
		os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		c := DefaultConfig()
		c.Save(p)
	}
	if c.Version != "1.0.5" {
		WARN("!!!!!!!!配置文件版本更新了，请重新配置配置文件")
		os.Rename(p, p+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		c := DefaultConfig()
		c.Save(p)
	}
	INFO("おはようございます。")
	c.InitConf()
	return &c
}

func (c *Yaml) Save(p string) {
	data, err := yaml.Marshal(c)
	if err != nil {
		ERROR("大失败！夜夜需要管理员权限")
	}
	WriteAllText(p, string(data))
}

func (conf *Yaml) InitConf() {
	for i, _ := range conf.BotConfs {
		for j, _ := range conf.BotConfs[i].WSSConf {
			conf.BotConfs[i].WSSConf[j].Status = 0
			conf.BotConfs[i].WSSConf[j].BotID = conf.BotConfs[i].Bot
			conf.BotConfs[i].WSSConf[j].Event = make(chan []byte, 100)
		}
		for k, _ := range conf.BotConfs[i].WSCConf {
			conf.BotConfs[i].WSCConf[k].Status = 0
			conf.BotConfs[i].WSCConf[k].BotID = conf.BotConfs[i].Bot
			conf.BotConfs[i].WSCConf[k].Event = make(chan []byte, 100)
		}
		for l, _ := range conf.BotConfs[i].HTTPConf {
			conf.BotConfs[i].HTTPConf[l].Status = 0
			conf.BotConfs[i].HTTPConf[l].BotID = conf.BotConfs[i].Bot
			conf.BotConfs[i].HTTPConf[l].Event = make(chan []byte, 100)
		}
	}
}

/*
func commandHandle(e XEvent) {
	if e.message == "/master" {
		if Conf.Master == 12345678 {
			Conf.Master = e.userID
			Conf.Save(AppPath + "config.yml")
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "登录完毕", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	} else if e.message == "/debug on" {
		if Conf.Master == e.userID {
			Conf.Debug = true
			Conf.Save(AppPath + "config.yml")
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!Debug On", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	} else if e.message == "/debug off" {
		if Conf.Master == e.userID {
			Conf.Debug = false
			Conf.Save(AppPath + "config.yml")
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!Debug Off", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	} else if e.message == "/夜夜" {
		if Conf.Master == e.userID {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "在！", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	}

	setWSCurl(e)
	setWSCon(e)
	setWSCoff(e)
	setWSCtoken(e)
}

func setWSCon(e XEvent) {
	if e.message == "/wsc enable" {
		if Conf.Master == e.userID {
			for i, conf := range Conf.BotConfs {
				if conf.Bot == e.selfID {
					Conf.BotConfs[i].WSCConf.Enable = true
					Conf.Save(AppPath + "config.yml")
					break
				}
				if i+1 == len(Conf.BotConfs) {
					newBotConf := DefaultBotConfig()
					newBotConf.Bot = e.selfID
					newBotConf.WSCConf.Enable = true
					Conf.BotConfs = append(Conf.BotConfs, newBotConf)
					Conf.Save(AppPath + "config.yml")
				}
			}
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!WebSocket Reverse Enable", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	}
}

func setWSCoff(e XEvent) {
	if e.message == "/wsc disable" {
		if Conf.Master == e.userID {
			for i, conf := range Conf.BotConfs {
				if conf.Bot == e.selfID {
					Conf.BotConfs[i].WSCConf.Enable = false
					Conf.Save(AppPath + "config.yml")
					break
				}
				if i+1 == len(Conf.BotConfs) {
					newBotConf := DefaultBotConfig()
					newBotConf.Bot = e.selfID
					newBotConf.WSCConf.Enable = false
					Conf.BotConfs = append(Conf.BotConfs, newBotConf)
					Conf.Save(AppPath + "config.yml")
				}
			}
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!WebSocket Reverse Disable", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	}
}

func setWSCurl(e XEvent) {
	wscUrlR := regexp.MustCompile(`\/wsc url (.*)`)
	if len(wscUrlR.FindStringSubmatch(e.message)) != 0 {
		if Conf.Master == e.userID {
			for i, conf := range Conf.BotConfs {
				if conf.Bot == e.selfID {
					Conf.BotConfs[i].WSCConf.Url = wscUrlR.FindStringSubmatch(e.message)[1]
					Conf.Save(AppPath + "config.yml")

					break
				}
				if i+1 == len(Conf.BotConfs) {
					newBotConf := DefaultBotConfig()
					newBotConf.Bot = e.selfID
					newBotConf.WSCConf.Url = wscUrlR.FindStringSubmatch(e.message)[1]
					Conf.BotConfs = append(Conf.BotConfs, newBotConf)
					Conf.Save(AppPath + "config.yml")
				}
			}
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!WebSocket Reverse Url Updated", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	}
}

func setWSCtoken(e XEvent) {
	wscTokenR := regexp.MustCompile(`\/wsc token (.*)`)
	if len(wscTokenR.FindStringSubmatch(e.message)) != 0 {
		if Conf.Master == e.userID {
			for i, conf := range Conf.BotConfs {
				if conf.Bot == e.selfID {
					Conf.BotConfs[i].WSCConf.AccessToken = wscTokenR.FindStringSubmatch(e.message)[1]
					Conf.Save(AppPath + "config.yml")

					break
				}
				if i+1 == len(Conf.BotConfs) {
					newBotConf := DefaultBotConfig()
					newBotConf.Bot = e.selfID
					newBotConf.WSCConf.AccessToken = wscTokenR.FindStringSubmatch(e.message)[1]
					Conf.BotConfs = append(Conf.BotConfs, newBotConf)
					Conf.Save(AppPath + "config.yml")
				}
			}
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!WebSocket Reverse Token Updated", 0)
		} else {
			core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
		}
	}
}
*/
/*
// 命令解析器
func commandParse(e XEvent) {
	commands := strings.Split(e.message, "/n")
	for _, command := range commands {
		if strings.Contains(command, "$") {

		}
	}
}

func selectCommandType(command string) bool {
	whereSapce := strings.Index(command, " ")
	if whereSapce == 0 {
		commandType := command[1:]
	} else {
		commandType := command[1:whereSapce]
	}
	switch commandType {
	case "夜夜":
		isYaYaOK(e)
	case "master":
		Conf := masterLogin(Conf, e)
	case "debug":
		commandParm := command[1:whereSapce]
		if
	case "wsc-url":
	}
}

// isYaYaOK 查看夜夜是否在线
func isYaYaOK(e XEvent) {
	core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "夜夜一直在你身边！", 0)
}

// masterLogin 变更主人
func masterLogin(conf *Yaml, e XEvent) *Yaml {
	if conf.Master == 12345678 {
		conf.Master = e.userID
		conf.Save(AppPath + "config.yml")
		core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "登录完毕", 0)
	} else {
		core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
	}
	return conf
}

// debugEnable 开启debug
func debugEnable(conf *Yaml, e XEvent) *Yaml {
	if e.userID == conf.Master {
		Conf.Debug = false
		Conf.Save(AppPath + "config.yml")
		core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!Debug Off", 0)
	} else {
		core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
	}
	return conf
}

// debugDisable 关闭debug
func debugDisable(conf *Yaml, e XEvent) *Yaml {
	if e.userID == conf.Master {
		Conf.Debug = false
		Conf.Save(AppPath + "config.yml")
		core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "!Debug Off", 0)
	} else {
		core.SendMsg(e.selfID, e.mseeageType, e.groupID, e.userID, "???", 0)
	}
	return conf
}

func isDebugCmd(command string) bool {
	if strings.ToLower(command) == "$master" {
		return true
	}
	return false
}

func isMasterCmd(command string) bool {
	if strings.ToLower(command) == "$master" {
		return true
	}
	return false
}
*/

// 兼容 Mirai CQHTTP
