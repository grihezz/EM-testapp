package service

import (
	"context"
	"fmt"
	"tt/internal/core"
)

type TaskRepo interface {
	GetAllTasks(ctx context.Context) (core.Tasks, error)
	AddTask(ctx context.Context, task core.ServiceTask) (int, error)
	GetTaskByUserID(ctx context.Context, userID int) (core.Tasks, error)
}

type TaskService struct {
	TaskRepository TaskRepo
}

func NewTaskService(repo TaskRepo) *TaskService {
	return &TaskService{TaskRepository: repo}
}

func (ts *TaskService) GetAllTasks(ctx context.Context) (core.Tasks, error) {
	const op = "service.GetAllTasks"
	tasks, err := ts.TaskRepository.GetAllTasks(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (ts *TaskService) AddTask(ctx context.Context, task core.ServiceTask) (int, error) {
	const op = "service.AddTask"
	taskID, err := ts.TaskRepository.AddTask(ctx, task)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return taskID, nil
}

func (ts *TaskService) GetTaskByUserID(ctx context.Context, userID int) (core.Tasks, error) {
	const op = "service.GetTaskByID"
	task, err := ts.TaskRepository.GetTaskByUserID(ctx, userID)
	if err != nil {
		return core.Tasks{}, fmt.Errorf("%s: %w", op, err)
	}
	return task, nil
}
