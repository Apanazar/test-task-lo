package services

import (
	"context"
	"errors"
	"test-task-lo/internal/logger"
	"test-task-lo/internal/models"
	"time"
)

type TaskRepository interface {
	Create(ctx context.Context, t models.Task) (models.Task, error)
	GetByID(ctx context.Context, id int) (models.Task, bool)
	List(ctx context.Context, status *models.TaskStatus) ([]models.Task, error)
}

type TaskService interface {
	Create(ctx context.Context, title string, status models.TaskStatus) (models.Task, error)
	GetByID(ctx context.Context, id int) (models.Task, bool)
	List(ctx context.Context, status *models.TaskStatus) ([]models.Task, error)
}

var (
	ErrEmptyTitle    = errors.New("title is required")
	ErrInvalidStatus = errors.New("invalid status")
)

type taskService struct {
	repo   TaskRepository
	logger *logger.Logger
}

func NewTaskService(repo TaskRepository, log *logger.Logger) TaskService {
	return &taskService{repo: repo, logger: log}
}

func (s *taskService) Create(ctx context.Context, title string, status models.TaskStatus) (models.Task, error) {
	if title == "" {
		return models.Task{}, ErrEmptyTitle
	}
	if !status.IsValid() {
		return models.Task{}, ErrInvalidStatus
	}

	t := models.NewTask(title, status)
	t.UpdatedAt = t.CreatedAt
	out, err := s.repo.Create(ctx, *t)

	if err != nil {
		if s.logger != nil {
			s.logger.Error("service.Create failed", map[string]interface{}{"error": err.Error()})
		}
		return models.Task{}, err
	}
	if s.logger != nil {
		s.logger.Info("service.Create ok", map[string]interface{}{"id": out.ID, "title": out.Title})
	}
	return out, nil
}

func (s *taskService) GetByID(ctx context.Context, id int) (models.Task, bool) {
	t, ok := s.repo.GetByID(ctx, id)
	return t, ok
}

func (s *taskService) List(ctx context.Context, status *models.TaskStatus) ([]models.Task, error) {
	ts, err := s.repo.List(ctx, status)
	if err != nil {
		return nil, err
	}
	_ = time.Now()
	return ts, nil
}
