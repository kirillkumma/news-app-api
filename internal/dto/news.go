package dto

import (
	"mime/multipart"
)

type (
	CreateNewsParams struct {
		MediaID int64  `json:"-"`
		Title   string `json:"title"`
		Text    string `json:"text"`
	}

	CreateOrUpdateAudioParams struct {
		NewsID  int64 `params:"news_id"`
		MediaID int64
		File    multipart.File
	}

	GetAudioParams struct {
		NewsID int64 `params:"news_id"`
	}

	GetNewsParams struct {
		NewsID int64 `params:"news_id"`
	}

	CreateOrUpdateImageParams struct {
		NewsID  int64 `params:"news_id"`
		MediaID int64
		File    multipart.File
	}

	GetImageParams struct {
		NewsID int64 `params:"news_id"`
	}

	ToggleFavoriteParams struct {
		NewsID int64 `params:"news_id"`
		UserID int64
	}

	ToggleFavoriteResult struct {
		IsFavorite bool `json:"isFavorite"`
	}
)
