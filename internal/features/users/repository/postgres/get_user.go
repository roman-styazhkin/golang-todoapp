package users_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_postgres_pool "github.com/roman-styazhkin/golang-todoapp/internal/core/repository/postgres/pool"
)

func (r *UsersRepository) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.pool.GetOpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	WHERE id=$1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id=%d not found, %w",
				id,
				err,
			)
		}

		return domain.User{}, fmt.Errorf("failed to scan user, %w", err)
	}

	userDomain := domainFromModel(userModel)
	return userDomain, nil
}
