package logflow

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
)

type LogResult struct {
	Total      int
	InfoCount  int
	WarnCount  int
	ErrorCount int
	ErrorRate  float64
}

type Task struct {
	ID     int
	Logs   []string
	Status TaskStatus
	Result *LogResult
}

func (t Task) IsCompleted() bool {
	return t.Status == StatusCompleted
}
