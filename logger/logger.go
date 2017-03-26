package logger

import (
	"fmt"
	"log"

	"github.com/davidqhr/sock5/client"
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

func Debug(client *client.Client, format string, formatParams ...interface{}) {
	log.Printf("\033[0;33m[DEBUG]\033[0m (%s): %s\n", client.Id, fmt.Sprintf(format, formatParams...))
}

func Error(client *client.Client, format string, formatParams ...interface{}) {
	log.Printf("\033[0;31m[ERROR]\033[0m (%s) :%v\n", client.Id, fmt.Sprintf(format, formatParams...))
}

func Info(client *client.Client, format string, formatParams ...interface{}) {
	log.Printf("\033[0;32m[INFO]\033[0m (%s) :%v\n", client.Id, fmt.Sprintf(format, formatParams...))
}
