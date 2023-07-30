package main

import ()

func init() {
	initLogger()
	initConfig()
}

func main() {
	// We need to make connection to sysboss
	// We need to send service states for the services we're expecting to run- that comes from config
	// Detailed information might come from module?
	// Again, another great application for RPC
}
