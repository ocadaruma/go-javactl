package logger

import (
	"fmt"
	"io"
	"log/syslog"
)

type systemLoggerStrategy struct {
	DryRun bool
}

func (this *systemLoggerStrategy) Log(level LogLevel, message string) {
	if this.DryRun {
		fmt.Printf("Would write to syslog: priority=%s, message=%s", level, message)
	} else {
		writer, err := syslog.New(syslog.Priority(level), "")
		if err != nil { return }

		io.WriteString(writer, message)
		writer.Close()
	}
}

func NewSystemLogger(dryRun bool) *Logger {
	return NewLogger(&systemLoggerStrategy{DryRun: dryRun})
}
