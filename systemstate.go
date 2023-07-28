package main

// This could be a lot better by having per-client reading/writing

import (
	"sync"
	"time"
)

type service struct {
	name           string
	status         bool
	lastConnection time.Time
}

type client struct {
	name           string
	ipaddress      string
	lastConnection time.Time
	services       []service
}

type systemState struct {
	clients map[string]*client
	mutex   sync.RWMutex
}

func NewSystemState() *systemState {
	ss := &systemState{
		clients: make(map[string]*client),
	}
	return ss
}

// Here's a problem, I don't want anyone to have direct access to any member
// It's all mutexed. So assignments generally copy, but not always. Gotta remember to make sure it copies. or "deep copies"
func (ss *systemState) UpdateClient(uC client) { // this isn't good enough We need to unmarshall json
	ss.mutex.Lock()
	defer ss.mutex.Unlock()
	ss.clients[uC.name] = &uC
	copy(uC.services, ss.clients[uC.name].services)
	ss.clients[uC.name].lastConnection = time.Now()
}

func (ss *systemState) GetClientsCopy() []client {
	ss.mutex.RLock()
	defer ss.mutex.RUnlock()
	clients := make([]client, len(ss.clients))
	i := 0
	for _, v := range ss.clients {
		clients[i] = *v
		copy(v.services, clients[i].services)
	}
	return clients
}

type clientList struct {
	names map[string]bool
	mutex sync.RWMutex
} // TODO POPULATE THIS

func NewClientList(names []string) *clientList {
	cL := make(map[string]bool)
	for _, s := range names {
		cL[s] = true
	}
	return &clientList{names: cL}
}

func (cl *clientList) UpdateClientList(names []string) {
	newList := make(map[string]bool)
	for _, s := range names {
		newList[s] = true
	}
	cl.mutex.Lock()
	defer cl.mutex.Unlock()
	cl.names = newList

}

func (cl *clientList) CheckClient(name string) bool {
	cl.mutex.RLock()
	defer cl.mutex.RUnlock()
	_, ok := cl.names[name]
	return ok
}

func (cl *clientList) GetAllClients() []string {
	cl.mutex.RLock()
	allNames := make([]string, len(cl.names))
	i := 0
	defer cl.mutex.RUnlock()
	for k, _ := range cl.names {
		allNames[i] = k
		i += 1
	}
	return allNames
}

// OKAY SO
// This guy will have a list of relevant clients, that comes from reading public keys
// He can compare it against the list of connected keys
// This is starting to look a lot like RPC, since they have to communicate with the same structure
// But they have to do it over an existing websockets connection? No, they don't.
