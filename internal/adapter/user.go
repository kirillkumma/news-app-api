package adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"news-app-api/internal/dto"
	"news-app-api/internal/entity"
)

type (
	UserRepository interface {
		GetUserByLogin(ctx context.Context, login string) (entity.User, error)
		GetUserByEmail(ctx context.Context, email string) (entity.User, error)
		GetUserByID(ctx context.Context, userID int64) (entity.User, error)
		CreateUser(ctx context.Context, p dto.RegisterUserParams) (entity.User, error)
	}

	userRepository struct {
		db *pgx.Conn
	}
)

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetUserByLogin(ctx context.Context, login string) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - GetUserByLogin: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetUserByLogin, login)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "Пользователь не найден",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - GetUserByEmail: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetUserByEmail, email)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "Пользователь не найден",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}

func (r *userRepository) CreateUser(ctx context.Context, p dto.RegisterUserParams) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - CreateUser: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryCreateUser, p.Login, p.Password, p.Name, p.Email)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	return
}

func (r *userRepository) GetUserByID(ctx context.Context, userID int64) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - GetUserByID: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetUserByID, userID)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "Пользователь не найден",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}
