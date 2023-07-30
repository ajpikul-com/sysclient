package main

import (
	"net/http"

	"golang.org/x/crypto/ssh"

	"github.com/ajpikul-com/wsssh/wsconn"
	"github.com/gorilla/websocket"
)

func ServeWSConn(w http.ResponseWriter, r *http.Request) {
	defaultLogger.Debug("Server: Incoming Req: " + r.Host + ", " + r.URL.Path)
	upgrader := &websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
		return true
	},
	}
	gorrilaconn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		defaultLogger.Error("Upgrade: " + err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	wsconn, err := wsconn.New(gorrilaconn)
	if err != nil {
		panic(err.Error()) // This actually is very severe
	}
	defer func() {
		defaultLogger.Debug("Closing WSConn")
		// Doesn't warn client, just closes
		if err := wsconn.CloseAll(); err != nil {
			defaultLogger.Error("wsconn.CloseAll(): " + err.Error())
		}
	}()

	sshconn, chans, reqs, err := GetServer(wsconn, globalConfig.PublicKeys, globalConfig.PrivateKey)
	if err != nil {
		defaultLogger.Error("GetServer(): " + err.Error())
		return
	}
	defaultLogger.Info("Welcome, " + sshconn.Permissions.Extensions["comment"])
	c := Client{
		Name:      sshconn.Permissions.Extensions["comment"],
		IPAddress: sshconn.RemoteAddr().String(),
	}
	globalState.UpdateClient(c)
	go ReadTexts(wsconn, c.Name)
	go ssh.DiscardRequests(reqs)
	for _ = range chans {
		// We do nothing for you
		// Otherwise it just stays open until client closes
	}
	defaultLogger.Info(sshconn.Permissions.Extensions["comment"] + " disconnected")
}
