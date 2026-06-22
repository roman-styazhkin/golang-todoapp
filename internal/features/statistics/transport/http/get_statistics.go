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
	TasksCreated               int      `json:"tasks_created" example:"10"`
	TasksCompleted             int      `json:"tasks_completed" example:"5"`
	TasksCompletedRate         *float64 `json:"tasks_completed_rate" example:"50"`
	TasksAverageCompletionTime *string  `json:"tasks_average_completion_time" example:"1m30s"`
}

// GetStatistics godoc
// @Summary Получение статистики
// @Description Получение статистики по задачам
// @Tags statistics
// @Produce json
// @Param user_id query int false "id пользователя с задачами"
// @Param from query string false "Начало промежутка рассмотрения статистики"
// @Param to query string false "Конец промежутка рассмотрения статистики"
// @Success 200 {object} StatisticsResponse "Статистика по задачам"
// @Failure 400 {object} core_http_response.ErrResponse "Bad request"
// @Failure 500 {object} core_http_response.ErrResponse "Internal server err"
// @Router /statistics [get]
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
