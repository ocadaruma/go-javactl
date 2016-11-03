package logger

import (
	"testing"
)

type logEntry struct {
	level LogLevel
	message string
}

type mockStrategy struct {
	buffer []logEntry
}

func (this *mockStrategy) Log(level LogLevel, message string) {
	this.buffer = append(this.buffer, logEntry{level, message})
}

func TestLogger(t *testing.T) {
	strategy := mockStrategy{}
	l := NewLogger(&strategy)

	l.Info("Info message.")
	l.Warn("Warn message.")
	l.Error("Error message.")

	testCases := []struct {
		actual, expect logEntry
	} {
		{strategy.buffer[0], logEntry{LOG_INFO, "[INFO] Info message."}},
		{strategy.buffer[1], logEntry{LOG_WARNING, "[WARN] Warn message."}},
		{strategy.buffer[2], logEntry{LOG_ERR, "[ERROR] Error message."}},
	}

	for i, c := range testCases {
		if c.actual != c.expect {
			t.Errorf("case %v : %v must equal to %v", i, c.actual, c.expect)
		}
	}
}
