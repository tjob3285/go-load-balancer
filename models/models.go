package models

import (
	"net/url"
	"sync"
)

type Server struct {
	URL       *url.URL
	IsHealthy bool
}

type LoadBalancer struct {
	Current int
	Mutex   sync.Mutex
}
