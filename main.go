package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var (
	maxRequests = 2
	interval    = 5 * time.Second
	requests    = []time.Time{}
	mu          sync.Mutex
)

type ResponseMessage string

const (
	RateLimitExceeded ResponseMessage = "Please come later, let me take rest 😅🙏"
)

func helloWorld(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	date := now.Format("02/01/2006")
	timestamp := now.Format("15:04:05")

	fmt.Printf("\nDate: %s, Timestamp: %s\n", date, timestamp)

	var newRequests []time.Time

	for _, timestamp := range requests {
		if now.Sub(timestamp) <= interval {
			newRequests = append(newRequests, timestamp)
		}
	}

	requests = newRequests

	if len(requests) >= maxRequests {
		fmt.Printf("Date: %s, Timestamp: %s, Response: %s\n", date, timestamp, RateLimitExceeded)
		http.Error(w, string(RateLimitExceeded), http.StatusTooManyRequests)
		return
	}

	requests = append(requests, now)

	fmt.Fprintln(w, "Hello World!")
}

func main() {
	http.HandleFunc("/hello", helloWorld)
	fmt.Println("Server running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
