package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/test/app/internal/generated/repo"
	"github.com/test/app/internal/types/models"
	"github.com/test/app/internal/types/structures"
)

// TaskHandler handles HTTP requests for Task resources.
type TaskHandler struct {
	repo repo.TaskRepository
}

// NewTaskHandler creates a new TaskHandler.
func NewTaskHandler(r repo.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: r}
}

// List handles GET /tasks
func (h *TaskHandler) List(c *gin.Context) {
	ctx := c.Request.Context()

	var projectID *string
	if v := c.Query("project_id"); v != "" {
		projectID = &v
	}

	items, err := h.repo.GetAll(ctx, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tasks"})
		return
	}

	responses := make([]structures.TaskResponse, len(items))
	for i, item := range items {
		responses[i] = toTaskResponse(item)
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetByID handles GET /tasks/:id
func (h *TaskHandler) GetByID(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	item, err := h.repo.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toTaskResponse(*item)})
}

// Create handles POST /tasks
func (h *TaskHandler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := h.repo.Create(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": toTaskResponse(*created)})
}

// Update handles PUT /tasks/:id
func (h *TaskHandler) Update(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	var input models.Task
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.repo.Update(ctx, id, &input)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toTaskResponse(*updated)})
}

// Delete handles DELETE /tasks/:id
func (h *TaskHandler) Delete(c *gin.Context) {
	ctx := c.Request.Context()
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Task ID is required"})
		return
	}

	err := h.repo.Delete(ctx, id)
	if err != nil {
		if err.Error() == "no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		}
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// RegisterRoutes registers all CRUD routes for Task.
func (h *TaskHandler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/tasks", h.List)
	rg.POST("/tasks", h.Create)
	rg.GET("/tasks/:id", h.GetByID)
	rg.PUT("/tasks/:id", h.Update)
	rg.DELETE("/tasks/:id", h.Delete)
}

// toTaskResponse converts a Task model to a TaskResponse DTO.
func toTaskResponse(m models.Task) structures.TaskResponse {
	return structures.TaskResponse{
		ID:        m.ID,
		Status:    m.Status,
		Title:     m.Title,
		ProjectID: m.ProjectID,
	}
}
