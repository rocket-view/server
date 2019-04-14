package server

import (
	"fmt"
	"io/ioutil"

	"../config"
)

// StartServer runs the server loop for a server group from the configuration file
func StartServer(cfg config.Server) {
	fmt.Printf("[%s] Starting up server...\n", cfg.Name)
	configChanged := make(chan chan string)
	changeConfig := make(chan string)
	callbacks := make([]chan string, 0)
	go func() {
		for {
			select {
			case callback := <-configChanged:
				callbacks = append(callbacks, callback)
			case config := <-changeConfig:
				ioutil.WriteFile(cfg.DataFile, []byte(config), 0644)
				for _, callback := range callbacks {
					callback <- config
				}
			}
		}
	}()
	for _, cxn := range cfg.Connections {
		openConnection(cfg, cxn, configChanged, changeConfig)
	}
	fmt.Printf("[%s] Loading server config from file...\n", cfg.Name)
	data, err := ioutil.ReadFile(cfg.DataFile)
	var config string
	if err != nil {
		fmt.Printf("[%s] Unable to load configuration.  Creating blank config.\n", cfg.Name)
		config = "{}"
	} else {
		config = string(data)
	}
	for _, callback := range callbacks {
		callback <- config
	}
	fmt.Printf("[%s] Server started up.\n", cfg.Name)
}
