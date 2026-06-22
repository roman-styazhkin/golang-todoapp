package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_errors "github.com/roman-styazhkin/golang-todoapp/internal/core/errors"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_request "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/request"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_http_types "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/types"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Title"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"Description for task"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean" example:"true"`
}

func (p *PatchTaskRequest) Validate() error {
	if p.Title.Set {
		if p.Title.Value == nil {
			return fmt.Errorf(
				"failed to validate, title is required, %w",
				core_errors.ErrInvalidArgument,
			)
		}

		titleLength := len([]rune(*p.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf(
				"failed to validate, title length must be between 1 and 100, len: %d, err: %w",
				titleLength,
				core_errors.ErrInvalidArgument,
			)
		}
	}

	if p.Description.Set {
		if p.Description.Value != nil {
			descriptionLength := len([]rune(*p.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf(
					"failed to validate, description length must be between 1 and 1000, len: %d, err: %w",
					descriptionLength,
					core_errors.ErrInvalidArgument,
				)
			}
		}
	}

	if p.Completed.Set {
		if p.Completed.Value == nil {
			return fmt.Errorf(
				"failed to validate, completed cannot be null, %w",
				core_errors.ErrInvalidArgument,
			)
		}
	}

	return nil
}

type PatchTaskResponse TaskDTO

// PatchTask godoc
// @Summary Изменение задачи
// @Description Изменение существующей в системе задачи
// @Description #### Логика по изменению задачи (Three-stage-logic)
// @Description **Title не задан** поле title никак не меняется в БД
// @Description **Title явно задан как title: "value"** поле title меняется в БД на value
// @Description **Description явно задан как description: null** поле description меняется в БД на null
// @Description **Ограничения:** title и completed не могут быть null
// @Description **Ограничения:** title и description не могут быть длиной меньше 1
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path int true "id изменеямой задачи"
// @Param request body PatchTaskRequest true "PatchTask тело задачи"
// @Success 200 {object} PatchTaskResponse "PatchTask измененной задачи"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 404 {object} core_http_response.ErrResponse "Task not found"
// @Failure 409 {object} core_http_response.ErrResponse "Conflict"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server error"
// @Router /tasks/{id} [patch]
func (h *TasksHttpHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	id, err := core_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse("failed to get id from path", err)
		return
	}

	var request PatchTaskRequest
	if err = core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse("failed to decode and validate request", err)
		return
	}

	taskPatch := domainFromPatch(request)
	taskDomain, err := h.tasksService.PatchTask(ctx, id, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse("failed to patch task", err)
		return
	}

	response := PatchTaskResponse(DTOFromDomain(taskDomain))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func domainFromPatch(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
