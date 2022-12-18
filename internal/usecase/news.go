package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"news-app-api/internal/adapter"
	"news-app-api/internal/dto"
	"news-app-api/internal/entity"
)

type (
	NewsUseCase interface {
		CreateNews(ctx context.Context, p dto.CreateNewsParams) (entity.News, error)
		CreateOrUpdateAudio(ctx context.Context, p dto.CreateOrUpdateAudioParams) error
		GetAudio(ctx context.Context, p dto.GetAudioParams) ([]byte, error)
		GetNews(ctx context.Context, p dto.GetNewsParams) (entity.NewsListItem, error)
	}

	newsUseCase struct {
		newsRepo      func() adapter.NewsRepository
		mediaRepo     adapter.MediaRepository
		audioFileRepo adapter.AudioFileRepository
	}
)

func NewNewsUseCase(
	newsRepo func() adapter.NewsRepository,
	mediaRepo adapter.MediaRepository,
	audioFileRepo adapter.AudioFileRepository,
) NewsUseCase {
	return &newsUseCase{newsRepo, mediaRepo, audioFileRepo}
}

func (u *newsUseCase) CreateNews(ctx context.Context, p dto.CreateNewsParams) (n entity.News, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - CreateNews: %w", err)
			}
		}
	}()

	r := u.newsRepo()

	err = r.Begin(ctx)
	if err != nil {
		return
	}
	defer r.Rollback(ctx)

	n, err = r.CreateNews(ctx, p)
	if err != nil {
		return
	}

	err = r.AddNewsToFeed(ctx, n.ID)
	if err != nil {
		return
	}

	err = r.Commit(ctx)
	return
}

func (u *newsUseCase) CreateOrUpdateAudio(ctx context.Context, p dto.CreateOrUpdateAudioParams) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - CreateOrUpdateAudio: %w", err)
			}
		}
	}()

	n, err := u.newsRepo().GetNews(ctx, p.NewsID)
	if err != nil {
		return
	}

	m, err := u.mediaRepo.GetMediaByRegistrationNumber(ctx, n.Media.RegistrationNumber)
	if err != nil {
		return
	}

	if m.ID != p.MediaID {
		return &dto.AppError{
			Code:    dto.ErrCodeUnauthorized,
			Message: "Недостаточно прав для совершения данной операции",
		}
	}

	data, err := io.ReadAll(p.File)
	if err != nil {
		return
	}

	err = u.audioFileRepo.Store(ctx, fmt.Sprint(p.NewsID), data)
	return
}

func (u *newsUseCase) GetAudio(ctx context.Context, p dto.GetAudioParams) (data []byte, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - GetAudio: %w", err)
			}
		}
	}()

	n, err := u.newsRepo().GetNews(ctx, p.NewsID)
	if err != nil {
		return
	}

	data, err = u.audioFileRepo.Get(ctx, fmt.Sprint(n.ID))
	return
}

func (u *newsUseCase) GetNews(ctx context.Context, p dto.GetNewsParams) (n entity.NewsListItem, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - GetNews: %w", err)
			}
		}
	}()

	n, err = u.newsRepo().GetNews(ctx, p.NewsID)
	return
}
