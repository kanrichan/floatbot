package onebot

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type GroupHonorInfo struct {
	GroupCode        string            `json:"gc"`
	Uin              string            `json:"uin"`
	Type             int64             `json:"type"`
	TalkativeList    []HonorMemberInfo `json:"talkativeList"`
	CurrentTalkative CurrentTalkative  `json:"currentTalkative"`
	ActorList        []HonorMemberInfo `json:"actorList"`
	LegendList       []HonorMemberInfo `json:"legendList"`
	StrongNewbieList []HonorMemberInfo `json:"strongnewbieList"`
	EmotionList      []HonorMemberInfo `json:"emotionList"`
}

type HonorMemberInfo struct {
	Uin    int64  `json:"uin"`
	Avatar string `json:"avatar"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
}

type CurrentTalkative struct {
	Uin      int64  `json:"uin"`
	DayCount int32  `json:"day_count"`
	Avatar   string `json:"avatar"`
	Name     string `json:"nick"`
}

func Base642Image(res string) string {
	data, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		ERROR("base64编码解码失败")
	}
	name := strings.ToUpper(fmt.Sprintf("%x", md5.Sum(data)))
	path := ImagePath + name + ".jpg"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("base64编码保存图片失败")
	} else {
		_, err = f.Write(data)
		if err != nil {
			ERROR("base64编码写入图片失败")
		}
	}
	return path
}

func Url2Image(url string) string {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Add("Net-Type", "Wifi")
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器图片%s下载失败", url)
		return "error"
	}
	resp, err := client.Do(reqest)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器图片%s下载失败", url)
		return "error"
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器图片%s下载失败", url)
		return "error"
	}
	name := byte2md5(data)
	path := ImagePath + name + ".jpg"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err == nil {
		_, err = f.Write(data)
		if err != nil {
			ERROR("[CQ码解析] 从TX服务器图片%s保存失败", url)
		}
	} else {
		ERROR("[CQ码解析] 从TX服务器图片%s保存失败", url)
	}
	return path
}

func byte2md5(data []byte) string {
	m := md5.New()
	m.Write(data)
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

func XmlEscape(c string) string {
	buf := new(bytes.Buffer)
	_ = xml.EscapeText(buf, []byte(c))
	return buf.String()
}

func groupHonor(groupID int64, honorType int64, cookie string) []byte {
	url := fmt.Sprintf("https://qun.qq.com/interactive/honorlist?gc=%d&type=%d", groupID, honorType)
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Add("Net-Type", "Wifi")
	reqest.Header.Add("Cookie", cookie)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器图片%s下载失败", url)
		return nil
	}
	resp, err := client.Do(reqest)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器图片%s下载失败", url)
		return nil
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器图片%s下载失败", url)
		return nil
	}
	return data
}
