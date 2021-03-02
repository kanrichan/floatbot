package xianqu

/*
char * eStrPtr2CStrPtr(char * str) {
	if (!str) {
		return NULL;
	}
	return str + 4;
}
*/
import "C"

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unsafe"

	sc "golang.org/x/text/encoding/simplifiedchinese"
)

func CString(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
}

func GoString(str *C.char) string {
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

func CPtr2GoStr(str *C.char) string {
	ptr := C.eStrPtr2CStrPtr(str)
	if ptr != nil {
		utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(ptr))
		return utf8str
	}
	return ""
}

func CBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

func GoBool(b C.int) bool {
	if b == 1 {
		return true
	}
	return false
}

func CByte(bt []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&bt))
}

func Str2Int(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

func Int2Str(val int64) string {
	str := strconv.FormatInt(val, 10)
	return str
}

func GoInt2CStr(val int64) *C.char {
	if val == 0 {
		return CString("")
	}
	return CString(Int2Str(val))
}

func CStr2GoInt(str *C.char) int64 {
	return Str2Int(GoString(str))
}

func Int2Bytes(val int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(val))
	return b
}

// EscapeEmoji 将 emoji的utf-8字符串 转化为 [emoji=FFFFFFFF]
func EscapeEmoji(text string) string {
	data := []byte(text)
	ret := []byte{}
	skip := 0
	for i := range data {
		if skip > 1 {
			skip -= 1
			continue
		}
		if data[i] == byte(240) && data[i+1] == byte(159) {
			code := hex.EncodeToString(data[i : i+4])
			ret = append(ret, []byte(fmt.Sprintf("[emoji=%s]", code))...)
			skip = 4
		} else {
			ret = append(ret, data[i])
		}
	}
	return string(ret)
}

// UnescapeEmoji 将 [emoji=FFFFFFFF] 转化为 emoji的utf-8字符串
func UnescapeEmoji(text string) string {
	data := []byte(text)
	ret := []byte{}
	skip := 0
	for i := range data {
		if skip > 1 {
			skip -= 1
			continue
		}
		if i+7 < len(data) && bytes.Equal(data[i:i+7], []byte("[emoji=")) {
			end := bytes.IndexByte(data[i:], byte(93))
			if end == -1 {
				return text
			}
			code, _ := hex.DecodeString(string(data[i+7 : end+i]))
			ret = append(ret, code...)
			skip = end + 1
		} else {
			ret = append(ret, data[i])
		}
	}
	return string(ret)
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

// PathExecute 返回当前运行目录
func PathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir + "\\"
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

func hashText(input string) string {
	m := md5.New()
	m.Write([]byte(input))
	return hex.EncodeToString(m.Sum(nil))
}
