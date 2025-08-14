package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"test-task-lo/internal/logger"
	"test-task-lo/internal/models"
	"test-task-lo/internal/services"
)

type TaskHandlers struct {
	svc    services.TaskService
	logger *logger.Logger
}

func NewTaskHandlers(svc services.TaskService, log *logger.Logger) *TaskHandlers {
	return &TaskHandlers{svc: svc, logger: log}
}

func (h *TaskHandlers) GetAllTasks(w http.ResponseWriter, r *http.Request) {
	statusParam := r.URL.Query().Get("status")
	var filter *models.TaskStatus
	if statusParam != "" {
		s := models.TaskStatus(statusParam)
		filter = &s
	}
	tasks, err := h.svc.List(r.Context(), filter)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandlers) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid id")
		return
	}
	task, ok := h.svc.GetByID(r.Context(), id)
	if !ok {
		writeJSONError(w, http.StatusNotFound, "task not found")
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var in struct {
		Title  string            `json:"title"`
		Status models.TaskStatus `json:"status"`
	}
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&in); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	task, err := h.svc.Create(r.Context(), in.Title, in.Status)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyTitle):
			writeJSONError(w, http.StatusBadRequest, "title is required")
			return
		case errors.Is(err, services.ErrInvalidStatus):
			writeJSONError(w, http.StatusBadRequest, "invalid status")
			return
		default:
			writeJSONError(w, http.StatusInternalServerError, "failed to create task")
			return
		}
	}
	writeJSON(w, http.StatusCreated, task)
}

type apiError struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, apiError{Error: msg})
}

func writeJSON(w http.ResponseWriter, code int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
