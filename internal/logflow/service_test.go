package logflow

import (
	"errors"
	"testing"
)

func validLogs() []string {
	return []string{
		"2026-10-16 10:00:00 INFO Server started",
		"2026-10-16 10:01:00 WARN High memory usage",
		"2026-10-16 10:02:00 ERROR Connection failed",
		"2026-10-16 10:03:00 ERROR Request timeout",
	}
}

func TestCreateTask(t *testing.T) {
	service := NewService(LogProcessor{})

	task, err := service.CreateTask(validLogs())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if task.ID != 1 {
		t.Errorf("expected ID 1, got %d", task.ID)
	}

	if task.Status != StatusPending {
		t.Errorf("expected pending, got %s", task.Status)
	}

	if task.Result != nil {
		t.Error("new task should not have a result")
	}

	if len(task.Logs) != 4 {
		t.Errorf("expected 4 logs, got %d", len(task.Logs))
	}
}

func TestInvalidInput(t *testing.T) {
	tests := []struct {
		name    string
		logs    []string
		wantErr error
	}{
		{
			name:    "empty input",
			logs:    nil,
			wantErr: ErrEmptyLogs,
		},
		{
			name:    "empty slice",
			logs:    []string{},
			wantErr: ErrEmptyLogs,
		},
		{
			name:    "whitespace only",
			logs:    []string{"   "},
			wantErr: ErrInvalidLog,
		},
		{
			name:    "invalid timestamp",
			logs:    []string{"invalid 10:00:00 INFO Started"},
			wantErr: ErrInvalidLog,
		},
		{
			name:    "unknown level",
			logs:    []string{"2026-10-16 10:00:00 DEBUG Started"},
			wantErr: ErrInvalidLog,
		},
		{
			name:    "missing message",
			logs:    []string{"2026-10-16 10:00:00 INFO"},
			wantErr: ErrInvalidLog,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := NewService(LogProcessor{})

			_, err := service.CreateTask(tt.logs)

			if !errors.Is(err, tt.wantErr) {
				t.Errorf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	service := NewService(LogProcessor{})

	created, err := service.CreateTask(validLogs())
	if err != nil {
		t.Fatal(err)
	}

	found, err := service.GetTask(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if found.ID != created.ID {
		t.Errorf("expected ID %d, got %d", created.ID, found.ID)
	}

	if found.Status != StatusPending {
		t.Errorf("expected pending, got %s", found.Status)
	}
}

func TestTaskNotFound(t *testing.T) {
	service := NewService(LogProcessor{})

	_, err := service.GetTask(999)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}

	_, err = service.ProcessTask(999)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Errorf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestProcessTask(t *testing.T) {
	service := NewService(LogProcessor{})

	created, err := service.CreateTask(validLogs())
	if err != nil {
		t.Fatal(err)
	}

	task, err := service.ProcessTask(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !task.IsCompleted() {
		t.Errorf("expected completed, got %s", task.Status)
	}

	if task.Result == nil {
		t.Fatal("expected processing result")
	}

	expected := LogResult{
		Total:      4,
		InfoCount:  1,
		WarnCount:  1,
		ErrorCount: 2,
		ErrorRate:  50,
	}

	if *task.Result != expected {
		t.Errorf("expected %+v, got %+v", expected, *task.Result)
	}

	stored, err := service.GetTask(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if stored.Status != StatusCompleted || stored.Result == nil {
		t.Error("processed result was not saved")
	}
}

func TestCannotProcessTwice(t *testing.T) {
	service := NewService(LogProcessor{})

	task, err := service.CreateTask(validLogs())
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.ProcessTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.ProcessTask(task.ID)
	if !errors.Is(err, ErrInvalidState) {
		t.Errorf("expected ErrInvalidState, got %v", err)
	}
}

func TestMultipleIndependentTasks(t *testing.T) {
	service := NewService(LogProcessor{})

	first, err := service.CreateTask(validLogs())
	if err != nil {
		t.Fatal(err)
	}

	second, err := service.CreateTask([]string{
		"2026-10-16 11:00:00 INFO Everything works",
	})
	if err != nil {
		t.Fatal(err)
	}

	if first.ID == second.ID {
		t.Fatal("task IDs must be unique")
	}

	_, err = service.ProcessTask(first.ID)
	if err != nil {
		t.Fatal(err)
	}

	unchanged, err := service.GetTask(second.ID)
	if err != nil {
		t.Fatal(err)
	}

	if unchanged.Status != StatusPending || unchanged.Result != nil {
		t.Error("processing first task affected second task")
	}

	processed, err := service.ProcessTask(second.ID)
	if err != nil {
		t.Fatal(err)
	}

	if processed.Result == nil {
		t.Fatal("expected second task result")
	}

	if processed.Result.InfoCount != 1 ||
		processed.Result.ErrorCount != 0 {
		t.Error("second task contains incorrect statistics")
	}
}

func TestTaskDataIsolation(t *testing.T) {
	service := NewService(LogProcessor{})

	input := validLogs()

	task, err := service.CreateTask(input)
	if err != nil {
		t.Fatal(err)
	}

	input[0] = "modified outside service"
	task.Logs[0] = "modified returned task"

	stored, err := service.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	if stored.Logs[0] != validLogs()[0] {
		t.Error("external modification affected stored logs")
	}

	processed, err := service.ProcessTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	if processed.Result == nil {
		t.Fatal("expected processing result")
	}

	processed.Result.Total = 999

	again, err := service.GetTask(task.ID)
	if err != nil {
		t.Fatal(err)
	}

	if again.Result == nil || again.Result.Total != 4 {
		t.Error("external modification affected stored result")
	}
}
