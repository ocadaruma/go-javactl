package logger

import (
	"fmt"
	"log/syslog"
)

type LogLevel syslog.Priority
const (
	LOG_INFO LogLevel = (LogLevel)(syslog.LOG_INFO)
	LOG_WARNING LogLevel = (LogLevel)(syslog.LOG_WARNING)
	LOG_ERR LogLevel = (LogLevel)(syslog.LOG_ERR)
)

var levelPrefix = map[LogLevel]string {
	LOG_INFO: "[INFO]",
	LOG_WARNING: "[WARN]",
	LOG_ERR: "[ERROR]",
}

func (this LogLevel) String() string {
	return levelPrefix[this]
}

type LogStrategy interface {
	Log(level LogLevel, message string)
}

type Logger struct {
	strategy LogStrategy
}

func NewLogger(strategy LogStrategy) Logger {
	l := Logger{ strategy: strategy }
	return l
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

	this.strategy.Log(level, msg)
}
