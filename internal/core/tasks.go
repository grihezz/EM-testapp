package core

import "time"

type Task struct {
	ID          int        `json:"task_id" db:"task_id"`
	UserID      int        `json:"user_id" db:"user_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	StartTime   *time.Time `json:"start_time" db:"start_time"`
	StopTime    *time.Time `json:"stop_time" db:"stop_time"`
}

type Tasks []Task

type ServiceTask struct {
	UserID      int        `json:"user_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	StopTime    *time.Time `json:"stop_time"`
}

type AddTaskResponse struct {
	TaskID int `json:"task_id"`
}
