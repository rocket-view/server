package main

import (
	"math/rand"
	"time"

	"./config"
	"./server"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	config.Load()
	for _, srv := range config.Servers {
		server.StartServer(srv)
	}
	select {}
}
