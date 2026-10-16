# LogFlow

**LogFlow** is a lightweight log processing engine built with Go. It validates structured server logs, analyzes severity levels, and manages processing tasks using an in-memory storage system.

This project is developed as part of the SIS course project for Golang Application Development at KBTU.

## Features

- Structured server log validation
- Severity-level analysis (INFO, WARN, ERROR)
- Error rate calculation
- Task creation with unique identifiers
- In-memory task storage and retrieval
- Task lifecycle management
- Error handling and validation
- Unit tests with Go's standard testing package

## Supported Log Format

Each log entry must follow this structure:

```text
YYYY-MM-DD HH:MM:SS LEVEL Message
```

Example:

```text
2026-10-16 10:00:00 INFO Server started
2026-10-16 10:01:00 ERROR Database connection failed
2026-10-16 10:02:00 WARN High memory usage
2026-10-16 10:03:00 INFO Request completed
2026-10-16 10:04:00 ERROR Request timeout
```

Supported severity levels: `INFO`, `WARN`, `ERROR`.

## Example Output

```text
Task created: ID=1, Status=pending
Task processed: ID=1, Status=completed

Log Analysis Results:
Total logs: 5
INFO: 2
WARN: 1
ERROR: 2
Error rate: 40.00%
```

## Task Lifecycle

Each task follows a defined processing lifecycle:

```text
Pending -> Processing -> Completed
```

- **Pending:** Task has been created and validated.
- **Processing:** Log analysis is being performed.
- **Completed:** Processing finished successfully and the result is stored.

A task can only be processed when its status is `Pending`.

Processing a completed task again is not allowed. If processing fails, the task returns to `Pending`, and an error is returned.

## Project Structure

```text
logflow/
├── go.mod
├── main.go
├── README.md
└── internal/
    └── logflow/
        ├── task.go
        ├── processor.go
        ├── service.go
        └── service_test.go
```

## Architecture

**Task**

Represents a log processing request containing a unique identifier, input logs, processing status, and optional result.

**Processor**

An interface defining validation and processing operations.

**LogProcessor**

Implements log format validation and calculates log statistics.

**Service**

Manages task creation, retrieval, processing, and in-memory storage using a Go map.

## Requirements

- Go 1.27

## Installation

Clone the repository:

```bash
git clone https://github.com/Oralbaev/logflow.git
cd logflow
```

Run the application:

```bash
go run .
```

Run automated tests:

```bash
go test ./... -v
```

Run static analysis:

```bash
go vet ./...
```

Check code formatting:

```bash
gofmt -l .
```

## Technology Stack

- Go
- Go standard library
- In-memory storage
- Standard Go testing framework

## Roadmap

- [x] Checkpoint 1: Core Go, synchronous task processing, validation and unit tests
- [ ] Checkpoint 2: Concurrent log processing using goroutines, channels and worker pools
- [ ] Checkpoint 3: HTTP API, context cancellation and graceful shutdown

## Current Limitations

- Tasks are stored in memory and are lost when the application exits.
- Processing is synchronous.
- The application does not currently expose an HTTP API.
- Supported log levels are limited to INFO, WARN, and ERROR.