package users_repository

import (
	"context"
	"fmt"
)

func (r *UsersRepository) DeleteUser(
	ctx context.Context,
	id int,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.pool.GetOpTimeout())
	defer cancel()

	query := `
	DELETE FROM todoapp.users
	WHERE id=$1;
	`

	cmdTag, err := r.pool.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete user, %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf(
			"user with id=%d not found, err: %w",
			id,
			err,
		)
	}

	return nil
}
