package service

import (
	"fmt"

	"github.com/Mukam21/go-task-api/logger"
	"github.com/Mukam21/go-task-api/model"
	"github.com/Mukam21/go-task-api/repository"
)

type TaskService struct {
	repo   *repository.TaskRepository
	logger *logger.Logger
}

func NewTaskService(repo *repository.TaskRepository, logger *logger.Logger) *TaskService {
	return &TaskService{repo: repo, logger: logger}
}

func (s *TaskService) CreateTask(t *model.Task) *model.Task {
	res := s.repo.Create(t)
	s.logger.Log(fmt.Sprintf("[CreateTask] id=%d title=%q status=%s", res.ID, res.Title, res.Status))
	return res
}

func (s *TaskService) GetTasks(status string) []*model.Task {
	res := s.repo.GetAll(status)
	s.logger.Log(fmt.Sprintf("[GetTasks] count=%d status=%s", len(res), status))
	return res
}

func (s *TaskService) GetTask(id int) (*model.Task, error) {
	t, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Log(fmt.Sprintf("[GetTask] not_found id=%d", id))
		return nil, err
	}
	s.logger.Log(fmt.Sprintf("[GetTask] id=%d", id))
	return t, nil
}
