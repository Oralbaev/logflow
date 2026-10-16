package logflow

import (
	"errors"
	"fmt"
)

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrInvalidState = errors.New("invalid task state")
)

type Service struct {
	tasks     map[int]Task
	nextID    int
	processor Processor
}

func NewService(processor Processor) *Service {
	return &Service{
		tasks:     make(map[int]Task),
		nextID:    1,
		processor: processor,
	}
}

func (s *Service) CreateTask(logs []string) (Task, error) {
	if err := s.processor.Validate(logs); err != nil {
		return Task{}, err
	}

	task := Task{
		ID:     s.nextID,
		Logs:   append([]string(nil), logs...),
		Status: StatusPending,
	}

	s.tasks[task.ID] = task
	s.nextID++

	return cloneTask(task), nil
}

func (s *Service) GetTask(id int) (Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return Task{}, fmt.Errorf(
			"%w: id=%d",
			ErrTaskNotFound,
			id,
		)
	}

	return cloneTask(task), nil
}

func cloneTask(task Task) Task {
	task.Logs = append([]string(nil), task.Logs...)

	if task.Result != nil {
		result := *task.Result
		task.Result = &result
	}

	return task
}

func (s *Service) ProcessTask(id int) (Task, error) {
	task, exists := s.tasks[id]
	if !exists {
		return Task{}, fmt.Errorf(
			"%w: id=%d",
			ErrTaskNotFound,
			id,
		)
	}

	if task.Status != StatusPending {
		return Task{}, fmt.Errorf(
			"%w: task %d has status %s",
			ErrInvalidState,
			id,
			task.Status,
		)
	}

	task.Status = StatusProcessing
	s.tasks[id] = task

	result, err := s.processor.Process(task.Logs)
	if err != nil {
		task.Status = StatusPending
		s.tasks[id] = task

		return Task{}, fmt.Errorf(
			"process task %d: %w",
			id,
			err,
		)
	}

	task.Result = &result
	task.Status = StatusCompleted

	s.tasks[id] = task

	return cloneTask(task), nil
}
