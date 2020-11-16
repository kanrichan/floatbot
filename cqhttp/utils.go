package cqhttp

import (
	"fmt"
	"io/ioutil"
	"os"
)

func INFO(s string, v ...interface{}) {
	OutPutLog("[INFO] " + fmt.Sprintf(s, v...))
}

func WARN(s string, v ...interface{}) {
	OutPutLog("[WARN] " + fmt.Sprintf(s, v...))
}

func DEBUG(s string, v ...interface{}) {
	if Conf.Debug {
		OutPutLog("[DEBUG] " + fmt.Sprintf(s, v...))
	}
}

func ERROR(s string, v ...interface{}) {
	OutPutLog("[ERROR] " + fmt.Sprintf(s, v...))
}

func TEST(s string, v ...interface{}) {
	if Conf.Debug {
		OutPutLog("[TEST] " + fmt.Sprintf(s, v...))
	}
}

func PathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		ERROR("判断当前运行路径失败")
	}
	return dir + "/"
}

func CreatePath(path string) {
	err := os.MkdirAll(path, 0644)
	if err != nil {
		ERROR("生成应用目录失败")
	}
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
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
		}
	}()
	entry()
}
