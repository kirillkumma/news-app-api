package usecase

import (
	"context"
	"errors"
	"fmt"
	"news-app-api/internal/adapter"
	"news-app-api/internal/dto"
)

type (
	FeedUseCase interface {
		GetFeed(ctx context.Context, p dto.GetFeedParams) (dto.GetFeedResult, error)
	}

	feedUseCase struct {
		newsRepo func() adapter.NewsRepository
	}
)

func NewFeedUseCase(newsRepo func() adapter.NewsRepository) FeedUseCase {
	return &feedUseCase{newsRepo}
}

func (u *feedUseCase) GetFeed(ctx context.Context, p dto.GetFeedParams) (res dto.GetFeedResult, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("FeedUseCase - GetFeed: %w", err)
			}
		}
	}()

	r := u.newsRepo()

	res.Items, err = r.GetFeedNewsList(ctx, p)
	if err != nil {
		return
	}

	res.Total, err = r.CountFeedNews(ctx, p.UserID)
	return
}
