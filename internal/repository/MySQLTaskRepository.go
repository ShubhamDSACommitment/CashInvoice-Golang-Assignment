package repository

import (
	"database/sql"
	"errors"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
)

type MySQLTaskRepository struct {
	db *sql.DB
}

// Constructor
func NewMySQLTaskRepository(db *sql.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

// Compile-time check
var _ TaskRepository = (*MySQLTaskRepository)(nil)

// -------------------- CREATE --------------------

func (r *MySQLTaskRepository) Create(task *models.Task) error {
	query := `
        INSERT INTO tasks (
            id,
            title,
            description,
            status,
            user_id,
            created_at,
            updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?)
    `

	_, err := r.db.Exec(
		query,
		task.ID,
		task.Title,
		task.Description,
		task.Status,
		task.UserID,
		task.CreatedAt,
		task.UpdatedAt,
	)

	return err
}

// -------------------- GET BY ID --------------------

func (r *MySQLTaskRepository) GetByID(id string) (*models.Task, error) {
	query := `
        SELECT 
            id,
            title,
            description,
            status,
            user_id,
            created_at,
            updated_at
        FROM tasks
        WHERE id = ?
    `

	var task models.Task

	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserID,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("task not found")
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

// -------------------- GET ALL --------------------

func (r *MySQLTaskRepository) GetAll(userID string, isAdmin bool) ([]models.Task, error) {
	var (
		rows *sql.Rows
		err  error
	)

	if isAdmin {
		rows, err = r.db.Query(`
            SELECT 
                id,
                title,
                description,
                status,
                user_id,
                created_at,
                updated_at
            FROM tasks
            ORDER BY created_at DESC
        `)
	} else {
		rows, err = r.db.Query(`
            SELECT 
                id,
                title,
                description,
                status,
                user_id,
                created_at,
                updated_at
            FROM tasks
            WHERE user_id = ?
            ORDER BY created_at DESC
        `, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := []models.Task{}

	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserID,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// -------------------- DELETE --------------------

func (r *MySQLTaskRepository) Delete(id string) error {
	result, err := r.db.Exec(
		"DELETE FROM tasks WHERE id = ?",
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}

// -------------------- UPDATE STATUS --------------------

func (r *MySQLTaskRepository) UpdateStatus(id string, status string) error {
	result, err := r.db.Exec(
		`
        UPDATE tasks
        SET status = ?, updated_at = NOW()
        WHERE id = ?
        `,
		status,
		id,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("task not found")
	}

	return nil
}
