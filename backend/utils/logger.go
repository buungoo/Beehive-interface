package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var apiLogger *Logger

type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
}

// Open/Create a logfile and return it with write permission and initialize the logger
func InitLogger() (*os.File, error) {
	logFile, err := os.OpenFile("/logs/logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	logFlags := log.Ldate | log.Ltime
	apiLogger = &Logger{
		Info:  log.New(logFile, "INFO: ", logFlags),
		Warn:  log.New(logFile, "WARN: ", logFlags),
		Error: log.New(logFile, "ERROR: ", logFlags),
		Fatal: log.New(logFile, "FATAL: ", logFlags),
	}

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
	if apiLogger != nil {
		apiLogger.Info.Println(callerInfo(), message)
	}
}

func LogWarn(message string) {
	if apiLogger != nil {
		apiLogger.Warn.Println(callerInfo(), message)
	}
}

func LogError(message string, err error) {
	if apiLogger != nil {
		apiLogger.Error.Println(callerInfo(), message, err)
	}
}

func LogFatal(message string, err error) {
	if apiLogger != nil {
		apiLogger.Fatal.Println(callerInfo(), message, err)
		os.Exit(1)
	}
}
