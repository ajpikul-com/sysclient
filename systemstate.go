package main

// This could be a lot better by having per-client reading/writing

import (
	"sync"
	"time"
)

type Service struct {
	Name           string
	IPAddress      string
	Status         bool
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

func (ss *systemState) GetServices() []Service {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	services := make([]Service, len(ss.Services))
	i := 0
	for _, v := range ss.Services {
		services[i] = *v
	}
	return services
}

func (ss *systemState) UpdateClientList(names []string) {
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	copy(ss.ExpectedClients, names)

}

func (ss *systemState) GetAllClients() []string {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	clients := make([]string, len(ss.ExpectedClients))
	copy(clients, ss.ExpectedClients)
	return clients
}
