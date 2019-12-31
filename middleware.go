package main

import (
	"flag"
	"net/http"
	"sync"

	"golang.org/x/time/rate"
)

/**
* Clients is a shared data structure to keep track of clients connections by ip.
**/
type Clients struct {
	mux *sync.RWMutex
	ips map[string]*rate.Limiter
}

/**
* Flags for command line configuration.
**/
var (
	max     = flag.Uint("max", 10, "Maximum connections per client per second.")
	burst   = flag.Int("burst", 1, "Maximum burst size events.")
	clients *Clients
)

/**
* init by parsing flags and creating shared data structure.
**/
func init() {
	flag.Parse()
	clients = &Clients{
		mux: &sync.RWMutex{},
		ips: make(map[string]*rate.Limiter),
	}
}

/**
* rateLimit uses time/rate to test if a connection from given ip is allowed.
**/
func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		if !connectionAllowed(ip) {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})

}

/**
* connectionAllowed is a thread safe function that uses
* a shared data structure to access a rate limiter per ip.
**/
func connectionAllowed(ip string) bool {
	clients.mux.Lock()
	defer clients.mux.Unlock()

	if _, ok := clients.ips[ip]; !ok {
		clients.ips[ip] = rate.NewLimiter(rate.Limit(*max), *burst)
	}

	return clients.ips[ip].Allow()
}
