package gateway

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	Bot      int64       `yaml:"bot"`
	WSSConf  []*WSSYaml  `yaml:"websocket"`
	WSCConf  []*WSCYaml  `yaml:"websocket_reverse"`
	HTTPConf []*HTTPYaml `yaml:"http"`
}

// Http & Post
type HTTPYaml struct {
	Name              string `yaml:"name"`
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	AccessToken       string `yaml:"token"`
	PostUrl           string `yaml:"post_url"`
	Secret            string `yaml:"secret"`
	PostMessageFormat string `yaml:"post_message_format"`
}

// 反向WS
type WSCYaml struct {
	Name              string `yaml:"name"`
	Enable            bool   `yaml:"enable"`
	Url               string `yaml:"url"`
	AccessToken       string `yaml:"access_token"`
	PostMessageFormat string `yaml:"post_message_format"`
}

// 正向WS
type WSSYaml struct {
	Name              string `yaml:"name"`
	Enable            bool   `yaml:"enable"`
	Host              string `yaml:"host"`
	Port              int64  `yaml:"port"`
	AccessToken       string `yaml:"access_token"`
	PostMessageFormat string `yaml:"post_message_format"`
}

// defaultConfig 默认配置文件
func defaultConfig() *Yaml {
	return &Yaml{
		Version:  "1.2.0",
		Master:   0,
		Debug:    false,
		BotConfs: []*BotYaml{defaultBotConfig()},
	}
}

func defaultBotConfig() *BotYaml {
	return &BotYaml{
		Bot: 0,
		WSSConf: []*WSSYaml{
			{
				Name:              "WSS EXAMPLE",
				Enable:            true,
				Host:              "127.0.0.1",
				Port:              6700,
				AccessToken:       "",
				PostMessageFormat: "string",
			},
		},
		WSCConf: []*WSCYaml{
			{
				Name:              "WSC EXAMPLE",
				Enable:            true,
				Url:               "ws://127.0.0.1:8080/ws",
				AccessToken:       "",
				PostMessageFormat: "string",
			},
			{
				Name:              "WSC EXAMPLE 2",
				Enable:            false,
				Url:               "ws://127.0.0.1:8081/ws",
				AccessToken:       "",
				PostMessageFormat: "string",
			},
		},
		HTTPConf: []*HTTPYaml{
			{
				Name:              "HTTP EXAMPLE",
				Enable:            true,
				Host:              "127.0.0.1",
				Port:              5700,
				AccessToken:       "",
				PostUrl:           "",
				Secret:            "",
				PostMessageFormat: "string",
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
	if !c.isExists(filepath.Dir(file)) {
		err := os.MkdirAll(filepath.Dir(file), 0644)
		if err != nil {
			return err
		}
	}
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

type ConfChecker struct {
	file string
	data time.Time
}

func NewConfChecker(file string) *ConfChecker {
	f, _ := os.Open(file)
	fi, _ := f.Stat()
	return &ConfChecker{
		file: file,
		data: fi.ModTime(),
	}
}

func (c *ConfChecker) Check() bool {
	f, _ := os.Open(c.file)
	fi, _ := f.Stat()
	if fi.ModTime() != c.data {
		c.data = fi.ModTime()
		return true
	}
	return false
}
