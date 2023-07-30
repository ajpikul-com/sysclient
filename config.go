package main

import (
	"encoding/json"
	"os"
)

type config struct {
	Hostname         string
	PublicKeys       string
	PrivateKey       string
	authorizedKeyMap map[string]string
}

var globalConfig config          // This is the sysboss config
var globalState *systemState     // This is the whol system state
var globalClientList *clientList // This is the list of possible clients (we need to know this!) I think it should be added to global state

func initConfig() {
	globalState = NewSystemState()
	globalClientList = NewClientList([]string{}) // This needs to be somewhere else

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
