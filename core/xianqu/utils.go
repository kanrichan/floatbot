package xianqu

//#include <string.h>
import "C"

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unsafe"

	sc "golang.org/x/text/encoding/simplifiedchinese"
)

// CString 将 GO 字符串 转为 C char指针
func CString(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
}

// GoString 将 C char指针 转为 GO 字符串
func GoString(str *C.char) string {
	if str == nil {
		return ""
	}
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

// CBool 将 GO 布尔型 转为 C 整数
func CBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

// GoBool 将 C 整数 转为 GO 布尔型
func GoBool(b C.int) bool {
	if b == 1 {
		return true
	}
	return false
}

// CByte 将 GO 字节数组 转为 C 字符串指针
func CByte(bt []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&bt))
}

// Str2Int 将string转为int64
func Str2Int(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

// Int2Str 将int64转为string
func Int2Str(val int64) string {
	str := strconv.FormatInt(val, 10)
	return str
}

// GoInt2CStr 将 GO int64 转为 C char指针
func GoInt2CStr(val int64) *C.char {
	if val == 0 {
		return CString("")
	}
	return CString(Int2Str(val))
}

// CStr2GoInt 将 C char指针 转为 GO int64
func CStr2GoInt(str *C.char) int64 {
	return Str2Int(GoString(str))
}

// Int2Bytes 将 int64 转为字节数组
func Int2Bytes(val int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(val))
	return b
}

// EscapeEmoji 将 emoji 转化为 先驱[emoji=FFFFFFFF]
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

// UnescapeEmoji 将 先驱[emoji=FFFFFFFF] 转化为 emoji
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

// XmlEscape XML 编码
func XmlEscape(c string) string {
	buf := new(bytes.Buffer)
	_ = xml.EscapeText(buf, []byte(c))
	return buf.String()
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
	if !PathExists(filepath.Dir(path)) {
		err := os.MkdirAll(filepath.Dir(path), 0644)
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

// FileSize 返回文件大小
func FileSize(file string) int64 {
	if fi, err := os.Stat(file); err == nil {
		return fi.Size()
	}
	return 0
}

// ReadAllText 返回文件字符串
func ReadAllText(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

// WriteAllText 向文件写入字符串
func WriteAllText(path, text string) {
	_ = ioutil.WriteFile(path, []byte(text), 0644)
}

// TextMD5 返回字符串的MD5值
func TextMD5(input string) string {
	m := md5.New()
	m.Write([]byte(input))
	return hex.EncodeToString(m.Sum(nil))
}

// FileMD5 返回文件的MD5值
func FileMD5(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	m := md5.New()
	m.Write(b)
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

// GetBnk 返回 tx cookie 的 bnk
func GetBnk(cookie string) (bnk int) {
	skey := cookie[strings.Index(cookie, "skey=")+5:]
	bnk = 5381
	for i := range skey {
		bnk += (bnk << 5) + int(skey[i])
	}
	return bnk & 2147483647
}

// escape 临时应对 nonebot CQ码转数组未反转义的问题
func escape(text string) string {
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&#44;", ",")
	text = strings.ReplaceAll(text, "&#91;", "[")
	text = strings.ReplaceAll(text, "&#93;", "]")
	return text
}
