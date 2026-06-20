package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type GetTasksResponse []TaskDTO

func (h *TasksHttpHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	userID, limit, offset, err := getUserIDLimitOffsetQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse("failed get query params", err)
		return
	}

	taskDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse("failed to get tasks", err)
		return
	}

	response := GetTasksResponse(DTOListFromDomains(taskDomains))
	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetQueryParams(r *http.Request) (*int, *int, *int, error) {
	const (
		queryUserID = "user_id"
		queryLimit  = "limit"
		queryOffset = "offset"
	)

	userID, err := core_utils.GetIntQueryParam(r, queryUserID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get user_id from query params, %w",
			err,
		)
	}

	limit, err := core_utils.GetIntQueryParam(r, queryLimit)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get limit from query params, %w",
			err,
		)
	}

	offset, err := core_utils.GetIntQueryParam(r, queryOffset)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get offset from query params, %w",
			err,
		)
	}

	return userID, limit, offset, nil
}
