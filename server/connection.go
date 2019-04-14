package server

import (
	"fmt"
	"math/rand"

	"../config"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func openConnection(srv config.Server, cxn config.Connection, configChanged chan chan string, changeConfig chan string) {
	fmt.Printf("[%s -> %s] Opening connection to server...\n", srv.Name, cxn.Name)
	var proto string
	if cxn.Ssl {
		proto = "ssl"
	} else {
		proto = "tcp"
	}
	opts := MQTT.NewClientOptions().AddBroker(fmt.Sprintf("%s://%s:%d", proto, cxn.Host, cxn.Port))
	if cxn.ClientID == "" {
		cs := make([]byte, 23-len("rocketui"))
		for i := range cs {
			cs[i] = alphabet[rand.Int63()%int64(len(alphabet))]
		}
		opts.SetClientID(string(cs))
	} else {
		opts.SetClientID(cxn.ClientID)
	}
	opts.SetDefaultPublishHandler(func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("[%s -> %s] Updating config.\n", srv.Name, cxn.Name)
		changeConfig <- string(msg.Payload())
	})
	c := MQTT.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	if token := c.Subscribe("rocket_view/display_config/set", 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	_configChanged := make(chan string)
	configChanged <- _configChanged
	go func() {
		for {
			select {
			case cfg := <-_configChanged:
				fmt.Printf("[%s -> %s] Sending out new config.\n", srv.Name, cxn.Name)
				c.Publish("rocket_view/display_config", 0, true, cfg)
			default:
			}
		}
	}()
	fmt.Printf("[%s -> %s] Opened connection to server.\n", srv.Name, cxn.Name)
}
