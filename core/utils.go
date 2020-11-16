package core

import "C"
import "strconv"
import sc "golang.org/x/text/encoding/simplifiedchinese"

func CString(str string) *C.char {
	gbstr, _ := sc.GB18030.NewEncoder().String(str)
	return C.CString(gbstr)
}

func GoString(str *C.char) string {
	utf8str, _ := sc.GB18030.NewDecoder().String(C.GoString(str))
	return utf8str
}

func CBool(b bool) C.int {
	if b {
		return 1
	}
	return 0
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
	return CString(Int2Str(val))
}

func CStr2GoInt(str *C.char) int64 {
	return Str2Int(GoString(str))
}
