package testutils

import (
	"log"

	"github.com/messagebird/internal/config"
)

// GetConfig returns ServerConfig with test configurations
func GetConfig() *config.ServerConfig {
	testConfig := []byte(`{
  "Version": "dev",
  "ServerAddr": "http://127.0.0.1:8070",
  "LogLevel": "debug",
  "AccessKey": "",
  "Originator": ""
}`)
	cnf, err := config.MakeServerConfig(testConfig)
	if err != nil {
		log.Fatal("Error creating config", err)
	}
	return cnf
}

