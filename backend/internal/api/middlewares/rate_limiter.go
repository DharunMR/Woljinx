package middlewares

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu        sync.Mutex
	visitor   map[string]int
	limit     int
	resetTime time.Duration
}

func (rl *rateLimiter) resetVisitorCount() {
	for {
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitor = make(map[string]int)
		rl.mu.Unlock()
	}
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitor:   make(map[string]int),
		limit:     limit,
		resetTime: resetTime,
	}

	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		rl.mu.Lock()
		defer rl.mu.Unlock()

		visitorIP := r.RemoteAddr
		rl.visitor[visitorIP]++

		if rl.limit < rl.visitor[visitorIP] {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
