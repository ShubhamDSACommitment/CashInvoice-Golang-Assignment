package repository

import (
	"database/sql"
	"errors"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	Create(user *models.User) error
}

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) *MySQLUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) GetByEmail(email string) (*models.User, error) {
	query := `
        SELECT id, email, password, role
        FROM users
        WHERE email = ?
    `

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Role,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}

	return &user, err
}

func (r *MySQLUserRepository) Create(user *models.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, email, password, role)
         VALUES (?, ?, ?, ?)`,
		user.ID,
		user.Email,
		user.Password,
		user.Role,
	)
	return err
}
