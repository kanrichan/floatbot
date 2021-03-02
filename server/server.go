package server

import "fmt"

var (
	LOG = func(s string, v ...interface{}) { fmt.Printf(s, v...) }
)

func ERROR(s interface{}, err interface{}) {
	switch s.(type) {
	case *HttpServer:
		this := s.(*HttpServer)
		LOG("[INFO][HTTP] ID: %d ADDR: %s POST_URL: %s ERROR: %v\n", this.id, this.addr, this.url, err)
	case *WSC:
		this := s.(*WSC)
		LOG("[ERROR][WSC] ID: %d URL: %s ERROR: %v\n", this.id, this.addr, err)
	case *WSS:
		this := s.(*WSS)
		LOG("[ERROR]][WSS] ID: %d ADDR: %s ERROR: %v\n", this.id, this.addr, err)
	}
}

func INFO(s interface{}, info interface{}) {
	switch s.(type) {
	case *HttpServer:
		this := s.(*HttpServer)
		LOG("[INFO][HTTP] ID: %d ADDR: %s POST_URL: %s INFO: %v\n", this.id, this.addr, this.url, info)
	case *WSC:
		this := s.(*WSC)
		LOG("[INFO][WSC] ID: %d URL: %s INFO: %v\n", this.id, this.addr, info)
	case *WSS:
		this := s.(*WSS)
		LOG("[INFO][WSS] ID: %d ADDR: %s INFO: %v\n", this.id, this.addr, info)
	}
}
