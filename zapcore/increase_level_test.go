package zapcore_test

import (
	"fmt"
	"testing"

	"github.com/jslyzt/zap"
	. "github.com/jslyzt/zap/zapcore"
	"github.com/jslyzt/zap/zaptest/observer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIncreaseLevel(t *testing.T) {
	tests := []struct {
		coreLevel     Level
		increaseLevel Level
		wantErr       bool
		with          []Field
	}{
		{
			coreLevel:     InfoLevel,
			increaseLevel: DebugLevel,
			wantErr:       true,
		},
		{
			coreLevel:     InfoLevel,
			increaseLevel: InfoLevel,
		},
		{
			coreLevel:     InfoLevel,
			increaseLevel: ErrorLevel,
		},
		{
			coreLevel:     InfoLevel,
			increaseLevel: ErrorLevel,
			with:          []Field{zap.String("k", "v")},
		},
		{
			coreLevel:     ErrorLevel,
			increaseLevel: DebugLevel,
			wantErr:       true,
		},
		{
			coreLevel:     ErrorLevel,
			increaseLevel: InfoLevel,
			wantErr:       true,
		},
		{
			coreLevel:     ErrorLevel,
			increaseLevel: WarnLevel,
			wantErr:       true,
		},
		{
			coreLevel:     ErrorLevel,
			increaseLevel: PanicLevel,
		},
	}

	for _, tt := range tests {
		msg := fmt.Sprintf("increase %v to %v", tt.coreLevel, tt.increaseLevel)
		t.Run(msg, func(t *testing.T) {
			logger, logs := observer.New(tt.coreLevel)

			filteredLogger, err := NewIncreaseLevelCore(logger, tt.increaseLevel)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "invalid increase level")
				return
			}

			if len(tt.with) > 0 {
				filteredLogger = filteredLogger.With(tt.with)
			}

			require.NoError(t, err)

			for l := DebugLevel; l <= FatalLevel; l++ {
				enabled := filteredLogger.Enabled(l)
				entry := Entry{Level: l}
				ce := filteredLogger.Check(entry, nil)
				ce.Write()
				entries := logs.TakeAll()

				if l >= tt.increaseLevel {
					assert.True(t, enabled, "expect %v to be enabled", l)
					assert.NotNil(t, ce, "expect non-nil Check")
					assert.NotEmpty(t, entries, "Expect log to be written")
				} else {
					assert.False(t, enabled, "expect %v to be disabled", l)
					assert.Nil(t, ce, "expect nil Check")
					assert.Empty(t, entries, "No logs should have been written")
				}

				// Write should always log the entry as per the Core interface
				require.NoError(t, filteredLogger.Write(entry, nil), "Write failed")
				require.NoError(t, filteredLogger.Sync(), "Sync failed")
				assert.NotEmpty(t, logs.TakeAll(), "Write should always log")
			}
		})
	}
}
