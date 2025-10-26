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

	// Agora que ParseLine está implementado, esperamos 2 entradas
	if len(entries) != 2 {
		t.Errorf("Expected 2 entries, got %d", len(entries))
	}

	// Verifica se a primeira entrada foi parseada corretamente
	if entries[0].Level != "INFO" {
		t.Errorf("Expected first entry level INFO, got %s", entries[0].Level)
	}

	if entries[1].Level != "ERROR" {
		t.Errorf("Expected second entry level ERROR, got %s", entries[1].Level)
	}
}

func TestParseLineValidFormat(t *testing.T) {
	// Testa ParseLine com formato válido
	entry := LogEntry{}
	err := entry.ParseLine("[2025-01-15 10:30:45] [INFO] User logged in successfully")

	if err != nil {
		t.Errorf("Unexpected error for valid format: %v", err)
	}

	// Verifica se os campos foram preenchidos corretamente
	expectedTime := time.Date(2025, 1, 15, 10, 30, 45, 0, time.UTC)
	if !entry.Timestamp.Equal(expectedTime) {
		t.Errorf("Expected timestamp %v, got %v", expectedTime, entry.Timestamp)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}

	if entry.Message != "User logged in successfully" {
		t.Errorf("Expected message 'User logged in successfully', got '%s'", entry.Message)
	}
}

func TestParseLineDifferentLevels(t *testing.T) {
	// Testa ParseLine com diferentes níveis de log
	testCases := []struct {
		line            string
		expectedLevel   string
		expectedMessage string
	}{
		{"[2025-01-15 10:00:00] [ERROR] Database connection failed", "ERROR", "Database connection failed"},
		{"[2025-01-15 10:05:00] [WARNING] High memory usage detected", "WARNING", "High memory usage detected"},
		{"[2025-01-15 10:10:00] [DEBUG] Processing request ID: 12345", "DEBUG", "Processing request ID: 12345"},
	}

	for _, tc := range testCases {
		entry := LogEntry{}
		err := entry.ParseLine(tc.line)

		if err != nil {
			t.Errorf("Unexpected error for line '%s': %v", tc.line, err)
		}

		if entry.Level != tc.expectedLevel {
			t.Errorf("Expected level %s, got %s", tc.expectedLevel, entry.Level)
		}

		if entry.Message != tc.expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", tc.expectedMessage, entry.Message)
		}
	}
}

func TestParseLineInvalidFormat(t *testing.T) {
	// Testa ParseLine com formatos inválidos
	invalidLines := []string{
		"This is not a valid log line",
		"[2025-01-15] [INFO] Missing time",
		"[2025-01-15 10:00:00] INFO Missing brackets",
		"[2025-01-15 10:00:00] [INFO]",
		"",
	}

	for _, line := range invalidLines {
		entry := LogEntry{}
		err := entry.ParseLine(line)

		if err == nil {
			t.Errorf("Expected error for invalid line '%s', got nil", line)
		}
	}
}

func TestParseLineInvalidTimestamp(t *testing.T) {
	// Testa ParseLine com timestamp inválido
	entry := LogEntry{}
	err := entry.ParseLine("[2025-13-45 25:70:99] [INFO] Invalid timestamp")

	if err == nil {
		t.Error("Expected error for invalid timestamp, got nil")
	}
}

func TestLoadRealLogFile(t *testing.T) {
	// Testa o carregamento do arquivo real app.log
	entries, err := LoadLogs("app.log")

	if err != nil {
		t.Errorf("Unexpected error loading real log file: %v", err)
	}

	// Verifica se carregou todas as 10 linhas do arquivo
	if len(entries) != 10 {
		t.Errorf("Expected 10 entries from app.log, got %d", len(entries))
	}

	// Verifica se os níveis foram parseados corretamente
	levels := make(map[string]int)
	for _, entry := range entries {
		levels[entry.Level]++
	}

	// Verifica se temos os níveis esperados
	expectedLevels := map[string]int{
		"INFO":    4,
		"WARNING": 2,
		"ERROR":   2,
		"DEBUG":   2,
	}

	for level, expectedCount := range expectedLevels {
		if levels[level] != expectedCount {
			t.Errorf("Expected %d %s entries, got %d", expectedCount, level, levels[level])
		}
	}
}

func TestLoadLogsWithMixedLines(t *testing.T) {
	// Testa o carregamento de arquivo com linhas válidas e inválidas misturadas
	tmpFile, err := os.CreateTemp("", "test_mixed.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Escreve dados mistos: válidos e inválidos
	testData := `[2025-01-15 10:00:00] [INFO] Valid log line 1
This is not a valid log line
[2025-01-15 10:05:00] [ERROR] Valid error line
Another invalid line
[2025-01-15 10:10:00] [WARNING] Valid warning line`

	if _, err := tmpFile.WriteString(testData); err != nil {
		t.Fatalf("Failed to write test data: %v", err)
	}
	tmpFile.Close()

	// Testa o carregamento do arquivo com linhas misturadas
	entries, err := LoadLogs(tmpFile.Name())

	if err != nil {
		t.Errorf("Unexpected error loading file with mixed lines: %v", err)
	}

	// Deve ter carregado apenas as 3 linhas válidas
	if len(entries) != 3 {
		t.Errorf("Expected 3 valid entries, got %d", len(entries))
	}

	// Verifica se as linhas válidas foram carregadas corretamente
	if entries[0].Level != "INFO" {
		t.Errorf("Expected first entry level INFO, got %s", entries[0].Level)
	}

	if entries[1].Level != "ERROR" {
		t.Errorf("Expected second entry level ERROR, got %s", entries[1].Level)
	}

	if entries[2].Level != "WARNING" {
		t.Errorf("Expected third entry level WARNING, got %s", entries[2].Level)
	}
}

func TestErrorPropagation(t *testing.T) {
	// Testa se os erros são propagados corretamente
	// Teste 1: Arquivo inexistente deve retornar erro
	_, err := LoadLogs("file_that_does_not_exist.log")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Verifica se a mensagem de erro contém informações úteis
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}

	// Teste 2: ParseLine com formato inválido deve retornar erro
	entry := LogEntry{}
	err = entry.ParseLine("Invalid line format")
	if err == nil {
		t.Error("Expected error for invalid line format, got nil")
	}

	// Verifica se a mensagem de erro contém informações úteis
	if err.Error() == "" {
		t.Error("Error message should not be empty")
	}
}
