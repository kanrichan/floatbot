package xianqu

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tidwall/gjson"
)

func update(version, path string) {
	defer func() {
		if err := recover(); err != nil {
			ApiOutPutLog(fmt.Sprintf("[E][更新] 更新失败: %v", err))
		}
	}()
	last, link, body, err := getLastRelease()
	if err != nil {
		panic(err)
	}
	if isNeedUpdate(version, last, path) {
		ret := ApiMessageBoxButton(
			fmt.Sprintf("发现新版本 %s\n更新内容: \n%s\n\n是否下载安装更新？", last, body),
		)
		if ret == 6 {
			ApiOutPutLog(fmt.Sprintf("[I][更新] 开始下载: %s", link))
			err := downLastRelease(link, path)
			if err != nil {
				panic(err)
			}
			ApiOutPutLog(fmt.Sprintf("[I][更新] 下载完毕"))
			err = install(path)
			if err != nil {
				panic(err)
			}
			ApiOutPutLog(fmt.Sprintf("[I][更新] 安装完毕"))
			ApiCallMessageBox("安装成功，重启先驱框架生效")
			return
		}
		ret = ApiMessageBoxButton("是否跳过此版本更新，不再提示？")
		if ret == 6 {
			ApiOutPutLog(fmt.Sprintf("[I][更新] 跳过版本 %s 的更新", last))
			skipUpdate(last, path)
		}
	}
}

func isNeedUpdate(ver, last, path string) bool {
	if last == "" {
		return false
	}
	if "v"+ver == last {
		return false
	}
	file := path + "OneBot\\SkipUpdate.txt"
	b, _ := ioutil.ReadFile(file)
	if string(b) == last {
		return false
	}
	return true
}

func skipUpdate(ver, path string) {
	file := path + "OneBot\\SkipUpdate.txt"
	f, _ := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	f.Write([]byte(ver))
	f.Close()
}

func getLastRelease() (last, link, body string, err error) {
	var api = "https://api.github.com/repos/Yiwen-Chan/OneBot-YaYa/releases/latest"
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", api, nil)
	reqest.Header.Set("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Set("Net-Type", "Wifi")
	resp, err := client.Do(reqest)
	if err != nil {
		return "", "", "", err
	}
	data, _ := ioutil.ReadAll(resp.Body)
	last = gjson.ParseBytes(data).Get("tag_name").Str
	link = gjson.ParseBytes(data).Get("assets.0.browser_download_url").Str
	body = gjson.ParseBytes(data).Get("body").Str
	return last, link, body, nil
}

func downLastRelease(link, path string) (err error) {
	link = strings.ReplaceAll(link, "github.com", "hub.fastgit.org")
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", link, nil)
	reqest.Header.Set("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Set("Net-Type", "Wifi")
	resp, err := client.Do(reqest)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	file := path + "OneBot\\OneBot-YaYa.XQ.dll"
	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.Write(data)
	f.Close()
	resp.Body.Close()
	return nil
}

func install(path string) (err error) {
	return os.Rename(path+"OneBot\\OneBot-YaYa.XQ.dll", path+"Plugin\\OneBot-YaYa.XQ.dll")
}
