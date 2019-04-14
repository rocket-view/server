package config

// Server represents a block in the configuration file
type Server struct {
	Name        string
	Connections []Connection
	DataFile    string
}
