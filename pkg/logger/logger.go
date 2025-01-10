package logger

import (
	"fmt"
	"time"
)

type Logger interface {
	Info(format string, args ...any)
	Error(format string, args ...any)
	Debug(format string, args ...any)
	Warn(format string, args ...any)
}

type DefaultLogger struct {
	level LogLevel
}

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

func NewDefaultLogger(level LogLevel) *DefaultLogger {
	return &DefaultLogger{
		level: level,
	}
}

const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorReset  = "\033[0m"
)

func (l *DefaultLogger) Info(format string, args ...any) {
	if l.level <= INFO {
		l.log("INFO", colorGreen, format, args...)
	}
}

func (l *DefaultLogger) Warn(format string, args ...any) {
	if l.level <= WARN {
		l.log("WARN", colorYellow, format, args...)
	}
}

func (l *DefaultLogger) Debug(format string, args ...any) {
	if l.level <= DEBUG {
		l.log("DEBUG", colorBlue, format, args...)
	}
}

func (l *DefaultLogger) Error(format string, args ...any) {
	if l.level <= ERROR {
		l.log("ERROR", colorRed, format, args...)
	}
}

func (l *DefaultLogger) log(level, color, format string, args ...any) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var message string
	if len(args) > 0 {
		message = fmt.Sprintf(format, args...)
	} else {
		message = format // handle literal string
	}

	fmt.Printf("%s[%s] %s: %s%s\n", color, timestamp, level, message, colorReset)
}
