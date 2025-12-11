package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/ArditZubaku/rate-limiting/limiter"
)

func main() {
	url := "http://127.0.0.1:8080/"

	rateLimiter := limiter.NewTokenBucket(1, 10)

	var wg sync.WaitGroup

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	timeout := time.After(5 * time.Second)

	for {
		select {
		case <-timeout:
			log.Println("5 seconds elapsed")
			return
		case <-ticker.C:
			for i := range 50 {
				if !rateLimiter.Allow() {
					fmt.Printf(
						"Request %2d: Status: %s Time: %d\n",
						i+1, "skipped", time.Now().Second(),
					)
					continue
				}

				wg.Add(1)
				go func(i int) {
					defer wg.Done()
					makeRequestAndReport(url, i)
				}(i)
			}

			wg.Wait()
		}
	}
}

func makeRequestAndReport(url string, i int) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf(
		"Request %2d: Status: %d Time: %d\n",
		i+1, resp.StatusCode, time.Now().Second(),
	)
	time.Sleep(100 * time.Millisecond)
}
