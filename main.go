package main

import (
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ajpikul-com/gitstatus"
	"github.com/ajpikul-com/wsssh/wsconn"
	gws "github.com/gorilla/websocket"
)

func init() {
	initLogger()
	initConfig()

	gitstatus.InitDataStore(globalConfig.GitService.DataStore)
	gitstatus.UpdateRepos()

}

func WriteText(conn *wsconn.WSConn) {
	gitI := 0
	for {
		// So what we're going to do here
		defaultLogger.Debug("Iterating over get services")
		for _, v := range globalConfig.GetServices {
			status := "Offline"
			if v.Module == "GET" {
				res, err := http.Get(v.URL)
				if err != nil {
					status = err.Error()
				} else if res.StatusCode == http.StatusOK {
					status = "Online"
				}
			}
			payload := map[string]*Service{
				"get": &Service{ // Is this not a protobuff
					Name:           v.Name,
					Status:         status,
					ParentService:  globalConfig.MyName,
					LastConnection: time.Now(),
				}}
			b, err := json.Marshal(payload)
			defaultLogger.Debug(string(b))
			if err != nil {
				defaultLogger.Error(err.Error())
				continue // skip this service
			}
			defaultLogger.Debug("Writing")
			_, err = conn.WriteText(b) // TODO Can we be sure this will write everything
			if err != nil {
				defaultLogger.Error("wsconn.WriteText(): " + err.Error())
				break
			}
		}
		if gitI == 0 {
			defaultLogger.Debug("In git")
			repostates := gitstatus.VerifyRepos() // payload is a map of gitstatus.RepoStates, we will occasionally send it

			payload := map[string]map[string]gitstatus.RepoState{
				"git": repostates}
			b, err := json.Marshal(payload)
			defaultLogger.Debug("About to send")
			defaultLogger.Debug(string(b))
			if err != nil {
				defaultLogger.Error(err.Error())
				continue // skip this service
			}
			defaultLogger.Debug("Writing")
			_, err = conn.WriteText(b) // TODO Can we be sure this will write everything
			if err != nil {
				defaultLogger.Error("wsconn.WriteText(): " + err.Error())
				break
			}
			gitI = 5
		}
		gitI -= 1
		time.Sleep(5 * time.Minute)
	}
}

func Pinger(conn *wsconn.WSConn) error {
	defaultLogger.Debug("Beggining Ping Loop")
	for {
		err := conn.WritePing([]byte("Pingaring'll Payload"))
		if err != nil {
			defaultLogger.Error("Pinger dead: " + err.Error())
			return err
		}
		time.Sleep(10000 * time.Millisecond)
	}
	defaultLogger.Debug("Ending Ping Loop, will never get here")
	return nil
}

func main() {

	var wg sync.WaitGroup
	var wsconn *wsconn.WSConn

	// Start goroutine to wait for signal to close
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	wg.Add(1)
	go func() {
		defaultLogger.Debug("Waiting for SIGINT")
		for _ = range c {
			// Here we exit the program
			defaultLogger.Info("Recieved SIGINT")
			break
		}
		wg.Done()
	}()
	go func() {
		for { // All this depends on ssh sitting on top of Read() which TODO Not sure it does
			var err error
			defaultLogger.Debug("Trying to reconnect")
			wsconn, err = Reconnect()
			if err != nil {
				defaultLogger.Error("Problem with reconnect: " + err.Error())
				time.Sleep(20 * time.Second)
				continue
			}
			go ReadTexts(wsconn) // No real reason to do this yet
			go WriteText(wsconn)
			err = Pinger(wsconn)
			defaultLogger.Debug("Pinger Error: " + err.Error())
			// Why bother, we can't do this if pinger failed!
			wsconn.Conn.WriteControl(gws.CloseMessage, []byte(""), time.Time{})
			// We're closing all to signal the end of some go routines
			wsconn.CloseAll()
		}
		defaultLogger.Debug("Trying to send myself interrupt")

		pid := os.Getpid()
		p, _ := os.FindProcess(pid)
		_ = p.Signal(os.Interrupt)

	}()

	wg.Wait()
	defaultLogger.Debug("passed wg.Wait()")
	if wsconn != nil {
		defaultLogger.Debug("Trying to close cleanly")
		wsconn.Conn.WriteControl(gws.CloseMessage, []byte(""), time.Time{})
		err := wsconn.CloseAll()
		if err != nil {
			defaultLogger.Info("Tried to close: " + err.Error())
		}
	}
}
