package statistics_service

import (
	"context"
	"fmt"
	"time"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func (s *StatisticsService) GetStatistics(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) (domain.Statistics, error) {
	if from != nil && to != nil {
		if to.Before(*from) || from.Equal(*to) {
			return domain.Statistics{}, fmt.Errorf(
				"failed to validate 'from' 'to' query params, %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	tasks, err := s.staticsRepository.GetTasks(ctx, userID, from, to)
	if err != nil {
		return domain.Statistics{}, fmt.Errorf(
			"failed to get statistics, %w",
			err,
		)
	}

	statistics := calculateStatistics(tasks)
	return statistics, nil
}

func calculateStatistics(tasks []domain.Task) domain.Statistics {
	if len(tasks) == 0 {
		return domain.Statistics{
			TasksCompleted:             0,
			TasksCreated:               0,
			TasksCompletedRate:         nil,
			TasksAverageCompletionTime: nil,
		}
	}

	created := len(tasks)
	completed := 0

	var totalCompletionDuration time.Duration

	for _, task := range tasks {
		if task.Completed {
			completed++
		}

		duration := task.CompletionDuration()
		if duration != nil {
			totalCompletionDuration += *duration
		}
	}

	completedRate := float64(completed) / float64(created) * 100
	var avgCompletionTime *time.Duration

	if totalCompletionDuration != 0 && completed > 0 {
		avg := totalCompletionDuration / time.Duration(completed)
		avgCompletionTime = &avg
	}

	return domain.Statistics{
		TasksCompleted:             completed,
		TasksCreated:               created,
		TasksCompletedRate:         &completedRate,
		TasksAverageCompletionTime: avgCompletionTime,
	}
}
