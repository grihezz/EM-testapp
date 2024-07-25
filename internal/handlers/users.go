package handlers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"regexp"
	"strconv"
	"tt/internal/core"
)

type UserService interface {
	GetAllUsers(ctx context.Context) (core.Users, error)
	AddUser(ctx context.Context, user core.ServiceUser) (int, error)
	UpdateUser(ctx context.Context, userID int, newData core.User) (core.User, error)
	GetUserByNumber(ctx context.Context, passportNumber string) (core.User, error)
	DeleteUser(ctx context.Context, userID int) error
}

type UserHandler struct {
	UserService UserService
	SlogLogger  slog.Logger
	Client      *http.Client
}

func NewUserHandler(userService UserService, slogLogger *slog.Logger, client *http.Client) *UserHandler {
	return &UserHandler{
		UserService: userService,
		SlogLogger:  *slogLogger,
		Client:      client,
	}
}

func (uh *UserHandler) GetAllUsersHandler(c *gin.Context) {
	// Получаем requestID из контекста
	reqID := c.Request.Context().Value("requestID")
	reqIDString := fmt.Sprintf("requestID: %s", reqID)

	c.Set("requestID", reqIDString)
	users, err := uh.UserService.GetAllUsers(context.Background())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Возвращаем список пользователей в формате JSON
	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) AddUserHandler(c *gin.Context) {
	var req core.ServiceUser

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passportRegex := regexp.MustCompile(`^\d{2} \d{2} \d{6}$`)
	if !passportRegex.MatchString(req.PassportNum) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport number format"})
		return
	}

	user := core.ServiceUser{
		PassportNum: req.PassportNum,
		Surname:     req.Surname,
		Name:        req.Name,
		Patronymic:  req.Patronymic,
		Address:     req.Address,
	}

	userID, err := uh.UserService.AddUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, core.AddUserResponse{UserID: userID})
}

func (uh *UserHandler) GetUserByNumberHandler(c *gin.Context) {
	var req core.RequestUsersPassport

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := uh.UserService.GetUserByNumber(c.Request.Context(), req.PassportNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) UpdateUserHandler(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req core.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.PassportNumber != "" {
		passportRegex := regexp.MustCompile(`^\d{2} \d{2} \d{6}$`)
		if !passportRegex.MatchString(req.PassportNumber) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid passport number format"})
			return
		}
	}

	user, err := uh.UserService.UpdateUser(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteUserHandler(c *gin.Context) {
	userIDStr := c.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	err = uh.UserService.DeleteUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
