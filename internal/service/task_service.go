package service

import (
	"context"
	"errors"
	"time"

	"todo-list/internal/model"
	"todo-list/internal/repository"

	"github.com/google/uuid"
)

type TaskService interface {
	Create(ctx context.Context, req CreateTaskRequest) (*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	GetAll(ctx context.Context, filters repository.TaskFilters) ([]*model.Task, error)
	Update(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status) (*model.Task, error)
}

type CreateTaskRequest struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	Priority    model.Priority `json:"priority"`
	Status      model.Status   `json:"status"`
	DueDate     *time.Time     `json:"due_date,omitempty"`
}

type UpdateTaskRequest struct {
	Title       *string         `json:"title,omitempty"`
	Description *string         `json:"description,omitempty"`
	Priority    *model.Priority `json:"priority,omitempty"`
	Status      *model.Status   `json:"status,omitempty"`
	DueDate     *time.Time      `json:"due_date,omitempty"`
}

type taskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(ctx context.Context, req CreateTaskRequest) (*model.Task, error) {
	if req.Title == "" {
		return nil, errors.New("title is required")
	}
	if !req.Priority.IsValid() {
		return nil, errors.New("invalid priority: must be low, medium, or high")
	}
	if !req.Status.IsValid() {
		return nil, errors.New("invalid status: must be todo, in_progress, or completed")
	}

	now := time.Now()
	task := &model.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		Status:      req.Status,
		DueDate:     req.DueDate,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) GetAll(ctx context.Context, filters repository.TaskFilters) ([]*model.Task, error) {
	return s.repo.GetAll(ctx, filters)
}

func (s *taskService) Update(ctx context.Context, id uuid.UUID, req UpdateTaskRequest) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		if *req.Title == "" {
			return nil, errors.New("title cannot be empty")
		}
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Priority != nil {
		if !req.Priority.IsValid() {
			return nil, errors.New("invalid priority: must be low, medium, or high")
		}
		task.Priority = *req.Priority
	}
	if req.Status != nil {
		if !req.Status.IsValid() {
			return nil, errors.New("invalid status: must be todo, in_progress, or completed")
		}
		task.Status = *req.Status
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	task.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *taskService) UpdateStatus(ctx context.Context, id uuid.UUID, status model.Status) (*model.Task, error) {
	if !status.IsValid() {
		return nil, errors.New("invalid status: must be todo, in_progress, or completed")
	}

	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	task.Status = status
	task.UpdatedAt = time.Now()

	if err := s.repo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return task, nil
}
