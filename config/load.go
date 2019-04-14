package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var (
	// Servers is all of the server blocks from the configuration file
	Servers []Server
)

// Load the configuration settings from the filesystem
func Load() {
	viper.SetConfigName("rocket-view")
	viper.AddConfigPath("/etc/rocket-view/")
	viper.AddConfigPath("$HOME/.rocket-view/")
	viper.AddConfigPath(".")
	viper.SetConfigType("hcl")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
	Servers = make([]Server, 0)
	for _, serverMap := range viper.Get("server").([]map[string]interface{}) {
		for serverName, serverConfig := range serverMap {
			srv := Server{
				Name:        serverName,
				Connections: make([]Connection, 0),
				DataFile:    "",
			}
			for _, serverConfigMap := range serverConfig.([]map[string]interface{}) {
				if cxns, ok := serverConfigMap["connection"]; ok {
					for _, cxnMap := range cxns.([]map[string]interface{}) {
						for cxnName, cxnConfigs := range cxnMap {
							cxn := Connection{
								Name:     cxnName,
								Host:     "",
								Port:     0,
								Username: "",
								Password: "",
								ClientID: "",
								Ssl:      true,
							}
							for _, cxnConfig := range cxnConfigs.([]map[string]interface{}) {
								if host, ok := cxnConfig["host"]; ok {
									cxn.Host = host.(string)
								}
								if port, ok := cxnConfig["port"]; ok {
									cxn.Port = uint16(port.(int))
								}
								if user, ok := cxnConfig["username"]; ok {
									cxn.Username = user.(string)
								}
								if pass, ok := cxnConfig["password"]; ok {
									cxn.Password = pass.(string)
								}
								if client, ok := cxnConfig["clientID"]; ok {
									cxn.ClientID = client.(string)
								}
								if ssl, ok := cxnConfig["ssl"]; ok {
									cxn.Ssl = ssl.(bool)
								}
							}
							if cxn.Host == "" {
								panic(fmt.Sprintf("Connection '%s' missing key 'host'.", cxn.Name))
							}
							if cxn.Port == 0 {
								panic(fmt.Sprintf("Connection '%s' missing key 'port'.", cxn.Name))
							}
							srv.Connections = append(srv.Connections, cxn)
						}
					}
				}
				if data, ok := serverConfigMap["dataFile"]; ok {
					if srv.DataFile != "" {
						panic("Can only have one data file per server configuration")
					}
					srv.DataFile = data.(string)
				}
			}
			if len(srv.Connections) == 0 {
				panic(fmt.Sprintf("Server '%s' does not have any connections.", srv.Name))
			}
			if srv.DataFile == "" {
				panic(fmt.Sprintf("Server '%s' missing key 'dataFile'.", srv.Name))
			}
			Servers = append(Servers, srv)
		}
	}
}
