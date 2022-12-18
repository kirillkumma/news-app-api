package adapter

import (
	"context"
	"errors"
	"fmt"
	"io"
	"news-app-api/internal/dto"
	"os"
)

type (
	AudioFileRepository interface {
		Store(ctx context.Context, filename string, data []byte) error
		Get(ctx context.Context, filename string) ([]byte, error)
	}

	audioFileRepository struct {
	}
)

func NewAudioFileRepository() AudioFileRepository {
	return &audioFileRepository{}
}

func (a audioFileRepository) Store(ctx context.Context, filename string, data []byte) (err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("AudioFileRepository - Store: %w", err)
			}
		}
	}()
	f, err := os.OpenFile(fmt.Sprintf("audio/%s", filename), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return
	}
	defer f.Close()
	_, err = f.Write(data)
	return
}

func (a audioFileRepository) Get(ctx context.Context, filename string) (data []byte, err error) {
	defer func() {
		if err != nil {
			var appErr *dto.AppError
			if !errors.As(err, &appErr) {
				err = fmt.Errorf("AudioFileRepository - Get: %w", err)
			}
		}
	}()
	f, err := os.Open(fmt.Sprintf("audio/%s", filename))
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
