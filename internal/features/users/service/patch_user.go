package users_service

import (
	"context"
	"fmt"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

func (s *UsersService) PatchUser(
	ctx context.Context,
	id int,
	patch domain.UserPatch,
) (domain.User, error) {
	user, err := s.usersRepository.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to get user by id=%d, %w",
			id,
			err,
		)
	}

	if err = user.ApplyPatch(patch); err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to apply patch, %w",
			err,
		)
	}

	userDomain, err := s.usersRepository.PatchUser(ctx, id, user)
	if err != nil {
		return domain.User{}, fmt.Errorf(
			"failed to patch user, %w",
			err,
		)
	}

	return userDomain, nil
}
