package middleware

import (
	"context"
	"net/http"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"sync/atomic"
)

type (
	RequestKey         struct{}
	ResponseHeadersKey struct{}
)

var requestPrefix [8]byte

func init() {
	if _, err := rand.Read(requestPrefix[:]); err != nil {
		requestPrefix = [8]byte{'s', 'u', 'p', 'e', 'r', 'k', 'i', 't'}
	}
}

var requestCounter atomic.Uint64

func WithRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), RequestKey{}, r)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = generateRequestID()
		}
		w.Header().Set("X-Request-ID", id)
		ctx := context.WithValue(r.Context(), RequestKey{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() string {
	var b [16]byte
	copy(b[:8], requestPrefix[:])
	binary.BigEndian.PutUint64(b[8:], requestCounter.Add(1))
	return hex.EncodeToString(b[:])
}
