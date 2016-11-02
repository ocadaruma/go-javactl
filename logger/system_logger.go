package logger

import (
	"fmt"
	"log/syslog"
)

type systemLoggerStrategy struct {
	DryRun bool
}

func (this systemLoggerStrategy) Log(level LogLevel, message string) {
	if this.DryRun {
		fmt.Printf("Would write to syslog: priority=%s, message=%s", level, message)
	} else {
		writer, err := syslog.New(level, "")
		if err != nil { return }

		writer.Write([]byte(message))
		writer.Close()
	}
}

func NewSystemLogger(dryRun bool) Logger {
	return NewLogger(systemLoggerStrategy{DryRun: dryRun})
}
