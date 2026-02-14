package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/test/app/internal/generated/repo"
	"github.com/test/app/internal/types/models"
	"github.com/test/app/internal/types/structures"
)

// OrganizationHandler handles HTTP requests for Organization resources.
type OrganizationHandler struct {
	repo repo.OrganizationRepository
}

// NewOrganizationHandler creates a new OrganizationHandler.
func NewOrganizationHandler(r repo.OrganizationRepository) *OrganizationHandler {
	return &OrganizationHandler{repo: r}
}

// List handles GET /organizations
func (h *OrganizationHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	items, err := h.repo.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	responses := make([]structures.OrganizationResponse, len(items))
	for i, item := range items {
		responses[i] = toOrganizationResponse(item)
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetByID handles GET /organizations/:id
func (h *OrganizationHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	item, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toOrganizationResponse(*item)})
}

// GetByCode handles GET /organizations/by-code/:code
func (h *OrganizationHandler) GetByCode(c *gin.Context) {
	ctx := c.Request.Context()
	code := c.Param("code")

	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization Code is required"})
		return
	}

	item, err := h.repo.GetByCode(ctx, code)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toOrganizationResponse(*item)})
}

// Create handles POST /organizations
func (h *OrganizationHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var input models.Organization
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.repo.Create(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toOrganizationResponse(*created)})
}

// Update handles PUT /organizations/:id
func (h *OrganizationHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	var input models.Organization
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repo.Update(ctx, id, &input)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toOrganizationResponse(*updated)})
}

// Delete handles DELETE /organizations/:id
func (h *OrganizationHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Organization ID is required"})
		return
	}

	err := h.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Organization not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// RegisterRoutes registers all CRUD routes for Organization.
func (h *OrganizationHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/organizations", h.List)
	rg.POST("/organizations", h.Create)
	rg.GET("/organizations/:id", h.GetByID)
	rg.PUT("/organizations/:id", h.Update)
	rg.DELETE("/organizations/:id", h.Delete)
	rg.GET("/organizations/by-code/:code", h.GetByCode)
}

// toOrganizationResponse converts a Organization model to a OrganizationResponse DTO.
func toOrganizationResponse(m models.Organization) structures.OrganizationResponse {
	return structures.OrganizationResponse{
		Code: m.Code,
		ID:   m.ID,
		Name: m.Name,
	}
}
