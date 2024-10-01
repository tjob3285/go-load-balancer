package main

import (
	"fmt"
	"math/rand"
	"net/http"
)

type RandomLB struct {
	servers []*Server
}

func NewRandomLB(servers []*Server) *RandomLB {
	return &RandomLB{
		servers: servers,
	}
}

// ServeHTTP distributes the incoming request to a random backend server.
func (lb *RandomLB) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server := lb.GetNextAvailableServer()
	if server != nil {
		proxy := NewReverseProxy(server.URL)
		fmt.Print("server is ", server.URL)
		proxy.ServerHttp(w, r)

	} else {
		//TODO:
		// Handle the case when no available server is found
	}
}

// Select two random servers and return the server with fewer connections
func (lb *RandomLB) GetNextAvailableServer() *Server {
	var availableServers []*Server

	for _, server := range lb.servers {
		server.mutex.Lock()

		if server.Alive {
			availableServers = append(availableServers, server)
		}
		server.mutex.Unlock()
	}

	if len(availableServers) == 0 {
		// No available servers found, return nil
		return nil
	}

	// return some random server from available servers
	rs, err := generateTwoRandomServers(availableServers)
	if err != nil {
		lb.GetNextAvailableServer()
	}

	return rs
}

func generateTwoRandomServers(servers []*Server) (*Server, error) {
	firstIndex := rand.Intn(len(servers))
	secondIndex := firstIndex

	for secondIndex == firstIndex {
		secondIndex = rand.Intn(len(servers))
	}

	s1 := servers[firstIndex]
	s2 := servers[secondIndex]

	s1.mutex.Lock()
	defer s1.mutex.Unlock()
	s2.mutex.Lock()
	defer s2.mutex.Unlock()

	// Compare URLs
	if s1.URL.String() == s2.URL.String() {
		return nil, fmt.Errorf("same urls") // URLs are the same
	}

	// Return the server with fewer connections
	if s1.Connections < s2.Connections {
		return s1, nil
	}
	return s2, nil
}
