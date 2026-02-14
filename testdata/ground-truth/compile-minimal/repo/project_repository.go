package repo

import (
	"context"

	"github.com/test/app/internal/types/models"
)

// ProjectRepository defines the data access contract for Project.
type ProjectRepository interface {
	GetAll(ctx context.Context, organizationID *string) ([]models.Project, error)
	GetByID(ctx context.Context, id string) (*models.Project, error)
	GetByCode(ctx context.Context, code string) (*models.Project, error)
	Create(ctx context.Context, input *models.Project) (*models.Project, error)
	Update(ctx context.Context, id string, input *models.Project) (*models.Project, error)
	Delete(ctx context.Context, id string) error
}
