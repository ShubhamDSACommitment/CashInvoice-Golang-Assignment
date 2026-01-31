package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/CashInvoice-Golang-Assignment/internal/models"
)

type MySQLTaskRepository struct {
	db *sql.DB
}

func NewMySQLTaskRepository(db *sql.DB) *MySQLTaskRepository {
	return &MySQLTaskRepository{db: db}
}

// Compile-time check
var _ TaskRepository = (*MySQLTaskRepository)(nil)

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
	var status string
	var createdAtStr, updatedAtStr string

	err := r.db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&status,
		&task.UserID,
		&createdAtStr,
		&updatedAtStr,
	)
	task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
	task.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("task not found")
	}

	if err != nil {
		return nil, err
	}

	return &task, nil
}

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
		var status string
		var createdAtStr, updatedAtStr string

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&status,
			&task.UserID,
			&createdAtStr,
			&updatedAtStr,
		)
		if err != nil {
			return nil, err
		}
		task.CreatedAt, _ = time.Parse("2006-01-02 15:04:05", createdAtStr)
		task.UpdatedAt, _ = time.Parse("2006-01-02 15:04:05", updatedAtStr)
		task.Status = models.TaskStatus(status)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

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
func (r *MySQLTaskRepository) AutoCompleteIfPending(id string) error {
	result, err := r.db.Exec(`
        UPDATE tasks
        SET status = 'completed', updated_at = NOW()
        WHERE id = ?
          AND status IN ('pending', 'in_progress')
    `, id)

	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()

	if rows == 0 {
		return nil
	}

	return nil
}
