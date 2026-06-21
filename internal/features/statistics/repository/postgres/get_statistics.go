package statistics_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (r *StatisticsRepository) GetTasks(
	ctx context.Context,
	userID *int,
	from *time.Time,
	to *time.Time,
) ([]domain.Task, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOpTimeout())
	defer cancel()

	var builder strings.Builder

	builder.WriteString(`
		SELECT id, version, title, description, completed, created_at, completed_at, author_user_id
		FROM todoapp.tasks
	`)

	args := []any{}
	condition := []string{}

	if userID != nil {
		condition = append(condition, fmt.Sprintf("author_user_id=$%d", len(args)+1))
		args = append(args, userID)
	}

	if from != nil {
		condition = append(condition, fmt.Sprintf("created_at>=$%d", len(args)+1))
		args = append(args, from)
	}

	if to != nil {
		condition = append(condition, fmt.Sprintf("created_at<$%d", len(args)+1))
		args = append(args, to)
	}

	if len(condition) > 0 {
		builder.WriteString(" WHERE " + strings.Join(condition, "AND"))
	}

	builder.WriteString(" ORDER BY id ASC ")

	rows, err := r.pool.Query(ctx, builder.String(), args...)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get tasks, err: %w",
			err,
		)
	}
	defer rows.Close()

	var taskModels []TaskModel

	for rows.Next() {
		var taskModel TaskModel

		err = rows.Scan(
			&taskModel.ID,
			&taskModel.Version,
			&taskModel.Title,
			&taskModel.Description,
			&taskModel.Completed,
			&taskModel.CreatedAt,
			&taskModel.CompletedAt,
			&taskModel.AuthorUserID,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan task, %w", err)
		}

		taskModels = append(taskModels, taskModel)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed with next rows, %w", err)
	}

	taskDomains := domainListFromModels(taskModels)
	return taskDomains, nil
}
