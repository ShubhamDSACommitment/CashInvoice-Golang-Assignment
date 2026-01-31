package repository

import "github.com/CashInvoice-Golang-Assignment/internal/models"

type TaskRepository interface {
	Create(task *models.Task) error
	GetByID(id string) (*models.Task, error)
	GetAll(userID string, isAdmin bool) ([]models.Task, error)
	Delete(id string) error
	UpdateStatus(id string, status string) error
	AutoCompleteIfPending(id string) error
}
