package config

import (
	"encoding/json"
	"io/ioutil"
)

// ServerConfig holds parameters needed to start http server.
type ServerConfig struct {
	Version    string
	ServerAddr string
	LogLevel   string
	AccessKey  string
	Originator  string
}

// MakeServerConfigFromFile creates server config json config file
func MakeServerConfigFromFile(jsonFile string) (*ServerConfig, error) {
	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		return nil, err
	}
	return MakeServerConfig(b)
}

// MakeServerConfig creates server config from bytes
func MakeServerConfig(b []byte) (*ServerConfig, error) {
	c := &ServerConfig{}
	err := json.Unmarshal(b, c)
	if err != nil {
		return nil, err
	}
	return c, nil
}