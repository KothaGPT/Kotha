package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	clients map[string][]time.Time
	limit   int
	window  time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string][]time.Time),
		limit:   limit,
		window:  window,
	}
	go rl.cleanup()
	return rl
}

func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for ip, timestamps := range rl.clients {
			var valid []time.Time
			for _, ts := range timestamps {
				if now.Sub(ts) < rl.window {
					valid = append(valid, ts)
				}
			}
			if len(valid) == 0 {
				delete(rl.clients, ip)
			} else {
				rl.clients[ip] = valid
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	var valid []time.Time
	for _, ts := range rl.clients[key] {
		if now.Sub(ts) < rl.window {
			valid = append(valid, ts)
		}
	}

	if len(valid) >= rl.limit {
		rl.clients[key] = valid
		return false
	}

	valid = append(valid, now)
	rl.clients[key] = valid
	return true
}

func (rl *RateLimiter) Middleware(keyFunc func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyFunc(r)
			if !rl.Allow(key) {
				w.Header().Set("Retry-After", "60")
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func CSRFMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get("X-CSRF-Token")
		cookie, err := r.Cookie("csrf_token")
		if err != nil || token == "" || token != cookie.Value {
			w.Header().Set("HX-Retarget", "#error-message")
			w.Header().Set("HX-Reswap", "innerHTML")
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprint(w, "Invalid CSRF token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func GetClientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.Split(fwd, ",")[0]
	}
	if fwd := r.Header.Get("X-Real-IP"); fwd != "" {
		return fwd
	}
	ip := r.RemoteAddr
	if idx := strings.LastIndex(ip, ":"); idx != -1 {
		return ip[:idx]
	}
	return ip
}

func (rl *RateLimiter) CleanupExpiredSessions(db interface{}) {
	_ = db
}
