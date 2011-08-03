package logger

import (
	"fmt"
	"log"
)

type Severity byte

const (
	DEBUG Severity = iota
	INFO
	WARN
	ERROR
	FATAL
	UNKNOWN
)

func (severity Severity) String() string {
	switch severity {
	case DEBUG:
		return "D"
	case INFO:
		return "I"
	case WARN:
		return "W"
	case ERROR:
		return "E"
	case FATAL:
		return "F"
	case UNKNOWN:
		return "U"
	}
	return "U"
}

type Logger interface {
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warn(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
	Unknown(format string, v ...interface{})
}

type base struct {
	LogLevel Severity
	addFunc  func(severity Severity, format string, v ...interface{})
}

func (logger *base) Debug(format string, v ...interface{}) {
	logger.Add(DEBUG, format, v...)
}

func (logger *base) Info(format string, v ...interface{}) {
	logger.Add(INFO, format, v...)
}

func (logger *base) Warn(format string, v ...interface{}) {
	logger.Add(WARN, format, v...)
}

func (logger *base) Error(format string, v ...interface{}) {
	logger.Add(ERROR, format, v...)
}

func (logger *base) Fatal(format string, v ...interface{}) {
	logger.Add(FATAL, format, v...)
}

func (logger *base) Unknown(format string, v ...interface{}) {
	logger.Add(UNKNOWN, format, v...)
}

func (logger *base) Add(severity Severity, format string, v ...interface{}) {
	if logger.addFunc == nil {
		log.Panic("Tried to use base logger, which has no ability to output. Use descendants instead!")
	}

	if severity < logger.LogLevel {
		return
	}
	logger.addFunc(severity, format, v...)
}

/******************************************************************************/

type ConsoleLogger struct {
	*base
}

func NewConsoleLogger(logLevel Severity) *ConsoleLogger {
	logger := &ConsoleLogger{}
	logger.base = &base{}
	logger.base.LogLevel = logLevel
	logger.base.addFunc = func(severity Severity, format string, v ...interface{}) {
		format = fmt.Sprintf("%s %s", severity, format)
		log.Printf(format, v...)
	}
	return logger
}
