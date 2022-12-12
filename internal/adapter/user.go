package adapter

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"news-app-api/internal/dto"
	"news-app-api/internal/entity"
)

type (
	UserRepository interface {
		GetByLogin(ctx context.Context, login string) (entity.User, error)
		GetByEmail(ctx context.Context, email string) (entity.User, error)
		Create(ctx context.Context, p dto.RegisterParams) (entity.User, error)
	}

	userRepository struct {
		db *pgx.Conn
	}
)

func NewUserRepository(db *pgx.Conn) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByLogin(ctx context.Context, login string) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - GetByLogin: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetByLogin, login)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = &dto.AppError{
				Message: "Пользователь не найден",
				Code:    dto.ErrCodeNotFound,
			}
			return
		}
		return
	}
	return
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - GetByEmail: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetByEmail, email)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = &dto.AppError{
				Message: "Пользователь не найден",
				Code:    dto.ErrCodeNotFound,
			}
			return
		}
		return
	}
	return
}

func (r *userRepository) Create(ctx context.Context, p dto.RegisterParams) (u entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - Create: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryCreate, p.Login, p.Password, p.Name, p.Email)
	err = row.Scan(
		&u.ID,
		&u.Login,
		&u.Password,
		&u.Name,
		&u.Email,
	)
	return
}
