package testutils

import (
	"github.com/akhilesharora/hire/internal/config"
	"log"
	"net/http"
)

var ServerAddr = "127.0.0.1:7070"
// GetConfig returns ServerConfig with test configurations
func GetConfig() *config.ServerConfig {
	return &config.ServerConfig{
		Version:    "dev",
		ServerAddr: ServerAddr,
		LogLevel:   "debug",
		AccessKey: "Nxra5cRS0exx0zxLMMJ32voC4",
		Originator: "HelloTest1",
	}
}

type TestSmsClient struct {}

func (cl *TestSmsClient) sendSms(msg Messages) error {
	return nil
}

type ServersTest struct {
	conf *config.ServerConfig
	cl *TestSmsClient
	q chan Messages
}

type Messages struct {
	Recipients []string
	Originator string
	Message string
}
func StartTestServerInstance(){
	// Run server
	log.Println("Listen and service on", ServerAddr)
	err := http.ListenAndServe(ServerAddr, nil)
	if err != nil {
		log.Fatal("Could not start server: ", err.Error())
	}
}
