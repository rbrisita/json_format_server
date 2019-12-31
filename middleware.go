package main

import (
	"flag"
	"fmt"
	"net/http"
)

var (
	clients map[string]uint

	max = flag.Uint("max", 10, "Maximum connections per client per second.")
)

func init() {
	flag.Parse()
	clients = make(map[string]uint)
}

func rateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Connecting: ", r.RemoteAddr)

		if connectionAllowed(r.RemoteAddr) == 0 {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})

}

func connectionAllowed(ip string) uint {
	if _, ok := clients[ip]; !ok {
		clients[ip] = 0
	}
	clients[ip]++

	total := clients[ip]
	if total > *max {
		clients[ip]--
		return 0
	}

	return total
}
