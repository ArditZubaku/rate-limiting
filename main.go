package main

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/ArditZubaku/rate-limiting/limiter"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		time.Sleep(3 * time.Second)
		_, _ = w.Write([]byte("hi there\n"))
	})

	// handler := rateLimiterMiddleware(mux, limiter.NewSlidingWindow(1, 2)) // 2 reqs per second

	handler := rateLimiterMiddleware(
		mux,
		limiter.NewTokenBucket(1, 10), // 2.3 tokens per second are regenerated - burst of 10
		false,
	)

	err := http.ListenAndServe("127.0.0.1:8080", handler)
	if err != nil {
		log.Fatal(err)
	}
}

func rateLimiterMiddleware(
	next http.Handler,
	rateLimiter limiter.RateLimiter,
	enable bool,
) http.Handler {
	// For each ip/client we will a rate limiter
	ipToLimiterMap := sync.Map{}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !enable {
			next.ServeHTTP(w, r)
			return
		}

		// Get the IP of the client
		ip := getClientIP(r)

		// Rate limit
		ipLimiter, _ := ipToLimiterMap.LoadOrStore(ip, rateLimiter)

		if !ipLimiter.(limiter.RateLimiter).Allow() {
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println("Failed to get client IP:", err)
	}

	return host
}
