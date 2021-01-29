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
	"regexp"
	"strings"
	"sync"

	"github.com/Yiwen-Chan/go-silk/silk"
)

type PicsCache struct {
	Lock sync.Mutex
	Max  int
	Md5  []string
}

func (pool *PicsCache) init(max int) {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()
	pool.Max = max
}

func (pool *PicsCache) add(md5 string) {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()
	if len(pool.Md5) >= pool.Max {
		start := int(0.1 * float64(pool.Max))
		if start == 0 || start > len(pool.Md5) {
			start = 1
		}
		pool.Md5 = pool.Md5[start:]
	}
	pool.Md5 = append(pool.Md5, md5)
}

func (pool *PicsCache) search(md5 string) bool {
	pool.Lock.Lock()
	defer pool.Lock.Unlock()
	length := len(pool.Md5)
	for i := length; i > 0; i-- {
		if pool.Md5[i-1] == md5 {
			return true
		}
	}
	return false
}

func (pool *PicsCache) addPicPool(text string) {
	pic := regexp.MustCompile(`\[pic={(.*?)-(.*?)-(.*?)-(.*?)-(.*?)}(\..*?)\]`)
	for _, p := range pic.FindAllStringSubmatch(text, -1) {
		md5 := strings.ToUpper(fmt.Sprintf("%s%s%s%s%s", p[1], p[2], p[3], p[4], p[5]))
		pool.add(md5)
	}
}

func hash2txfile(md5, type_ string) string {
	return fmt.Sprintf(
		"{%s-%s-%s-%s-%s}%s",
		md5[:8],
		md5[8:12],
		md5[12:16],
		md5[16:20],
		md5[20:],
		type_,
	)
}

type picDownloader struct {
	file     string
	url      string
	res      string
	type_    string
	suffix   string
	savePath string
	iscache  bool
}

func (pic *picDownloader) path() string {
	if pic.url != "" {
		return pic.file
	}
	switch {
	default:
		pic.type_ = "file"
		pic.res = pic.file
	case strings.Contains(pic.file, "base64://"):
		pic.type_ = "base64"
		pic.res = pic.file[9:]
	case strings.Contains(pic.file, "file:///"):
		pic.type_ = "file"
		pic.res = pic.file[8:]
	case strings.Contains(pic.file, "http://"):
		pic.type_ = "http"
		pic.res = pic.file
	case strings.Contains(pic.file, "https://"):
		pic.type_ = "http"
		pic.res = pic.file
	}

	path := pic.cache()
	data, _ := ioutil.ReadFile(path)
	md5 := strings.ToUpper(fmt.Sprintf("%x", md5.Sum(data)))
	if PicPool.search(md5) {
		return hash2txfile(md5, pic.suffix)
	}
	return path
}

func (pic *picDownloader) cache() string {
	switch pic.type_ {
	default:
		return pic.res
	case "file":
		return pic.res
	case "http":
		return pic.urlCache()
	case "base64":
		return pic.base64Cache()
	}
}

// urlCache 缓存并返回缓存路径
func (pic *picDownloader) urlCache() string {
	// TODO 文件名为url的hash值
	name := hashText(pic.res)
	path := pic.savePath + name + pic.suffix
	// TODO 判断是否使用缓存以及缓存是否存在
	if PathExists(path) && pic.iscache {
		return path
	}
	// TODO 模拟QQ客户端请求
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", pic.res, nil)
	reqest.Header.Add("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Add("Net-Type", "Wifi")

	resp, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	// TODO 写入文件
	data, _ := ioutil.ReadAll(resp.Body)
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)
	return path
}

func (pic *picDownloader) base64Cache() string {
	data, err := base64.StdEncoding.DecodeString(pic.res)
	if err != nil {
		panic(err)
	}
	name := strings.ToUpper(fmt.Sprintf("%x", md5.Sum(data)))
	path := pic.savePath + name + pic.suffix

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)
	return path
}

type recDownloader struct {
	file     string
	url      string
	res      string
	type_    string
	suffix   string
	savePath string
	iscache  bool
}

