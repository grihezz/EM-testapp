package router

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tt/internal/handlers"
	middlewares "tt/internal/middleware"
	"tt/internal/service"
)

func SetupRouter(userSrv *service.UserService, taskSrv *service.TaskService, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.AccessLog(logger))
	r.Use(middlewares.GzipMiddleware())
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimitingMiddleware())

	uh := handlers.NewUserHandler(userSrv, logger, &http.Client{})
	th := handlers.NewTaskHandler(taskSrv, logger, &http.Client{})

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.GET("/", func(c *gin.Context) { uh.GetAllUsersHandler(c) })
			users.POST("/", func(c *gin.Context) { uh.AddUserHandler(c) })
			users.POST("/search", func(c *gin.Context) { uh.GetUserByNumberHandler(c) })
			users.PATCH("/:user_id", func(c *gin.Context) { uh.UpdateUserHandler(c) })
			users.DELETE("/:user_id", func(c *gin.Context) { uh.DeleteUserHandler(c) })
		}
		tasks := api.Group("/tasks")
		{
			tasks.GET("/", func(c *gin.Context) { th.GetAllTasksHandler(c) })
			tasks.POST("/", func(c *gin.Context) { th.AddTaskHandler(c) })
			tasks.GET("/:user_id", func(c *gin.Context) { th.GetTaskByUserIDHandler(c) })
		}
	}

	return r
}
