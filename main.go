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

func helloWorld(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	now := time.Now()
	var newRequests []time.Time
	for _, timestamp := range requests {
		if now.Sub(timestamp) <= interval {
			newRequests = append(newRequests, timestamp)
		}
	}

	requests = newRequests

	if len(requests) >= maxRequests {
		http.Error(w, "Please come later, let me take rest 😅🙏", http.StatusTooManyRequests)
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