func (rec *recDownloader) path() string {
	if rec.url != "" {
		return rec.file
	}
	switch {
	default:
		rec.type_ = "file"
		rec.res = rec.file
	case strings.Contains(rec.file, "base64://"):
		rec.type_ = "base64"
		rec.res = rec.file[9:]
	case strings.Contains(rec.file, "file:///"):
		rec.type_ = "file"
		rec.res = rec.file[8:]
	case strings.Contains(rec.file, "http://"):
		rec.type_ = "http"
		rec.res = rec.file
	case strings.Contains(rec.file, "https://"):
		rec.type_ = "http"
		rec.res = rec.file
	}

	path := rec.cache()
	return path
}

func (rec *recDownloader) cache() string {
	switch rec.type_ {
	default:
		return rec.res
	case "file":
		return rec.res
	case "http":
		return rec.urlCache()
	case "base64":
		return rec.base64Cache()
	}
}

// urlCache 缓存并返回缓存路径
func (rec *recDownloader) urlCache() string {
	// TODO 文件名为url的hash值
	name := hashText(rec.res)
	path := rec.savePath + name + rec.suffix
	// TODO 判断是否使用缓存以及缓存是否存在
	if PathExists(path) && rec.iscache {
		return path
	}
	// TODO 模拟QQ客户端请求
	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", rec.res, nil)
	reqest.Header.Add("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Add("Net-Type", "Wifi")

	resp, err := client.Do(reqest)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, _ := ioutil.ReadAll(resp.Body)
	// TODO 写入文件
	f, _ := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)
	// TODO 不使用缓存或文件不存在则转silk
	if !PathExists(rec.savePath+name+".silk") || !rec.iscache {
		path = rec2silk(path)
	}
	return path
}

func (rec *recDownloader) base64Cache() string {
	data, err := base64.StdEncoding.DecodeString(rec.res)
	if err != nil {
		panic(err)
	}
	// TODO 文件名为res的hash值
	name := strings.ToUpper(fmt.Sprintf("%x", md5.Sum(data)))
	path := rec.savePath + name + rec.suffix
	// TODO 写入文件
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	f.Write(data)
	// TODO 不使用缓存或文件不存在则转silk
	if !PathExists(rec.savePath+name+".silk") || !rec.iscache {
		path = rec2silk(path)
	}
	return path
}

func rec2silk(path string) string {
	silkEncoder := &silk.Encoder{}
	err := silkEncoder.Init("OneBot/record", "OneBot/codec")
	if err != nil {
		ERROR("[CQ码解析] %s", err)
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		ERROR("[CQ码解析] %s", err)
	}
	name := "not found"
	if strings.LastIndex(path, "\\") > strings.LastIndex(path, "/") {
		name = path[strings.LastIndex(path, "\\")+1 : strings.LastIndex(path, ".")]
	} else {
		name = path[strings.LastIndex(path, "/")+1 : strings.LastIndex(path, ".")]
	}
	_, err = silkEncoder.EncodeToSilk(data, name, true)
	if err != nil {
		ERROR("[CQ码解析] %s", err)
	}
	return RecordPath + name + ".silk"
}

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

func Base642Record(res string) string {
	data, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		ERROR("base64编码解码失败")
	}
	name := strings.ToUpper(fmt.Sprintf("%x", md5.Sum(data)))
	path := RecordPath + name + ".mp3"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		ERROR("base64编码保存语音失败")
	} else {
		_, err = f.Write(data)
		if err != nil {
			ERROR("base64编码写入语音失败")
		}
	}
	return path
}

func Url2Record(url string) string {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Add("Net-Type", "Wifi")
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器语音%s下载失败", url)
		return "error"
	}
	resp, err := client.Do(reqest)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器语音%s下载失败", url)
		return "error"
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ERROR("[CQ码解析] 从TX服务器语音%s下载失败", url)
		return "error"
	}
	name := byte2md5(data)
	path := RecordPath + name + ".mp3"
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	defer f.Close()
	if err == nil {
		_, err = f.Write(data)
		if err != nil {
			ERROR("[CQ码解析] 从TX服务器语音%s保存失败", url)
		}
	} else {
		ERROR("[CQ码解析] 从TX服务器语音%s保存失败", url)
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

func Base642ImageBytes(res string) []byte {
	data, err := base64.StdEncoding.DecodeString(res)
	if err != nil {
		ERROR("base64编码解码失败")
	}
	return data
}

func Url2ImageBytes(url string) []byte {
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", url, nil)
	reqest.Header.Add("User-Agent", "QQ/8.2.0.1296 CFNetwork/1126")
	reqest.Header.Add("Net-Type", "Wifi")
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

func Path2ImageBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}
