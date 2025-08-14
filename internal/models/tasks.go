package models

import (
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
)

func (s TaskStatus) IsValid() bool {
	switch s {
	case StatusPending, StatusInProgress, StatusCompleted:
		return true
	default:
		return false
	}
}

type Task struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewTask(title string, status TaskStatus) *Task {
	if !status.IsValid() {
		status = StatusPending
	}

	now := time.Now()
	return &Task{
		ID:        generateID(),
		Title:     title,
		Status:    status,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func generateID() int {
	return int(time.Now().UnixNano())
}
