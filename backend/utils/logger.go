// Package utils contains functions, methods and structs that can be used in multiple places throughtout the api.
//
// This package contains a custom logger that can be used with different levels of logging and be saved to a logfile.
// Utils package also has the Json response fuctions in the utils.go file.

package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

// Pointers to custom loggers
var infoLogger *log.Logger
var warnLogger *log.Logger
var errorLogger *log.Logger
var fatalLogger *log.Logger

// Opens/Creates a logfile and return it with write permission and initialize the logger.
func InitLogger(filePath string) (*os.File, error) {
	// Ensure the parent directory exists
	dir := filepath.Dir(filePath)
	err := os.MkdirAll(dir, 0755) // Create directory and parents if they don't exist
	if err != nil {
		return nil, err
	}

	logFile, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
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

// LogInfo writes to a logfile and is prefixed with the INFO tag.
func LogInfo(message string) {
	if infoLogger != nil {
		infoLogger.Println(callerInfo(), message)
	}
}

// LogWarn writes to a logfile and is prefixed with the WARN tag.
func LogWarn(message string) {
	if warnLogger != nil {
		warnLogger.Println(callerInfo(), message)
	}
}

// LogError writes to a logfile and is prefixed with the Error tag. It also includes the error message.
func LogError(message string, err error) {
	if errorLogger != nil {
		errorLogger.Println(callerInfo(), message, err)
	}
}

// LogInfo writes to a logfile and is prefixed with the FATAL tag. This also call terminates the program.
func LogFatal(message string, err error) {
	if fatalLogger != nil {
		fatalLogger.Println(callerInfo(), message, err)
		os.Exit(1)
	}
}
