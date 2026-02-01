package service

import (
	"errors"
	"log"
	"time"

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

	// Non-blocking send
	select {
	case s.queue <- task.ID:
		// enqueued
	case <-time.After(500 * time.Millisecond):
		log.Println("Auto-complete queue timeout, skipping:", task.ID)
	}
	return nil
}

func (s *TaskService) GetAllTasks(userID string, role string) ([]models.Task, error) {
	if userID == "" {
		return nil, errors.New("unauthorized")
	}

	isAdmin := role == "admin"

	return s.repo.GetAll(userID, isAdmin)
}

func (s *TaskService) GetTaskByID(taskID, userID, role string) (*models.Task, error) {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return nil, err
	}

	// Authorization: user can only access own task
	if role != "admin" && task.UserID != userID {
		return nil, errors.New("forbidden")
	}

	return task, nil
}

func (s *TaskService) DeleteTask(taskID, userID, role string) error {
	task, err := s.repo.GetByID(taskID)
	if err != nil {
		return err
	}

	// Authorization: user can only delete own task
	if role != "admin" && task.UserID != userID {
		return errors.New("forbidden")
	}

	return s.repo.Delete(taskID)
}
