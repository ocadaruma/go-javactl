package logger

import (
	"io"

	"github.com/natefinch/lumberjack"
)

func NewConsoleLogger(path string, maxBytes int64, backupCount int) io.Writer {
	writer := lumberjack.Logger{
		Filename: path,
		MaxSize: int(maxBytes / 1024 / 1024),
		MaxBackups: backupCount,
		LocalTime: true,
	}

	strategy := consoleStrategy{ underlying: &writer }

	result := consoleLogger{ underlying: &strategy }

	return &result
}

type consoleStrategy struct {
	underlying *lumberjack.Logger
	writeCallback func(int, error)
}

func (this *consoleStrategy) Log(level LogLevel, message string) {
	n, err := this.underlying.Write([]byte(message))
	this.writeCallback(n, err)
}

type consoleLogger struct {
	underlying *consoleStrategy
}

func (this *consoleLogger) Write(p []byte) (n int, err error) {
	this.underlying.writeCallback = func(nCallback int, errCallback error) {
		n = nCallback
		err = errCallback
	}
	this.underlying.Log(LOG_INFO, string(p))
	return
}
