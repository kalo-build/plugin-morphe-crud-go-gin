package repo

import (
	"context"

	"github.com/test/app/internal/types/models"
)

// TaskRepository defines the data access contract for Task.
type TaskRepository interface {
	GetAll(ctx context.Context, projectID *string) ([]models.Task, error)
	GetByID(ctx context.Context, id string) (*models.Task, error)
	Create(ctx context.Context, input *models.Task) (*models.Task, error)
	Update(ctx context.Context, id string, input *models.Task) (*models.Task, error)
	Delete(ctx context.Context, id string) error
}
