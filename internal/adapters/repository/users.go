package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"strings"
	"tt/internal/core"
)

var ErrUsrNotExists = errors.New("user not exists")

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (ur *UserRepository) GetAllUsers(ctx context.Context) (core.Users, error) {
	const op = "repository.GetAllUsers"
	baseQuery := "SELECT user_id, passport_number, surname, name, patronymic, address FROM users ORDER BY user_id;"

	rows, err := ur.db.Query(ctx, baseQuery)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var users core.Users
	for rows.Next() {
		var user core.User

		err := rows.Scan(
			&user.ID,
			&user.PassportNumber,
			&user.Surname,
			&user.Name,
			&user.Patronymic,
			&user.Address,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return users, nil
}

func (ur *UserRepository) AddUser(ctx context.Context, user core.ServiceUser) (int, error) {
	const op = "repository.AddUser"
	baseQuery := "INSERT INTO users (passport_number, surname, name, patronymic, address) VALUES ($1, $2, $3, $4, $5) RETURNING user_id;"
	var userID int

	err := ur.db.QueryRow(ctx, baseQuery, user.PassportNum, user.Surname, user.Name, user.Patronymic, user.Address).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (ur *UserRepository) GetUserByNumber(ctx context.Context, passportNumber string) (core.User, error) {
	const op = "repository.GetUserByNumber"
	baseQuery := "SELECT user_id,passport_number, surname, name, patronymic, address FROM users WHERE passport_number = $1;"
	var user core.User

	err := ur.db.QueryRow(ctx, baseQuery, passportNumber).Scan(
		&user.ID,
		&user.PassportNumber,
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.Address,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return core.User{}, ErrUsrNotExists
		}
		return core.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, userID int, newData core.User) (core.User, error) {
	const op = "repository.UpdateUser"
	var user core.User
	baseQuery := "UPDATE users SET "
	var args []interface{}
	var setClauses []string
	paramIndex := 1

	if newData.PassportNumber != "" {
		setClauses = append(setClauses, "passport_number = $"+strconv.Itoa(paramIndex))
		args = append(args, newData.PassportNumber)
		paramIndex++
	}

	if newData.Surname != "" {
		setClauses = append(setClauses, "surname = $"+strconv.Itoa(paramIndex))
		args = append(args, newData.Surname)
		paramIndex++
	}

	if newData.Name != "" {
		setClauses = append(setClauses, "name = $"+strconv.Itoa(paramIndex))
		args = append(args, newData.Name)
		paramIndex++
	}

	if newData.Patronymic != "" {
		setClauses = append(setClauses, "patronymic = $"+strconv.Itoa(paramIndex))
		args = append(args, newData.Patronymic)
		paramIndex++
	}

	if newData.Address != "" {
		setClauses = append(setClauses, "address = $"+strconv.Itoa(paramIndex))
		args = append(args, newData.Address)
		paramIndex++
	}

	if len(setClauses) == 0 {
		err := errors.New("no fields to update")
		return core.User{}, fmt.Errorf("%s: %w", op, err)
	}

	baseQuery += strings.Join(setClauses, ", ")
	baseQuery += " WHERE id = $" + strconv.Itoa(paramIndex)
	args = append(args, userID)

	_, err := ur.db.Exec(ctx, baseQuery, args...)
	if err != nil {
		return core.User{}, fmt.Errorf("%s: %w", op, err)
	}

	err = ur.db.QueryRow(ctx, "SELECT user_id,passport_number, surname, name, patronymic, address FROM users WHERE user_id = $1", userID).
		Scan(&user.ID, &user.PassportNumber, &user.Surname, &user.Name, &user.Patronymic, &user.Address)

	if err != nil {
		return core.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (ur *UserRepository) DeleteUser(ctx context.Context, userID int) error {
	const op = "repository.DeleteUser"
	baseURL := "DELETE FROM users WHERE user_id = $1"

	result, err := ur.db.Exec(ctx, baseURL, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrUsrNotExists
	}

	return nil
}
