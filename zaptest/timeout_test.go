package zaptest

import (
	"testing"
	"time"

	"github.com/jslyzt/zap/internal/ztest"
	"github.com/stretchr/testify/assert"
)

func TestTimeout(t *testing.T) {
	defer ztest.Initialize("2")()
	assert.Equal(t, time.Duration(100), Timeout(50), "Expected to scale up timeout.")
}

func TestSleep(t *testing.T) {
	defer ztest.Initialize("2")()
	const sleepFor = 50 * time.Millisecond
	now := time.Now()
	Sleep(sleepFor)
	elapsed := time.Since(now)
	assert.True(t, 2*sleepFor <= elapsed, "Expected to scale up timeout.")
}
