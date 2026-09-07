package metrics

import (
	"net/http"
	"sync"
	"time"
)

type RequestMetrics struct {
	mu             sync.RWMutex
	RequestCount   int64
	StatusCounts   map[int]int64
	Latencies      []time.Duration
	maxLatencies   int
}

func NewRequestMetrics(maxLatencies int) *RequestMetrics {
	if maxLatencies <= 0 {
		maxLatencies = 1000
	}
	return &RequestMetrics{
		StatusCounts: make(map[int]int64),
		Latencies:    make([]time.Duration, 0, maxLatencies),
		maxLatencies: maxLatencies,
	}
}

func (m *RequestMetrics) Record(status int, latency time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RequestCount++
	m.StatusCounts[status]++
	if len(m.Latencies) < m.maxLatencies {
		m.Latencies = append(m.Latencies, latency)
	}
}

type MetricsSummary struct {
	RequestCount int64
	StatusCounts map[int]int64
	P50          time.Duration
	P95          time.Duration
	P99          time.Duration
}

func (m *RequestMetrics) Summary() MetricsSummary {
	m.mu.RLock()
	defer m.mu.RUnlock()
	s := MetricsSummary{
		RequestCount: m.RequestCount,
		StatusCounts: make(map[int]int64),
	}
	for k, v := range m.StatusCounts {
		s.StatusCounts[k] = v
	}
	if len(m.Latencies) > 0 {
		s.P50 = percentile(m.Latencies, 50)
		s.P95 = percentile(m.Latencies, 95)
		s.P99 = percentile(m.Latencies, 99)
	}
	return s
}

func percentile(durations []time.Duration, p int) time.Duration {
	if len(durations) == 0 {
		return 0
	}
	sorted := make([]time.Duration, len(durations))
	copy(sorted, durations)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j] < sorted[i] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	idx := (p * len(sorted)) / 100
	if idx >= len(sorted) {
		idx = len(sorted) - 1
	}
	return sorted[idx]
}

var globalMetrics = NewRequestMetrics(1000)

func RecordRequest(status int, latency time.Duration) {
	globalMetrics.Record(status, latency)
}

func GetMetrics() MetricsSummary {
	return globalMetrics.Summary()
}

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(ww, r)
		latency := time.Since(start)
		RecordRequest(ww.status, latency)
	})
}

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (w *responseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}
