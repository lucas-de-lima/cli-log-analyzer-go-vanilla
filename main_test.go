package main

import (
	"os"
	"testing"
	"time"
)

func TestLogEntryCreation(t *testing.T) {
	// Testa a criação de uma LogEntry com dados válidos
	timestamp := time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC)
	entry := LogEntry{
		Timestamp: timestamp,
		Level:     "INFO",
		Message:   "Test message",
	}

	// Verifica se os campos foram definidos corretamente
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
	// Testa a criação de uma LogEntry vazia
	entry := LogEntry{}

	// Verifica os valores padrão
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
	// Testa se o arquivo de log de exemplo existe e é legível
	_, err := os.Stat("app.log")
	if err != nil {
		t.Errorf("Sample log file 'app.log' not found: %v", err)
	}
}

func TestLogFileContent(t *testing.T) {
	// Testa se o arquivo de log tem conteúdo
	content, err := os.ReadFile("app.log")
	if err != nil {
		t.Errorf("Failed to read log file: %v", err)
	}

	if len(content) == 0 {
		t.Error("Log file is empty")
	}

	// Verifica se o arquivo contém o formato de log esperado
	contentStr := string(content)
	if len(contentStr) < 50 {
		t.Error("Log file seems too short")
	}
}

func TestLoadLogsFileNotFound(t *testing.T) {
	// Testa o carregamento de um arquivo inexistente
	entries, err := LoadLogs("nonexistent.log")

	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	if entries != nil {
		t.Error("Expected nil entries for non-existent file")
	}
}

func TestLoadLogsEmptyFile(t *testing.T) {
	// Cria um arquivo temporário vazio
	tmpFile, err := os.CreateTemp("", "test_empty.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Testa o carregamento de arquivo vazio
	entries, err := LoadLogs(tmpFile.Name())

	if err != nil {
		t.Errorf("Unexpected error loading empty file: %v", err)
	}

	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for empty file, got %d", len(entries))
	}
}

func TestLoadLogsValidFile(t *testing.T) {
	// Cria um arquivo temporário com dados de teste
	tmpFile, err := os.CreateTemp("", "test_valid.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Escreve dados de teste
	testData := `[2025-01-15 10:00:00] [INFO] Test message 1
[2025-01-15 10:05:00] [ERROR] Test error message`

	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Testa o carregamento de arquivo válido
	entries, err := LoadLogs(tmpFile.Name())

	if err != nil {
		t.Errorf("Unexpected error loading valid file: %v", err)
	}

	// Nota: Este teste falhará até implementarmos o método ParseLine
	// Por enquanto, esperamos 0 entradas porque ParseLine falhará
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries (ParseLine not implemented), got %d", len(entries))
	}
}

func TestParseLineNotImplemented(t *testing.T) {
	// Testa o método ParseLine (implementação placeholder)
	entry := LogEntry{}
	err := entry.ParseLine("[2025-01-15 10:00:00] [INFO] Test message")

	if err == nil {
		t.Error("Expected error for ParseLine not implemented, got nil")
	}

	expectedError := "ParseLine not implemented yet"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestParseLineWithEmptyString(t *testing.T) {
	// Testa ParseLine com string vazia
	entry := LogEntry{}
	err := entry.ParseLine("")

	if err == nil {
		t.Error("Expected error for ParseLine with empty string, got nil")
	}
}

func TestParseLineWithInvalidFormat(t *testing.T) {
	// Testa ParseLine com formato inválido
	entry := LogEntry{}
	err := entry.ParseLine("This is not a valid log line")

	if err == nil {
		t.Error("Expected error for ParseLine with invalid format, got nil")
	}
}
