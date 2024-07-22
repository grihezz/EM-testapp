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
	baseQuery := "SELECT id, passport_number, surname, name, patronymic, address FROM users ORDER BY id;"

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
	baseQuery := "INSERT INTO users (passport_number, surname, name, patronymic, address) VALUES ($1, $2, $3, $4, $5) RETURNING id;"
	var userID int

	err := ur.db.QueryRow(ctx, baseQuery, user.PassportNum, user.Surname, user.Name, user.Patronymic, user.Address).Scan(&userID)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return userID, nil
}

func (ur *UserRepository) GetUserByID(ctx context.Context, userID int) (core.ServiceUser, error) {
	const op = "repository.GetUserByID"
	baseQuery := "SELECT passport_number, surname, name, patronymic, address FROM users WHERE id = $1;"
	var user core.ServiceUser

	err := ur.db.QueryRow(ctx, baseQuery, userID).Scan(
		&user.PassportNum,
		&user.Surname,
		&user.Name,
		&user.Patronymic,
		&user.Address,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return core.ServiceUser{}, ErrUsrNotExists
		}
		return core.ServiceUser{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (ur *UserRepository) UpdateUser(ctx context.Context, userID int, newData core.ServiceUser) error {
	const op = "repository.UpdateUser"
	baseQuery := "UPDATE users SET "
	var args []interface{}
	var setClauses []string
	paramIndex := 1

	if newData.PassportNum != "" {
		setClauses = append(setClauses, "passport_number = $"+strconv.Itoa(paramIndex))
		args = append(args, newData.PassportNum)
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

	// Ensure we have at least one field to update
	if len(setClauses) == 0 {
		err := errors.New("no fields to update")
		return fmt.Errorf("%s: %w", op, err)
	}

	// Combine the SET clauses into the query
	baseQuery += strings.Join(setClauses, ", ")
	baseQuery += " WHERE id = $" + strconv.Itoa(paramIndex)
	args = append(args, userID)

	_, err := ur.db.Exec(ctx, baseQuery, args...)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}