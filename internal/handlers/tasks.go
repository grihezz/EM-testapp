package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"strconv"
	"tt/internal/core"
)

type TaskService interface {
	GetAllTasks(ctx context.Context) (core.Tasks, error)
	AddTask(ctx context.Context, task core.ServiceTask) (int, error)
	GetTaskByUserID(ctx context.Context, userID int) (core.Tasks, error)
}

type TaskHandler struct {
	TaskService TaskService
	Logger      slog.Logger
	Client      *http.Client
}

func NewTaskHandler(taskService TaskService, logger *slog.Logger, client *http.Client) *TaskHandler {
	return &TaskHandler{
		TaskService: taskService,
		Logger:      *logger,
		Client:      client,
	}
}

func (th *TaskHandler) GetAllTasksHandler(c *gin.Context) {
	reqID := c.Request.Context().Value("requestID")
	reqIDString := fmt.Sprintf("requestID: %s", reqID)

	c.Set("requestID", reqIDString)
	tasks, err := th.TaskService.GetAllTasks(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (th *TaskHandler) AddTaskHandler(c *gin.Context) {
	var req core.ServiceTask

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task := core.ServiceTask{
		UserID:      req.UserID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		StopTime:    req.StopTime,
	}

	taskID, err := th.TaskService.AddTask(c.Request.Context(), task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, core.AddTaskResponse{TaskID: taskID})
}

func (th *TaskHandler) GetTaskByUserIDHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	task, err := th.TaskService.GetTaskByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, task)
}
