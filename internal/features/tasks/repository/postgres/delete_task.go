package tasks_repository

import (
	"context"
	"fmt"
)

func (r *TaskRepository) DeleteTask(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.tasks
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to exec, %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"failed to delete task, task with id=%d, err:%w",
			id,
			err,
		)
	}

	return nil
}
