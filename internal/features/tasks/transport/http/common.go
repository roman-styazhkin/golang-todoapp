package tasks_transport_http

import (
	"time"

	"github.com/roman-styazhkin/golang-todoapp/internal/core/domain"
)

type TaskDTO struct {
	ID           int        `json:"id"`
	Version      int        `json:"version"`
	Title        string     `json:"title"`
	Description  *string    `json:"description"`
	Completed    bool       `json:"completed"`
	CreatedAt    time.Time  `json:"create_at"`
	CompletedAt  *time.Time `json:"completed_at"`
	AuthorUserID int        `json:"author_user_id"`
}

func NewTaskDTO(
	id int,
	version int,
	title string,
	description *string,
	completed bool,
	createdAt time.Time,
	completedAt *time.Time,
	authorUserID int,
) TaskDTO {
	return TaskDTO{
		ID:           id,
		Version:      version,
		Title:        title,
		Description:  description,
		Completed:    completed,
		CreatedAt:    createdAt,
		CompletedAt:  completedAt,
		AuthorUserID: authorUserID,
	}
}

func DTOFromDomain(domain domain.Task) TaskDTO {
	return NewTaskDTO(
		domain.ID,
		domain.Version,
		domain.Title,
		domain.Description,
		domain.Completed,
		domain.CreatedAt,
		domain.CompletedAt,
		domain.AuthorUserID,
	)
}

func DTOListFromDomains(domains []domain.Task) []TaskDTO {
	dtoList := make([]TaskDTO, len(domains))

	for i, domain := range domains {
		dtoList[i] = DTOFromDomain(domain)
	}

	return dtoList
}
