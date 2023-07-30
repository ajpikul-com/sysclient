package main

// This could be a lot better by having per-client reading/writing

import (
	"sync"
	"time"
)

type Service struct { // Is this not a protobuff
	Name           string
	IPAddress      string
	Status         string
	ParentService  string
	LastConnection time.Time
}

type systemState struct {
	Services        map[string]*Service
	ExpectedClients []string
	mutex           sync.RWMutex
}

func NewSystemState() *systemState {
	ss := &systemState{
		Services:        make(map[string]*Service),
		ExpectedClients: make([]string, 0),
	}
	return ss
}

func (ss *systemState) UpdateService(uS Service) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.Services[uS.Name] = &uS
	if ss.Services[uS.Name].LastConnection.IsZero() {
		ss.Services[uS.Name].LastConnection = time.Now()
	}
}

func (ss *systemState) UpdateTime(serviceName string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.Services[serviceName].LastConnection = time.Now()
}

func (ss *systemState) UpdateClientList(names []string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	copy(ss.ExpectedClients, names)

}
func (ss *systemState) ReadLock() {
	ss.mutex.RLock()
}
func (ss *systemState) ReadUnlock() {
	ss.mutex.RUnlock()
}

// TODO: maybe if we do a JSONMarshal function we don't need to call these manually
