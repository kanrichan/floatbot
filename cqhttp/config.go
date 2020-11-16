package cqhttp

import (
	"gopkg.in/yaml.v2"
	"os"
	"strconv"
	"time"
)

var Conf *Yaml

type Yaml struct {
	Version  string     `yaml:"version"`
	Master   int64      `yaml:"master"`
	Debug    bool       `yaml:"debug"`
	BotConfs []*BotYaml `yaml:"bots"`
}

type BotYaml struct {
	Bot           int64          `yaml:"bot"`
	CacheImage    bool           `yaml:"cacheImage"`
	CacheRecord   bool           `yaml:"cacheRrcord"`
	HeratBeatConf *HeratBeatYaml `yaml:"heratbeat"`
	WSSConf       *WSSYaml       `yaml:"websocket"`
	WSCConf       *WSCYaml       `yaml:"websocket_reverse"`
	HTTPConf      *HTTPYaml      `yaml:"http"`
}

type HeratBeatYaml struct {
	Enable   bool  `yaml:"enable"`
	Interval int64 `yaml:"interval"`
}

type HTTPYaml struct {
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	PostUrl           string `yaml:"post_url"`
	Secret            string `yaml:"secret"`
	TimeOut           int64  `yaml:"time_out"`
	PostMessageFormat string `yaml:"post_message_format"`
}

type WSCYaml struct {
	Enable             bool   `yaml:"enable"`
	Url                string `yaml:"url"`
	ApiUrl             string `yaml:"api_url"`
	EventUrl           string `yaml:"event_url"`
	UseUniversalClient bool   `yaml:"use_universal_client"`
	AccessToken        string `yaml:"access_token"`
	PostMessageFormat  string `yaml:"post_message_format"`
	ReconnectInterval  int64  `yaml:"reconnect_interval"`
}

type WSSYaml struct {
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	AccessToken       string `yaml:"access_token"`
	PostMessageFormat string `yaml:"post_message_format"`
}

func DefaultConfig() *Yaml {
	return &Yaml{
		Version: "1.0.1",
		Master:  12345678,
		Debug:   false,
		BotConfs: []*BotYaml{
			{
				Bot:         0,
				CacheImage:  false,
				CacheRecord: false,
				HeratBeatConf: &HeratBeatYaml{
					Enable:   true,
					Interval: 10000,
				},
				WSSConf: &WSSYaml{
					Enable:            true,
					Host:              "127.0.0.1",
					Port:              6700,
					AccessToken:       "",
					PostMessageFormat: "string",
				},
				WSCConf: &WSCYaml{
					Enable:             true,
					Url:                "ws://127.0.0.1:8080/ws",
					ApiUrl:             "ws://127.0.0.1:8080/api",
					EventUrl:           "ws://127.0.0.1:8080/event",
					UseUniversalClient: false,
					AccessToken:        "",
					PostMessageFormat:  "string",
					ReconnectInterval:  3000,
				},
				HTTPConf: &HTTPYaml{
					Enable:            false,
					Host:              "127.0.0.1",
					Port:              5700,
					PostUrl:           "http://127.0.0.1:5705/",
					Secret:            "",
					TimeOut:           0,
					PostMessageFormat: "string",
				},
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
		return nil
	}
	INFO("おはようございます。")
	return &c
}

func (c *Yaml) Save(p string) {
	data, err := yaml.Marshal(c)
	if err != nil {
		ERROR("大失败！夜夜需要管理员权限")
	}
	WriteAllText(p, string(data))
	INFO("将以默认配置启动~")
}
