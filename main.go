package main

import (
	"fmt"
	"time"
)

// LogEntry represents a single log entry with timestamp, level, and message
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

func main() {
	fmt.Println("CLI Log Analyzer - Go Vanilla")
}
