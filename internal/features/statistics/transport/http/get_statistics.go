package statistics_transport_http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
	core_logger "github.com/roman-styazhkin/golang-todoapp/internal/core/logger"
	core_http_response "github.com/roman-styazhkin/golang-todoapp/internal/core/transport/http/response"
	core_utils "github.com/roman-styazhkin/golang-todoapp/internal/core/utils"
)

type StatisticsResponse struct {
	TasksCreated               int      `json:"tasks_created"`
	TasksCompleted             int      `json:"tasks_completed"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time"`
}

func (h *StatisticsHttpHandler) GetStatistics(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(rw, log)

	userID, from, to, err := getUserIDFromToQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse("failed to get params", err)
		return
	}

	statisticsDomain, err := h.statisticsService.GetStatistics(ctx, userID, from, to)
	if err != nil {
		responseHandler.ErrorResponse("failed to get statistics", err)
		return
	}

	response := dtoFromDomain(statisticsDomain)
	responseHandler.JSONResponse(response, http.StatusOK)
}

func dtoFromDomain(domain domain.Statistics) StatisticsResponse {
	var avgTime *string

	if domain.TasksAverageCompletionTime != nil {
		duration := domain.TasksAverageCompletionTime.String()
		avgTime = &duration
	}

	return StatisticsResponse{
		TasksCompleted:             domain.TasksCompleted,
		TasksCreated:               domain.TasksCreated,
		TasksCompletedRate:         domain.TasksCompletedRate,
		TasksAverageCompletionTime: avgTime,
	}
}

func getUserIDFromToQueryParams(r *http.Request) (*int, *time.Time, *time.Time, error) {
	const (
		queryUserID = "user_id"
		queryFrom   = "from"
		queryTo     = "to"
	)

	userID, err := core_utils.GetIntQueryParam(r, queryUserID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get 'user_id' query param, %w",
			err,
		)
	}

	from, err := core_utils.GetDateQueryParam(r, queryFrom)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get 'from' query param, %w",
			err,
		)
	}

	to, err := core_utils.GetDateQueryParam(r, queryTo)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(
			"failed to get 'to' query param, %w",
			err,
		)
	}

	return userID, from, to, nil
}
