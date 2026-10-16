package logflow

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var (
	ErrEmptyLogs  = errors.New("logs cannot be empty")
	ErrInvalidLog = errors.New("invalid log format")
)

type Processor interface {
	Validate(logs []string) error
	Process(logs []string) (LogResult, error)
}

type LogProcessor struct{}

var _ Processor = LogProcessor{}

func (p LogProcessor) Validate(logs []string) error {
	if len(logs) == 0 {
		return ErrEmptyLogs
	}

	for i, line := range logs {
		fields := strings.Fields(line)

		if len(fields) < 4 {
			return fmt.Errorf(
				"%w at line %d: missing required fields",
				ErrInvalidLog,
				i+1,
			)
		}

		timestamp := fields[0] + " " + fields[1]

		if _, err := time.Parse(
			"2006-01-02 15:04:05",
			timestamp,
		); err != nil {
			return fmt.Errorf(
				"%w at line %d: invalid timestamp: %v",
				ErrInvalidLog,
				i+1,
				err,
			)
		}

		switch fields[2] {
		case "INFO", "WARN", "ERROR":
			// Supported log level.
		default:
			return fmt.Errorf(
				"%w at line %d: unsupported level %q",
				ErrInvalidLog,
				i+1,
				fields[2],
			)
		}
	}

	return nil
}

func (p LogProcessor) Process(logs []string) (LogResult, error) {
	if err := p.Validate(logs); err != nil {
		return LogResult{}, err
	}

	result := LogResult{
		Total: len(logs),
	}

	for _, line := range logs {
		fields := strings.Fields(line)

		switch fields[2] {
		case "INFO":
			result.InfoCount++

		case "WARN":
			result.WarnCount++

		case "ERROR":
			result.ErrorCount++
		}
	}

	result.ErrorRate = float64(result.ErrorCount) /
		float64(result.Total) * 100

	return result, nil
}
