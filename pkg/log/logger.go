package log

import (
	"fmt"
	"github.com/haevg-rz/git-file-downloader/pkg/exit"
	goLog "log"
	"os"
	"sync"
	"time"
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
	logger = NewLogger()

	GracefulShutdown sync.WaitGroup
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

	logFile, err := os.OpenFile(fmt.Sprintf("%s/%s-log.txt", outputPath, time.Now().Format(FilenameFormat)), os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	goLog.SetOutput(logFile)

	GracefulShutdown.Add(1)
	go func() {
		V(3).Printf("logging to file %s with v=%d\n", logFile.Name(), logLevel)
		<-doneCh
		V(3).Printf("closing writer on logfile %s\n", logFile.Name())

		err = logFile.Close()
		if err != nil {
			exit.Code = exit.InternalError
		}
		GracefulShutdown.Done()
	}()

	return nil
}

func V(level int) *Logger {
	logger.level = level
	return logger
}

func (l *Logger) Println(v ...interface{}) {
	if l.level <= Level {
		goLog.Println(v...)
	}
}

func (l *Logger) Printf(format string, v ...interface{}) {
	if l.level <= Level {
		goLog.Printf(format, v...)
	}
}
