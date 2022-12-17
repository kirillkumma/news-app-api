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
		GetSubscriptionList(ctx context.Context, p dto.GetSubscriptionListParams) ([]entity.MediaListItem, error)
		CountSubscriptions(ctx context.Context, userID int64) (int64, error)
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

func (r *userRepository) GetSubscriptionList(
	ctx context.Context,
	p dto.GetSubscriptionListParams,
) (list []entity.MediaListItem, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - GetSubscriptionList: %w", err)
			}
		}
	}()
	list = make([]entity.MediaListItem, 0, p.Limit.Int64)
	rows, err := r.db.Query(ctx, queryGetSubscriptionList, p.UserID, p.Limit, p.Offset)
	if err != nil {
		return
	}
	for rows.Next() {
		item := entity.MediaListItem{}
		err = rows.Scan(
			&item.ID,
			&item.RegistrationNumber,
			&item.Name,
			&item.Email,
			&item.Editor.LastName,
			&item.Editor.FirstName,
			&item.SubscriptionCount,
		)
		if err != nil {
			return
		}
		list = append(list, item)
	}
	return
}

func (r *userRepository) CountSubscriptions(ctx context.Context, userID int64) (v int64, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserRepository - CountSubscriptions: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryCountSubscriptions, userID)
	err = row.Scan(&v)
	return
}
