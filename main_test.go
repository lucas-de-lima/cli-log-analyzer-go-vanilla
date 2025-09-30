package main

import (
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
