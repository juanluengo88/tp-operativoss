package log

import (
	"fmt"
	"io"
	"log"
	"os"
)

type LogLevel string

const (
	DEBUG LogLevel = "DEBUG"
	INFO  LogLevel = "INFO"
	ERROR LogLevel = "ERROR"
)

type LoggerStruct struct {
	Logger      *log.Logger
	LogFile     *os.File
	environment string
}

// Log formats and logs a message at the given severity level.
func (l *LoggerStruct) Log(message string, level LogLevel) {
	if l.environment == "prod" && level == "DEBUG" {
		return
	}

	format := fmt.Sprintf("[%s]: %s", level, message)
	l.Logger.Println(format)
}

func (l *LoggerStruct) CloseLogger() {
	l.LogFile.Close()
}

// Initializes a Logger to write to both stdout and a specified file.
func ConfigureLogger(filename, env string) *LoggerStruct {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logger := LoggerStruct{Logger: log.New(mw, "", log.LstdFlags), LogFile: logFile, environment: env}
	return &logger
}
