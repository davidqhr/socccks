package socks5

import (
	"fmt"
	"log"
)

// # Reset
// Color_Off='\033[0m'       # Text Reset
//
// # Regular Colors
// Black='\033[0;30m'        # Black
// Red='\033[0;31m'          # Red
// Green='\033[0;32m'        # Green
// Yellow='\033[0;33m'       # Yellow
// Blue='\033[0;34m'         # Blue
// Purple='\033[0;35m'       # Purple
// Cyan='\033[0;36m'         # Cyan
// White='\033[0;37m'        # White

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
