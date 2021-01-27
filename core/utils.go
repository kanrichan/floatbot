package core

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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strconv"
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
