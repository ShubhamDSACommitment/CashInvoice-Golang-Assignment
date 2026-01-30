package service

import (
	"errors"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
	"github.com/CashInvoice-Golang-Assignment/internal/repository"
)

type TaskService struct {
	repo  repository.TaskRepository
	queue chan string
}

func NewTaskService(r repository.TaskRepository, q chan string) *TaskService {
	return &TaskService{repo: r, queue: q}
}

func (s *TaskService) CreateTask(task *models.Task) error {
	if err := s.repo.Create(task); err != nil {
		return err
	}

	//// Send task ID to background worker (non-blocking if buffered channel)
	//select {
	//case s.queue <- task.ID:
	//default:
	//	// Queue full, skip auto-complete to avoid blocking API
	//}

	return nil
}
func (s *TaskService) GetAllTasks(userID string, role string) ([]models.Task, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	isAdmin := role == "admin"

	return s.repo.GetAll(userID, isAdmin)
}
