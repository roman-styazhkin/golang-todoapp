package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type GetTaskResponse TaskDTO

// GetTask godoc
// @Summary Получение задачи
// @Description Получение существующей в системе задачи
// @Tags tasks
// @Produce json
// @Param id path int true "id получаемой задачи"
// @Success 200 {object} GetTaskResponse "Тело получаемой задачи"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrResponse "Task not found"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server err"
// @Router /tasks/{id} [get]
func (h *TasksHttpHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	id, err := core_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, id)
	if err != nil {
		responseHandler.ErrorResponse("failed to get task", err)
		return
	}

	response := GetTaskResponse(DTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}
