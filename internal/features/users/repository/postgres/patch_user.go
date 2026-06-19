package users_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func (r *UsersRepository) PatchUser(
	ctx context.Context,
	id int,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.pool.GetOpTimeout())
	defer cancel()

	query := `
	UPDATE todoapp.users
	SET full_name=$1, phone_number=$2, version=version+1
	WHERE id=$3 AND version=$4
	RETURNING id, version, full_name, phone_number;
	`

	row := r.pool.QueryRow(ctx, query,
		user.FullName,
		user.PhoneNumber,
		id,
		user.Version,
	)

	var userModel UserModel

	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf(
				"user with id=%d not found, %w",
				id,
				core_errors.ErrNotFound,
			)
		}

		return domain.User{}, fmt.Errorf(
			"failed to scan user, %w",
			err,
		)
	}

	userDomain := domainFromModel(userModel)
	return userDomain, nil
}
