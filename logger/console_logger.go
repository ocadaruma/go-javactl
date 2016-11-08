package logger

import (
	"fmt"
	"io"
	"time"
	"github.com/natefinch/lumberjack"
)

func NewConsoleLogger(path string, maxBytes int64, backupCount int) io.WriteCloser {
	mBytes := int(maxBytes / 1024 / 1024)

	writer := lumberjack.Logger{
		Filename: path,
		LocalTime: true,
	}
	if mBytes > 0 { writer.MaxSize = mBytes }
	if backupCount > 0 { writer.MaxBackups = backupCount }

	result := consoleLogger{ underlying: &writer }

	return &result
}

type consoleLogger struct {
	underlying *lumberjack.Logger
}

func (this *consoleLogger) Write(p []byte) (int, error) {
	formatted := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15-04-05,000"), string(p))

	return io.WriteString(this.underlying, formatted)
}

func (this *consoleLogger) Close() error {
	return this.underlying.Close()
}
