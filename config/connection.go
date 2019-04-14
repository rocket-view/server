package config

// Connection represents a block in the configuration file
type Connection struct {
	Name     string
	Host     string
	Port     uint16
	Username string
	Password string
	Ssl      bool
}
