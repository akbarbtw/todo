package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"todo-list/internal/model"

	"github.com/google/uuid"
)

type TaskRepository interface {
	Create(ctx context.Context, task *model.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetAll(ctx context.Context, filters TaskFilters) ([]*model.Task, error)
	Update(ctx context.Context, task *model.Task) error
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status) error
}

type TaskFilters struct {
	Status   *model.Status
	Priority *model.Priority
	DueDate  *string
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks (id, title, description, priority, status, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Priority, task.Status,
		task.DueDate, task.CreatedAt, task.UpdatedAt,
	)
	return err
}

func (r *taskRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	query := `
		SELECT id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var task model.Task
	err := row.Scan(
		&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status,
		&task.DueDate, &task.CreatedAt, &task.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("task not found")
	}
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) GetAll(ctx context.Context, filters TaskFilters) ([]*model.Task, error) {
	query := `
		SELECT id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks
	`
	args := []interface{}{}
	conditions := []string{}

	if filters.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", len(args)+1))
		args = append(args, *filters.Status)
	}
	if filters.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", len(args)+1))
		args = append(args, *filters.Priority)
	}
	if filters.DueDate != nil {
		conditions = append(conditions, fmt.Sprintf("due_date::date = $%d::date", len(args)+1))
		args = append(args, *filters.DueDate)
	}

	if len(conditions) > 0 {
		query += " WHERE " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}

	query += " ORDER BY created_at DESC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Priority, &task.Status,
			&task.DueDate, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}

func (r *taskRepository) Update(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE tasks
		SET title = $2, description = $3, priority = $4, status = $5, due_date = $6, updated_at = $7
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Priority, task.Status,
		task.DueDate, task.UpdatedAt,
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (r *taskRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func (r *taskRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status) error {
	query := `UPDATE tasks SET status = $2, updated_at = $3 WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id, status, time.Now())
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}
