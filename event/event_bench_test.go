package event

import (
	"context"
	"testing"
)

func BenchmarkEventEmit(b *testing.B) {
	Subscribe("bench.topic", func(_ context.Context, _ any) {})
	defer Unsubscribe(stream.subs["bench.topic"][0])

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Emit("bench.topic", i)
	}
}

func BenchmarkEventEmitConcurrent(b *testing.B) {
	Subscribe("bench.topic2", func(_ context.Context, _ any) {})

	b.ReportAllocs()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Emit("bench.topic2", 1)
		}
	})
}
