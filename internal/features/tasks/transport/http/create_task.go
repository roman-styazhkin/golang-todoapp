package tasks_transport_http

import (
	"net/http"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_request "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
)

type CreateTaskRequest struct {
	Title        string  `json:"title" validate:"required,min=1,max=100" example:"Сделать домашку"`
	Description  *string `json:"description" validate:"omitempty,min=1,max=1000" example:"Сделать домашку по математике"`
	AuthorUserID int     `json:"author_user_id" validate:"required" example:"1"`
}

type CreateTaskResponse TaskDTO

// CreateTask godoc
// @Summary Создание новой задачи
// @Description Создание новой задачи в системе
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body CreateTaskRequest true "Тело созданной задачи"
// @Success 201 {object} CreateTaskResponse "Успешно созданная задача"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrResponse "Author not found"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server error"
// @Router /tasks [post]
func (h *TasksHttpHandler) CreateTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	var request CreateTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse("failed to decode and validate request", err)
		return
	}

	taskDomain := domainFromRequest(request)
	taskDomain, err := h.tasksService.CreateTask(ctx, taskDomain)
	if err != nil {
		responseHandler.ErrorResponse("failed to create task", err)
		return
	}

	response := CreateTaskResponse(DTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func domainFromRequest(request CreateTaskRequest) domain.Task {
	return domain.NewTaskUninitialized(
		request.Title,
		request.Description,
		request.AuthorUserID,
	)
}
