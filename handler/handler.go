package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Mukam21/go-task-api/model"
	"github.com/Mukam21/go-task-api/service"
)

type TaskHandler struct {
	service *service.TaskService
}

func NewTaskHandler(service *service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

func (h *TaskHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/tasks", h.handleTasks)
	mux.HandleFunc("/tasks/", h.handleTaskByID)
}

func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		status := r.URL.Query().Get("status")
		tasks := h.service.GetTasks(status)
		writeJSON(w, http.StatusOK, tasks)
	case http.MethodPost:
		var task model.Task
		if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON"})
			return
		}
		if task.Title == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
			return
		}
		if task.Status == "" {
			task.Status = "new"
		}
		created := h.service.CreateTask(&task)
		writeJSON(w, http.StatusCreated, created)
	default:
		writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func (h *TaskHandler) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid ID"})
		return
	}

	task, err := h.service.GetTask(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "task not found"})
		return
	}
	writeJSON(w, http.StatusOK, task)
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
