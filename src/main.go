package main

import (
	"net/http/httputil"
	"net/http"
	"net/url"
	"os"
	"fmt"
)

type Server interface {
	Address() string
	isAlive() bool
	Serve(rw http.ResponseWriter, r *http.Request)
}

type simpleServer struct {
	addr string
	proxy httputil.ReverseProxy
}

func newSimpleServer(addr string) *simpleServer {
	serverUrl, err := url.Parse(addr)
	handleErr(err)
	return &simpleServer{
		addr: addr,
		proxy: *httputil.NewSingleHostReverseProxy(serverUrl),
	}
}

type loadBalancer struct {
	port string
	roundRobinIndex int
	servers []Server
}

func newLoadBalancer(port string, servers []Server) *loadBalancer {
	return &loadBalancer{
		port: port,
		roundRobinIndex: 0,
		servers: servers,
	}
}

func handleErr(err error){
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func (s *simpleServer) Address() string {return s.addr}

func (s *simpleServer)isAlive() bool {return true}

func (s *simpleServer) Serve(rw http.ResponseWriter, req *http.Request) {
	s.proxy.ServeHTTP(rw, req)
}

func (lb *loadBalancer) getNextAvailableServer() Server {
    totalServers := len(lb.servers)
    for i := 0; i < totalServers; i++ {
        server := lb.servers[lb.roundRobinIndex%totalServers]
        lb.roundRobinIndex++
        if server.isAlive() {
            return server
        }
    }
    return nil // or handle case when no servers are alive
}

func (lb *loadBalancer) serveProxy(rw http.ResponseWriter, req *http.Request) {
    targetServer := lb.getNextAvailableServer()
    if targetServer == nil {
        http.Error(rw, "Service Unavailable", http.StatusServiceUnavailable)
        return
    }
    fmt.Printf("Redirecting request to %s\n", targetServer.Address())
    targetServer.Serve(rw, req)
}

func main() {
	servers := []Server{
		newSimpleServer("https://www.google.com"),
		newSimpleServer("https://www.wikipedia.org"),
		newSimpleServer("https://www.github.com"),
	}
	lb := newLoadBalancer(":8000", servers)
	handleRedirect := func(rw http.ResponseWriter, req *http.Request){
		lb.serveProxy(rw, req)
	}
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Load balancer is running on port %s\n", lb.port)
	handleErr(http.ListenAndServe(lb.port, nil))

}