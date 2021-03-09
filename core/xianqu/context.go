package xianqu

import (
	"fmt"
	"strconv"
)

type value map[string]interface{}

func Parse(m map[string]interface{}) value {
	return m
}

func (v value) Get(path string) value {
	return Parse(v[path].(map[string]interface{}))
}

func (v value) Int(path string) int64 {
	if v[path] == nil {
		return 0
	}
	temp := v[path]
	switch temp.(type) {
	case float64:
		return int64(temp.(float64))
	case int64:
		return temp.(int64)
	case string:
		r, _ := strconv.ParseInt(temp.(string), 10, 64)
		return r
	case bool:
		if temp.(bool) {
			return 1
		}
		return 0
	case interface{}:
		return 0
	}
	return 0
}

func (v value) Str(path string) string {
	if v[path] == nil {
		return ""
	}
	temp := v[path]
	switch temp.(type) {
	case string:
		return temp.(string)
	case float64:
		return strconv.FormatFloat(temp.(float64), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(temp.(bool))
	case interface{}:
		return fmt.Sprint(temp)
	}
	return ""
}

func (v value) Bool(path string) bool {
	if v[path] == nil {
		return false
	}
	temp := v[path]
	switch temp.(type) {
	case bool:
		return temp.(bool)
	case float64:
		if temp.(float64) != 0 {
			return true
		}
		return false
	case string:
		r, _ := strconv.ParseBool(temp.(string))
		return r
	case interface{}:
		return false
	}
	return false
}

func (v value) Array(path string) []value {
	temp := []value{}
	switch v[path].(type) {
	case []map[string]interface{}:
		for _, e := range v[path].([]map[string]interface{}) {
			temp = append(temp, e)
		}
	case []interface{}:
		for _, e := range v[path].([]interface{}) {
			temp = append(temp, e.(map[string]interface{}))
		}
	}
	return temp
}

func (v value) Exist(path string) bool {
	if v[path] == nil {
		return false
	}
	return true
}
