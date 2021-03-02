package gateway

import (
	"io/ioutil"
	"onebot/server"
	"os"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"
)

// 整个OneBot的配置
type Yaml struct {
	Version  string     `yaml:"version"`
	Master   int64      `yaml:"master"`
	Debug    bool       `yaml:"debug"`
	BotConfs []*BotYaml `yaml:"bots"`
}

// Bot所有连接
type BotYaml struct {
	Bot      int64                     `yaml:"bot"`
	WSSConf  *WSSYaml                  `yaml:"websocket"`
	WSCConf  []*server.WebSocketClient `yaml:"websocket_reverse"`
	HTTPConf []*server.HttpServer      `yaml:"http"`
}

// 正向WS
type WSSYaml struct {
	Name              string `yaml:"name"`
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	AccessToken       string `yaml:"access_token"`
	PostMessageFormat string `yaml:"post_message_format"`
	HeartBeatInterval int64  `yaml:"heartbeat_interval"`
}

// defaultConfig 默认配置文件
func defaultConfig() *Yaml {
	return &Yaml{
		Version:  "1.2.0",
		Master:   12345678,
		Debug:    true,
		BotConfs: []*BotYaml{defaultBotConfig()},
	}
}

func defaultBotConfig() *BotYaml {
	return &BotYaml{
		Bot: 0,
		WSSConf: &WSSYaml{
			Name:              "WSS EXAMPLE",
			Enable:            true,
			Host:              "127.0.0.1",
			Port:              6700,
			AccessToken:       "",
			PostMessageFormat: "string",
			HeartBeatInterval: 10000,
		},
		WSCConf: []*server.WebSocketClient{
			{
				Name:               "WSC EXAMPLE",
				Enable:             true,
				Url:                "ws://127.0.0.1:8080/ws",
				ApiUrl:             "ws://127.0.0.1:8080/api",
				EventUrl:           "ws://127.0.0.1:8080/event",
				UseUniversalClient: true,
				AccessToken:        "",
				PostMessageFormat:  "string",
				HeartBeatInterval:  10000,
				ReconnectInterval:  3000,
			},
		},
		HTTPConf: []*server.HttpServer{
			{
				Name:              "HTTP EXAMPLE",
				Enable:            true,
				Host:              "127.0.0.1",
				Port:              5700,
				AccessToken:       "",
				PostUrl:           "",
				Secret:            "",
				TimeOut:           0,
				PostMessageFormat: "string",
				HeartBeatInterval: 10000,
			},
		},
	}
}

// ConfigLoad 加载配置文件
func configLoad(file string) (conf *Yaml) {
	conf = defaultConfig()
	if !conf.isExists(file) {
		conf.save(file)
		INFO("[配置] 检测到无默认配置文件，已自动生成 %s", file)
		return conf
	}
	err := conf.read(file)
	if err != nil {
		ERROR("[配置] 配置文件解析失败，将备份并重新生成默认配置")
		os.Rename(file, file+".backup"+strconv.FormatInt(time.Now().Unix(), 10))
		conf := defaultConfig()
		conf.save(file)
		return conf
	}
	INFO("[配置] 加载完毕！")
	conf.save(file)
	return conf
}

func (c *Yaml) read(file string) (err error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Yaml) save(file string) (err error) {
	data, _ := yaml.Marshal(c)
	err = ioutil.WriteFile(file, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (c *Yaml) isExists(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}
