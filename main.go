package main

import (
	"./config"
	"./server"
)

func main() {
	config.Load()
	for _, srv := range config.Servers {
		server.StartServer(srv)
	}
	select {}
}
