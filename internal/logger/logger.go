package logger

import (
	"log"
	"os"
)

var (
	InfoLog      = log.New(os.Stdout, "[INFO]: ", log.LstdFlags)
	ErrorLog     = log.New(os.Stderr, "[ERROR]: ", log.LstdFlags)
	DebugLog     = log.New(os.Stderr, "[DEBUG]: ", log.LstdFlags)
	DebugEnabled = false
)

func Debug(msg string, args ...interface{}) {
	if DebugEnabled {
		DebugLog.Printf(msg, args...)
	}
}

func Info(msg string, args ...interface{}) {
	InfoLog.Printf(msg, args...)
}

func Error(msg string, args ...interface{}) {
	ErrorLog.Printf(msg, args...)
}
