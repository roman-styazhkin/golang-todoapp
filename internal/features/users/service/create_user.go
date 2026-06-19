package users_service

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (s *UsersService) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	if err := user.Validate(); err != nil {
		return domain.User{}, fmt.Errorf("failed to validate user, %w", err)
	}

	userDomain, err := s.usersRepository.CreateUser(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user, %w", err)
	}

	return userDomain, nil
}
