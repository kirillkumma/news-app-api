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
	MediaRepository interface {
		GetMediaByRegistrationNumber(ctx context.Context, registrationNumber int64) (entity.Media, error)
		GetMediaByName(ctx context.Context, name string) (entity.Media, error)
		GetMediaByEmail(ctx context.Context, email string) (entity.Media, error)
		CreateMedia(ctx context.Context, p dto.RegisterMediaParams) (entity.Media, error)
		GetMediaByID(ctx context.Context, mediaID int64) (entity.Media, error)
		GetMediaList(ctx context.Context, p dto.GetMediaListParams) ([]entity.MediaListItem, error)
		CountMedia(ctx context.Context) (int64, error)
		IsSubscriptionExists(ctx context.Context, mediaID, userID int64) (bool, error)
		CreateSubscription(ctx context.Context, mediaID, userID int64) error
		DeleteSubscription(ctx context.Context, mediaID, userID int64) error
	}

	mediaRepository struct {
		db *pgx.Conn
	}
)

func NewMediaRepository(db *pgx.Conn) MediaRepository {
	return &mediaRepository{db}
}

func (r *mediaRepository) GetMediaByRegistrationNumber(
	ctx context.Context,
	registrationNumber int64,
) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - GetMediaByRegistrationNumber: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetMediaByRegistrationNumber, registrationNumber)
	err = row.Scan(
		&m.ID,
		&m.RegistrationNumber,
		&m.Name,
		&m.Email,
		&m.Editor.LastName,
		&m.Editor.FirstName,
		&m.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "СМИ не найдено",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}

func (r *mediaRepository) GetMediaByName(ctx context.Context, name string) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - GetMediaByName: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetMediaByName, name)
	err = row.Scan(
		&m.ID,
		&m.RegistrationNumber,
		&m.Name,
		&m.Email,
		&m.Editor.LastName,
		&m.Editor.FirstName,
		&m.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "СМИ не найдено",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}

func (r *mediaRepository) GetMediaByEmail(ctx context.Context, email string) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - GetMediaByEmail: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetMediaByEmail, email)
	err = row.Scan(
		&m.ID,
		&m.RegistrationNumber,
		&m.Name,
		&m.Email,
		&m.Editor.LastName,
		&m.Editor.FirstName,
		&m.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "СМИ не найдено",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}

func (r *mediaRepository) CreateMedia(ctx context.Context, p dto.RegisterMediaParams) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - CreateMedia: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(
		ctx,
		queryCreateMedia,
		p.RegistrationNumber,
		p.Name,
		p.Email,
		p.Editor.LastName,
		p.Editor.FirstName,
		p.Password,
	)
	err = row.Scan(
		&m.ID,
		&m.RegistrationNumber,
		&m.Name,
		&m.Email,
		&m.Editor.LastName,
		&m.Editor.FirstName,
		&m.Password,
	)
	return
}

func (r *mediaRepository) GetMediaByID(ctx context.Context, mediaID int64) (m entity.Media, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - GetMediaByID: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryGetMediaByID, mediaID)
	err = row.Scan(
		&m.ID,
		&m.RegistrationNumber,
		&m.Name,
		&m.Email,
		&m.Editor.LastName,
		&m.Editor.FirstName,
		&m.Password,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			err = &dto.AppError{
				Message: "СМИ не найдено",
				Code:    dto.ErrCodeNotFound,
			}
		}
		return
	}
	return
}

func (r *mediaRepository) GetMediaList(
	ctx context.Context,
	p dto.GetMediaListParams,
) (list []entity.MediaListItem, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - GetMediaList: %w", err)
			}
		}
	}()
	list = make([]entity.MediaListItem, 0, p.Limit.Int64)
	rows, err := r.db.Query(ctx, queryGetMediaList, p.Limit, p.Offset)
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

func (r *mediaRepository) CountMedia(ctx context.Context) (v int64, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - CountMedia: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryCountMedia)
	err = row.Scan(&v)
	return
}

func (r *mediaRepository) IsSubscriptionExists(ctx context.Context, mediaID, userID int64) (v bool, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - IsSubscriptionExists: %w", err)
			}
		}
	}()
	row := r.db.QueryRow(ctx, queryIsSubscriptionExists, mediaID, userID)
	err = row.Scan(&v)
	return
}

func (r *mediaRepository) CreateSubscription(ctx context.Context, mediaID, userID int64) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - CreateSubscription: %w", err)
			}
		}
	}()
	_, err = r.db.Exec(ctx, queryCreateSubscription, mediaID, userID)
	return
}

func (r *mediaRepository) DeleteSubscription(ctx context.Context, mediaID, userID int64) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("MediaRepository - DeleteSubscription: %w", err)
			}
		}
	}()
	_, err = r.db.Exec(ctx, queryDeleteSubscription, mediaID, userID)
	return
}
