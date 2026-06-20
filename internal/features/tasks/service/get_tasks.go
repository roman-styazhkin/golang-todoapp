package tasks_service

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func (s *TasksService) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	if limit != nil && *limit < 0 {
		return nil, fmt.Errorf(
			"failed to validate limit, limit must be positive integer, %w",
			core_errors.ErrInvalidArgument,
		)
	}

	if offset != nil && *offset < 0 {
		return nil, fmt.Errorf(
			"failed to validate offset, offset must be positive integer, %w",
			core_errors.ErrInvalidArgument,
		)
	}

	taskDomains, err := s.tasksRepository.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get tasks, %v, %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return taskDomains, nil
}
