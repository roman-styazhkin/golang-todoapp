package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удаление существующей в системе задачи
// @Tags tasks
// @Param id path int true "id удаляемой задачи"
// @Success 204 "Успешное удаление задачи"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrResponse "Task not found"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server err"
// @Router /tasks/{id} [delete]
func (h *TasksHttpHandler) DeleteTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	id, err := core_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	err = h.tasksService.DeleteTask(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse("failed to delete task", err)
		return
	}

	responseHandler.NoContentResponse()
}
