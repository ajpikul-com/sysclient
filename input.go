package main

import (
	"bytes"
	"encoding/json"

	"github.com/ajpikul-com/wsssh/wsconn"
)

func ReadTexts(conn *wsconn.WSConn, name string) {
	defaultLogger.Debug("Starting to read texts")
	channel, _ := conn.SubscribeToTexts()
	buffer := bytes.NewBuffer([]byte{})
	commandDecoder := json.NewDecoder(buffer)
	go func() {
		var command interface{}
		err := commandDecoder.Decode(command) // what happens if the buffer isn't complete json currenlty
		if err != nil {
			panic(err.Error())
		}
		if commandService, ok := command.(Service); ok {
			_ = commandService
			// TODO okay write new service
			//globalState.UpdateService(name, commandService.Name, commandService.Status, commandService.LastConnection)
			// We should be able to test the service at this point
		} else {
			panic("Command wasn't a Service type")
		}
	}()
	for s := range channel {
		buffer.WriteString(s)
	}
	defaultLogger.Debug("ReadTexts Channel Closed")
}
