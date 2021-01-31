package onebot

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"

	"yaya/core"
)

func INFO(s string, v ...interface{}) {
	core.OutPutLog("[INFO] " + fmt.Sprintf(s, v...))
}

func WARN(s string, v ...interface{}) {
	core.OutPutLog("[WARN] " + fmt.Sprintf(s, v...))
}

func DEBUG(s string, v ...interface{}) {
	if Conf.Debug {
		core.OutPutLog("[DEBUG] " + fmt.Sprintf(s, v...))
	}
}

func ERROR(s string, v ...interface{}) {
	core.OutPutLog("[ERROR] " + fmt.Sprintf(s, v...))
}

func META(s string, v ...interface{}) {
	if Conf.Meta {
		core.OutPutLog("[META] " + fmt.Sprintf(s, v...))
	}
}

func TEST(s string, v ...interface{}) {
	if Conf.Debug {
		core.OutPutLog("[TEST] " + fmt.Sprintf(s, v...))
	}
}

// PathExecute 返回当前运行目录
func PathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir + "/"
}

// CreatePath 生成路径或文件所对应的目录
func CreatePath(path string) {
	length := len(path)
	switch {
	case path[length:] != "/":
		path = path[:strings.LastIndex(path, "/")]
	case path[length:] != "\\":
		path = path[:strings.LastIndex(path, "\\")]
	default:
		//
	}
	if !PathExists(path) {
		err := os.MkdirAll(path, 0644)
		if err != nil {
			panic(err)
		}
	}
}

// PathExists 判断路径或文件是否存在
func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// FileSize 获取文件大小
func FileSize(file string) int64 {
	if fi, err := os.Stat(file); err == nil {
		return fi.Size()
	}
	return 0
}

func ReadAllText(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func WriteAllText(path, text string) {
	_ = ioutil.WriteFile(path, []byte(text), 0644)
}

func ProtectRun(entry func(), label string) {
	defer func() {
		err := recover()
		if err != nil {
			ERROR("[协程] %v协程发生了不可预知的错误，请在GitHub提交issue：%v", label, err)
			buf := make([]byte, 1<<16)
			runtime.Stack(buf, true)
			ERROR("traceback:\n%v", string(buf))
		}
	}()
	entry()
}

func unicode2chinese(text string) string {
	if !strings.Contains(text, "\\u") {
		return text
	}
	t := strings.Split(text, "\\u")
	var out string
	for _, v := range t {
		if len(v) < 1 {
			continue
		}
		if len(v) == 4 {
			temp, _ := strconv.ParseInt(v, 16, 32)
			out += fmt.Sprintf("%c", temp)
		} else {
			temp, _ := strconv.ParseInt(v[:3], 16, 32)
			out += fmt.Sprintf("%c", temp)
			out += fmt.Sprintf("%s", v[4:])
		}
	}
	return out
}

func hashText(input string) string {
	m := md5.New()
	m.Write([]byte(input))
	return hex.EncodeToString(m.Sum(nil))
}
