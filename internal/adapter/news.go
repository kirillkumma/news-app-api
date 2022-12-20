package adapter

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/guregu/null.v3"
	"news-app-api/internal/dto"
	"news-app-api/internal/entity"
)

type (
	NewsRepository interface {
		Transactor
		CreateNews(ctx context.Context, p dto.CreateNewsParams) (entity.News, error)
		AddNewsToFeed(ctx context.Context, newsID int64) error
		GetNews(ctx context.Context, newsID int64) (entity.NewsListItem, error)
		GetFeedNewsList(ctx context.Context, p dto.GetFeedParams) ([]entity.NewsListItem, error)
		CountFeedNews(ctx context.Context, userID int64, since null.Int) (int64, error)
		IsFavorite(ctx context.Context, userID, newsID int64) (bool, error)
		AddToFavorite(ctx context.Context, userID, newsID int64) error
		RemoveFromFavorite(ctx context.Context, userID, newsID int64) error
		GetFavoriteList(ctx context.Context, p dto.GetFavoriteListParams) ([]entity.NewsListItem, error)
		CountFavorites(ctx context.Context, userID int64) (int64, error)
	}

	newsRepository struct {
		db *pgxpool.Pool
		q  Querier
	}
)

func NewNewsRepository(db *pgxpool.Pool) NewsRepository {
	return &newsRepository{db, db}
}

func (r *newsRepository) Begin(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - Begin: %w", err)
			}
		}
	}()
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return
	}
	r.q = tx
	return
}

func (r *newsRepository) Commit(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - Commit: %w", err)
			}
		}
	}()
	tx, ok := r.q.(pgx.Tx)
	if !ok {
		return ErrTxNotStarted
	}
	return tx.Commit(ctx)
}

func (r *newsRepository) Rollback(ctx context.Context) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - Rollback: %w", err)
			}
		}
	}()
	tx, ok := r.q.(pgx.Tx)
	if !ok {
		return ErrTxNotStarted
	}
	return tx.Rollback(ctx)
}

func (r *newsRepository) CreateNews(ctx context.Context, p dto.CreateNewsParams) (n entity.News, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - CreateNews: %w", err)
			}
		}
	}()
	row := r.q.QueryRow(ctx, queryCreateNews, p.MediaID, p.Title, p.Text)
	err = row.Scan(
		&n.ID,
		&n.MediaRegistrationNumber,
		&n.Title,
		&n.Text,
		&n.CreatedAt,
	)
	return
}

func (r *newsRepository) AddNewsToFeed(ctx context.Context, newsID int64) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - AddNewsToFeed: %w", err)
			}
		}
	}()
	_, err = r.q.Exec(ctx, queryAddNewsToFeed, newsID)
	return
}

func (r *newsRepository) GetNews(ctx context.Context, newsID int64) (n entity.NewsListItem, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - GetNews: %w", err)
			}
		}
	}()
	row := r.q.QueryRow(ctx, queryGetNews, newsID)
	err = row.Scan(
		&n.ID,
		&n.Media.ID,
		&n.Media.RegistrationNumber,
		&n.Media.Name,
		&n.Media.Email,
		&n.Media.Editor.FirstName,
		&n.Media.Editor.LastName,
		&n.Media.SubscriptionCount,
		&n.Title,
		&n.Text,
		&n.CreatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Code:    dto.ErrCodeNotFound,
				Message: "Новость не найдена",
			}
		}
		return
	}
	return
}

func (r *newsRepository) GetFeedNewsList(
	ctx context.Context,
	p dto.GetFeedParams,
) (list []entity.NewsListItem, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - GetFeedNewsList: %w", err)
			}
		}
	}()
	list = make([]entity.NewsListItem, 0, p.Limit.Int64)
	rows, err := r.q.Query(ctx, queryGetFeedNewsList, p.UserID, p.Since, p.Limit, p.Offset)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := entity.NewsListItem{}
		err = rows.Scan(
			&item.ID,
			&item.Media.ID,
			&item.Media.RegistrationNumber,
			&item.Media.Name,
			&item.Media.Email,
			&item.Media.Editor.FirstName,
			&item.Media.Editor.LastName,
			&item.Media.SubscriptionCount,
			&item.Title,
			&item.Text,
			&item.IsFavorite,
			&item.CreatedAt,
		)
		if err != nil {
			return
		}
		list = append(list, item)
	}
	return
}

func (r *newsRepository) CountFeedNews(ctx context.Context, userID int64, since null.Int) (v int64, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - CountFeedNews: %w", err)
			}
		}
	}()
	row := r.q.QueryRow(ctx, queryCountFeedNews, userID, since)
	err = row.Scan(&v)
	return
}

func (r *newsRepository) IsFavorite(ctx context.Context, userID, newsID int64) (v bool, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - IsFavorite: %w", err)
			}
		}
	}()
	row := r.q.QueryRow(ctx, queryIsFavorite, userID, newsID)
	err = row.Scan(&v)
	return
}

func (r *newsRepository) AddToFavorite(ctx context.Context, userID, newsID int64) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - AddToFavorite: %w", err)
			}
		}
	}()
	_, err = r.q.Exec(ctx, queryAddToFavorite, userID, newsID)
	return
}

func (r *newsRepository) RemoveFromFavorite(ctx context.Context, userID, newsID int64) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - RemoveFromFavorite: %w", err)
			}
		}
	}()
	_, err = r.q.Exec(ctx, queryRemoveFromFavorite, userID, newsID)
	return
}

func (r *newsRepository) GetFavoriteList(
	ctx context.Context,
	p dto.GetFavoriteListParams,
) (list []entity.NewsListItem, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - GetFavoriteList: %w", err)
			}
		}
	}()
	list = make([]entity.NewsListItem, 0, p.Limit.Int64)
	rows, err := r.q.Query(ctx, queryGetFavoriteList, p.UserID, p.Limit, p.Offset)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := entity.NewsListItem{}
		err = rows.Scan(
			&item.ID,
			&item.Media.ID,
			&item.Media.RegistrationNumber,
			&item.Media.Name,
			&item.Media.Email,
			&item.Media.Editor.FirstName,
			&item.Media.Editor.LastName,
			&item.Media.SubscriptionCount,
			&item.Title,
			&item.Text,
			&item.IsFavorite,
			&item.CreatedAt,
		)
		if err != nil {
			return
		}
		list = append(list, item)
	}
	return
}

func (r *newsRepository) CountFavorites(ctx context.Context, userID int64) (v int64, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("NewsRepository - CountFavorites: %w", err)
			}
		}
	}()
	row := r.q.QueryRow(ctx, queryCountFavorites, userID)
	err = row.Scan(&v)
	return
}
