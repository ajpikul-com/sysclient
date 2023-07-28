package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Hostname   string
	PublicKeys string
	PrivateKey string
}

// b, err := json.Marshal(instance of config)
var globalConfig config
var globalState *systemState
var globalClientList *clientList

func initConfig() {
	globalState = NewSystemState()
	globalClientList = NewClientList([]string{})
	// globalClientList = NewClientList() // We actually don't know yet
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
