package adapter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"news-app-api/internal/dto"
	"os"
	"strings"
)

type (
	VideoFileRepository interface {
		Store(ctx context.Context, filename string, data []byte) error
		Get(ctx context.Context, filename string) ([]byte, error)
	}

	videoFileRepository struct {
	}
)

func NewVideoFileRepository() (VideoFileRepository, error) {
	err := os.Mkdir("video", 0750)
	if err != nil && !strings.Contains(err.Error(), "exists") {
		return nil, err
	}

	return &videoFileRepository{}, nil
}

func (v videoFileRepository) Store(ctx context.Context, filename string, data []byte) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("VideoFileRepository - Store: %w", err)
			}
		}
	}()
	f, err := os.OpenFile(fmt.Sprintf("video/%s", filename), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	return
}

func (v videoFileRepository) Get(ctx context.Context, filename string) (data []byte, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("VideoFileRepository - Get: %w", err)
			}
		}
	}()
	f, err := os.Open(fmt.Sprintf("video/%s", filename))
	if err != nil {
		fmt.Println(err)
		err = &dto.AppError{
			Code:    dto.ErrCodeNotFound,
			Message: "Файл не найден",
		}
		return
	}
	defer f.Close()
	data, err = io.ReadAll(f)
	return
}
