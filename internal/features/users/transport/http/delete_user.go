package users_transport_http

import (
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

// DeleteUser godoc
// @Summary Удаление пользователя
// @Description Удаление существующего в системе пользователя
// @Tags users
// @Param id path int true "id удаляемого пользователя"
// @Success 204 "Успешное удаление пользователя"
// @Failure 404 {object} core_http_response.ErrResponse "Пользователь с таким id не найден"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server error"
// @Router /users/{id} [delete]
func (h *UsersHttpHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	id, err := core_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	if err = h.usersService.DeleteUser(ctx, id); err != nil {
		responseHandler.ErrorResponse("failed delete user", err)
		return
	}

	responseHandler.NoContentResponse()
}
