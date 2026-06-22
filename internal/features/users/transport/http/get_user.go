package users_transport_http

import (
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type GetUserResponse UserDTO

// GetUser godoc
// @Summary Получение пользователя
// @Description Получение пользователя из системы
// @Tags users
// @Produce json
// @Param id path int true "id получаемого пользователя"
// @Success 200 {object} GetUserResponse "Тело успешно полученного пользователя"
// @Failure 404 {object} core_http_response.ErrResponse "User with id not found"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server error"
// @Router /users/{id} [get]
func (h *UsersHttpHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	id, err := core_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	userDomain, err := h.usersService.GetUser(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	response := GetUserResponse(DTOFromDomain(userDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
