package service

import (
	"context"
	"fmt"
	"tt/internal/core"
)

type UserRepo interface {
	GetAllUsers(ctx context.Context) (core.Users, error)
	AddUser(ctx context.Context, user core.ServiceUser) (int, error)
	UpdateUser(ctx context.Context, userID int, newData core.User) (core.User, error)
	GetUserByNumber(ctx context.Context, passportNumber string) (core.User, error)
	DeleteUser(ctx context.Context, userID int) error
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

func (us *UserService) UpdateUser(ctx context.Context, userID int, newData core.User) (core.User, error) {
	const op = "service.UpdateUser"
	user, err := us.UserRepository.UpdateUser(ctx, userID, newData)
	if err != nil {
		return core.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (us *UserService) GetUserByNumber(ctx context.Context, passportNumber string) (core.User, error) {
	const op = "service.GetUserByNumber"
	user, err := us.UserRepository.GetUserByNumber(ctx, passportNumber)
	if err != nil {
		return core.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (us *UserService) DeleteUser(ctx context.Context, userID int) error {
	const op = "service.DeleteUser"
	err := us.UserRepository.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
