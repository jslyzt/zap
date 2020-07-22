package zapcore

import (
	"testing"

	"github.com/jslyzt/zap/internal/ztest"
)

func BenchmarkMultiWriteSyncer(b *testing.B) {
	b.Run("2", func(b *testing.B) {
		w := NewMultiWriteSyncer(
			&ztest.Discarder{},
			&ztest.Discarder{},
		)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.Write([]byte("foobarbazbabble"))
			}
		})
	})
	b.Run("4", func(b *testing.B) {
		w := NewMultiWriteSyncer(
			&ztest.Discarder{},
			&ztest.Discarder{},
			&ztest.Discarder{},
			&ztest.Discarder{},
		)
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				w.Write([]byte("foobarbazbabble"))
			}
		})
	})
}
