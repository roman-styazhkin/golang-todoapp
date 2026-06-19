package users_transport_http

import (
	"net/http"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_request "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"`
}

type CreateUserResponse UserDTO

func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse("failed to decode request", err)
		return
	}

	userDomain := domainFromDTO(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse("failed create user", err)
		return
	}

	response := CreateUserResponse(DTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(
		dto.FullName,
		dto.PhoneNumber,
	)
}
