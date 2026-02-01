package service

import (
	"context"
	"tasks-api/models"
	"tasks-api/repo"
)

type Service struct {
	r *repo.Repo
}

func New(r *repo.Repo) *Service {
	return &Service{
		r: r,
	}
}

func (s *Service) CreateTask(ctx context.Context, task models.CreateTask) error {
	return s.r.Create(ctx, task)
}

func (s *Service) GetTasks(ctx context.Context) ([]models.TaskResponse, error) {
	return s.r.RequestTasks(ctx)
}

func (s *Service) GetTask(ctx context.Context, taskId string) (models.TaskResponse, error) {
	return s.r.RequestTask(ctx, taskId)
}

func (s *Service) DeleteTask(ctx context.Context, taskId string) error {
	return s.r.Delete(ctx, taskId)
}

func (s *Service) ChangeTask(ctx context.Context, task models.CreateTask, taskId string) error {
	return s.r.Change(ctx, task, taskId)
}
