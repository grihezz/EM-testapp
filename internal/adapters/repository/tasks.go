package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"tt/internal/core"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrTaskExists   = errors.New("task with the same title already exists")
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (tr *TaskRepository) GetAllTasks(ctx context.Context) (core.Tasks, error) {
	const op = "repository.GetAllTasks"
	baseQuery := "SELECT task_id, user_id, title, description, start_time, stop_time FROM tasks ORDER BY task_id;"

	rows, err := tr.db.Query(ctx, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks core.Tasks
	for rows.Next() {
		var task core.Task
		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.StartTime,
			&task.StopTime,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return tasks, nil

}
func (tr *TaskRepository) AddTask(ctx context.Context, task core.ServiceTask) (int, error) {
	const op = "repository.AddTask"

	// Проверка на существование пользователя
	var userExists bool
	userQuery := "SELECT EXISTS (SELECT 1 FROM users WHERE user_id=$1)"
	err := tr.db.QueryRow(ctx, userQuery, task.UserID).Scan(&userExists)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if !userExists {
		return 0, fmt.Errorf("%s: user with id %d does not exist", op, task.UserID)
	}

	// Проверка на существование задачи с таким же названием
	var taskExists bool
	taskQuery := "SELECT EXISTS (SELECT 1 FROM tasks WHERE title=$1)"
	err = tr.db.QueryRow(ctx, taskQuery, task.Title).Scan(&taskExists)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	if taskExists {
		return 0, ErrTaskExists
	}

	// Вставка новой задачи
	baseQuery := "INSERT INTO tasks (user_id, title, description, start_time, stop_time) VALUES ($1, $2, $3, $4, $5) RETURNING task_id;"
	var taskID int
	err = tr.db.QueryRow(ctx, baseQuery, task.UserID, task.Title, task.Description, task.StartTime, task.StopTime).Scan(&taskID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return taskID, nil
}

func (tr *TaskRepository) GetTaskByUserID(ctx context.Context, userID int) (core.Tasks, error) {
	const op = "repository.GetTaskByUserID"
	baseQuery := "SELECT task_id, user_id, title, description, start_time, stop_time FROM tasks WHERE user_id = $1;"
	var tasks core.Tasks

	rows, err := tr.db.Query(ctx, baseQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task core.Task

		err := rows.Scan(
			&task.ID,
			&task.UserID,
			&task.Title,
			&task.Description,
			&task.StartTime,
			&task.StopTime)

		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}
	return tasks, nil
}
