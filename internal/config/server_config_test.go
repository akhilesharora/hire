package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	sampleConfig = []byte(`{
	"Version": "dev",
	"ServerAddr": "localhost:7070",
	"LogLevel": "debug",
	"AccessKey": "",
	"Originator":""
}`)
)

func TestMakeServerConfigFromWrongFile(t *testing.T) {
	_, err := MakeServerConfigFromFile("fake.json")
	assert.Error(t, err)
}

func TestMakeServerConfig(t *testing.T) {
	config, err := MakeServerConfig(sampleConfig)

	expected := &ServerConfig{
		Version:    "dev",
		ServerAddr: "localhost:7070",
		LogLevel:   "debug",
		AccessKey: "",
		Originator: "",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, config)
}

func TestMakeServerConfigFromFile(t *testing.T) {
	config, err := MakeServerConfigFromFile("config.test.json")

	expected := &ServerConfig{
		Version:    "dev",
		ServerAddr: "localhost:7070",
		LogLevel:   "debug",
		AccessKey: "",
		Originator: "",
	}

	assert.NoError(t, err)
	assert.Equal(t, expected, config)
}