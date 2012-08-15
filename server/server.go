package main

import (
	"github.com/kisielk/gcfg3/gcfg3"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

func main() {
	server := gcfg3.Gcfg3Server{"."}
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("Listen error:", e)
	}
	http.Serve(l, nil)
}
