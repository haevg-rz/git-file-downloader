package log

import (
	"log"
)

var (

	//	Level:
	//
	//	0 = No logging
	//
	//	1 = Minimal
	//
	//	2 = Logging of major steps
	//
	//	3 = Everything
	Level = 0

	// Shared logger instance used throughout the project. Acts dependent on the global log level.
	logger = NewLogger()
)

type ILogger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

type Logger struct {
	level int
}

func NewLogger() *Logger {
	return &Logger{level: 0}
}

func V(level int) *Logger {
	logger.level = level
	return logger
}

func (l *Logger) Println(v ...interface{}) {
	if l.level <= Level {
		log.Println(v...)
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.level <= Level {
		log.Printf(format, v...)
	}
}
