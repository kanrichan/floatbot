package main

import (
	"encoding/json"

	"yaya/cqhttp"
)

// 插件信息
type AppInfo struct {
	Name   string `json:"name"`   // 插件名字
	Pver   string `json:"pver"`   // 插件版本
	Sver   int    `json:"sver"`   // 框架版本
	Author string `json:"author"` // 作者名字
	Desc   string `json:"desc"`   // 插件说明
}

func newAppInfo() *AppInfo {
	return &AppInfo{
		Name:   "OneBot-YaYa",
		Pver:   "1.0.7 beta",
		Sver:   3,
		Author: "kanri",
		Desc:   "http://github.com/Yiwen-Chan/OneBot-YaYa",
	}
}

func init() {
	data, _ := json.Marshal(newAppInfo())
	cqhttp.AppInfoJson = string(data)
}

func main() { cqhttp.Main() }
