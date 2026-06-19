package users_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type GetUserResponse []UserDTO

func (h *UsersHttpHandler) GetUsers(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	limit, offset, err := getLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			"failed to get limit, offset from query params",
			err,
		)
		return
	}

	userDomains, err := h.usersService.GetUsers(ctx, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse("failed to get users", err)
		return
	}

	response := GetUserResponse(dtoListFromDomains(userDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getLimitOffsetQueryParams(r *http.Request) (*int, *int, error) {
	const (
		queryLimit  = "limit"
		queryOffset = "offset"
	)

	limit, err := core_utils.GetIntQueryParam(r, queryLimit)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get limit from query params, %w",
			err,
		)
	}

	offset, err := core_utils.GetIntQueryParam(r, queryOffset)
	if err != nil {
		return nil, nil, fmt.Errorf(
			"failed to get offset from query params, %w",
			err,
		)
	}

	return limit, offset, nil
}
