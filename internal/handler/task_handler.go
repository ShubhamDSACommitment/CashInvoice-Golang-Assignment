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
