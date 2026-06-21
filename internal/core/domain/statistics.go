package domain

import "time"

type Statistics struct {
	TasksCompleted             int
	TasksCreated               int
	TasksCompletedRate         *float64
	TasksAverageCompletionTime *time.Duration
}

func NewStatistics(
	tasksCompleted int,
	tasksCreated int,
	tasksCompletedRate *float64,
	tasksAverageCompletionTime *time.Duration,
) Statistics {
	return Statistics{
		TasksCompleted:             tasksCompleted,
		TasksCreated:               tasksCreated,
		TasksCompletedRate:         tasksCompletedRate,
		TasksAverageCompletionTime: tasksAverageCompletionTime,
	}
}
