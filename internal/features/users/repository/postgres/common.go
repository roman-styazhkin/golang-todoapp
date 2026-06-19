package users_repository

import "github.com/roman-styazhkin/golang-todoapp/internal/core/domain"

type UserModel struct {
	ID          int
	Version     int
	FullName    string
	PhoneNumber *string
}

func NewUserModel(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) UserModel {
	return UserModel{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func domainFromModel(model UserModel) domain.User {
	return domain.NewUser(
		model.ID,
		model.Version,
		model.FullName,
		model.PhoneNumber,
	)
}

func domainListFromModels(models []UserModel) []domain.User {
	domainList := make([]domain.User, len(models))

	for i, model := range models {
		domainList[i] = domainFromModel(model)
	}

	return domainList
}
