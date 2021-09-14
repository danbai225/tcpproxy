package main

import (
	"flag"
	"fmt"
	logs "github.com/danbai225/go-logs"
)

var server bool
var (
	addr string
	port string
	pass string
)

func main() {
	flag.BoolVar(&server, "s", false, "服务端模式")
	flag.StringVar(&port, "port", "18889", "监听端口")
	flag.StringVar(&pass, "ps", "", "密码")
	flag.StringVar(&addr, "addr", "", "目标地址")
	flag.Parse()
	if addr == "" {
		logs.Err("缺少目标地址")
		return
	}
	if server {
		Server{}.New(pass, addr, fmt.Sprintf(":%s", port)).Start()
	} else {
		Client{}.New(pass, addr, fmt.Sprintf(":%s", port)).Start()
	}
}
