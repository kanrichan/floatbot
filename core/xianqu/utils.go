package xianqu

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

	"github.com/wdvxdr1123/ZeroBot/utils/helper"
	sc "golang.org/x/text/encoding/simplifiedchinese"
)

// CString 将 GO 字符串 转为 C char指针
func cString(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
}

// GoString 将 C char指针 转为 GO 字符串
func goString(str *C.char) string {
	if str == nil {
		return ""
	}
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

// CBool 将 GO 布尔型 转为 C 整数
func cBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

// GoBool 将 C 整数 转为 GO 布尔型
func goBool(b C.int) bool {
	if b == 1 {
		return true
	}
	return false
}

// CByte 将 GO 字节数组 转为 C 字符串指针
func cByte(bt []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&bt))
}

// Str2Int 将string转为int64
func str2Int(str string) int64 {
	val, _ := strconv.ParseInt(str, 10, 64)
	return val
}

// Int2Str 将int64转为string
func int2Str(val int64) string {
	str := strconv.FormatInt(val, 10)
	return str
}

// GoInt2CStr 将 GO int64 转为 C char指针
func goInt2CStr(val int64) *C.char {
	if val == 0 {
		return cString("")
	}
	return cString(int2Str(val))
}

// CStr2GoInt 将 C char指针 转为 GO int64
func cStr2GoInt(str *C.char) int64 {
	return str2Int(goString(str))
}

// Int2Bytes 将 int64 转为字节数组
func int2Bytes(val int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(val))
	return b
}

// EscapeEmoji 将 emoji 转化为 先驱[emoji=FFFFFFFF]
func escapeEmoji(text string) string {
	data := helper.StringToBytes(text)
	ret := []byte{}
	skip := 0
	for i := range data {
		if skip > 1 {
			skip -= 1
			continue
		}
		if data[i] == byte(240) && data[i+1] == byte(159) {
			code := hex.EncodeToString(data[i : i+4])
			ret = append(ret, helper.StringToBytes(fmt.Sprintf("[emoji=%s]", code))...)
			skip = 4
		} else {
			ret = append(ret, data[i])
		}
	}
	return helper.BytesToString(ret)
}

// UnescapeEmoji 将 先驱[emoji=FFFFFFFF] 转化为 emoji
func unescapeEmoji(text string) string {
	data := helper.StringToBytes(text)
	ret := []byte{}
	skip := 0
	for i := range data {
		if skip > 1 {
			skip -= 1
			continue
		}
		if i+7 < len(data) && bytes.Equal(data[i:i+7], helper.StringToBytes("[emoji=")) {
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
	return helper.BytesToString(ret)
}

// XmlEscape XML 编码
func xmlEscape(c string) string {
	buf := new(bytes.Buffer)
	_ = xml.EscapeText(buf, helper.StringToBytes(c))
	return buf.String()
}

// PathExecute 返回当前运行目录
func pathExecute() string {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return dir + "\\"
}

// CreatePath 生成路径或文件所对应的目录
func createPath(path string) {
	if !pathExists(filepath.Dir(path)) {
		err := os.MkdirAll(filepath.Dir(path), 0644)
		if err != nil {
			panic(err)
		}
	}
}

// PathExists 判断路径或文件是否存在
func pathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// FileSize 返回文件大小
func fileSize(file string) int64 {
	if fi, err := os.Stat(file); err == nil {
		return fi.Size()
	}
	return 0
}

// TextMD5 返回字符串的MD5值
func textMD5(input string) string {
	m := md5.New()
	m.Write(helper.StringToBytes(input))
	return hex.EncodeToString(m.Sum(nil))
}

// FileMD5 返回文件的MD5值
func fileMD5(path string) string {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return ""
	}
	m := md5.New()
	m.Write(b)
	return strings.ToUpper(hex.EncodeToString(m.Sum(nil)))
}

// GetBnk 返回 tx cookie 的 bnk
func getBnk(cookie string) (bnk int) {
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
