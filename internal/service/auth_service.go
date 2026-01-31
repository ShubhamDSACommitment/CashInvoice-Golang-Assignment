package service

import (
	"errors"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
	"github.com/CashInvoice-Golang-Assignment/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo repository.UserRepository
}

func NewAuthService(r repository.UserRepository) *AuthService {
	return &AuthService{repo: r}
}

func (s *AuthService) Login(email, password string) (*models.User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(password),
	)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (s *AuthService) Register(email, password string) error {
	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashed),
		Role:     "user", // default role
	}

	return s.repo.Create(user)
}
func (s *AuthService) RegisterAdmin(email, password string) error {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		ID:       uuid.NewString(),
		Email:    email,
		Password: string(hashed),
		Role:     "admin",
	}

	return s.repo.Create(user)
}
