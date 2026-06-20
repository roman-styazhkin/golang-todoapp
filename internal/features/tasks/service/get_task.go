package tasks_service

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (s *TasksService) GetTask(
	ctx context.Context,
	id int,
) (domain.Task, error) {
	taskDomain, err := s.tasksRepository.GetTask(ctx, id)

	if err != nil {
		return domain.Task{}, fmt.Errorf("failed to get task, %w", err)
	}

	return taskDomain, nil
}
