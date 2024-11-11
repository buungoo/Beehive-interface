package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

// Pointers to custom loggers
var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger
var fatalLogger *log.Logger

// Open/Create a logfile and return it with write permission and initialize the logger
func InitLogger() (*os.File, error) {
	logFile, err := os.OpenFile("/logs/logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	logFlags := log.Ldate | log.Ltime
	// Initialize the loggers
	infoLogger = log.New(logFile, "INFO: ", logFlags)
	warnLogger = log.New(logFile, "WARN: ", logFlags)
	errorLogger = log.New(logFile, "ERROR: ", logFlags)
	fatalLogger = log.New(logFile, "FATAL: ", logFlags)

	return logFile, nil
}

// Find calling function from the callstack
func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
func LogInfo(message string) {
	if infoLogger != nil {
		infoLogger.Println(callerInfo(), message)
	}
}

func LogWarn(message string) {
	if warnLogger != nil {
		warnLogger.Println(callerInfo(), message)
	}
}

func LogError(message string, err error) {
	if errorLogger != nil {
		errorLogger.Println(callerInfo(), message, err)
	}
}

func LogFatal(message string, err error) {
	if fatalLogger != nil {
		fatalLogger.Println(callerInfo(), message, err)
		os.Exit(1)
	}
}
