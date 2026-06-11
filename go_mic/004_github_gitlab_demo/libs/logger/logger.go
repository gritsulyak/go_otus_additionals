package logger

import "fmt"

func Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

func Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}
