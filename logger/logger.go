package logger

import (
	"fmt"
)

type LogLevel int
const (
	LOG_INFO LogLevel = 6
	LOG_WARNING LogLevel = 4
	LOG_ERR LogLevel = 3
)

var levelPrefix = map[LogLevel]string {
	LOG_INFO: "[INFO]",
	LOG_WARNING: "[WARN]",
	LOG_ERR: "[ERROR]",
}

type Logger interface {
	Log(level LogLevel, message string)
}

func (this Logger) Info(message string) {
	this.logMessage(LOG_INFO, message)
}

func (this Logger) Warn(message string) {
	this.logMessage(LOG_WARNING, message)
}

func (this Logger) Error(message string) {
	this.logMessage(LOG_ERR, message)
}

func (this Logger) logMessage(level LogLevel, message string) {
	prefix := levelPrefix[level]

	msg := fmt.Sprintf("%s%s", prefix, message)

	this.Log(level, msg)
}
