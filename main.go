package main

import (
	"fmt"
	"log"

	"github.com/Oralbaev/logflow/internal/logflow"
)

func main() {
	service := logflow.NewService(logflow.LogProcessor{})

	logs := []string{
		"2026-10-16 10:00:00 INFO Server started",
		"2026-10-16 10:01:00 ERROR Database connection failed",
		"2026-10-16 10:02:00 WARN High memory usage",
		"2026-10-16 10:03:00 INFO Request completed",
		"2026-10-16 10:04:00 ERROR Request timeout",
	}

	task, err := service.CreateTask(logs)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task created: ID=%d, Status=%s\n",
		task.ID, task.Status)

	task, err = service.ProcessTask(task.ID)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Task processed: ID=%d, Status=%s\n",
		task.ID, task.Status)

	result := task.Result

	fmt.Printf("\nLog Analysis Results:\n")
	fmt.Printf("Total logs: %d\n", result.Total)
	fmt.Printf("INFO: %d\n", result.InfoCount)
	fmt.Printf("WARN: %d\n", result.WarnCount)
	fmt.Printf("ERROR: %d\n", result.ErrorCount)
	fmt.Printf("Error rate: %.2f%%\n", result.ErrorRate)
}
