package main

import (
	"bufio"
	"fmt"
	"os"
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
	// Por enquanto, vamos implementar uma versão simples que sempre retorna erro
	// Na próxima etapa implementaremos o regex para extrair os dados
	return fmt.Errorf("ParseLine not implemented yet")
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
