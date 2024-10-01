package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// layer 4
// tcp/udp
// focuses on distributing traffic based on IP addresses, ports, and transport protocols. Layer 4 load balancers make routing decisions based on network-level information without inspecting application-layer protocols.

// layer 7

// LoadBalancer defines the interface for a load balancer.
// LoadBalancer defines the interface for a load balancer.
type LoadBalancer interface {
	ServeHttp(w http.ResponseWriter, r *http.Request)
	GetNextAvailableServer() *Server
}

type ReverseProxy struct {
	backendURL *url.URL
	proxy      *httputil.ReverseProxy
}

func NewReverseProxy(backendURL *url.URL) *ReverseProxy {

	return &ReverseProxy{
		backendURL: backendURL,
		proxy:      httputil.NewSingleHostReverseProxy(backendURL),
	}

}

// Forwards the incoming request to backend server
func (rp *ReverseProxy) ServerHttp(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Forwarding request to %s : %s\n", rp.backendURL, r.URL.Path)
	rp.proxy.ServeHTTP(w, r)
}
