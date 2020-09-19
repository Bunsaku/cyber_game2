package main

import (
	"flag"
	"math/rand"
	"time"

	"20dojo-online/pkg/server"
)

var (
	// Listenするアドレス+ポート
	addr string
)

func init() {
	flag.StringVar(&addr, "addr", ":8080", "tcp host:port to connect")
	flag.Parse()
	//乱数設定
	rand.Seed(time.Now().UnixNano())
}

func main() {
	server.Serve(addr)
}
