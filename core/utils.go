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
	"encoding/binary"
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
