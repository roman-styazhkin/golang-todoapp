package users_transport_http

import "github.com/roman-styazhkin/golang-todoapp/internal/core/domain"

type UserDTO struct {
	ID          int     `json:"id" example:"10"`
	Version     int     `json:"version" example:"1"`
	FullName    string  `json:"full_name" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" example:"+78005553535"`
}

func NewUserDTO(
	id int,
	version int,
	fullName string,
	phoneNumber *string,
) UserDTO {
	return UserDTO{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func DTOFromDomain(domain domain.User) UserDTO {
	return NewUserDTO(
		domain.ID,
		domain.Version,
		domain.FullName,
		domain.PhoneNumber,
	)
}

func dtoListFromDomains(domains []domain.User) []UserDTO {
	dtoList := make([]UserDTO, len(domains))

	for i, domain := range domains {
		dtoList[i] = DTOFromDomain(domain)
	}

	return dtoList
}
