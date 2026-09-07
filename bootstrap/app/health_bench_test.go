package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khulnasoft/superkit/kit"
)

func BenchmarkHandleHealth(b *testing.B) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w.Result().Body.Close()
		w = httptest.NewRecorder()
		k := &kit.Kit{
			Response: w,
			Request:  req,
		}
		HandleHealth(k)
	}
}
