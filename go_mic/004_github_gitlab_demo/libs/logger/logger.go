package logger

import "fmt"

func Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}
