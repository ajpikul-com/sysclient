package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/ajpikul-com/wsssh/wsconn"
	"github.com/gorilla/mux"
)

func WriteText(conn *wsconn.WSConn) {
	for {
		_, err := conn.WriteText([]byte("Test Message")) // TODO Can we be sure this will write everything
		if err != nil {
			defaultLogger.Error("wsconn.WriteText(): " + err.Error())
			break
		}
		time.Sleep(1000 * time.Millisecond)
	}
}

func ReadTexts(conn *wsconn.WSConn) {
	defaultLogger.Debug("Starting to read texts")
	channel, _ := conn.SubscribeToTexts()
	for s := range channel {
		defaultLogger.Info("ReadTexts: " + s)
	}
	defaultLogger.Debug("ReadTexts Channel Closed")
	// The channel has been closed by someone else
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
	availableClients := globalClientList.GetAllClients()
	connectedClients := globalState.GetClientsCopy()
	_ = connectedClients
	jsonWrite := json.NewEncoder(w)
	jsonWrite.SetEscapeHTML(true)
	jsonWrite.Encode(availableClients)
	w.(http.Flusher).Flush()
}
