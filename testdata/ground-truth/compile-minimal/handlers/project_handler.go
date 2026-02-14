package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/test/app/internal/generated/repo"
	"github.com/test/app/internal/types/models"
	"github.com/test/app/internal/types/structures"
)

// ProjectHandler handles HTTP requests for Project resources.
type ProjectHandler struct {
	repo repo.ProjectRepository
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(r repo.ProjectRepository) *ProjectHandler {
	return &ProjectHandler{repo: r}
}

// List handles GET /projects
func (h *ProjectHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var organizationID *string
	if v := c.Query("organization_id"); v != "" {
		organizationID = &v
	}

	items, err := h.repo.GetAll(ctx, organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch projects"})
		return
	}

	responses := make([]structures.ProjectResponse, len(items))
	for i, item := range items {
		responses[i] = toProjectResponse(item)
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetByID handles GET /projects/:id
func (h *ProjectHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	item, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toProjectResponse(*item)})
}

// GetByCode handles GET /projects/by-code/:code
func (h *ProjectHandler) GetByCode(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Code is required"})
		return
	}

	item, err := h.repo.GetByCode(ctx, code)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch project"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toProjectResponse(*item)})
}

// Create handles POST /projects
func (h *ProjectHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.repo.Create(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toProjectResponse(*created)})
}

// Update handles PUT /projects/:id
func (h *ProjectHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	var input models.Project
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repo.Update(ctx, id, &input)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toProjectResponse(*updated)})
}

// Delete handles DELETE /projects/:id
func (h *ProjectHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project ID is required"})
		return
	}

	err := h.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project"})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// RegisterRoutes registers all CRUD routes for Project.
func (h *ProjectHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/projects", h.List)
	rg.POST("/projects", h.Create)
	rg.GET("/projects/:id", h.GetByID)
	rg.PUT("/projects/:id", h.Update)
	rg.DELETE("/projects/:id", h.Delete)
	rg.GET("/projects/by-code/:code", h.GetByCode)
}

// toProjectResponse converts a Project model to a ProjectResponse DTO.
func toProjectResponse(m models.Project) structures.ProjectResponse {
	return structures.ProjectResponse{
		Code:           m.Code,
		Description:    m.Description,
		ID:             m.ID,
		Name:           m.Name,
		OrganizationID: m.OrganizationID,
	}
}
