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
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
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
