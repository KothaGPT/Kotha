package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func BenchmarkRateLimiterAllow(b *testing.B) {
	rl := NewRateLimiter(1000, time.Minute)
	keys := []string{"ip-1", "ip-2", "ip-3", "ip-4", "ip-5"}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rl.Allow(keys[i%len(keys)])
	}
}

func BenchmarkRateLimiterConcurrentAllow(b *testing.B) {
	rl := NewRateLimiter(1000, time.Minute)
	keys := []string{"ip-1", "ip-2", "ip-3", "ip-4", "ip-5",
		"ip-6", "ip-7", "ip-8", "ip-9", "ip-10"}

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			rl.Allow(keys[i%len(keys)])
			i++
		}
	})
}

func BenchmarkRateLimiterMiddleware(b *testing.B) {
	rl := NewRateLimiter(1000, time.Minute)
	keyFunc := func(r *http.Request) string {
		return "test-client"
	}
	mw := rl.Middleware(keyFunc)
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := mw(next)

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}
}
