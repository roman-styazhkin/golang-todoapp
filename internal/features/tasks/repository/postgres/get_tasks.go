package tasks_repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (r *TaskRepository) GetTasks(
	ctx context.Context,
	userID *int,
	limit *int,
	offset *int,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOpTimeout())
	defer cancel()

	var builder strings.Builder
	builder.WriteString(`
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
	`)

	args := make([]any, 0)

	if userID != nil {
		builder.WriteString(fmt.Sprintf("WHERE author_user_id=$%d ", len(args)+1))
		args = append(args, userID)
	}

	builder.WriteString("ORDER BY id ASC ")

	if limit != nil {
		builder.WriteString(fmt.Sprintf("LIMIT $%d ", len(args)+1))
		args = append(args, limit)
	}

	if offset != nil {
		builder.WriteString(fmt.Sprintf("OFFSET $%d ", len(args)+1))
		args = append(args, offset)
	}

	rows, err := r.pool.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks, %w", err)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		if err = rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		); err != nil {
			return nil, fmt.Errorf("failed to scan task, %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed with rows, %w", err)
	}

	taskDomains := domainListFromModels(taskModels)

	return taskDomains, nil
}
