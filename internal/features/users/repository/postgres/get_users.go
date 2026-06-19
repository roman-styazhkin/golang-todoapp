package users_repository

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (r *UsersRepository) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.GetOpTimeout())
	defer cancel()

	query := `
	SELECT id, version, full_name, phone_number
	FROM todoapp.users
	ORDER BY id ASC
	LIMIT $1
	OFFSET $2;
	`

	rows, err := r.pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get users, %w", err)
	}

	var userModels []UserModel

	for rows.Next() {
		var userModel UserModel
		err = rows.Scan(
			&userModel.ID,
			&userModel.Version,
			&userModel.FullName,
			&userModel.PhoneNumber,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to scan user, %w", err)
		}

		userModels = append(userModels, userModel)
	}

	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed rows user, %w", err)
	}

	domainList := domainListFromModels(userModels)
	return domainList, nil
}
