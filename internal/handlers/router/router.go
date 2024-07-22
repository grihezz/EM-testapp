package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"tt/internal/service"
)

func SetupRouter(userSrv *service.UserService, logger *slog.Logger) *gin.Engine {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to the Gin framework!",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	router.GET("/api/users", func(c *gin.Context) {
		handleGetAllUsers(c, userSrv)
	})

	return router
}

func handleGetAllUsers(c *gin.Context, userService *service.UserService) {
	users, err := userService.GetAllUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
