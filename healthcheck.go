package main

import (
	"fmt"
	"net/http"
	"time"
)

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
