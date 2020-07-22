package zap

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTakeStacktrace(t *testing.T) {
	trace := takeStacktrace()
	lines := strings.Split(trace, "\n")
	require.True(t, len(lines) > 0, "Expected stacktrace to have at least one frame.")
	assert.Contains(
		t,
		lines[0],
		"testing.",
		"Expected stacktrace to start with the test runner (zap frames are filtered out) %s.", lines[0],
	)
}

func TestIsZapFrame(t *testing.T) {
	zapFrames := []string{
		"github.com/jslyzt/zap.Stack",
		"github.com/jslyzt/zap.(*SugaredLogger).log",
		"github.com/jslyzt/zap/zapcore.(ArrayMarshalerFunc).MarshalLogArray",
		"github.com/uber/tchannel-go/vendor/github.com/jslyzt/zap.Stack",
		"github.com/uber/tchannel-go/vendor/github.com/jslyzt/zap.(*SugaredLogger).log",
		"github.com/uber/tchannel-go/vendor/github.com/jslyzt/zap/zapcore.(ArrayMarshalerFunc).MarshalLogArray",
	}
	nonZapFrames := []string{
		"github.com/uber/tchannel-go.NewChannel",
		"go.uber.org/not-zap.New",
		"github.com/jslyzt/zapext.ctx",
		"github.com/jslyzt/zap_ext/ctx.New",
	}

	t.Run("zap frames", func(t *testing.T) {
		for _, f := range zapFrames {
			require.True(t, isZapFrame(f), f)
		}
	})
	t.Run("non-zap frames", func(t *testing.T) {
		for _, f := range nonZapFrames {
			require.False(t, isZapFrame(f), f)
		}
	})
}

func BenchmarkTakeStacktrace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		takeStacktrace()
	}
}
