package core

//#include <xqapi.h>
import "C"

import (
	"fmt"
	"strconv"
	"strings"
)

func SendMsg(selfID int64, messageType int64, groupID int64, userID int64, message string, bubble int64) {
	C.S3_Api_SendMsg(
		GoInt2CStr(selfID), C.int(messageType), GoInt2CStr(groupID), GoInt2CStr(userID), CString(message), C.int(bubble),
	)
}

func OutPutLog(text string) {
	C.S3_Api_OutPutLog(
		CString(text),
	)
}

func GetGroupList(selfID int64) string {
	return GoString(C.S3_Api_GetGroupList(
		GoInt2CStr(selfID),
	))
}

// sender
func GetNick(selfID int64, userID int64) string {
	context := GoString(C.S3_Api_GetNick(
		GoInt2CStr(selfID),
		GoInt2CStr(userID),
	))
	OutPutLog(fmt.Sprint(context))
	sUnicodev := strings.Split(fmt.Sprint(context), `\u`)
	var nickname string
	for _, v := range sUnicodev {
		if len(v) < 1 {
			continue
		}
		temp, err := strconv.ParseInt(v, 16, 32)
		if err != nil {
			OutPutLog("昵称解码失败")
		}
		nickname += fmt.Sprintf("%c", temp)
	}
	return "nickname"
}

func GetGender(selfID int64, userID int64) string {
	sex := int64(C.S3_Api_GetGender(
		GoInt2CStr(selfID),
		GoInt2CStr(userID),
	))
	switch sex {
	case 1:
		return "male"
	case 2:
		return "female"
	default:
		return "unknow"
	}
}

func GetAge(selfID int64, userID int64) string {
	age := int64(C.S3_Api_GetAge(
		GoInt2CStr(selfID),
		GoInt2CStr(userID),
	))
	return Int2Str(age)
}
