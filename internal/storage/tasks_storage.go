package storage

import (
	"context"
	"sync"
	"test-task-lo/internal/logger"
	"test-task-lo/internal/models"
)

type TaskStorage struct {
	mu     sync.RWMutex
	tasks  map[int]models.Task
	logger *logger.Logger

	nextID int
}

func NewTaskStorage(log *logger.Logger) *TaskStorage {
	return &TaskStorage{
		tasks:  make(map[int]models.Task),
		logger: log,
		nextID: 1,
	}
}

func (s *TaskStorage) Create(_ context.Context, t models.Task) (models.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextID
	s.nextID++

	t.ID = id
	s.tasks[id] = t

	if s.logger != nil {
		s.logger.Info("repo.Create", map[string]interface{}{"id": id, "title": t.Title, "status": t.Status})
	}
	out := t
	return out, nil
}

func (s *TaskStorage) GetByID(_ context.Context, id int) (models.Task, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	if !ok {
		return models.Task{}, false
	}
	return t, true
}

func (s *TaskStorage) List(_ context.Context, status *models.TaskStatus) ([]models.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]models.Task, 0, len(s.tasks))
	if status == nil {
		for _, t := range s.tasks {
			out = append(out, t)
		}
		return out, nil
	}
	for _, t := range s.tasks {
		if t.Status == *status {
			out = append(out, t)
		}
	}
	return out, nil
}
