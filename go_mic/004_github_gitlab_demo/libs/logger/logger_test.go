package logger

import "testing"

func TestInfo(t *testing.T) {
	// smoke test — just ensure it doesn't panic
	Info("test message")
}
