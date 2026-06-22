package users_transport_http

import (
	"net/http"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_request "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"Ivan Ivanov"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+78005553535"`
}

type CreateUserResponse UserDTO

// CreateUser godoc
// @Summary Создание нового пользователя
// @Description Создание нового пользователя в системе
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "CreateUser тело запроса"
// @Success 201 {object} CreateUserResponse "Успешно созданный пользователь"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server error"
// @Router /users [post]
func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	var request CreateUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse("failed to decode request", err)
		return
	}

	userDomain := domainFromRequest(request)
	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse("failed create user", err)
		return
	}

	response := CreateUserResponse(DTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusCreated)
}

func domainFromRequest(request CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(
		request.FullName,
		request.PhoneNumber,
	)
}
