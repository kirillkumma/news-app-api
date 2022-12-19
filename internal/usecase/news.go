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
		CreateOrUpdateImage(ctx context.Context, p dto.CreateOrUpdateImageParams) error
		GetAudio(ctx context.Context, p dto.GetAudioParams) ([]byte, error)
		GetNews(ctx context.Context, p dto.GetNewsParams) (entity.NewsListItem, error)
		GetImage(ctx context.Context, p dto.GetImageParams) ([]byte, error)
		ToggleFavorite(ctx context.Context, p dto.ToggleFavoriteParams) (dto.ToggleFavoriteResult, error)
		GetFavoriteList(ctx context.Context, p dto.GetFavoriteListParams) (dto.GetFavoriteListResult, error)
	}

	newsUseCase struct {
		newsRepo      func() adapter.NewsRepository
		mediaRepo     adapter.MediaRepository
		audioFileRepo adapter.AudioFileRepository
		imageFileRepo adapter.ImageFileRepository
	}
)

func NewNewsUseCase(
	newsRepo func() adapter.NewsRepository,
	mediaRepo adapter.MediaRepository,
	audioFileRepo adapter.AudioFileRepository,
	imageFileRepo adapter.ImageFileRepository,
) NewsUseCase {
	return &newsUseCase{newsRepo, mediaRepo, audioFileRepo, imageFileRepo}
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

	err = u.audioFileRepo.Store(ctx, fmt.Sprintf("%d.wav", p.NewsID), data)
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

	data, err = u.audioFileRepo.Get(ctx, fmt.Sprintf("%d.wav", n.ID))
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

func (u *newsUseCase) CreateOrUpdateImage(ctx context.Context, p dto.CreateOrUpdateImageParams) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - CreateOrUpdateImage: %w", err)
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

	return u.imageFileRepo.Store(ctx, fmt.Sprintf("%d.png", p.NewsID), data)
}

func (u *newsUseCase) GetImage(ctx context.Context, p dto.GetImageParams) (data []byte, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - GetImage: %w", err)
			}
		}
	}()

	n, err := u.newsRepo().GetNews(ctx, p.NewsID)
	if err != nil {
		return
	}

	return u.imageFileRepo.Get(ctx, fmt.Sprintf("%d.png", n.ID))
}

func (u *newsUseCase) ToggleFavorite(
	ctx context.Context,
	p dto.ToggleFavoriteParams,
) (res dto.ToggleFavoriteResult, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - ToggleFavorite: %w", err)
			}
		}
	}()

	r := u.newsRepo()

	isFavorite, err := r.IsFavorite(ctx, p.UserID, p.NewsID)
	if err != nil {
		return
	}

	if !isFavorite {
		err = r.AddToFavorite(ctx, p.UserID, p.NewsID)
	} else {
		err = r.RemoveFromFavorite(ctx, p.UserID, p.NewsID)
	}
	if err != nil {
		return
	}

	res.IsFavorite = !isFavorite

	return
}

func (u *newsUseCase) GetFavoriteList(
	ctx context.Context,
	p dto.GetFavoriteListParams,
) (res dto.GetFavoriteListResult, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsUseCase - GetFavoriteList: %w", err)
			}
		}
	}()

	r := u.newsRepo()

	res.Items, err = r.GetFavoriteList(ctx, p)
	if err != nil {
		return
	}

	res.Total, err = r.CountFavorites(ctx, p.UserID)
	return
}
