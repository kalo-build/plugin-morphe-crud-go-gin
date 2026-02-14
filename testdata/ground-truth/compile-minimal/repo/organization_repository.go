package repo

import (
	"context"

	"github.com/test/app/internal/types/models"
)

// OrganizationRepository defines the data access contract for Organization.
type OrganizationRepository interface {
	GetAll(ctx context.Context) ([]models.Organization, error)
	GetByID(ctx context.Context, id string) (*models.Organization, error)
	GetByCode(ctx context.Context, code string) (*models.Organization, error)
	Create(ctx context.Context, input *models.Organization) (*models.Organization, error)
	Update(ctx context.Context, id string, input *models.Organization) (*models.Organization, error)
	Delete(ctx context.Context, id string) error
}
