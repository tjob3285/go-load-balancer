package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/tjob3285/go-load-balancer/config"
)

// Server represents a backend server.
type Server struct {
	URL         *url.URL
	Alive       bool
	Connections int
	mutex       sync.Mutex // using it to protect concurrent access to alive and connections field
}

// health check function that runs in given interval to check health of servers
func healthCheck(s *Server, healthCheckInterval time.Duration) {
	for range time.Tick(healthCheckInterval) {
		res, err := http.Head(s.URL.String())
		s.mutex.Lock()
		if err != nil || res.StatusCode != http.StatusOK {
			fmt.Printf("%s is down, alert someone!\n", s.URL)
			s.Alive = false
		} else {
			fmt.Printf("%s is alive and running\n", s.URL)
			s.Alive = true
		}
		s.mutex.Unlock()
	}
}

func main() {
	// config file for servers and backend server
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	healthCheckInterval, err := time.ParseDuration(config.HealthInterval)
	if err != nil {
		log.Fatalf("Invalid health check interval: %s", err.Error())
	}

	// map config servers to Server type and run health checks
	var servers []*Server
	for _, serverUrl := range config.URLs {
		u, _ := url.Parse(serverUrl)

		server := &Server{URL: u, Alive: true}
		servers = append(servers, server)
		go healthCheck(server, healthCheckInterval)
	}

	// switch lb based on config algorithm
	var loadBalancer http.Handler
	switch config.Algorithm {
	case "round-robin":
		loadBalancer = NewRoundRobinLB(servers)
	case "least-connection":
		loadBalancer = NewLeastConnectionLB(servers)
	case "rdm":
		loadBalancer = NewRandomLB(servers)
	default:
		log.Fatalf("Invalid algorithm type")
	}

	// Register the load balancers as HTTP handlers
	http.Handle("/", loadBalancer)

	// Start the server
	fmt.Println("Load balancers started.")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}
}
