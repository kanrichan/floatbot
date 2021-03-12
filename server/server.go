package server

import "fmt"

var (
	WSSHandler      = func(bot int64, data []byte) []byte { fmt.Println(string(data)); return []byte("ok") }
	WSCHandler      = func(bot int64, data []byte) []byte { fmt.Println(string(data)); return []byte("ok") }
	HttpPostHandler = func(bot int64, send []byte, data []byte) { fmt.Println(string(data)) }
	HttpHandler     = func(bot int64, path string, data []byte) []byte { fmt.Println(string(data)); return []byte("ok") }

	CoreInfo  = func(s string, v ...interface{}) { fmt.Printf(s, v...) }
	CoreDebug = func(s string, v ...interface{}) { fmt.Printf(s, v...) }
)

type Server interface {
	Run()
	Close()
	Send(data []byte)
}

func (s *WSC) INFO(text interface{}) {
	CoreInfo("[I][WSC][%d][%s] 信息: %v", s.ID, s.Addr, text)
}

func (s *WSC) DEBUG(text interface{}) {
	CoreDebug("[D][WSC][%d][%s] 信息: %v", s.ID, s.Addr, text)
}

func (s *WSC) ERROR(text interface{}) {
	CoreInfo("[E][WSC][%d][%s] 错误: %v", s.ID, s.Addr, text)
}

func (s *WSC) PANIC(err interface{}, traceback []byte) {
	CoreInfo("[P][WSC][%d][%s] \n[错误]\n%v\n[TRACEBACK]\n%s", s.ID, s.Addr, err, traceback)
}

func (s *WSS) INFO(text interface{}) {
	CoreInfo("[I][WSS][%d][%s] 信息: %v", s.ID, s.Addr, text)
}

func (s *WSS) DEBUG(text interface{}) {
	CoreDebug("[D][WSS][%d][%s] 信息: %v", s.ID, s.Addr, text)
}

func (s *WSS) ERROR(text interface{}) {
	CoreInfo("[E][WSS][%d][%s] 错误: %v ", s.ID, s.Addr, text)
}

func (s *WSS) PANIC(err interface{}, traceback []byte) {
	CoreInfo("[P][WSS][%d][%s] \n[错误]\n%v\n[TRACEBACK]\n%s", s.ID, s.Addr, err, traceback)
}

func (s *HTTP) INFO(text interface{}) {
	CoreInfo("[I][HTTP][%d][%s] 信息: %v", s.ID, s.Addr, text)
}

func (s *HTTP) DEBUG(text interface{}) {
	CoreDebug("[D][HTTP][%d][%s] 信息: %v", s.ID, s.Addr, text)
}

func (s *HTTP) ERROR(text interface{}) {
	CoreInfo("[E][HTTP][%d][%s] 错误: %v", s.ID, s.Addr, text)
}

func (s *HTTP) PANIC(err interface{}, traceback []byte) {
	CoreInfo("[P][HTTP][%d][%s] \n[错误]\n%v\n[TRACEBACK]\n%s", s.ID, s.Addr, err, traceback)
}
