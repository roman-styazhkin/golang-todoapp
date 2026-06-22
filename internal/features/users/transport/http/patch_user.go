package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_request "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/types"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Ivan Ivanov"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+71112223344"`
}

func (p *PatchUserRequest) Validate() error {
	if p.FullName.Set {
		if p.FullName.Value == nil {
			return fmt.Errorf(
				"failed to validate, full_name cannot be null, %w",
				core_errors.ErrInvalidArgument,
			)
		}

		fullNameLength := len([]rune(*p.FullName.Value))

		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf(
				"failed to validate, full_name length is incorrect, len: %d, %w",
				fullNameLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.PhoneNumber.Set {
		if p.PhoneNumber.Value != nil {
			phoneNumberLength := len([]rune(*p.PhoneNumber.Value))
			if phoneNumberLength < 10 || phoneNumberLength > 15 {
				return fmt.Errorf(
					"failed to validate, phone_number length is incorrect, len: %d, %w",
					phoneNumberLength,
					core_errors.ErrInvalidArgument,
				)
			}

			if !strings.HasPrefix(*p.PhoneNumber.Value, "+") {
				return fmt.Errorf(
					"failed to validate, phone_number must starts with +, err: %w",
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTO

// PatchUser godoc
// @Summary Изменение пользователя
// @Description Изменение существующего в системе пользователя
// @Description #### Логика обновления полей (Three state logic):
// @Description **1. Поле не передано** `phone_number is empty` игнорируется, значение phone_number в БД не меняется
// @Description **2. Поле явно передано** `phone_number:+791112223344` значение phone_number меняется на новое в БД
// @Description **3. Поле явно передано как null** `phone_number:null` значение phone_number меняется на null
// @Description **Ограничения: ** full_name не может быть передан как null
// @Description **Ограничения: ** full_name не может быть передан как "" (empty string)
// @Description **Ограничения: ** phone_number не может быть передан как "" (empty string)
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "id изменяемого пользователя"
// @Param request body PatchUserRequest true "PatchUser тело запроса"
// @Success 200 {object} PatchUserResponse "Успешно измененный пользователь"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrResponse "User not found"
// @Failure 409 {object} core_http_response.ErrResponse "Err conflict"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server error"
// @Router /users/{id} [patch]
func (h *UsersHttpHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	var request PatchUserRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse("failed to decode and validate request", err)
		return
	}

	id, err := core_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	patchDomain := domainFromPatch(request)
	userDomain, err := h.usersService.PatchUser(ctx, id, patchDomain)
	if err != nil {
		responseHandler.ErrorResponse("failed to patch user", err)
		return
	}

	response := PatchUserResponse(DTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func domainFromPatch(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
