package logger

import (
	"fmt"
	"log"
	"os"
)

type messageType int

const (
	INFO messageType = 0 + iota
	GRPC
	WARNING
	ERROR
	FATAL
)

const (
	grpcColor    = "\033[1;35m%s\033[0m"
	infoColor    = "\033[1;34m%s\033[0m"
	warningColor = "\033[1;33m%s\033[0m"
	errorColor   = "\033[1;31m%s\033[0m"
	fatalColor   = "\033[41;1;37m%s\033[0m"
)

type logOutType string

const (
	STDOUT logOutType = "stdout"
	FILE   logOutType = "file"
)

// Write saves the message on log gile or via os.Stdout
func Write(messageType messageType, message string, logTarget logOutType) {

	target, err := os.OpenFile(".log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	if logTarget == "stdout" {
		target = os.Stdout
	}

	switch messageType {
	case INFO:
		logger := log.New(target, fmt.Sprintf(infoColor, "INFO: "), log.Ldate|log.Ltime)
		logger.Println(message)
	case WARNING:
		logger := log.New(target, fmt.Sprintf(warningColor, "WARNING: "), log.Ldate|log.Ltime)
		logger.Println(message)
	case ERROR:
		logger := log.New(target, fmt.Sprintf(errorColor, "ERROR: "), log.Ldate|log.Ltime)
		logger.Println(message)
	case FATAL:
		logger := log.New(target, fmt.Sprintf(fatalColor, "FATAL: "), log.Ldate|log.Ltime)
		logger.Fatal(message)
	case GRPC:
		logger := log.New(target, fmt.Sprintf(grpcColor, "gRPC: "), log.Ldate|log.Ltime)
		logger.Println(message)
	}
}
