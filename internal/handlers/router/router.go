package router

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tt/internal/handlers"
	middlewares "tt/internal/middleware"
	"tt/internal/service"
)

func SetupRouter(userSrv *service.UserService, logger *slog.Logger) *gin.Engine {
	r := gin.Default()

	r.Use(middlewares.AccessLog(logger))
	r.Use(middlewares.GzipMiddleware())
	r.Use(middlewares.ErrorHandler())
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.RateLimitingMiddleware())

	uh := handlers.NewUserHandler(userSrv, logger, &http.Client{})

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Gin framework!",
		})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/api/users", func(c *gin.Context) {
		uh.GetAllUsersHandler(c)
	})

	r.POST("/api/users", func(c *gin.Context) {
		uh.AddUserHandler(c)
	})

	r.POST("/api/users/search", func(c *gin.Context) {
		uh.GetUserByNumberHandler(c)
	})

	return r
}
