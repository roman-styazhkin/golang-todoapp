package users_service

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (s *UsersService) GetUser(
	ctx context.Context,
	id int,
) (domain.User, error) {
	userDomain, err := s.usersRepository.GetUser(ctx, id)

	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to get user, %w",
			err,
		)
	}

	return userDomain, nil
}
