package main

import (
	"fmt"

	"./config"
)

func main() {
	config.Load()
	fmt.Println(config.Servers)
}
