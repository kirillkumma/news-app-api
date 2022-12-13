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
	MediaUseCase interface {
		Register(ctx context.Context, p dto.RegisterMediaParams) (entity.Media, error)
		LoginMedia(ctx context.Context, p dto.LoginMediaParams) (entity.Media, error)
		GetMediaByID(ctx context.Context, mediaID int64) (entity.Media, error)
	}

	mediaUseCase struct {
		mediaRepo adapter.MediaRepository
	}
)

func NewMediaUseCase(mediaRepo adapter.MediaRepository) MediaUseCase {
	return &mediaUseCase{mediaRepo}
}

func (u *mediaUseCase) Register(ctx context.Context, p dto.RegisterMediaParams) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaUseCase - RegisterUser: %w", err)
			}
		}
	}()

	err = p.Validate()
	if err != nil {
		return
	}

	_, err = u.mediaRepo.GetMediaByEmail(ctx, p.Email)
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

	_, err = u.mediaRepo.GetMediaByName(ctx, p.Name)
	if err == nil {
		err = &dto.AppError{
			Message: "Название уже занято",
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

	_, err = u.mediaRepo.GetMediaByRegistrationNumber(ctx, p.RegistrationNumber)
	if err == nil {
		err = &dto.AppError{
			Message: "Регистрационный номер уже занят",
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

	return u.mediaRepo.CreateMedia(ctx, p)
}

func (u *mediaUseCase) LoginMedia(ctx context.Context, p dto.LoginMediaParams) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaUseCase - LoginMedia: %w", err)
			}
		}
	}()

	m, err = u.mediaRepo.GetMediaByRegistrationNumber(ctx, p.RegistrationNumber)
	if err != nil {
		return
	}

	if p.Password != m.Password {
		err = &dto.AppError{
			Message: "Неверный пароль",
			Code:    dto.ErrCodeUnauthorized,
		}
		return
	}

	return
}

func (u *mediaUseCase) GetMediaByID(ctx context.Context, mediaID int64) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaUseCase - GetMediaByID: %w", err)
			}
		}
	}()
	return u.mediaRepo.GetMediaByID(ctx, mediaID)
}
