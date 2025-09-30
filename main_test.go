package main

import (
	"os"
	"testing"
	"time"
)

func TestLogEntryCreation(t *testing.T) {
	// Test creating a LogEntry with valid data
	timestamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	entry := LogEntry{
		Timestamp: timestamp,
		Level:     "INFO",
		Message:   "Test message",
	}

	// Verify the fields are set correctly
	if entry.Timestamp != timestamp {
		t.Errorf("Expected timestamp %v, got %v", timestamp, entry.Timestamp)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}

	if entry.Message != "Test message" {
		t.Errorf("Expected message 'Test message', got '%s'", entry.Message)
	}
}

func TestLogEntryEmpty(t *testing.T) {
	// Test creating an empty LogEntry
	entry := LogEntry{}

	// Verify default values
	if !entry.Timestamp.IsZero() {
		t.Errorf("Expected zero timestamp, got %v", entry.Timestamp)
	}

	if entry.Level != "" {
		t.Errorf("Expected empty level, got %s", entry.Level)
	}

	if entry.Message != "" {
		t.Errorf("Expected empty message, got '%s'", entry.Message)
	}
}

func TestLogFileExists(t *testing.T) {
	// Test that the sample log file exists and is readable
	_, err := os.Stat("app.log")
	if err != nil {
		t.Errorf("Sample log file 'app.log' not found: %v", err)
	}
}

func TestLogFileContent(t *testing.T) {
	// Test that the log file has content
	content, err := os.ReadFile("app.log")
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file is empty")
	}

	// Check if file contains expected log format
	contentStr := string(content)
	if len(contentStr) < 50 {
		t.Error("Log file seems too short")
	}
}
