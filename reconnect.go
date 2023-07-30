package main

import (
	"github.com/ajpikul-com/wsssh/wsconn"
	gws "github.com/gorilla/websocket"
)

func Reconnect() (*wsconn.WSConn, error) {
	dialer := gws.Dialer{}
	conn1, _, err := dialer.Dial(globalConfig.Hostname, nil) // TODO: maybe dump response
	if err != nil {
		defaultLogger.Error("websocket.Dialier.Dial: Dial fail: " + err.Error())
		return nil, err
	}
	wssshConn, err := wsconn.New(conn1)
	if err != nil {
		defaultLogger.Error("Problem with wsconn: " + err.Error())
		return nil, err
	}
	_, err = GetClient(wssshConn, "ajp", globalConfig.PrivateKey, globalConfig.HostKey)
	if err != nil {
		defaultLogger.Error("Couldn't auth: " + err.Error())
		return nil, err
	}
	return wssshConn, nil
}
