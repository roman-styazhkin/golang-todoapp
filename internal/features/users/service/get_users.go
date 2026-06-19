package users_service

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
)

func (s *UsersService) GetUsers(
	ctx context.Context,
	limit *int,
	offset *int,
) ([]domain.User, error) {
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

	userDomains, err := s.usersRepository.GetUsers(ctx, limit, offset)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to get users from repository, %w",
			core_errors.ErrInvalidArgument,
		)
	}

	return userDomains, nil
}
