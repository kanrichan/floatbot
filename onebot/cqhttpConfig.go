package onebot

import (
	"gopkg.in/yaml.v3"
	"os"
	"strconv"
	"strings"
	"time"

	"yaya/core"

	"database/sql"
	"github.com/gorilla/websocket"
)

var Conf *Yaml

type Yaml struct {
	Version       string         `yaml:"version"`
	Master        int64          `yaml:"master"`
	Debug         bool           `yaml:"debug"`
	Meta          bool           `yaml:"-"`
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
	DB       *sql.DB     `yaml:"-"`
	WSSConf  []*WSSYaml  `yaml:"websocket"`
	WSCConf  []*WSCYaml  `yaml:"websocket_reverse"`
	HTTPConf []*HTTPYaml `yaml:"http"`
}

type HTTPYaml struct {
	Name              string      `yaml:"name"`
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
	Heart             chan []byte `yaml:"-"`
}

type WSCYaml struct {
	Name               string          `yaml:"name"`
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
	Heart              chan []byte     `yaml:"-"`
}

type WSSYaml struct {
	Name              string            `yaml:"name"`
	Enable            bool              `yaml:"enable"`
	Host              string            `yaml:"host"`
	Port              int64             `yaml:"port"`
	AccessToken       string            `yaml:"access_token"`
	PostMessageFormat string            `yaml:"post_message_format"`
	BotID             int64             `yaml:"-"`
	Status            int64             `yaml:"-"`
	Conn              []*websocket.Conn `yaml:"-"`
	Event             chan []byte       `yaml:"-"`
	Heart             chan []byte       `yaml:"-"`
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

func DefaultQQ() int64 {
	botList := strings.Split(core.GetQQList(), "/n")
	if len(botList) < 0 {
		return 0
	}
	return core.Str2Int(botList[0])
}

func DefaultBotConfig() *BotYaml {
	return &BotYaml{
		Bot: DefaultQQ(),
		WSSConf: []*WSSYaml{
			&WSSYaml{
				Name:              "WSS EXAMPLE",
				Enable:            false,
				Host:              "127.0.0.1",
				Port:              6700,
				AccessToken:       "",
				PostMessageFormat: "string",
			},
		},
		WSCConf: []*WSCYaml{
			&WSCYaml{
				Name:               "WSC EXAMPLE",
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
				Name:              "HTTP EXAMPLE",
				Enable:            false,
				Host:              "127.0.0.1",
				Port:              5700,
				AccessToken:       "",
				PostUrl:           "http://127.0.0.1/plugin",
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
	c.Save(p)
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
	conf.Meta = false
	for i, _ := range conf.BotConfs {
		for j, _ := range conf.BotConfs[i].WSSConf {
			conf.BotConfs[i].WSSConf[j].Status = 0
			conf.BotConfs[i].WSSConf[j].BotID = conf.BotConfs[i].Bot
			conf.BotConfs[i].WSSConf[j].Event = make(chan []byte, 100)
			conf.BotConfs[i].WSSConf[j].Heart = make(chan []byte, 1)
		}
		for k, _ := range conf.BotConfs[i].WSCConf {
			conf.BotConfs[i].WSCConf[k].Status = 0
			conf.BotConfs[i].WSCConf[k].BotID = conf.BotConfs[i].Bot
			conf.BotConfs[i].WSCConf[k].Event = make(chan []byte, 100)
			conf.BotConfs[i].WSCConf[k].Heart = make(chan []byte, 1)
		}
		for l, _ := range conf.BotConfs[i].HTTPConf {
			conf.BotConfs[i].HTTPConf[l].Status = 0
			conf.BotConfs[i].HTTPConf[l].BotID = conf.BotConfs[i].Bot
			conf.BotConfs[i].HTTPConf[l].Event = make(chan []byte, 100)
			conf.BotConfs[i].HTTPConf[l].Heart = make(chan []byte, 1)
		}
	}
}
