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
// Retorna erro se o arquivo não puder ser aberto ou lido
func LoadLogs(filename string) ([]LogEntry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo %s: %w", filename, err)
	}
	defer file.Close()

	var entries []LogEntry
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()

		// Pula linhas vazias
		if line == "" {
			continue
		}

		// Tenta fazer o parse da linha
		var entry LogEntry
		if err := entry.ParseLine(line); err != nil {
			// Se o parse falhar, pula a linha mas continua
			// Em um sistema de produção, você pode querer registrar esses erros
			// Por enquanto, apenas ignoramos linhas mal formatadas
			continue
		}

		entries = append(entries, entry)
	}

	// Verifica se houve erro ao ler o arquivo
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Retorna sucesso mesmo que algumas linhas tenham sido ignoradas
	return entries, nil
}

func main() {
	fmt.Println("CLI Log Analyzer - Go Vanilla")
}
