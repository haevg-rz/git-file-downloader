package log

import (
	"fmt"
	goLog "log"
	"os"
	"sync"
	"time"

	"github.com/haevg-rz/git-file-downloader/pkg/exit"
)

const (
	FilenameFormat = "2006-01-02-15-04-05"
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
	logger = NewLoggerHandler()

	FileLogWg sync.WaitGroup
)

// Logger describes the default logging behaviour
type Logger interface {
	Println(v ...interface{})
	Printf(format string, v ...interface{})
}

// LoggerHandler is a specific implementation of Logger. Defines a log-level which controls whether a log will be printed.
type LoggerHandler struct {
	level int
}

// NewLoggerHandler creates a new instance of LoggerHandler
func NewLoggerHandler() *LoggerHandler {
	return &LoggerHandler{}
}

// InitFileLog initiates logging to file. Spawns a goroutine which closes the fileHandle once no longer needed.
func InitFileLog(outputPath string, logLevel int, doneCh chan bool) error {
	var err error

	defer func() {
		if err != nil {
			exit.Code = exit.InternalError
		}
	}()

	if _, err = os.Stat(outputPath); err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir(outputPath, os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s-log.log", outputPath, time.Now().Format(FilenameFormat)), os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	goLog.SetOutput(logFile)

	FileLogWg.Add(1)
	go func() {
		V(3).Printf("logging to file %s with v=%d\n", logFile.Name(), logLevel)
		<-doneCh
		V(3).Printf("closing writer on logfile %s\n", logFile.Name())

		err = logFile.Close()
		if err != nil {
			exit.Code = exit.InternalError
		}
		FileLogWg.Done()
	}()

	return nil
}

// V returns a local shared logger-instance with given level.
func V(level int) *LoggerHandler {
	logger.level = level
	return logger
}

// Println prints the given arguments.
func (l *LoggerHandler) Println(v ...interface{}) {
	if l.level <= Level {
		goLog.Println(v...)
	}
}

// Printf prints the format string with the given arguments.
func (l *LoggerHandler) Printf(format string, v ...interface{}) {
	if l.level <= Level {
		goLog.Printf(format, v...)
	}
}
