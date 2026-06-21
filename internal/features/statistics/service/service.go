package statistics_service

import (
	"context"
	"time"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

type StatisticsService struct {
	staticsRepository StatisticsRepository
}

type StatisticsRepository interface {
	GetTasks(
		ctx context.Context,
		userID *int,
		from *time.Time,
		to *time.Time,
	) ([]domain.Task, error)
}

func NewStatisticsService(staticsRepository StatisticsRepository) *StatisticsService {
	return &StatisticsService{
		staticsRepository: staticsRepository,
	}
}
