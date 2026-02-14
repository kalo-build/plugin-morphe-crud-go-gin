package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/test/app/internal/generated/repo"
)

// Repos holds all repository implementations needed by the generated handlers.
type Repos struct {
	Organization repo.OrganizationRepository
	Project      repo.ProjectRepository
	Task         repo.TaskRepository
}

// RegisterAllRoutes registers all generated CRUD handlers on the given router group.
func RegisterAllRoutes(rg *gin.RouterGroup, repos Repos) {
	NewOrganizationHandler(repos.Organization).RegisterRoutes(rg)
	NewProjectHandler(repos.Project).RegisterRoutes(rg)
	NewTaskHandler(repos.Task).RegisterRoutes(rg)
}
