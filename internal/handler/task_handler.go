package handler

import (
	"net/http"
	"time"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
	"github.com/CashInvoice-Golang-Assignment/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{service: s}
}

type CreateTaskRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
}

// Create CreateTask godoc
// @Summary Create a task
// @Tags Tasks
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param task body CreateTaskRequest true "Task payload"
// @Success 201 {object} models.Task
// @Failure 401 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	var req CreateTaskRequest

	// Validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "title is required",
		})
		return
	}

	userID := c.GetString("user_id") // from JWT middleware

	task := &models.Task{
		ID:          uuid.NewString(),
		Title:       req.Title,
		Description: req.Description,
		Status:      models.StatusPending,
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.service.CreateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to create task",
		})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// GetAllTask GetAllTasks godoc
// @Summary Get all tasks
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tasks [get]
func (h *TaskHandler) GetAllTask(c *gin.Context) {
	// Extract auth context (set by JWT middleware)
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "unauthorized",
		})
		return
	}

	// Call service layer
	tasks, err := h.service.GetAllTasks(userID, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch tasks",
		})
		return
	}

	// Success
	c.JSON(http.StatusOK, gin.H{
		"count": len(tasks),
		"tasks": tasks,
	})
}

// GetByID GetTaskByID godoc
// @Summary Get task by ID
// @Tags Tasks
// @Produce json
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task id required"})
		return
	}

	task, err := h.service.GetTaskByID(taskID, userID, role)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Delete DeleteTask godoc
// @Summary Delete task
// @Tags Tasks
// @Security BearerAuth
// @Param id path string true "Task ID"
// @Success 204
// @Failure 403 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	taskID := c.Param("id")
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if taskID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "task id required"})
		return
	}

	err := h.service.DeleteTask(taskID, userID, role)
	if err != nil {
		if err.Error() == "forbidden" {
			c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "task deleted successfully",
	})
}
