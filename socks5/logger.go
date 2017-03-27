package socks5

import (
	"fmt"
	"log"
)

type Logger struct{}

func (_ *Logger) Debug(args ...interface{}) {
	log.Printf("\033[0;33m[DEBUG]\033[0m: %s", fmt.Sprintln(args...))
}

func (_ *Logger) Error(args ...interface{}) {
	log.Printf("\033[0;31m[ERROR]\033[0m: %v", fmt.Sprintln(args...))
}

func (_ *Logger) Info(args ...interface{}) {
	log.Printf("\033[0;32m[INFO ]\033[0m: %v", fmt.Sprintln(args...))
}

var logger = &Logger{}
