package logger

import "fmt"

func Info(msg string) {
	fmt.Printf("[INFO] %s\n", msg)
}

// new
func Error(msg string) {
	fmt.Printf("[ERROR] %s\n", msg)
}

func Debug(msg string) {
	fmt.Printf("[DEBUG] %s\n", msg)
}
