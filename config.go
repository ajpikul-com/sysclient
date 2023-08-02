package main

import (
	"encoding/json"
	"os"
)

type GetServices struct {
	Name   string
	URL    string
	Module string
}

type GitServices struct {
	DataStore string
	Module    string
}

type config struct {
	MyName      string
	Hostname    string
	PrivateKey  string
	HostKey     string
	GetServices []GetServices
	GitService  GitServices
}

var globalConfig config // This is the sysboss config

func initConfig() {

	configFile, err := os.Open("/home/ajp/systems/sysclient/sysclient.json")
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
