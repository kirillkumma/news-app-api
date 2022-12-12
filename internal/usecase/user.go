package usecase

import (
	"context"
	"errors"
	"fmt"
	"news-app-api/internal/adapter"
	"news-app-api/internal/dto"
	"news-app-api/internal/entity"
)

type (
	UserUseCase interface {
		Register(ctx context.Context, p dto.RegisterParams) (entity.User, error)
		Login(ctx context.Context, p dto.LoginParams) (entity.User, error)
	}

	userUseCase struct {
		userRepo adapter.UserRepository
	}
)

func NewUserUseCase(userRepo adapter.UserRepository) UserUseCase {
	return &userUseCase{userRepo}
}

func (u *userUseCase) Register(ctx context.Context, p dto.RegisterParams) (user entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserUseCase - Register: %w", err)
			}
		}
	}()

	err = p.Validate()
	if err != nil {
		return
	}

	_, err = u.userRepo.GetByLogin(ctx, p.Login)
	if err == nil {
		err = &dto.AppError{
			Message: "Логин уже занят",
			Code:    dto.ErrCodeConflict,
		}
		return
	}
	if err != nil {
		var appErr *dto.AppError
		if !errors.As(err, &appErr) || appErr.Code != dto.ErrCodeNotFound {
			return
		}
	}

	_, err = u.userRepo.GetByEmail(ctx, p.Email)
	if err == nil {
		err = &dto.AppError{
			Message: "Email уже занят",
			Code:    dto.ErrCodeConflict,
		}
		return
	}
	if err != nil {
		var appErr *dto.AppError
		if !errors.As(err, &appErr) || appErr.Code != dto.ErrCodeNotFound {
			return
		}
	}

	return u.userRepo.Create(ctx, p)
}

func (u *userUseCase) Login(ctx context.Context, p dto.LoginParams) (user entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserUseCase - Login: %w", err)
			}
		}
	}()

	user, err = u.userRepo.GetByLogin(ctx, p.Login)
	if err != nil {
		return
	}

	if user.Password != p.Password {
		err = &dto.AppError{
			Message: "Неверный пароль",
			Code:    dto.ErrCodeUnauthorized,
		}
		return
	}

	return
}
