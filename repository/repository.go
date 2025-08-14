package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/Mukam21/go-task-api/model"
)

type TaskRepository struct {
	mu     sync.RWMutex
	tasks  map[int]*model.Task
	lastID int
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{
		tasks: make(map[int]*model.Task),
	}
}

func (r *TaskRepository) Create(task *model.Task) *model.Task {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lastID++
	task.ID = r.lastID
	task.CreatedAt = time.Now().UTC()
	task.UpdatedAt = task.CreatedAt
	r.tasks[task.ID] = task
	return task
}

func (r *TaskRepository) GetAll(status string) []*model.Task {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []*model.Task
	for _, t := range r.tasks {
		if status == "" || t.Status == status {
			result = append(result, t)
		}
	}
	return result
}

func (r *TaskRepository) GetByID(id int) (*model.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if t, ok := r.tasks[id]; ok {
		return t, nil
	}
	return nil, errors.New("task not found")
}
