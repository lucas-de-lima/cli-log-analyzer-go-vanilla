package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"time"
)

// LogEntry representa uma entrada de log com timestamp, nível e mensagem
type LogEntry struct {
	Timestamp time.Time
	Level     string
	Message   string
}

// ParseLine parses a log line and populates the LogEntry fields
// Este método recebe uma linha bruta do log e extrai timestamp, level e message
func (le *LogEntry) ParseLine(rawLine string) error {
	// Define o padrão regex com grupos de captura nomeados
	// Formato esperado: [YYYY-MM-DD HH:MM:SS] [LEVEL] Mensagem
	pattern := `^\[(?P<timestamp>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2})\] \[(?P<level>\w+)\] (?P<message>.*)$`

	// Compila o regex para melhor performance
	regex := regexp.MustCompile(pattern)

	// Encontra correspondências na linha
	matches := regex.FindStringSubmatch(rawLine)
	if matches == nil {
		return fmt.Errorf("linha não corresponde ao formato esperado: %s", rawLine)
	}

	// Extrai os grupos nomeados usando os índices
	// matches[0] = linha completa
	// matches[1] = timestamp
	// matches[2] = level
	// matches[3] = message
	timestampStr := matches[1]
	level := matches[2]
	message := matches[3]

	// Converte a string de timestamp para time.Time
	// Usa o layout específico do nosso formato
	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		return fmt.Errorf("erro ao converter timestamp '%s': %w", timestampStr, err)
	}

	// Preenche os campos da struct
	le.Timestamp = timestamp
	le.Level = level
	le.Message = message

	return nil
}

// LoadLogs lê um arquivo de log e retorna um slice de LogEntry
func LoadLogs(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue // Pula linhas vazias
		}

		var entry LogEntry
		if err := entry.ParseLine(line); err != nil {
			// Por enquanto, vamos pular linhas mal formatadas e continuar
			// Em um sistema de produção, você pode querer registrar esses erros
			continue
		}

		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return entries, nil
}

func main() {
	fmt.Println("CLI Log Analyzer - Go Vanilla")
}
