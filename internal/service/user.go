package service

import (
	"context"
	"fmt"
	"tt/internal/core"
)

type UserRepo interface {
	GetAllUsers(ctx context.Context) (core.Users, error)
	AddUser(ctx context.Context, user core.ServiceUser) (int, error)
	UpdateUser(ctx context.Context, userID int, newData core.ServiceUser) error
}

type UserService struct {
	UserRepository UserRepo
}

func NewUserService(repo UserRepo) *UserService {
	return &UserService{UserRepository: repo}
}

func (us *UserService) GetAllUsers(ctx context.Context) (core.Users, error) {
	const op = "service.GetAllUsers"
	users, err := us.UserRepository.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return users, nil
}

func (us *UserService) AddUser(ctx context.Context, user core.ServiceUser) (int, error) {
	const op = "service.AddUser"
	userID, err := us.UserRepository.AddUser(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (us *UserService) UpdateUser(ctx context.Context, userID int, newData core.ServiceUser) error {
	const op = "service.UpdateUser"
	err := us.UserRepository.UpdateUser(ctx, userID, newData)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
