package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajpikul-com/wsssh/wsconn"
	"github.com/gorilla/mux"
)

func init() {
	initLogger()
	initConfig()
}

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

func main() {
	initConfig()
	m := mux.NewRouter()
	// Add two gets here. Why are we using mux? I don't know
	m.HandleFunc("/", ServeWSConn)
	m.HandleFunc("/status", func(http.ResponseWriter, *http.Request) {})
	m.HandleFunc("/sysstatus", writeSystemStatus)
	s := &http.Server{
		Addr:           globalConfig.Hostname,
		Handler:        m,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	defaultLogger.Info("Initiating server")
	err := s.ListenAndServe()
	if err != nil {
		defaultLogger.Error("http.Server.ListenAndServe: " + err.Error())
	}
	defaultLogger.Info("Server exiting")
}

func writeSystemStatus(w http.ResponseWriter, r *http.Request) {
	// If r.Method = GET/POST
	jsonWrite := json.NewEncoder(w)
	jsonWrite.SetEscapeHTML(true)
	globalState.ReadLock()
	jsonWrite.Encode(globalState)
	globalState.ReadUnlock()
	w.(http.Flusher).Flush()
}
