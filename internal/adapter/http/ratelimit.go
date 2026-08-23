package httpadapter

import (
	"net"
	"net/http"
	"sync"
	"time"
)

type bucket struct {
	available float64
	updated   time.Time
}

type RateLimiter struct {
	mu      sync.Mutex
	buckets map[string]bucket
	rate    float64
	burst   float64
}

func NewRateLimiter(rate, burst int) *RateLimiter {
	if rate <= 0 {
		rate = 100
	}
	if burst <= 0 {
		burst = rate
	}
	return &RateLimiter{buckets: make(map[string]bucket), rate: float64(rate), burst: float64(burst)}
}

func (l *RateLimiter) Allow(key string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	value, exists := l.buckets[key]
	if !exists {
		l.buckets[key] = bucket{available: l.burst - 1, updated: now}
		return true
	}
	elapsed := now.Sub(value.updated).Seconds()
	value.available += elapsed * l.rate
	if value.available > l.burst {
		value.available = l.burst
	}
	value.updated = now
	if value.available < 1 {
		l.buckets[key] = value
		return false
	}
	value.available--
	l.buckets[key] = value
	return true
}

func (l *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			host = r.RemoteAddr
		}
		if !l.Allow(host, time.Now()) {
			writeError(w, APIError{Code: "rate_limited", Message: "request rate exceeded", Status: http.StatusTooManyRequests})
			return
		}
		next.ServeHTTP(w, r)
	})
}
