package main

import (
	"github.com/ajpikul-com/ilog"
	"github.com/ajpikul-com/wsssh/wsconn"
)

var defaultLogger ilog.LoggerInterface

func init() {
	defaultLogger = &ilog.SimpleLogger{}
	err := defaultLogger.Init()
	if err != nil {
		panic(err)
	}
	packageLogger := &ilog.SimpleLogger{}
	packageLogger.Level(ilog.INFO)
	packageLogger.Init()
	wsconn.SetDefaultLogger(packageLogger)
}
