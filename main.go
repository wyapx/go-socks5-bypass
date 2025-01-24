package main

import (
	socks5 "bypass/src"
	"fmt"
)

func main() {
	bindAddr := "0.0.0.0:7900"
	fmt.Printf("Socks5 server running at %s\n", bindAddr)
	conf := &socks5.Config{}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", bindAddr); err != nil {
		panic(err)
	}
}
