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
		RegisterUser(ctx context.Context, p dto.RegisterUserParams) (entity.User, error)
		LoginUser(ctx context.Context, p dto.LoginUserParams) (entity.User, error)
		GetUserByID(ctx context.Context, userID int64) (entity.User, error)
		GetSubscriptionList(ctx context.Context, p dto.GetSubscriptionListParams) (dto.GetSubscriptionListResult, error)
	}

	userUseCase struct {
		userRepo adapter.UserRepository
	}
)

func NewUserUseCase(userRepo adapter.UserRepository) UserUseCase {
	return &userUseCase{userRepo}
}

func (u *userUseCase) RegisterUser(ctx context.Context, p dto.RegisterUserParams) (user entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserUseCase - RegisterUser: %w", err)
			}
		}
	}()

	err = p.Validate()
	if err != nil {
		return
	}

	_, err = u.userRepo.GetUserByLogin(ctx, p.Login)
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

	_, err = u.userRepo.GetUserByEmail(ctx, p.Email)
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

	return u.userRepo.CreateUser(ctx, p)
}

func (u *userUseCase) LoginUser(ctx context.Context, p dto.LoginUserParams) (user entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserUseCase - LoginUser: %w", err)
			}
		}
	}()

	user, err = u.userRepo.GetUserByLogin(ctx, p.Login)
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

func (u *userUseCase) GetUserByID(ctx context.Context, userID int64) (user entity.User, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserUseCase - GetUserByID: %w", err)
			}
		}
	}()
	return u.userRepo.GetUserByID(ctx, userID)
}

func (u *userUseCase) GetSubscriptionList(
	ctx context.Context,
	p dto.GetSubscriptionListParams,
) (res dto.GetSubscriptionListResult, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("UserUseCase - GetSubscriptionList: %w", err)
			}
		}
	}()

	res.Items, err = u.userRepo.GetSubscriptionList(ctx, p)
	if err != nil {
		return
	}

	res.Total, err = u.userRepo.CountSubscriptions(ctx, p.UserID)

	return
}
