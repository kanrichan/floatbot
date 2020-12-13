package core

import "C"

var Create func(version string) string
var Event func(selfID int64, mseeageType int64, subType int64, groupID int64, userID int64, noticID int64, message string, messageNum int64, messageID int64, rawMessage []byte, time int64, ret int64) int64
var DestroyPlugin func() int64
var SetUp func() int64

//export GO_Create
func GO_Create(version *C.char) *C.char {
	return CString(Create(GoString(version)))
}

//export GO_Event
func GO_Event(selfID *C.char, mseeageType C.int, subType C.int, groupID *C.char, userID *C.char, noticID *C.char, message *C.char, messageNum *C.char, messageID *C.char, rawMessage *C.char, time *C.char, ret *C.char) C.int {
	return C.int(Event(CStr2GoInt(selfID),
		int64(mseeageType),
		int64(subType),
		CStr2GoInt(groupID),
		CStr2GoInt(userID),
		CStr2GoInt(noticID),
		GoString(message),
		CStr2GoInt(messageNum),
		CStr2GoInt(messageID),
		[]byte(GoString(rawMessage)),
		CStr2GoInt(time),
		CStr2GoInt(ret),
	))
}

//export GO_DestroyPlugin
func GO_DestroyPlugin() C.int {
	return C.int(DestroyPlugin())
}

//export GO_SetUp
func GO_SetUp() C.int {
	return C.int(SetUp())
}

func main() {
	//
}
