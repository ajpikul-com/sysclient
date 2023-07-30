package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Hostname          string
	PublicKeys        string
	PrivateKey        string
	authorizedKeysMap map[string]string
}

var globalConfig config      // This is the sysboss config
var globalState *systemState // This is the whol system state

func initConfig() {
	globalState = NewSystemState()

	configFile, err := os.Open("/home/ajp/systems/ajpikul.com_system/configs/sysboss.json")
	if err != nil {
		panic(err.Error())
	}
	defer configFile.Close()
	configDecoder := json.NewDecoder(configFile)
	if err != nil {
		panic(err.Error())
	}
	//err := json.Unmarshal(bytes, &config)
	err = configDecoder.Decode(&globalConfig)
	if err != nil {
		panic(err.Error())
	}
}
