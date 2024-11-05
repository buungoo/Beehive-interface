package utils

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

var ApiLogger *Logger

type Logger struct {
	Info  *log.Logger
	Warn  *log.Logger
	Error *log.Logger
	Fatal *log.Logger
}

func InitLogger() (*os.File, error) {
	logFile, err := os.OpenFile("/logs/logFile.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	logFlags := log.Ldate | log.Ltime
	ApiLogger = &Logger{
		Info:  log.New(logFile, "INFO: ", logFlags),
		Warn:  log.New(logFile, "WARN: ", logFlags),
		Error: log.New(logFile, "ERROR: ", logFlags),
		Fatal: log.New(logFile, "FATAL: ", logFlags),
	}

	return logFile, nil
}

func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	return fmt.Sprintf("%s:%d", file, line)
}
func LogInfo(message string) {
	if ApiLogger != nil {
		ApiLogger.Info.Println(callerInfo(), message)
	}
}

func LogWarn(message string) {
	if ApiLogger != nil {
		ApiLogger.Warn.Println(callerInfo(), message)
	}
}

func LogError(message string, err error) {
	if ApiLogger != nil {
		ApiLogger.Error.Println(callerInfo(), message, err)
	}
}

func LogFatal(message string, err error) {
	if ApiLogger != nil {
		ApiLogger.Fatal.Println(callerInfo(), message, err)
		os.Exit(1)
	}
}
