package zap

import (
	"bytes"
	"testing"

	"github.com/jslyzt/zap/zapcore"
	"github.com/jslyzt/zap/zaptest/observer"
	"github.com/stretchr/testify/assert"
)

func newLoggedEntry(level zapcore.Level, msg string, fields ...zapcore.Field) observer.LoggedEntry {
	if len(fields) == 0 {
		fields = []zapcore.Field{}
	}
	return observer.LoggedEntry{
		Entry:   zapcore.Entry{Level: level, Message: msg},
		Context: fields,
	}
}

func TestIncreaseLevelTryDecrease(t *testing.T) {
	errorOut := &bytes.Buffer{}
	opts := []Option{
		ErrorOutput(zapcore.AddSync(errorOut)),
	}
	withLogger(t, WarnLevel, opts, func(logger *Logger, logs *observer.ObservedLogs) {
		logger.Warn("original warn log")

		debugLogger := logger.WithOptions(IncreaseLevel(DebugLevel))
		debugLogger.Debug("ignored debug log")
		debugLogger.Warn("increase level warn log")
		debugLogger.Error("increase level error log")

		assert.Equal(t, []observer.LoggedEntry{
			newLoggedEntry(WarnLevel, "original warn log"),
			newLoggedEntry(WarnLevel, "increase level warn log"),
			newLoggedEntry(ErrorLevel, "increase level error log"),
		}, logs.AllUntimed(), "unexpected logs")
		assert.Equal(t,
			"failed to IncreaseLevel: invalid increase level, as level \"info\" is allowed by increased level, but not by existing core\n",
			errorOut.String(),
			"unexpected error output",
		)
	})
}

func TestIncreaseLevel(t *testing.T) {
	errorOut := &bytes.Buffer{}
	opts := []Option{
		ErrorOutput(zapcore.AddSync(errorOut)),
	}
	withLogger(t, WarnLevel, opts, func(logger *Logger, logs *observer.ObservedLogs) {
		logger.Warn("original warn log")

		errorLogger := logger.WithOptions(IncreaseLevel(ErrorLevel))
		errorLogger.Debug("ignored debug log")
		errorLogger.Warn("ignored warn log")
		errorLogger.Error("increase level error log")

		withFields := errorLogger.With(String("k", "v"))
		withFields.Debug("ignored debug log with fields")
		withFields.Warn("ignored warn log with fields")
		withFields.Error("increase level error log with fields")

		assert.Equal(t, []observer.LoggedEntry{
			newLoggedEntry(WarnLevel, "original warn log"),
			newLoggedEntry(ErrorLevel, "increase level error log"),
			newLoggedEntry(ErrorLevel, "increase level error log with fields", String("k", "v")),
		}, logs.AllUntimed(), "unexpected logs")

		assert.Empty(t, errorOut.String(), "expect no error output")
	})
}
